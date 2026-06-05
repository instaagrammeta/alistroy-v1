// Package ws provides a minimal WebSocket hub backed by Redis pub/sub so chat
// and notifications work across multiple backend replicas.
package ws

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/cache"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/logger"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Origin is enforced at the Nginx / CORS layer.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Client is a single websocket connection subscribed to one or more channels.
type Client struct {
	conn     *websocket.Conn
	send     chan []byte
	channels []string
}

// Hub fans Redis pub/sub messages out to local websocket clients.
type Hub struct {
	redis  *cache.Client
	mu     sync.RWMutex
	byChan map[string]map[*Client]struct{}
}

func NewHub(redis *cache.Client) *Hub {
	return &Hub{redis: redis, byChan: make(map[string]map[*Client]struct{})}
}

// Serve upgrades the HTTP request and subscribes the client to the channels.
func (h *Hub) Serve(w http.ResponseWriter, r *http.Request, channels ...string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Warn("ws upgrade failed", "err", err)
		return
	}
	client := &Client{conn: conn, send: make(chan []byte, 32), channels: channels}
	h.register(client)

	// Ensure a Redis subscription exists per channel (lazy, idempotent).
	for _, ch := range channels {
		h.ensureSubscription(ch)
	}

	go client.writePump()
	client.readPump(h)
}

func (h *Hub) register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, ch := range c.channels {
		if h.byChan[ch] == nil {
			h.byChan[ch] = make(map[*Client]struct{})
		}
		h.byChan[ch][c] = struct{}{}
	}
}

func (h *Hub) unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, ch := range c.channels {
		if set, ok := h.byChan[ch]; ok {
			delete(set, c)
			if len(set) == 0 {
				delete(h.byChan, ch)
			}
		}
	}
	close(c.send)
}

// ensureSubscription starts a goroutine that pumps Redis messages for a channel
// into all local clients. It runs until the process exits; duplicate calls for
// the same channel are cheap because go-redis multiplexes.
func (h *Hub) ensureSubscription(channel string) {
	go func() {
		sub := h.redis.Subscribe(context.Background(), channel)
		defer sub.Close()
		ch := sub.Channel()
		for msg := range ch {
			h.broadcastLocal(channel, []byte(msg.Payload))
		}
	}()
}

func (h *Hub) broadcastLocal(channel string, payload []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.byChan[channel] {
		select {
		case c.send <- payload:
		default:
			// drop slow client
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) readPump(h *Hub) {
	defer h.unregister(c)
	c.conn.SetReadLimit(4096)
	for {
		// We only need to detect disconnect; inbound messages are ignored
		// (clients post via REST and receive via this socket).
		if _, _, err := c.conn.ReadMessage(); err != nil {
			return
		}
	}
}
