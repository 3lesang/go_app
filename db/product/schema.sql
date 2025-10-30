CREATE TABLE IF NOT EXISTS categories (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  slug TEXT UNIQUE NOT NULL,
  parent_id BIGINT REFERENCES categories (id) ON DELETE SET NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS collections (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  slug TEXT UNIQUE NOT NULL,
  image_url TEXT,
  is_featured BOOLEAN DEFAULT FALSE,
  meta_title TEXT,
  meta_description TEXT,
  meta_keywords TEXT,
  canonical_url TEXT,
  og_title TEXT,
  og_description TEXT,
  og_image TEXT,
  layout TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS products (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  slug TEXT UNIQUE NOT NULL,
  origin_price INT NOT NULL DEFAULT 0,
  sale_price INT NOT NULL DEFAULT 0,
  meta_title TEXT NOT NULL DEFAULT '',
  meta_description TEXT NOT NULL DEFAULT '',
  meta_keywords TEXT NOT NULL DEFAULT '',
  canonical_url TEXT NOT NULL DEFAULT '',
  og_title TEXT NOT NULL DEFAULT '',
  og_description TEXT NOT NULL DEFAULT '',
  og_image TEXT NOT NULL DEFAULT '',
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  category_id BIGINT REFERENCES categories (id) ON DELETE SET NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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
  total_amount INT NOT NULL DEFAULT 0,
  discount_amount INT NOT NULL DEFAULT 0,
  shipping_address_id BIGINT REFERENCES addresses (id) ON DELETE SET NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS order_items (
  id BIGSERIAL PRIMARY KEY,
  quantity INT NOT NULL DEFAULT 0,
  sale_price INT NOT NULL DEFAULT 0,
  order_id BIGINT NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
  product_id BIGINT NOT NULL REFERENCES products (id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS addresses (
  id BIGSERIAL PRIMARY KEY,
  user_id INT,
  full_name TEXT NOT NULL,
  address_line TEXT NOT NULL,
  city TEXT NOT NULL,
  state TEXT,
  country TEXT NOT NULL,
  postal_code TEXT,
  phone TEXT
);

CREATE TABLE IF NOT EXISTS reviews (
  id BIGSERIAL PRIMARY KEY,
  product_id BIGINT REFERENCES products (id) ON DELETE CASCADE,
  rating INT CHECK (rating BETWEEN 1 AND 5),
  comment TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS coupons (
  id BIGSERIAL PRIMARY KEY,
  code TEXT UNIQUE NOT NULL,
  description TEXT,
  discount_percent NUMBERIC(5, 2) CHECK (discount_percent BETWEEN 0 AND 100),
  valid_from TIMESTAMP,
  valid_until TIMESTAMP,
  is_active BOOLEAN DEFAULT TRUE
);
