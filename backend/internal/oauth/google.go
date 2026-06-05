// Package oauth wraps Google OAuth2 for customer sign-in.
package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/config"
)

type GoogleProfile struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type Google struct {
	cfg     *oauth2.Config
	enabled bool
}

func NewGoogle(c config.GoogleOAuthConfig) *Google {
	if !c.Enabled() {
		return &Google{enabled: false}
	}
	return &Google{
		enabled: true,
		cfg: &oauth2.Config{
			ClientID:     c.ClientID,
			ClientSecret: c.ClientSecret,
			RedirectURL:  c.RedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func (g *Google) Enabled() bool { return g.enabled }

func (g *Google) AuthURL(state string) string {
	return g.cfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// Exchange swaps the auth code for a token and fetches the user's profile.
func (g *Google) Exchange(ctx context.Context, code string) (*GoogleProfile, error) {
	if !g.enabled {
		return nil, fmt.Errorf("google oauth disabled")
	}
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	tok, err := g.cfg.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	client := g.cfg.Client(ctx, tok)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google userinfo status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var p GoogleProfile
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, err
	}
	return &p, nil
}
