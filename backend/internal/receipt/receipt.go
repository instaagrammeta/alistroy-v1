// Package receipt renders a printable HTML receipt for an order. The frontend
// (or browser print-to-PDF) turns this into a PDF; keeping it as HTML avoids a
// heavy native PDF dependency while remaining fully printable.
package receipt

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
)

type Data struct {
	BrandName string
	Subtitle  string
	Order     *models.Order
	Now       string
}

var tpl = template.Must(template.New("receipt").Funcs(template.FuncMap{
	"money": func(v float64) string { return fmt.Sprintf("%.2f", v) },
	"add1":  func(i int) int { return i + 1 },
}).Parse(receiptHTML))

// RenderHTML produces a standalone, printable HTML document for the order.
func RenderHTML(brandName, subtitle string, o *models.Order) ([]byte, error) {
	if brandName == "" {
		brandName = "AliStroy CRM"
	}
	data := Data{
		BrandName: brandName,
		Subtitle:  subtitle,
		Order:     o,
		Now:       time.Now().Format("02.01.2006, 15:04"),
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

const receiptHTML = `<!doctype html>
<html lang="tg">
<head>
<meta charset="utf-8">
<title>{{.BrandName}} — {{.Order.Number}}</title>
<style>
  * { box-sizing: border-box; }
  body { font-family: 'Segoe UI', Arial, sans-serif; color:#1F2937; margin:0; padding:24px; }
  .receipt { max-width: 420px; margin: 0 auto; border:1px solid #E5E7EB; border-radius:12px; padding:20px; }
  .brand { font-size:22px; font-weight:800; color:#FF661A; text-align:center; }
  .sub { text-align:center; color:#6B7280; font-size:12px; margin-bottom:8px; }
  .meta { font-size:13px; margin:12px 0; }
  .meta div { display:flex; justify-content:space-between; padding:2px 0; }
  table { width:100%; border-collapse:collapse; margin-top:8px; font-size:13px; }
  th, td { text-align:left; padding:6px 4px; border-bottom:1px solid #F3F4F6; }
  td.r, th.r { text-align:right; }
  .total { font-weight:800; font-size:16px; margin-top:10px; display:flex; justify-content:space-between; }
  .pay { display:flex; justify-content:space-between; color:#FF661A; font-weight:700; }
  .thanks { text-align:center; margin-top:14px; color:#6B7280; font-size:12px; }
  @media print { .no-print { display:none; } body { padding:0; } }
  .btn { display:inline-block; margin-top:16px; background:#FF661A; color:#fff; border:none; padding:10px 18px; border-radius:8px; cursor:pointer; }
</style>
</head>
<body>
  <div class="receipt">
    <div class="brand">{{.BrandName}}</div>
    <div class="sub">{{.Subtitle}}</div>
    <div class="sub">{{.Now}}</div>
    <div class="meta">
      <div><span>Заявка №:</span><strong>{{.Order.Number}}</strong></div>
      <div><span>Клиент:</span><strong>{{.Order.CustomerName}}</strong></div>
      <div><span>Тел:</span><strong>{{.Order.CustomerPhone}}</strong></div>
      <div><span>Расондан:</span><strong>{{.Order.DeliveryAddress}}</strong></div>
      <div><span>Статус:</span><strong>{{.Order.Status}}</strong></div>
    </div>
    <table>
      <thead>
        <tr><th>#</th><th>Маҳсулот</th><th class="r">Шумора</th><th class="r">Нарх</th><th class="r">Ҷамъ</th></tr>
      </thead>
      <tbody>
        {{range $i, $it := .Order.Items}}
        <tr>
          <td>{{add1 $i}}</td>
          <td>{{$it.NameSnap}}</td>
          <td class="r">{{$it.Quantity}}</td>
          <td class="r">{{money $it.SalePrice}}</td>
          <td class="r">{{money $it.LineTotal}}</td>
        </tr>
        {{end}}
      </tbody>
    </table>
    <div class="total"><span>Маблағи умумӣ:</span><span>{{money .Order.Total}} {{.Order.Currency}}</span></div>
    <div class="pay"><span>Барои пардохт:</span><span>{{money .Order.Total}} {{.Order.Currency}}</span></div>
    <div class="thanks">Раҳмат барои хариди шумо!</div>
    <div class="no-print" style="text-align:center;">
      <button class="btn" onclick="window.print()">Чоп / Печать</button>
    </div>
  </div>
</body>
</html>`
