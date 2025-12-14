package models

import (
	"errors"
	"kasir/config"
	"time"
)

// User model
type User struct {
	ID          int
	Username    string
	Password    string
	Role        string // "admin" atau "user"
	WarehouseID *int   // nil untuk admin (akses semua gudang)
	CreatedAt   time.Time
}

// CurrentUser menyimpan user yang sedang login
var CurrentUser *User

// Login melakukan validasi login
func Login(username, password string) (*User, error) {
	var u User
	var warehouseID *int

	err := config.DB.QueryRow(`
		SELECT id, username, password, role, warehouse_id, created_at 
		FROM users 
		WHERE username = $1 AND password = $2
	`, username, password).Scan(&u.ID, &u.Username, &u.Password, &u.Role, &warehouseID, &u.CreatedAt)

	if err != nil {
		return nil, errors.New("username atau password salah")
	}

	u.WarehouseID = warehouseID
	CurrentUser = &u
	return &u, nil
}

// Logout menghapus current user session
func Logout() {
	CurrentUser = nil
}

// Register membuat user baru
func Register(username, password, role string, warehouseID *int) (*User, error) {
	var u User
	err := config.DB.QueryRow(`
		INSERT INTO users (username, password, role, warehouse_id) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, username, password, role, warehouse_id, created_at
	`, username, password, role, warehouseID).Scan(&u.ID, &u.Username, &u.Password, &u.Role, &u.WarehouseID, &u.CreatedAt)

	if err != nil {
		return nil, errors.New("gagal membuat user, username mungkin sudah digunakan")
	}
	return &u, nil
}

// GetAllUsers mengambil semua user
func GetAllUsers() ([]User, error) {
	rows, err := config.DB.Query(`
		SELECT u.id, u.username, u.role, u.warehouse_id, u.created_at,
			   COALESCE(w.name, '-') as warehouse_name
		FROM users u
		LEFT JOIN warehouses w ON u.warehouse_id = w.id
		ORDER BY u.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		var warehouseName string
		err := rows.Scan(&u.ID, &u.Username, &u.Role, &u.WarehouseID, &u.CreatedAt, &warehouseName)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// GetUserByID mengambil user berdasarkan ID
func GetUserByID(id int) (*User, error) {
	var u User
	err := config.DB.QueryRow(`
		SELECT id, username, password, role, warehouse_id, created_at 
		FROM users 
		WHERE id = $1
	`, id).Scan(&u.ID, &u.Username, &u.Password, &u.Role, &u.WarehouseID, &u.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &u, nil
}

// DeleteUser menghapus user
func DeleteUser(id int) error {
	result, err := config.DB.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("user tidak ditemukan")
	}
	return nil
}

// IsAdmin mengecek apakah user adalah admin
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// GetWarehouseID mengembalikan warehouse_id user (0 jika admin/nil)
func (u *User) GetWarehouseID() int {
	if u.WarehouseID == nil {
		return 0
	}
	return *u.WarehouseID
}
