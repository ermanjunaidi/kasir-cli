package models

import (
	"fmt"
	"kasir/config"
	"time"
)

// Product model
type Product struct {
	ID            int
	Name          string
	PurchasePrice float64 // Harga Beli
	SellingPrice  float64 // Harga Jual
	Stock         int
	WarehouseID   int
	CreatedAt     time.Time
}

// GetAllProducts mengambil semua produk (filter by warehouse jika user biasa)
func GetAllProducts() ([]Product, error) {
	var query string
	var args []interface{}

	if CurrentUser != nil && !CurrentUser.IsAdmin() && CurrentUser.WarehouseID != nil {
		query = `
			SELECT id, name, purchase_price, selling_price, stock, warehouse_id, created_at 
			FROM products 
			WHERE warehouse_id = $1
			ORDER BY id
		`
		args = append(args, *CurrentUser.WarehouseID)
	} else {
		query = `
			SELECT id, name, purchase_price, selling_price, stock, warehouse_id, created_at 
			FROM products 
			ORDER BY id
		`
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.PurchasePrice, &p.SellingPrice, &p.Stock, &p.WarehouseID, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

// GetProductsByWarehouse mengambil produk berdasarkan warehouse
func GetProductsByWarehouse(warehouseID int) ([]Product, error) {
	rows, err := config.DB.Query(`
		SELECT id, name, purchase_price, selling_price, stock, warehouse_id, created_at 
		FROM products 
		WHERE warehouse_id = $1
		ORDER BY id
	`, warehouseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.PurchasePrice, &p.SellingPrice, &p.Stock, &p.WarehouseID, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

// GetProductByID mengambil produk berdasarkan ID
func GetProductByID(id int) (*Product, error) {
	var p Product
	err := config.DB.QueryRow(`
		SELECT id, name, purchase_price, selling_price, stock, warehouse_id, created_at 
		FROM products 
		WHERE id = $1
	`, id).Scan(&p.ID, &p.Name, &p.PurchasePrice, &p.SellingPrice, &p.Stock, &p.WarehouseID, &p.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &p, nil
}

// CreateProduct membuat produk baru
func CreateProduct(name string, purchasePrice, sellingPrice float64, stock, warehouseID int) (*Product, error) {
	var p Product
	err := config.DB.QueryRow(`
		INSERT INTO products (name, purchase_price, selling_price, stock, warehouse_id) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, name, purchase_price, selling_price, stock, warehouse_id, created_at
	`, name, purchasePrice, sellingPrice, stock, warehouseID).Scan(
		&p.ID, &p.Name, &p.PurchasePrice, &p.SellingPrice, &p.Stock, &p.WarehouseID, &p.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &p, nil
}

// UpdateProduct mengupdate produk
func UpdateProduct(id int, name string, purchasePrice, sellingPrice float64, stock int) error {
	result, err := config.DB.Exec(`
		UPDATE products 
		SET name = $1, purchase_price = $2, selling_price = $3, stock = $4 
		WHERE id = $5
	`, name, purchasePrice, sellingPrice, stock, id)

	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("produk dengan ID %d tidak ditemukan", id)
	}
	return nil
}

// DeleteProduct menghapus produk
func DeleteProduct(id int) error {
	result, err := config.DB.Exec(`DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("produk dengan ID %d tidak ditemukan", id)
	}
	return nil
}

// UpdateStock mengupdate stok produk
func UpdateStock(id int, quantity int) error {
	result, err := config.DB.Exec(`
		UPDATE products 
		SET stock = stock - $1 
		WHERE id = $2 AND stock >= $1
	`, quantity, id)

	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("stok tidak mencukupi atau produk tidak ditemukan")
	}
	return nil
}

// GetProfit menghitung profit per item
func (p *Product) GetProfit() float64 {
	return p.SellingPrice - p.PurchasePrice
}
