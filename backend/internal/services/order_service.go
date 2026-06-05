package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/repositories"
)

type OrderService struct {
	orders    *repositories.OrderRepository
	products  *repositories.ProductRepository
	customers *repositories.CustomerRepository
	drivers   *repositories.DriverRepository
	txs       *repositories.TransactionRepository
	cart      *repositories.CartRepository
	notifier  *NotificationService
}

func NewOrderService(
	o *repositories.OrderRepository,
	p *repositories.ProductRepository,
	c *repositories.CustomerRepository,
	d *repositories.DriverRepository,
	t *repositories.TransactionRepository,
	cart *repositories.CartRepository,
	n *NotificationService,
) *OrderService {
	return &OrderService{orders: o, products: p, customers: c, drivers: d, txs: t, cart: cart, notifier: n}
}

type OrderItemInput struct {
	ProductID uuid.UUID
	Quantity  int
}

type OrderInput struct {
	CustomerID      *uuid.UUID
	CustomerName    string
	CustomerPhone   string
	DeliveryAddress string
	DeliveryDate    *time.Time
	DiscountPercent float64
	Notes           string
	Status          string
	Items           []OrderItemInput
}

// Create builds an order from explicit items (admin manual create or API).
func (s *OrderService) Create(ctx context.Context, in OrderInput) (*models.Order, error) {
	if len(in.Items) == 0 {
		return nil, ErrValidation
	}
	order := &models.Order{
		CustomerID:      in.CustomerID,
		CustomerName:    in.CustomerName,
		CustomerPhone:   in.CustomerPhone,
		DeliveryAddress: in.DeliveryAddress,
		DeliveryDate:    in.DeliveryDate,
		DiscountPercent: in.DiscountPercent,
		Notes:           in.Notes,
		Status:          defaultStr(in.Status, models.OrderStatusNew),
		Currency:        "TJS",
	}
	// Resolve customer snapshot (fallback delivery address from profile).
	if in.CustomerID != nil {
		if cust, err := s.customers.FindByID(ctx, *in.CustomerID); err == nil {
			order.DeliveryAddress = firstNonEmpty(order.DeliveryAddress, cust.Address)
		}
	}

	var items []models.OrderItem
	for _, it := range in.Items {
		if it.Quantity <= 0 {
			continue
		}
		p, err := s.products.FindByID(ctx, it.ProductID)
		if err != nil {
			return nil, ErrValidation
		}
		line := p.SalePrice * float64(it.Quantity)
		profit := (p.SalePrice - p.CostPrice) * float64(it.Quantity)
		items = append(items, models.OrderItem{
			ProductID: p.ID,
			NameSnap:  p.NameRU,
			Unit:      p.Unit,
			Quantity:  it.Quantity,
			CostPrice: p.CostPrice,
			SalePrice: p.SalePrice,
			LineTotal: line,
			Profit:    profit,
		})
	}
	if len(items) == 0 {
		return nil, ErrValidation
	}
	computeTotals(order, items)
	order.Items = items

	seq, err := s.orders.NextSequenceForYear(ctx, time.Now().UTC().Year())
	if err != nil {
		return nil, err
	}
	order.Number = fmt.Sprintf("ALS-%d-%04d", time.Now().UTC().Year(), seq)

	if err := s.orders.Create(ctx, order); err != nil {
		return nil, err
	}
	return s.orders.FindByID(ctx, order.ID)
}

// CreateFromCart turns a customer's cart into an order then clears the cart.
func (s *OrderService) CreateFromCart(ctx context.Context, customerID uuid.UUID, deliveryAddr, notes string, deliveryDate *time.Time) (*models.Order, error) {
	cartItems, err := s.cart.List(ctx, customerID)
	if err != nil {
		return nil, err
	}
	if len(cartItems) == 0 {
		return nil, ErrValidation
	}
	in := OrderInput{
		CustomerID:      &customerID,
		DeliveryAddress: deliveryAddr,
		DeliveryDate:    deliveryDate,
		Notes:           notes,
		Status:          models.OrderStatusNew,
	}
	for _, ci := range cartItems {
		in.Items = append(in.Items, OrderItemInput{ProductID: ci.ProductID, Quantity: ci.Quantity})
	}
	order, err := s.Create(ctx, in)
	if err != nil {
		return nil, err
	}
	_ = s.cart.Clear(ctx, customerID)
	return order, nil
}

