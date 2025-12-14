-- Script untuk menambahkan 1 juta produk
-- Menggunakan generate_series untuk performa optimal

INSERT INTO products (name, price, stock)
SELECT 
    'Produk-' || i AS name,
    (random() * 99000 + 1000)::DECIMAL(10,2) AS price,
    (random() * 100 + 1)::INT AS stock
FROM generate_series(1, 1000000) AS i;
