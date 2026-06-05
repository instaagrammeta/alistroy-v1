// Package exporter builds Excel (.xlsx) files for the various admin lists.
package exporter

import (
	"bytes"
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/repositories"
)

const sheet = "Sheet1"

func newFile(headers []string) (*excelize.File, error) {
	f := excelize.NewFile()
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		if err := f.SetCellValue(sheet, cell, h); err != nil {
			return nil, err
		}
	}
	style, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	first, _ := excelize.CoordinatesToCellName(1, 1)
	last, _ := excelize.CoordinatesToCellName(len(headers), 1)
	_ = f.SetCellStyle(sheet, first, last, style)
	return f, nil
}

func write(f *excelize.File, row int, vals ...any) {
	for i, v := range vals {
		cell, _ := excelize.CoordinatesToCellName(i+1, row)
		_ = f.SetCellValue(sheet, cell, v)
	}
}

func toBytes(f *excelize.File) ([]byte, error) {
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Products(items []models.Product) ([]byte, error) {
	f, err := newFile([]string{"ID", "Name TJ", "Name RU", "SKU", "Category", "Seller", "Cost", "Sale", "Stock", "Status", "Created"})
	if err != nil {
		return nil, err
	}
	for i, p := range items {
		cat, seller := "", ""
		if p.Category != nil {
			cat = p.Category.NameRU
		}
		if p.Seller != nil {
			seller = p.Seller.FullName
		}
		write(f, i+2, p.ID.String(), p.NameTJ, p.NameRU, p.SKU, cat, seller,
			p.CostPrice, p.SalePrice, p.StockQuantity, p.Status, p.CreatedAt.Format(time.RFC3339))
	}
	return toBytes(f)
}

func Categories(items []models.Category) ([]byte, error) {
	f, err := newFile([]string{"ID", "Name TJ", "Name RU", "Slug", "Sort", "Active"})
	if err != nil {
		return nil, err
	}
	for i, c := range items {
		write(f, i+2, c.ID.String(), c.NameTJ, c.NameRU, c.Slug, c.SortOrder, c.Active)
	}
	return toBytes(f)
}

func Orders(items []models.Order) ([]byte, error) {
	f, err := newFile([]string{"Number", "Customer", "Phone", "Status", "Subtotal", "Total", "Cost", "Profit", "Created"})
	if err != nil {
		return nil, err
	}
	for i, o := range items {
		write(f, i+2, o.Number, o.CustomerName, o.CustomerPhone, o.Status,
			o.Subtotal, o.Total, o.CostTotal, o.Profit, o.CreatedAt.Format(time.RFC3339))
	}
	return toBytes(f)
}

func Customers(items []models.Customer) ([]byte, error) {
	f, err := newFile([]string{"ID", "Company", "City", "Address", "Created"})
	if err != nil {
		return nil, err
	}
	for i, c := range items {
		write(f, i+2, c.ID.String(), c.Company, c.City, c.Address, c.CreatedAt.Format(time.RFC3339))
	}
	return toBytes(f)
}

func Sellers(items []models.Seller) ([]byte, error) {
	f, err := newFile([]string{"ID", "Full Name", "Company", "Market", "Phone", "City", "Active"})
	if err != nil {
		return nil, err
	}
	for i, s := range items {
		write(f, i+2, s.ID.String(), s.FullName, s.CompanyName, s.MarketName, s.Phone, s.City, s.Active)
	}
	return toBytes(f)
}

func Drivers(items []models.Driver) ([]byte, error) {
	f, err := newFile([]string{"ID", "Full Name", "Phone", "Vehicle", "Active", "On Duty"})
	if err != nil {
		return nil, err
	}
	for i, d := range items {
		write(f, i+2, d.ID.String(), d.FullName, d.Phone, d.Vehicle, d.Active, d.OnDuty)
	}
	return toBytes(f)
}

func Transactions(items []models.FinancialTransaction, summary *repositories.Summary) ([]byte, error) {
	f, err := newFile([]string{"Kind", "Amount", "Currency", "Description", "Occurred"})
	if err != nil {
		return nil, err
	}
	for i, t := range items {
		write(f, i+2, t.Kind, t.Amount, t.Currency, t.Description, t.OccurredAt.Format(time.RFC3339))
	}
	if summary != nil {
		base := len(items) + 3
		write(f, base, "Income", summary.Income)
		write(f, base+1, "Purchase", summary.Purchase)
		write(f, base+2, "Expense", summary.Expense)
		write(f, base+3, "Profit", summary.Profit)
	}
	return toBytes(f)
}

// FileName builds a timestamped download name.
func FileName(prefix string) string {
	return fmt.Sprintf("%s_%s.xlsx", prefix, time.Now().UTC().Format("20060102_150405"))
}