type StatusUpdate struct {
	Status   string
	DriverID *uuid.UUID
}

func (s *OrderService) UpdateStatus(ctx context.Context, id uuid.UUID, upd StatusUpdate) (*models.Order, error) {
	order, err := s.orders.FindByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}
	if !validOrderStatus(upd.Status) {
		return nil, ErrValidation
	}
	prev := order.Status
	order.Status = upd.Status
	now := time.Now().UTC()

	if upd.DriverID != nil {
		if _, err := s.drivers.FindByID(ctx, *upd.DriverID); err != nil {
			return nil, ErrValidation
		}
		order.DriverID = upd.DriverID
		order.AssignedAt = &now
		if order.Status == models.OrderStatusNew {
			order.Status = models.OrderStatusAssigned
		}
	}

	if order.Status == models.OrderStatusCompleted && prev != models.OrderStatusCompleted {
		order.CompletedAt = &now
		if err := s.recordCompletionFinance(ctx, order); err != nil {
			return nil, err
		}
	}
	if err := s.orders.Save(ctx, order); err != nil {
		return nil, err
	}

	// Notify driver on assignment.
	if upd.DriverID != nil && s.notifier != nil {
		if d, err := s.drivers.FindByID(ctx, *upd.DriverID); err == nil {
			_ = s.notifier.NotifyUser(ctx, d.UserID, models.NotifKindOrder,
				"Фармоиши нав", "Новый заказ",
				"Ба шумо фармоиш таъин шуд: "+order.Number,
				"Вам назначен заказ: "+order.Number, "/driver/orders")
		}
	}
	return s.orders.FindByID(ctx, order.ID)
}

func (s *OrderService) Assign(ctx context.Context, id, driverID uuid.UUID) (*models.Order, error) {
	return s.UpdateStatus(ctx, id, StatusUpdate{Status: models.OrderStatusAssigned, DriverID: &driverID})
}

func (s *OrderService) recordCompletionFinance(ctx context.Context, o *models.Order) error {
	txs := []models.FinancialTransaction{
		{Kind: models.TxKindIncome, OrderID: &o.ID, Amount: o.Total, Currency: o.Currency,
			Description: "Order " + o.Number + " income", OccurredAt: time.Now().UTC()},
		{Kind: models.TxKindPurchase, OrderID: &o.ID, Amount: o.CostTotal, Currency: o.Currency,
			Description: "Order " + o.Number + " cost", OccurredAt: time.Now().UTC()},
	}
	return s.txs.CreateMany(ctx, txs)
}

func (s *OrderService) Get(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	o, err := s.orders.FindByID(ctx, id)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return o, nil
}

func (s *OrderService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.orders.Delete(ctx, id)
}

func (s *OrderService) List(ctx context.Context, p repositories.ListOrdersParams) ([]models.Order, int64, error) {
	return s.orders.List(ctx, p)
}

func computeTotals(o *models.Order, items []models.OrderItem) {
	var subtotal, cost, profit float64
	for _, it := range items {
		subtotal += it.LineTotal
		cost += it.CostPrice * float64(it.Quantity)
		profit += it.Profit
	}
	o.Subtotal = subtotal
	disc := subtotal * o.DiscountPercent / 100.0
	o.Total = subtotal - disc
	o.CostTotal = cost
	o.Profit = o.Total - cost
}

func validOrderStatus(s string) bool {
	switch s {
	case models.OrderStatusNew, models.OrderStatusProcessing, models.OrderStatusAssigned,
		models.OrderStatusOnDelivery, models.OrderStatusCompleted, models.OrderStatusCancelled:
		return true
	}
	return false
}
