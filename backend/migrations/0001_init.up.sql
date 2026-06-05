CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =========================================================
-- users (auth root)
-- =========================================================
CREATE TABLE IF NOT EXISTS users (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email          VARCHAR(255) NOT NULL DEFAULT '',
    phone          VARCHAR(32)  NOT NULL DEFAULT '',
    login          VARCHAR(64)  NOT NULL DEFAULT '',
    password_hash  VARCHAR(255) NOT NULL DEFAULT '',
    name           VARCHAR(160) NOT NULL,
    role           VARCHAR(20)  NOT NULL DEFAULT 'customer',
    status         VARCHAR(20)  NOT NULL DEFAULT 'active',
    locale         VARCHAR(5)   NOT NULL DEFAULT 'tg',
    avatar_url     VARCHAR(500) NOT NULL DEFAULT '',
    google_id      VARCHAR(128) NOT NULL DEFAULT '',
    last_login_at  TIMESTAMPTZ,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at     TIMESTAMPTZ
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE email <> '' AND deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_phone ON users(phone) WHERE phone <> '' AND deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_login ON users(login) WHERE login <> '' AND deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_role        ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_status      ON users(status);
CREATE INDEX IF NOT EXISTS idx_users_google_id   ON users(google_id) WHERE google_id <> '';
CREATE INDEX IF NOT EXISTS idx_users_deleted_at  ON users(deleted_at);

-- =========================================================
-- customers
-- =========================================================
CREATE TABLE IF NOT EXISTS customers (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    company     VARCHAR(160) NOT NULL DEFAULT '',
    address     VARCHAR(500) NOT NULL DEFAULT '',
    city        VARCHAR(100) NOT NULL DEFAULT '',
    notes       TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_customers_deleted_at ON customers(deleted_at);

-- addresses (per customer)
CREATE TABLE IF NOT EXISTS addresses (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id  UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    label        VARCHAR(80)  NOT NULL DEFAULT '',
    phone        VARCHAR(32)  NOT NULL DEFAULT '',
    city         VARCHAR(100) NOT NULL DEFAULT '',
    street       VARCHAR(500) NOT NULL DEFAULT '',
    is_primary   BOOLEAN NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_addresses_customer ON addresses(customer_id);

-- =========================================================
-- categories
-- =========================================================
CREATE TABLE IF NOT EXISTS categories (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug                VARCHAR(160) NOT NULL UNIQUE,
    name_tj             VARCHAR(200) NOT NULL,
    name_ru             VARCHAR(200) NOT NULL,
    description_tj      TEXT NOT NULL DEFAULT '',
    description_ru      TEXT NOT NULL DEFAULT '',
    icon_url            VARCHAR(500) NOT NULL DEFAULT '',
    banner_url          VARCHAR(500) NOT NULL DEFAULT '',
    seo_title_tj        VARCHAR(255) NOT NULL DEFAULT '',
    seo_title_ru        VARCHAR(255) NOT NULL DEFAULT '',
    seo_description_tj  TEXT NOT NULL DEFAULT '',
    seo_description_ru  TEXT NOT NULL DEFAULT '',
    sort_order          INTEGER NOT NULL DEFAULT 0,
    active              BOOLEAN NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_categories_sort       ON categories(sort_order);
CREATE INDEX IF NOT EXISTS idx_categories_active     ON categories(active);
CREATE INDEX IF NOT EXISTS idx_categories_deleted_at ON categories(deleted_at);

-- =========================================================
-- subcategories
-- =========================================================
CREATE TABLE IF NOT EXISTS subcategories (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id  UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    slug         VARCHAR(160) NOT NULL UNIQUE,
    name_tj      VARCHAR(200) NOT NULL,
    name_ru      VARCHAR(200) NOT NULL,
    icon_url     VARCHAR(500) NOT NULL DEFAULT '',
    sort_order   INTEGER NOT NULL DEFAULT 0,
    active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_subcategories_category   ON subcategories(category_id);
CREATE INDEX IF NOT EXISTS idx_subcategories_sort       ON subcategories(sort_order);
CREATE INDEX IF NOT EXISTS idx_subcategories_active     ON subcategories(active);
CREATE INDEX IF NOT EXISTS idx_subcategories_deleted_at ON subcategories(deleted_at);

-- =========================================================
-- brands
-- =========================================================
CREATE TABLE IF NOT EXISTS brands (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug        VARCHAR(160) NOT NULL UNIQUE,
    name        VARCHAR(200) NOT NULL,
    logo_url    VARCHAR(500) NOT NULL DEFAULT '',
    active      BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order  INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_brands_active     ON brands(active);
CREATE INDEX IF NOT EXISTS idx_brands_deleted_at ON brands(deleted_at);

-- =========================================================
-- sellers
-- =========================================================
CREATE TABLE IF NOT EXISTS sellers (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    full_name              VARCHAR(160) NOT NULL,
    company_name           VARCHAR(200) NOT NULL DEFAULT '',
    market_name            VARCHAR(160) NOT NULL DEFAULT '',
    slug                   VARCHAR(200) UNIQUE,
    phone                  VARCHAR(32)  NOT NULL DEFAULT '',
    phone_alt              VARCHAR(32)  NOT NULL DEFAULT '',
    whats_app              VARCHAR(32)  NOT NULL DEFAULT '',
    telegram               VARCHAR(32)  NOT NULL DEFAULT '',
    telegram_username      VARCHAR(64)  NOT NULL DEFAULT '',
    address                VARCHAR(500) NOT NULL DEFAULT '',
    city                   VARCHAR(100) NOT NULL DEFAULT '',
    business_category      UUID REFERENCES categories(id) ON DELETE SET NULL,
    logo_url               VARCHAR(500) NOT NULL DEFAULT '',
    notes                  TEXT NOT NULL DEFAULT '',
    active                 BOOLEAN NOT NULL DEFAULT TRUE,
    is_featured            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at             TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_sellers_active            ON sellers(active);
CREATE INDEX IF NOT EXISTS idx_sellers_featured          ON sellers(is_featured);
CREATE INDEX IF NOT EXISTS idx_sellers_business_category ON sellers(business_category);
CREATE INDEX IF NOT EXISTS idx_sellers_deleted_at        ON sellers(deleted_at);

-- =========================================================
-- drivers
-- =========================================================
CREATE TABLE IF NOT EXISTS drivers (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    full_name   VARCHAR(160) NOT NULL,
    age         INTEGER NOT NULL DEFAULT 0,
    phone       VARCHAR(32)  NOT NULL DEFAULT '',
    phone_alt   VARCHAR(32)  NOT NULL DEFAULT '',
    whats_app   VARCHAR(32)  NOT NULL DEFAULT '',
    telegram    VARCHAR(32)  NOT NULL DEFAULT '',
    vehicle     VARCHAR(160) NOT NULL DEFAULT '',
    photo_url   VARCHAR(500) NOT NULL DEFAULT '',
    notes       TEXT NOT NULL DEFAULT '',
    active      BOOLEAN NOT NULL DEFAULT TRUE,
    on_duty     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_drivers_active     ON drivers(active);
CREATE INDEX IF NOT EXISTS idx_drivers_deleted_at ON drivers(deleted_at);

-- =========================================================
-- products
-- =========================================================
CREATE TABLE IF NOT EXISTS products (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seller_id          UUID NOT NULL REFERENCES sellers(id) ON DELETE CASCADE,
    category_id        UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    subcategory_id     UUID REFERENCES subcategories(id) ON DELETE SET NULL,
    brand_id           UUID REFERENCES brands(id) ON DELETE SET NULL,
    slug               VARCHAR(255) NOT NULL UNIQUE,
    sku                VARCHAR(64)  NOT NULL DEFAULT '',
    name_tj            VARCHAR(255) NOT NULL,
    name_ru            VARCHAR(255) NOT NULL,
    description_tj     TEXT NOT NULL DEFAULT '',
    description_ru     TEXT NOT NULL DEFAULT '',
    unit               VARCHAR(20) NOT NULL DEFAULT 'pcs',
    currency           VARCHAR(8)  NOT NULL DEFAULT 'TJS',
    cost_price         NUMERIC(14,2) NOT NULL DEFAULT 0,
    sale_price         NUMERIC(14,2) NOT NULL DEFAULT 0,
    discount_percent   NUMERIC(5,2)  NOT NULL DEFAULT 0,
    stock_quantity     INTEGER NOT NULL DEFAULT 0,
    minimum_stock      INTEGER NOT NULL DEFAULT 0,
    is_available       BOOLEAN NOT NULL DEFAULT TRUE,
    is_featured        BOOLEAN NOT NULL DEFAULT FALSE,
    contact_owner      VARCHAR(10) NOT NULL DEFAULT 'admin',
    contact_phone      VARCHAR(32) NOT NULL DEFAULT '',
    contact_whats_app  VARCHAR(32) NOT NULL DEFAULT '',
    contact_telegram   VARCHAR(32) NOT NULL DEFAULT '',
    status             VARCHAR(20) NOT NULL DEFAULT 'draft',
    rejection_note     TEXT NOT NULL DEFAULT '',
    views_count        BIGINT NOT NULL DEFAULT 0,
    phone_clicks       BIGINT NOT NULL DEFAULT 0,
    whats_app_clicks   BIGINT NOT NULL DEFAULT 0,
    telegram_clicks    BIGINT NOT NULL DEFAULT 0,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at         TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_products_seller       ON products(seller_id);
CREATE INDEX IF NOT EXISTS idx_products_category     ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_subcategory  ON products(subcategory_id);
CREATE INDEX IF NOT EXISTS idx_products_brand        ON products(brand_id);
CREATE INDEX IF NOT EXISTS idx_products_status       ON products(status);
CREATE INDEX IF NOT EXISTS idx_products_available    ON products(is_available);
CREATE INDEX IF NOT EXISTS idx_products_featured     ON products(is_featured);
CREATE INDEX IF NOT EXISTS idx_products_sku          ON products(sku);
CREATE INDEX IF NOT EXISTS idx_products_created_at   ON products(created_at);
CREATE INDEX IF NOT EXISTS idx_products_deleted_at   ON products(deleted_at);

-- product images
CREATE TABLE IF NOT EXISTS product_images (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id   UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    url          VARCHAR(500) NOT NULL,
    alt          VARCHAR(255) NOT NULL DEFAULT '',
    sort_order   INTEGER NOT NULL DEFAULT 0,
    is_cover     BOOLEAN NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_product_images_product ON product_images(product_id);
CREATE INDEX IF NOT EXISTS idx_product_images_sort    ON product_images(sort_order);

-- =========================================================
-- banners (CMS)
-- =========================================================
CREATE TABLE IF NOT EXISTS banners (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    position        VARCHAR(32) NOT NULL,
    title_tj        VARCHAR(255) NOT NULL DEFAULT '',
    title_ru        VARCHAR(255) NOT NULL DEFAULT '',
    description_tj  TEXT NOT NULL DEFAULT '',
    description_ru  TEXT NOT NULL DEFAULT '',
    desktop_url     VARCHAR(500) NOT NULL DEFAULT '',
    tablet_url      VARCHAR(500) NOT NULL DEFAULT '',
    mobile_url      VARCHAR(500) NOT NULL DEFAULT '',
    link_url        VARCHAR(500) NOT NULL DEFAULT '',
    sort_order      INTEGER NOT NULL DEFAULT 0,
    active          BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_banners_position   ON banners(position);
CREATE INDEX IF NOT EXISTS idx_banners_active     ON banners(active);
CREATE INDEX IF NOT EXISTS idx_banners_sort       ON banners(sort_order);
CREATE INDEX IF NOT EXISTS idx_banners_deleted_at ON banners(deleted_at);

-- =========================================================
-- orders + order_items
-- =========================================================
CREATE TABLE IF NOT EXISTS orders (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    number            VARCHAR(32) NOT NULL UNIQUE,
    customer_id       UUID REFERENCES customers(id) ON DELETE SET NULL,
    customer_name     VARCHAR(160) NOT NULL DEFAULT '',
    customer_phone    VARCHAR(32)  NOT NULL DEFAULT '',
    delivery_address  VARCHAR(500) NOT NULL DEFAULT '',
    delivery_date     TIMESTAMPTZ,
    status            VARCHAR(20)  NOT NULL DEFAULT 'new',
    discount_percent  NUMERIC(5,2) NOT NULL DEFAULT 0,
    subtotal          NUMERIC(14,2) NOT NULL DEFAULT 0,
    total             NUMERIC(14,2) NOT NULL DEFAULT 0,
    cost_total        NUMERIC(14,2) NOT NULL DEFAULT 0,
    profit            NUMERIC(14,2) NOT NULL DEFAULT 0,
    currency          VARCHAR(8)   NOT NULL DEFAULT 'TJS',
    notes             TEXT NOT NULL DEFAULT '',
    driver_id         UUID REFERENCES drivers(id) ON DELETE SET NULL,
    assigned_at       TIMESTAMPTZ,
    completed_at      TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_orders_customer    ON orders(customer_id);
CREATE INDEX IF NOT EXISTS idx_orders_status      ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_driver      ON orders(driver_id);
CREATE INDEX IF NOT EXISTS idx_orders_created_at  ON orders(created_at);
CREATE INDEX IF NOT EXISTS idx_orders_deleted_at  ON orders(deleted_at);

CREATE TABLE IF NOT EXISTS order_items (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id      UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id    UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    name_snap     VARCHAR(255) NOT NULL,
    unit          VARCHAR(20) NOT NULL DEFAULT '',
    quantity      INTEGER NOT NULL DEFAULT 1,
    cost_price    NUMERIC(14,2) NOT NULL DEFAULT 0,
    sale_price    NUMERIC(14,2) NOT NULL DEFAULT 0,
    line_total    NUMERIC(14,2) NOT NULL DEFAULT 0,
    profit        NUMERIC(14,2) NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_order_items_order   ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_product ON order_items(product_id);

-- =========================================================
-- cart items
-- =========================================================
CREATE TABLE IF NOT EXISTS cart_items (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id  UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    product_id   UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity     INTEGER NOT NULL DEFAULT 1,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_cart_customer_product ON cart_items(customer_id, product_id) WHERE deleted_at IS NULL;

-- =========================================================
-- chat
-- =========================================================
CREATE TABLE IF NOT EXISTS chat_rooms (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id       UUID NOT NULL UNIQUE REFERENCES customers(id) ON DELETE CASCADE,
    last_message_at   TIMESTAMPTZ,
    unread_admin      INTEGER NOT NULL DEFAULT 0,
    unread_customer   INTEGER NOT NULL DEFAULT 0,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_chat_rooms_last_message ON chat_rooms(last_message_at);

CREATE TABLE IF NOT EXISTS chat_messages (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id       UUID NOT NULL REFERENCES chat_rooms(id) ON DELETE CASCADE,
    sender_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sender_role   VARCHAR(20) NOT NULL,
    body          TEXT NOT NULL DEFAULT '',
    read_at       TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_chat_messages_room    ON chat_messages(room_id);
CREATE INDEX IF NOT EXISTS idx_chat_messages_sender  ON chat_messages(sender_id);

CREATE TABLE IF NOT EXISTS chat_attachments (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id  UUID NOT NULL REFERENCES chat_messages(id) ON DELETE CASCADE,
    url         VARCHAR(500) NOT NULL,
    mime_type   VARCHAR(100) NOT NULL DEFAULT '',
    size_bytes  BIGINT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_chat_attachments_message ON chat_attachments(message_id);

-- =========================================================
-- notifications
-- =========================================================
CREATE TABLE IF NOT EXISTS notifications (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    kind        VARCHAR(32) NOT NULL,
    title_tj    VARCHAR(255) NOT NULL DEFAULT '',
    title_ru    VARCHAR(255) NOT NULL DEFAULT '',
    body_tj     TEXT NOT NULL DEFAULT '',
    body_ru     TEXT NOT NULL DEFAULT '',
    link_url    VARCHAR(500) NOT NULL DEFAULT '',
    read_at     TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_notifications_user    ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_kind    ON notifications(kind);
CREATE INDEX IF NOT EXISTS idx_notifications_read    ON notifications(read_at);

-- =========================================================
-- financial transactions
-- =========================================================
CREATE TABLE IF NOT EXISTS financial_transactions (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kind         VARCHAR(20) NOT NULL,
    order_id     UUID REFERENCES orders(id) ON DELETE SET NULL,
    product_id   UUID REFERENCES products(id) ON DELETE SET NULL,
    seller_id    UUID REFERENCES sellers(id) ON DELETE SET NULL,
    amount       NUMERIC(14,2) NOT NULL DEFAULT 0,
    currency     VARCHAR(8) NOT NULL DEFAULT 'TJS',
    description  VARCHAR(500) NOT NULL DEFAULT '',
    occurred_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_tx_kind     ON financial_transactions(kind);
CREATE INDEX IF NOT EXISTS idx_tx_order    ON financial_transactions(order_id);
CREATE INDEX IF NOT EXISTS idx_tx_seller   ON financial_transactions(seller_id);
CREATE INDEX IF NOT EXISTS idx_tx_occurred ON financial_transactions(occurred_at);

-- =========================================================
-- settings
-- =========================================================
CREATE TABLE IF NOT EXISTS settings (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key         VARCHAR(80) NOT NULL UNIQUE,
    value       TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

-- =========================================================
-- tracking events
-- =========================================================
CREATE TABLE IF NOT EXISTS tracking_events (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id   UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    seller_id    UUID NOT NULL REFERENCES sellers(id) ON DELETE CASCADE,
    user_id      UUID REFERENCES users(id) ON DELETE SET NULL,
    event        VARCHAR(32) NOT NULL,
    ip           VARCHAR(64)  NOT NULL DEFAULT '',
    user_agent   VARCHAR(500) NOT NULL DEFAULT '',
    referrer     VARCHAR(500) NOT NULL DEFAULT '',
    occurred_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_events_product   ON tracking_events(product_id);
CREATE INDEX IF NOT EXISTS idx_events_seller    ON tracking_events(seller_id);
CREATE INDEX IF NOT EXISTS idx_events_event     ON tracking_events(event);
CREATE INDEX IF NOT EXISTS idx_events_occurred  ON tracking_events(occurred_at);

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
CREATE UNIQUE INDEX IF NOT EXISTS idx_fav_user_product ON favorites(user_id, product_id) WHERE deleted_at IS NULL;
