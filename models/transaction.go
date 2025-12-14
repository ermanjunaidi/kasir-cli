package models

import (
	"kasir/config"
	"time"
)

// Transaction model
type Transaction struct {
	ID          int
	UserID      int
	WarehouseID int
	Total       float64
	Profit      float64
	Payment     float64
	Change      float64
	CreatedAt   time.Time
	Items       []TransactionItem
}

// TransactionItem model
type TransactionItem struct {
	ID            int
	TransactionID int
	ProductID     int
	ProductName   string
	Quantity      int
	PurchasePrice float64
	SellingPrice  float64
	Subtotal      float64
	Profit        float64
}

// CartItem untuk keranjang belanja
type CartItem struct {
	Product  *Product
	Quantity int
}

// CreateTransaction membuat transaksi baru dengan items
func CreateTransaction(items []CartItem, payment float64) (*Transaction, error) {
	// Hitung total dan profit
	var total float64
	var totalProfit float64
	for _, item := range items {
		subtotal := item.Product.SellingPrice * float64(item.Quantity)
		profit := (item.Product.SellingPrice - item.Product.PurchasePrice) * float64(item.Quantity)
		total += subtotal
		totalProfit += profit
	}
	change := payment - total

	// Get user dan warehouse info
	userID := 0
	warehouseID := 0
	if CurrentUser != nil {
		userID = CurrentUser.ID
		if CurrentUser.WarehouseID != nil {
			warehouseID = *CurrentUser.WarehouseID
		} else if len(items) > 0 {
			warehouseID = items[0].Product.WarehouseID
		}
	}

	// Mulai transaction database
	tx, err := config.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert transaksi
	var transactionID int
	var createdAt time.Time
	err = tx.QueryRow(`
		INSERT INTO transactions (user_id, warehouse_id, total, profit, payment, change) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id, created_at
	`, userID, warehouseID, total, totalProfit, payment, change).Scan(&transactionID, &createdAt)
	if err != nil {
		return nil, err
	}

	// Insert items dan update stok
	transaction := &Transaction{
		ID:          transactionID,
		UserID:      userID,
		WarehouseID: warehouseID,
		Total:       total,
		Profit:      totalProfit,
		Payment:     payment,
		Change:      change,
		CreatedAt:   createdAt,
	}

	for _, item := range items {
		subtotal := item.Product.SellingPrice * float64(item.Quantity)
		profit := (item.Product.SellingPrice - item.Product.PurchasePrice) * float64(item.Quantity)

		// Insert transaction item
		_, err = tx.Exec(`
			INSERT INTO transaction_items 
			(transaction_id, product_id, product_name, quantity, purchase_price, selling_price, subtotal, profit) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, transactionID, item.Product.ID, item.Product.Name, item.Quantity,
			item.Product.PurchasePrice, item.Product.SellingPrice, subtotal, profit)
		if err != nil {
			return nil, err
		}

		// Update stok
		result, err := tx.Exec(`
			UPDATE products 
			SET stock = stock - $1 
			WHERE id = $2 AND stock >= $1
		`, item.Quantity, item.Product.ID)
		if err != nil {
			return nil, err
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			return nil, err
		}

		transaction.Items = append(transaction.Items, TransactionItem{
			ProductID:     item.Product.ID,
			ProductName:   item.Product.Name,
			Quantity:      item.Quantity,
			PurchasePrice: item.Product.PurchasePrice,
			SellingPrice:  item.Product.SellingPrice,
			Subtotal:      subtotal,
			Profit:        profit,
		})
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// GetTransactionsByDate mengambil transaksi berdasarkan tanggal
func GetTransactionsByDate(date time.Time) ([]Transaction, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var query string
	var args []interface{}

	if CurrentUser != nil && !CurrentUser.IsAdmin() && CurrentUser.WarehouseID != nil {
		query = `
			SELECT id, user_id, warehouse_id, total, profit, payment, change, created_at 
			FROM transactions 
			WHERE created_at >= $1 AND created_at < $2 AND warehouse_id = $3
			ORDER BY created_at DESC
		`
		args = []interface{}{startOfDay, endOfDay, *CurrentUser.WarehouseID}
	} else {
		query = `
			SELECT id, user_id, warehouse_id, total, profit, payment, change, created_at 
			FROM transactions 
			WHERE created_at >= $1 AND created_at < $2 
			ORDER BY created_at DESC
		`
		args = []interface{}{startOfDay, endOfDay}
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.ID, &t.UserID, &t.WarehouseID, &t.Total, &t.Profit, &t.Payment, &t.Change, &t.CreatedAt)
		if err != nil {
			return nil, err
		}

		// Get items for this transaction
		t.Items, err = GetTransactionItems(t.ID)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, t)
	}
	return transactions, nil
}

// GetTransactionItems mengambil item-item transaksi
func GetTransactionItems(transactionID int) ([]TransactionItem, error) {
	rows, err := config.DB.Query(`
		SELECT id, transaction_id, product_id, product_name, quantity, purchase_price, selling_price, subtotal, profit 
		FROM transaction_items 
		WHERE transaction_id = $1
	`, transactionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []TransactionItem
	for rows.Next() {
		var item TransactionItem
		err := rows.Scan(&item.ID, &item.TransactionID, &item.ProductID,
			&item.ProductName, &item.Quantity, &item.PurchasePrice, &item.SellingPrice, &item.Subtotal, &item.Profit)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// GetDailyTotal mengambil total penjualan harian
func GetDailyTotal(date time.Time) (float64, float64, int, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var query string
	var args []interface{}

	if CurrentUser != nil && !CurrentUser.IsAdmin() && CurrentUser.WarehouseID != nil {
		query = `
			SELECT COALESCE(SUM(total), 0), COALESCE(SUM(profit), 0), COUNT(*) 
			FROM transactions 
			WHERE created_at >= $1 AND created_at < $2 AND warehouse_id = $3
		`
		args = []interface{}{startOfDay, endOfDay, *CurrentUser.WarehouseID}
	} else {
		query = `
			SELECT COALESCE(SUM(total), 0), COALESCE(SUM(profit), 0), COUNT(*) 
			FROM transactions 
			WHERE created_at >= $1 AND created_at < $2
		`
		args = []interface{}{startOfDay, endOfDay}
	}

	var total, profit float64
	var count int
	err := config.DB.QueryRow(query, args...).Scan(&total, &profit, &count)

	return total, profit, count, err
}
