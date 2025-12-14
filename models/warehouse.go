package models

import (
	"errors"
	"kasir/config"
	"time"
)

// Warehouse model
type Warehouse struct {
	ID        int
	Name      string
	Address   string
	CreatedAt time.Time
}

// GetAllWarehouses mengambil semua gudang
func GetAllWarehouses() ([]Warehouse, error) {
	rows, err := config.DB.Query(`
		SELECT id, name, COALESCE(address, ''), created_at 
		FROM warehouses 
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var warehouses []Warehouse
	for rows.Next() {
		var w Warehouse
		err := rows.Scan(&w.ID, &w.Name, &w.Address, &w.CreatedAt)
		if err != nil {
			return nil, err
		}
		warehouses = append(warehouses, w)
	}
	return warehouses, nil
}

// GetWarehouseByID mengambil gudang berdasarkan ID
func GetWarehouseByID(id int) (*Warehouse, error) {
	var w Warehouse
	err := config.DB.QueryRow(`
		SELECT id, name, COALESCE(address, ''), created_at 
		FROM warehouses 
		WHERE id = $1
	`, id).Scan(&w.ID, &w.Name, &w.Address, &w.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &w, nil
}

// CreateWarehouse membuat gudang baru
func CreateWarehouse(name, address string) (*Warehouse, error) {
	var w Warehouse
	err := config.DB.QueryRow(`
		INSERT INTO warehouses (name, address) 
		VALUES ($1, $2) 
		RETURNING id, name, COALESCE(address, ''), created_at
	`, name, address).Scan(&w.ID, &w.Name, &w.Address, &w.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &w, nil
}

// UpdateWarehouse mengupdate gudang
func UpdateWarehouse(id int, name, address string) error {
	result, err := config.DB.Exec(`
		UPDATE warehouses 
		SET name = $1, address = $2 
		WHERE id = $3
	`, name, address, id)

	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("gudang tidak ditemukan")
	}
	return nil
}

// DeleteWarehouse menghapus gudang
func DeleteWarehouse(id int) error {
	result, err := config.DB.Exec(`DELETE FROM warehouses WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("gudang tidak ditemukan")
	}
	return nil
}
