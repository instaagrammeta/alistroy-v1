# AliStroy — Multi-Vendor Construction Materials Marketplace

AliStroy is a multi-vendor construction materials marketplace and lead-generation platform for Tajikistan. The platform is **not** an online-payment marketplace: customers browse products and contact the seller (or marketplace) directly via phone or WhatsApp.

## Tech stack

- **Frontend:** Nuxt 4, Vue 3, TypeScript, TailwindCSS, Pinia, @nuxtjs/i18n, SSR
- **Backend:** Go (Gin), PostgreSQL, Redis, GORM, JWT auth, REST API
- **Infrastructure:** Docker, Docker Compose, Nginx

## Languages

- Tajik (`tg`) — default
- Russian (`ru`)

All product, category, and page content is stored bilingually.

## Repository layout

```
alistroy-v1/
├── backend/             # Go API (Gin + GORM + Redis)
├── frontend/            # Nuxt 4 SSR app
├── nginx/               # Nginx reverse proxy config
├── docker-compose.yml   # Production compose stack
├── .env.example         # Environment template
└── README.md
```

## Quick start (production VPS)

1. Clone the repository on the VPS:
   ```bash
   git clone https://github.com/instaagrammeta/alistroy-v1.git
   cd alistroy-v1
   ```

2. Copy and edit the environment file:
   ```bash
   cp .env.example .env
   nano .env
   ```
   Set strong secrets for `POSTGRES_PASSWORD`, `JWT_SECRET`, the admin seed account, and your domain.

3. Build and launch:
   ```bash
   docker compose up -d --build
   ```

4. Visit `http://your-domain/`.
   - Public site: `/`
   - Seller dashboard: `/seller`
   - Admin panel: `/admin`
   - REST API: `/api/v1/...`

5. Default admin (created on first boot from `.env`):
   - Email: `ADMIN_EMAIL`
   - Password: `ADMIN_PASSWORD`

   **Change the password immediately after first login.**

## Local development

```bash
# Backend
cd backend
cp .env.example .env
go run ./cmd/api

# Frontend
cd frontend
cp .env.example .env
npm install
npm run dev
```

## Business model

- Customers browse products. They do **not** place online orders.
- Each product page exposes a **Call** and a **WhatsApp** button.
- Each product is configured (by admin) to route inquiries either to:
  - the **marketplace** (admin) phone/WhatsApp, or
  - the **seller** phone/WhatsApp.
- All button clicks are tracked for analytics.

## Roles

- **Customer** — browses products, saves favorites, leaves reviews.
- **Seller** — manages own products, views inquiry counts.
- **Admin** — moderates products, manages users/sellers/categories/reviews/settings, sees analytics.

## Moderation workflow

1. Seller submits a product (status `pending`).
2. Admin reviews and either approves, rejects, or requests changes.
3. On approval, admin sets `contact_type`, `phone_number`, `whatsapp_number`.
4. Approved products become public.

## License

Proprietary — © AliStroy.
