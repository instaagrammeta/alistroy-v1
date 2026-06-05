-- AliStroy initial schema
-- Idempotent: safe to apply multiple times.

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =========================================================
-- users
-- =========================================================
CREATE TABLE IF NOT EXISTS users (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email            VARCHAR(255) NOT NULL UNIQUE,
    password_hash    VARCHAR(255) NOT NULL,
    name             VARCHAR(120) NOT NULL,
    phone            VARCHAR(32)  NOT NULL DEFAULT '',
    role             VARCHAR(20)  NOT NULL DEFAULT 'customer',
    locale           VARCHAR(5)   NOT NULL DEFAULT 'tg',
    is_active        BOOLEAN      NOT NULL DEFAULT TRUE,
    last_login_at    TIMESTAMPTZ,
    reset_token      VARCHAR(128) NOT NULL DEFAULT '',
    reset_expires_at TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_users_role        ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at  ON users(deleted_at);
CREATE INDEX IF NOT EXISTS idx_users_reset_token ON users(reset_token);

-- =========================================================
-- sellers
-- =========================================================
CREATE TABLE IF NOT EXISTS sellers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    name            VARCHAR(160) NOT NULL,
    slug            VARCHAR(200) NOT NULL UNIQUE,
    description_tj  TEXT NOT NULL DEFAULT '',
    description_ru  TEXT NOT NULL DEFAULT '',
    logo_url        VARCHAR(500) NOT NULL DEFAULT '',
    phone           VARCHAR(32)  NOT NULL DEFAULT '',
    whats_app       VARCHAR(32)  NOT NULL DEFAULT '',
    address         VARCHAR(255) NOT NULL DEFAULT '',
    city            VARCHAR(100) NOT NULL DEFAULT '',
    status          VARCHAR(20)  NOT NULL DEFAULT 'pending',
    is_featured     BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_sellers_status     ON sellers(status);
CREATE INDEX IF NOT EXISTS idx_sellers_featured   ON sellers(is_featured);
CREATE INDEX IF NOT EXISTS idx_sellers_deleted_at ON sellers(deleted_at);

-- =========================================================
-- categories
-- =========================================================
CREATE TABLE IF NOT EXISTS categories (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug        VARCHAR(120) NOT NULL UNIQUE,
    title_tj    VARCHAR(160) NOT NULL,
    title_ru    VARCHAR(160) NOT NULL,
    icon_url    VARCHAR(500) NOT NULL DEFAULT '',
    sort_order  INTEGER      NOT NULL DEFAULT 0,
    is_active   BOOLEAN      NOT NULL DEFAULT TRUE,
    parent_id   UUID REFERENCES categories(id) ON DELETE SET NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_categories_sort       ON categories(sort_order);
CREATE INDEX IF NOT EXISTS idx_categories_active     ON categories(is_active);
CREATE INDEX IF NOT EXISTS idx_categories_parent     ON categories(parent_id);
CREATE INDEX IF NOT EXISTS idx_categories_deleted_at ON categories(deleted_at);

-- =========================================================
-- products
-- =========================================================
CREATE TABLE IF NOT EXISTS products (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seller_id        UUID NOT NULL REFERENCES sellers(id) ON DELETE CASCADE,
    category_id      UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    slug             VARCHAR(220) NOT NULL UNIQUE,
    sku              VARCHAR(64)  NOT NULL DEFAULT '',
    title_tj         VARCHAR(255) NOT NULL,
    title_ru         VARCHAR(255) NOT NULL,
    description_tj   TEXT         NOT NULL DEFAULT '',
    description_ru   TEXT         NOT NULL DEFAULT '',
    price            NUMERIC(14,2) NOT NULL DEFAULT 0,
    currency         VARCHAR(8)   NOT NULL DEFAULT 'TJS',
    unit             VARCHAR(20)  NOT NULL DEFAULT 'pcs',
    stock_quantity   INTEGER      NOT NULL DEFAULT 0,
    is_available     BOOLEAN      NOT NULL DEFAULT TRUE,
    contact_type     VARCHAR(10)  NOT NULL DEFAULT 'admin',
    phone_number     VARCHAR(32)  NOT NULL DEFAULT '',
    whats_app_number VARCHAR(32)  NOT NULL DEFAULT '',
    status           VARCHAR(20)  NOT NULL DEFAULT 'draft',
    rejection_note   TEXT         NOT NULL DEFAULT '',
    is_featured      BOOLEAN      NOT NULL DEFAULT FALSE,
    views_count      BIGINT       NOT NULL DEFAULT 0,
    phone_clicks     BIGINT       NOT NULL DEFAULT 0,
    whats_app_clicks BIGINT       NOT NULL DEFAULT 0,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_products_seller     ON products(seller_id);
CREATE INDEX IF NOT EXISTS idx_products_category   ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_status     ON products(status);
CREATE INDEX IF NOT EXISTS idx_products_available  ON products(is_available);
CREATE INDEX IF NOT EXISTS idx_products_featured   ON products(is_featured);
CREATE INDEX IF NOT EXISTS idx_products_sku        ON products(sku);
CREATE INDEX IF NOT EXISTS idx_products_created_at ON products(created_at);
CREATE INDEX IF NOT EXISTS idx_products_deleted_at ON products(deleted_at);
-- Trigram-style search relies on btree-indexed text; use GIN on lower(...) only when pg_trgm exists.

-- =========================================================
-- product_images
-- =========================================================
CREATE TABLE IF NOT EXISTS product_images (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id  UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    url         VARCHAR(500) NOT NULL,
    alt         VARCHAR(255) NOT NULL DEFAULT '',
    sort_order  INTEGER      NOT NULL DEFAULT 0,
    is_cover    BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_product_images_product ON product_images(product_id);
CREATE INDEX IF NOT EXISTS idx_product_images_sort    ON product_images(sort_order);

-- =========================================================
-- reviews
-- =========================================================
CREATE TABLE IF NOT EXISTS reviews (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id  UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rating      INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment     TEXT NOT NULL DEFAULT '',
    status      VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_reviews_product ON reviews(product_id);
CREATE INDEX IF NOT EXISTS idx_reviews_user    ON reviews(user_id);
CREATE INDEX IF NOT EXISTS idx_reviews_status  ON reviews(status);

-- =========================================================
-- favorites
-- =========================================================
CREATE TABLE IF NOT EXISTS favorites (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id  UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_product ON favorites(user_id, product_id);

-- =========================================================
-- settings (key/value)
-- =========================================================
CREATE TABLE IF NOT EXISTS settings (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key        VARCHAR(80) NOT NULL UNIQUE,
    value      TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- =========================================================
-- tracking_events
-- =========================================================
CREATE TABLE IF NOT EXISTS tracking_events (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id  UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    seller_id   UUID NOT NULL REFERENCES sellers(id) ON DELETE CASCADE,
    user_id     UUID REFERENCES users(id) ON DELETE SET NULL,
    event       VARCHAR(32) NOT NULL,
    ip          VARCHAR(64) NOT NULL DEFAULT '',
    user_agent  VARCHAR(500) NOT NULL DEFAULT '',
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_events_product   ON tracking_events(product_id);
CREATE INDEX IF NOT EXISTS idx_events_seller    ON tracking_events(seller_id);
CREATE INDEX IF NOT EXISTS idx_events_event     ON tracking_events(event);
CREATE INDEX IF NOT EXISTS idx_events_occurred  ON tracking_events(occurred_at);
