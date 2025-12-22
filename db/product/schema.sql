CREATE TABLE IF NOT EXISTS categories (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  slug TEXT UNIQUE NOT NULL,
  parent_id BIGINT REFERENCES categories (id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS collections (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  slug TEXT UNIQUE NOT NULL,
  file TEXT,
  is_featured BOOLEAN DEFAULT FALSE,
  meta_title TEXT,
  meta_description TEXT,
  meta_keywords TEXT,
  canonical_url TEXT,
  og_title TEXT,
  og_description TEXT,
  og_image TEXT,
  layout TEXT,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS products (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  slug TEXT UNIQUE NOT NULL,
  origin_price INT NOT NULL DEFAULT 0,
  sale_price INT NOT NULL DEFAULT 0,
  stock INT DEFAULT 0,
  sku TEXT,
  weight INT DEFAULT 0,
  long INT DEFAULT 0,
  wide INT DEFAULT 0,
  high INT DEFAULT 0,
  meta_title TEXT NOT NULL DEFAULT '',
  meta_description TEXT NOT NULL DEFAULT '',
  meta_keywords TEXT NOT NULL DEFAULT '',
  canonical_url TEXT NOT NULL DEFAULT '',
  og_title TEXT NOT NULL DEFAULT '',
  og_description TEXT NOT NULL DEFAULT '',
  og_image TEXT NOT NULL DEFAULT '',
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  category_id BIGINT REFERENCES categories (id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS product_files (
  name TEXT,
  no INT NOT NULL,
  is_primary BOOLEAN NOT NULL DEFAULT FALSE,
  product_id BIGINT NOT NULL REFERENCES products (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS product_tags (
  name TEXT NOT NULL,
  product_id BIGINT NOT NULL REFERENCES products (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS options (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL DEFAULT '',
  no INT NOT NULL,
  product_id BIGINT NOT NULL REFERENCES products (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS option_values (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL DEFAULT '',
  no INT,
  option_id BIGINT NOT NULL REFERENCES options (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS variants (
  id BIGSERIAL PRIMARY KEY,
  origin_price INT NOT NULL DEFAULT 0,
  sale_price INT NOT NULL DEFAULT 0,
  file TEXT,
  stock INT NOT NULL DEFAULT 0,
  sku TEXT NOT NULL DEFAULT '',
  no INT NOT NULL DEFAULT 0,
  product_id BIGINT NOT NULL REFERENCES products (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS variant_options (
  variant_id BIGINT NOT NULL REFERENCES variants (id) ON DELETE CASCADE,
  option_value_id BIGINT NOT NULL REFERENCES option_values (id) ON DELETE CASCADE,
  option_id BIGINT NOT NULL REFERENCES options (id) ON DELETE CASCADE,
  PRIMARY KEY (variant_id, option_id)
);

CREATE TABLE IF NOT EXISTS product_collections (
  product_id BIGINT NOT NULL REFERENCES products (id) ON DELETE CASCADE,
  collection_id BIGINT NOT NULL REFERENCES collections (id) ON DELETE CASCADE,
  PRIMARY KEY (product_id, collection_id)
);

CREATE TABLE IF NOT EXISTS orders (
  id BIGSERIAL PRIMARY KEY,
  code TEXT NOT NULL DEFAULT '',
  total_amount INT NOT NULL DEFAULT 0,
  discount_amount INT NOT NULL DEFAULT 0,
  shipping_fee_amount INT DEFAULT 0,
  shipping_address_id BIGINT REFERENCES addresses (id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS order_items (
  id BIGSERIAL PRIMARY KEY,
  quantity INT NOT NULL DEFAULT 0,
  sale_price INT NOT NULL DEFAULT 0,
  order_id BIGINT NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
  product_id BIGINT REFERENCES products (id) ON DELETE SET NULL,
  variant_id BIGINT REFERENCES variants (id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS addresses (
  id BIGSERIAL PRIMARY KEY,
  user_id INT,
  full_name TEXT NOT NULL,
  address_line TEXT NOT NULL,
  city TEXT,
  state TEXT,
  country TEXT,
  postal_code TEXT,
  phone TEXT,
  email TEXT
);

CREATE TABLE IF NOT EXISTS customers (
  id BIGSERIAL PRIMARY KEY,
  name TEXT UNIQUE NOT NULL,
  phone TEXT UNIQUE NOT NULL,
  phone_verified  BOOLEAN DEFAULT FALSE,
  zns_otp TEXT,
  avatar TEXT,
  email TEXT,
  password TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS reviews (
  id BIGSERIAL PRIMARY KEY,
  product_id BIGINT REFERENCES products (id) ON DELETE CASCADE,
  rating INT CHECK (rating BETWEEN 1 AND 5),
  comment TEXT,
  has_file BOOLEAN NOT NULL DEFAULT FALSE,
  customer_id BIGINT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS review_files (
  name TEXT NOT NULL,
  review_id BIGINT NOT NULL REFERENCES reviews (id) on DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS files (
  id BIGSERIAL PRIMARY KEY,
  name TEXT UNIQUE NOT NULL,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pages (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  slug TEXT UNIQUE NOT NULL,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS menus (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  position TEXT UNIQUE NOT NULL,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE discounts (
  id BIGSERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT,
  code TEXT UNIQUE, -- nullable for automatic discounts
  discount_type TEXT NOT NULL, -- 'code' | 'automatic'
  status TEXT NOT NULL, -- 'draft' | 'active' | 'scheduled' | 'expired'
  usage_limit INTEGER, -- limit times used
  usage_count INTEGER DEFAULT 0,
  per_customer_limit INTEGER, -- optional per-customer limit
  starts_at TIMESTAMPTZ NOT NULL,
  ends_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE discount_conditions (
  id BIGSERIAL PRIMARY KEY,
  discount_id BIGINT NOT NULL,
  condition_type TEXT NOT NULL, -- 'specific_products' | 'specific_collections' | 'order_amount' | 'quantity' | 'customer'
  operator TEXT NOT NULL, -- 'eq' | 'gt' | 'gte' | 'lt' | 'lte' | 'in' | 'not_in'
  value TEXT NOT NULL, -- JSON (product IDs, min_subtotal, min_qty, etc)
  FOREIGN KEY (discount_id) REFERENCES discounts (id) ON DELETE CASCADE
);

CREATE TABLE discount_effects (
  id BIGSERIAL PRIMARY KEY,
  discount_id BIGINT NOT NULL,
  effect_type TEXT NOT NULL, -- 'percent' | 'fixed' | 'free_shipping' | 'bogo'
  value TEXT, -- percent: "20", fixed: "5", bogo: JSON {"buy":1, "get":1}
  applies_to TEXT NOT NULL, -- 'entire_order' | 'specific_products' | 'specific_collections'
  FOREIGN KEY (discount_id) REFERENCES discounts (id) ON DELETE CASCADE
);

CREATE TABLE discount_targets (
  id BIGSERIAL PRIMARY KEY,
  discount_id BIGINT NOT NULL,
  target_type TEXT NOT NULL, -- 'specific_products' | 'specific_collections'
  target_id INTEGER NOT NULL,
  FOREIGN KEY (discount_id) REFERENCES discounts (id) ON DELETE CASCADE
);

CREATE TABLE discount_customer_usages (
  id BIGSERIAL PRIMARY KEY,
  discount_id BIGINT NOT NULL REFERENCES discounts(id) ON DELETE CASCADE,
  customer_id BIGINT NOT NULL,
  used_count INTEGER DEFAULT 0,
  UNIQUE (discount_id, customer_id)
);

CREATE TABLE hotspots (
  id BIGSERIAL PRIMARY KEY,
  file TEXT NOT NULL
);

CREATE TABLE product_hotspots (
  id BIGSERIAL PRIMARY KEY,
  product_id BIGINT NOT NULL,
  hotspot_id BIGINT NOT NULL,
  x REAL NOT NULL,
  y REAL NOT NULL,
  FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
  FOREIGN KEY (hotspot_id) REFERENCES hotspots(id) ON DELETE CASCADE
);

CREATE TABLE shipping_fees (
  id BIGSERIAL PRIMARY KEY,
  name TEXT,
  min_weight INTEGER NOT NULL DEFAULT 0,
  max_weight INTEGER NOT NULL DEFAULT 0,
  fee_amount INTEGER NOT NULL DEFAULT 0,
  min_order_value INTEGER,
  free_shipping BOOLEAN DEFAULT FALSE,
  shipping_method TEXT,
  effective_from DATE,
  effective_to DATE,
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
