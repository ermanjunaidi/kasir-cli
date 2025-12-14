-- Database Schema untuk Sistem Kasir v2.0
-- Dengan fitur: multi-gudang, user auth, harga beli/jual

-- Hapus tabel jika sudah ada (untuk fresh install)
DROP TABLE IF EXISTS transaction_items CASCADE;
DROP TABLE IF EXISTS transactions CASCADE;
DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS warehouses CASCADE;

-- Tabel Gudang
CREATE TABLE warehouses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabel User
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    warehouse_id INT REFERENCES warehouses(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabel Produk dengan harga beli & jual
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    purchase_price DECIMAL(10,2) NOT NULL,
    selling_price DECIMAL(10,2) NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    warehouse_id INT REFERENCES warehouses(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabel Transaksi
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    warehouse_id INT REFERENCES warehouses(id),
    total DECIMAL(10,2) NOT NULL,
    profit DECIMAL(10,2) NOT NULL DEFAULT 0,
    payment DECIMAL(10,2) NOT NULL,
    change DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabel Detail Transaksi
CREATE TABLE transaction_items (
    id SERIAL PRIMARY KEY,
    transaction_id INT REFERENCES transactions(id) ON DELETE CASCADE,
    product_id INT REFERENCES products(id),
    product_name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    purchase_price DECIMAL(10,2) NOT NULL,
    selling_price DECIMAL(10,2) NOT NULL,
    subtotal DECIMAL(10,2) NOT NULL,
    profit DECIMAL(10,2) NOT NULL DEFAULT 0
);

-- Index untuk performa
CREATE INDEX idx_transactions_created_at ON transactions(created_at);
CREATE INDEX idx_transactions_warehouse_id ON transactions(warehouse_id);
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transaction_items_transaction_id ON transaction_items(transaction_id);
CREATE INDEX idx_products_warehouse_id ON products(warehouse_id);
CREATE INDEX idx_users_warehouse_id ON users(warehouse_id);

-- Sample data gudang
INSERT INTO warehouses (name, address) VALUES
    ('Gudang Pusat', 'Jl. Utama No. 1'),
    ('Gudang Cabang A', 'Jl. Cabang A No. 10'),
    ('Gudang Cabang B', 'Jl. Cabang B No. 20');

-- Sample admin user (password: admin123)
INSERT INTO users (username, password, role, warehouse_id) VALUES
    ('admin', 'admin123', 'admin', NULL);

-- Sample user per gudang (password: user123)
INSERT INTO users (username, password, role, warehouse_id) VALUES
    ('kasir1', 'user123', 'user', 1),
    ('kasir2', 'user123', 'user', 2),
    ('kasir3', 'user123', 'user', 3);

-- Sample data produk per gudang
INSERT INTO products (name, purchase_price, selling_price, stock, warehouse_id) VALUES
    ('Indomie Goreng', 2500, 3500, 100, 1),
    ('Aqua 600ml', 2500, 4000, 50, 1),
    ('Teh Botol Sosro', 3500, 5000, 30, 1),
    ('Roti Tawar', 12000, 15000, 20, 1),
    ('Susu Ultra 250ml', 4500, 6000, 40, 1),
    ('Indomie Goreng', 2500, 3500, 80, 2),
    ('Aqua 600ml', 2500, 4000, 40, 2),
    ('Teh Botol Sosro', 3500, 5000, 25, 2),
    ('Indomie Goreng', 2500, 3500, 60, 3),
    ('Aqua 600ml', 2500, 4000, 35, 3);
