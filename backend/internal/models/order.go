package models

import (
	"time"

	"github.com/google/uuid"
)

// Order is the customer purchase request. AliStroy is contact-driven, so an
// "order" represents the lead/quote and any later delivery assignment.
type Order struct {
	BaseModel
	Number          string     `gorm:"size:32;uniqueIndex;not null" json:"number"`
	CustomerID      *uuid.UUID `gorm:"type:uuid;index" json:"customer_id,omitempty"`
	CustomerName    string     `gorm:"size:160" json:"customer_name"`
	CustomerPhone   string     `gorm:"size:32" json:"customer_phone"`
	DeliveryAddress string     `gorm:"size:500" json:"delivery_address"`
	DeliveryDate    *time.Time `json:"delivery_date,omitempty"`
	Status          string     `gorm:"size:20;not null;default:new;index" json:"status"`
	DiscountPercent float64    `gorm:"type:numeric(5,2);not null;default:0" json:"discount_percent"`
	Subtotal        float64    `gorm:"type:numeric(14,2);not null;default:0" json:"subtotal"`
	Total           float64    `gorm:"type:numeric(14,2);not null;default:0" json:"total"`
	CostTotal       float64    `gorm:"type:numeric(14,2);not null;default:0" json:"cost_total"`
	Profit          float64    `gorm:"type:numeric(14,2);not null;default:0" json:"profit"`
	Currency        string     `gorm:"size:8;not null;default:TJS" json:"currency"`
	Notes           string     `gorm:"type:text" json:"notes"`
	DriverID        *uuid.UUID `gorm:"type:uuid;index" json:"driver_id,omitempty"`
	AssignedAt      *time.Time `json:"assigned_at,omitempty"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`

	Items    []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
	Customer *Customer   `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Driver   *Driver     `gorm:"foreignKey:DriverID" json:"driver,omitempty"`
}

// OrderItem snapshots a product at the moment of order.
type OrderItem struct {
	BaseModel
	OrderID   uuid.UUID `gorm:"type:uuid;index;not null" json:"order_id"`
	ProductID uuid.UUID `gorm:"type:uuid;index;not null" json:"product_id"`
	NameSnap  string    `gorm:"size:255;not null" json:"name_snapshot"`
	Unit      string    `gorm:"size:20" json:"unit"`
	Quantity  int       `gorm:"not null;default:1" json:"quantity"`
	CostPrice float64   `gorm:"type:numeric(14,2);not null;default:0" json:"cost_price"`
	SalePrice float64   `gorm:"type:numeric(14,2);not null;default:0" json:"sale_price"`
	LineTotal float64   `gorm:"type:numeric(14,2);not null;default:0" json:"line_total"`
	Profit    float64   `gorm:"type:numeric(14,2);not null;default:0" json:"profit"`

	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}
