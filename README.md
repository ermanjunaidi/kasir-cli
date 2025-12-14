# 🛒 Sistem Kasir CLI v2.0

Aplikasi kasir berbasis Command Line Interface (CLI) menggunakan **Go** dan **PostgreSQL**.

## 📋 Fitur

- ✅ **Autentikasi** - Login/Register dengan role admin/user
- ✅ **Multi-Gudang** - User hanya akses gudang tertentu
- ✅ **Harga Beli/Jual** - Track profit per transaksi
- ✅ **Manajemen Produk** - CRUD dengan warehouse filter
- ✅ **Transaksi Penjualan** - Keranjang & struk pembayaran
- ✅ **Laporan** - Penjualan harian dengan profit
- ✅ **Lihat Stok Semua Gudang** - Admin bisa lihat ringkasan stok
- ✅ **Export/Import Excel** - Export & import data produk ke Excel

## 🔧 Prasyarat

- [Go](https://golang.org/dl/) versi 1.19+
- [PostgreSQL](https://www.postgresql.org/download/) versi 12+

## 🚀 Instalasi

### 1. Clone & Setup Database

```bash
# Clone repository
git clone <repository-url>
cd kasir

# Buat database
sudo -u postgres psql -c "CREATE DATABASE kasir;"

# Jalankan migration
PGPASSWORD=123123 psql -U postgres -d kasir -h localhost -f migrations/init.sql
```

### 2. Konfigurasi (Opsional)

Set environment variables jika berbeda dari default:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=123123
export DB_NAME=kasir
```

### 3. Jalankan Aplikasi

```bash
go run main.go
```

## 🔑 Akun Default

| Username | Password | Role | Gudang |
|----------|----------|------|--------|
| admin | admin123 | admin | Semua |
| kasir1 | user123 | user | Gudang Pusat |
| kasir2 | user123 | user | Gudang Cabang A |
| kasir3 | user123 | user | Gudang Cabang B |

## 📖 Role & Permissions

### Admin
- Transaksi (semua gudang)
- Manajemen Produk (semua gudang)
- Laporan (semua gudang)
- Manajemen User
- Manajemen Gudang

### User (Kasir)
- Transaksi (gudang sendiri)
- Lihat Produk (gudang sendiri)
- Laporan (gudang sendiri)

## 📁 Struktur Proyek

```
kasir/
├── config/database.go      # Konfigurasi DB
├── handlers/
│   ├── auth.go             # Login & user management
│   ├── warehouse.go        # Warehouse management
│   ├── product.go          # Product CRUD
│   ├── transaction.go      # Sales transactions
│   └── report.go           # Sales reports
├── migrations/init.sql     # Database schema
├── models/
│   ├── user.go             # User model
│   ├── warehouse.go        # Warehouse model
│   ├── product.go          # Product model
│   └── transaction.go      # Transaction model
└── main.go                 # Entry point
```

## 🛠️ Troubleshooting

| Error | Solusi |
|-------|--------|
| relation does not exist | Jalankan `migrations/init.sql` |
| connection refused | Pastikan PostgreSQL berjalan |
| password authentication failed | Cek password di env variable |

## 📝 License

MIT License
