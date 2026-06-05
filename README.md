# AliStroy — Enterprise Construction Marketplace (Backend)

AliStroy is a multi-vendor construction-materials marketplace + CRM for Tajikistan.
This repository currently contains the **Backend (Phase 1)**. The Nuxt 4 frontend
will be added in Phase 2.

The platform is contact-driven: customers browse products and reach the seller
(or the marketplace) by phone / WhatsApp / Telegram. Admin moderates every
product, sets the public sale price, and decides whose contacts are shown.

## Tech stack

- **Go 1.22** + **Gin**
- **PostgreSQL** (GORM) + **Redis**
- **JWT** auth (access + refresh), role-based access control
- Clean architecture: **handlers → services → repositories → models**
- **WebSockets** (chat + notifications) via Redis pub/sub
- **Excel export** (excelize), **printable receipts** (HTML/print-to-PDF)
- **Google OAuth2** customer login (optional)
- **Docker**, **Docker Compose**, **Nginx**

## Languages

All public content is bilingual: **Tajik (default)** and **Russian**
(`name_tj` / `name_ru`, `description_tj` / `description_ru`, etc.).

## Repository layout

```
alistroy-v1/
├── backend/
│   ├── cmd/api/                 # main entrypoint
│   ├── internal/
│   │   ├── config/              # env config
│   │   ├── database/            # postgres connection
│   │   ├── cache/               # redis client
│   │   ├── jwt/                 # token manager
│   │   ├── httpx/               # response envelopes
│   │   ├── middleware/          # auth, cors, rate-limit, security headers
│   │   ├── validation/          # validator
│   │   ├── logger/              # levelled logger
│   │   ├── models/              # GORM models
│   │   ├── repositories/        # data access
│   │   ├── services/            # business logic
│   │   ├── handlers/            # HTTP handlers
│   │   ├── server/              # router wiring
│   │   ├── seed/                # bootstrap admin/settings/categories
│   │   ├── exporter/            # Excel exports
│   │   ├── receipt/             # printable receipt
│   │   ├── oauth/               # Google OAuth
│   │   └── ws/                  # websocket hub
│   ├── migrations/              # embedded SQL migrations
│   └── Dockerfile
├── nginx/
├── docker-compose.yml
└── .env.example
```

## Domain features

- **CMS Banners** — device-specific images (desktop/tablet/mobile), link, sort,
  active flag, grouped by position (hero, side, mid_large, mid_small, …).
- **Categories / Subcategories / Brands** — bilingual, icons, SEO fields, sort,
  enable/disable.
- **Sellers** — created by admin (login + password), full contact set.
- **Customers** — phone+password or Google login; admin can also create them.
- **Products** — multi-image, drag-sort (client sends ordered list), SKU, unit,
  cost/sale price, discount, stock + minimum-stock (low-stock detection).
- **Approval flow** — seller submits → status `pending`, the seller-entered price
  becomes **cost price**; admin sets the **sale price** + **contact owner**
  (admin or seller) + phone/WhatsApp/Telegram, then approves/rejects.
- **Orders** — cart checkout creates an order automatically; admin can also
  create orders manually. Statuses: new → processing → assigned → on_delivery →
  completed → cancelled. Profit computed automatically per item & per order.
- **Receipts** — printable HTML receipt at `GET /api/v1/admin/orders/:id/receipt`.
- **Drivers + delivery** — admin creates drivers (login), assigns orders; driver
  dashboard lists new/assigned/completed orders.
- **CRM chat** — realtime customer ↔ admin chat (text + image/video attachments)
  over WebSockets, with unread counters + notifications.
- **Accounting** — completed orders generate income + purchase transactions;
  profit auto-calculated. Reports filterable by today/week/month/custom.
- **Export** — Excel export honouring the active filters for products,
  categories, orders, customers, sellers, drivers and financial reports.
- **Analytics** — product views, phone/WhatsApp/Telegram clicks; admin dashboard.
- **Visual board** — XMind-style tree (`GET /api/v1/admin/board`):
  Category → Subcategory → Products.
- **Logo management** — site logo + favicon stored in settings, applied site-wide.
- **SEO** — `/robots.txt`, `/sitemap.xml`, SEO slugs + meta fields on categories.
- **Security** — JWT, bcrypt password hashing, request validation, per-IP rate
  limiting, secure headers, CORS.

## Quick start (VPS / production)

```bash
git clone https://github.com/instaagrammeta/alistroy-v1.git
cd alistroy-v1

cp .env.example .env
nano .env   # set JWT_SECRET (>=32 chars), POSTGRES_PASSWORD, ADMIN_*, MARKETPLACE_*

docker compose up -d --build
```

Then:

- API base: `http://YOUR_SERVER/api/v1`
- Health: `http://YOUR_SERVER/health` → `{"status":"ok"}`
- Robots: `http://YOUR_SERVER/robots.txt`
- Sitemap: `http://YOUR_SERVER/sitemap.xml`

The backend auto-runs migrations and seeds the admin account (`ADMIN_EMAIL` /
`ADMIN_PASSWORD`), default settings, and 15 root categories on first boot.

### Verify

```bash
docker compose ps
docker compose logs backend | grep -E "listening|seed"
curl http://localhost/health
curl -s http://localhost/api/v1/categories
```

## Local development

```bash
cd backend
cp .env.example .env   # POSTGRES_HOST=localhost, REDIS_HOST=localhost
go run ./cmd/api
```

(Requires a local PostgreSQL 16 + Redis 7, or point the env at remote ones.)

## Auth model

- `POST /api/v1/auth/register` — customer (name, phone, password, address).
- `POST /api/v1/auth/login` — identifier (phone/email/login) + password.
- `POST /api/v1/auth/refresh` — exchange refresh token.
- `GET /api/v1/auth/google/url` + `POST /api/v1/auth/google/callback` — Google
  login; response includes `needs_profile=true` when phone/address must still be
  collected.
- `Authorization: Bearer <access_token>` for protected routes.

Roles: `customer`, `seller`, `driver`, `admin`.

## API surface (high level)

| Area      | Prefix |
|-----------|--------|
| Public catalog | `GET /api/v1/products`, `/categories`, `/brands`, `/sellers`, `/banners`, `/settings/public` |
| Customer | `/api/v1/customer/*` (cart, checkout, orders, chat) |
| Seller   | `/api/v1/seller/*` (products CRUD, profile, stats) |
| Driver   | `/api/v1/driver/*` (assigned orders, status) |
| Admin    | `/api/v1/admin/*` (users, sellers, customers, drivers, products, moderation, orders, banners, categories, brands, reviews, chat, reports, exports, board, settings) |

## Migrations

SQL files in `backend/migrations/*.up.sql` are embedded and applied on boot
(idempotent — `CREATE TABLE IF NOT EXISTS`, etc.).

## Production notes

- Put TLS in front (Nginx + certbot) and set `PUBLIC_URL`/`CORS_ALLOWED_ORIGINS`.
- Back up PostgreSQL:
  ```bash
  docker compose exec postgres pg_dump -U $POSTGRES_USER $POSTGRES_DB > backup.sql
  ```
- Change the seeded admin password immediately after first login.

## License

Proprietary — © AliStroy.
