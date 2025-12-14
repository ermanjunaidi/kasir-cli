package handlers

import (
	"bufio"
	"fmt"
	"kasir/models"
	"os"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

// LoginMenu menampilkan menu login
func LoginMenu() bool {
	fmt.Println("\n╔══════════════════════════════════════╗")
	fmt.Println("║              LOGIN                   ║")
	fmt.Println("╚══════════════════════════════════════╝")

	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	user, err := models.Login(username, password)
	if err != nil {
		fmt.Println("❌", err)
		return false
	}

	warehouseInfo := "Semua Gudang"
	if user.WarehouseID != nil {
		warehouse, _ := models.GetWarehouseByID(*user.WarehouseID)
		if warehouse != nil {
			warehouseInfo = warehouse.Name
		}
	}

	fmt.Printf("\n✅ Login berhasil! Selamat datang, %s (%s)\n", user.Username, user.Role)
	fmt.Printf("📦 Gudang: %s\n", warehouseInfo)
	return true
}

// UserMenu menampilkan menu manajemen user (admin only)
func UserMenu() {
	for {
		fmt.Println("\n╔══════════════════════════════════════╗")
		fmt.Println("║        MANAJEMEN USER                ║")
		fmt.Println("╠══════════════════════════════════════╣")
		fmt.Println("║  1. Lihat Daftar User                ║")
		fmt.Println("║  2. Tambah User Baru                 ║")
		fmt.Println("║  3. Hapus User                       ║")
		fmt.Println("║  0. Kembali ke Menu Utama            ║")
		fmt.Println("╚══════════════════════════════════════╝")

		fmt.Print("Pilihan: ")
		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		switch choice {
		case "1":
			listUsers()
		case "2":
			registerUser()
		case "3":
			deleteUser()
		case "0":
			return
		default:
			fmt.Println("❌ Pilihan tidak valid!")
		}
	}
}

func listUsers() {
	users, err := models.GetAllUsers()
	if err != nil {
		fmt.Println("❌ Error:", err)
		return
	}

	fmt.Println("\n┌─────┬────────────────────┬──────────┬────────────────────────┐")
	fmt.Println("│ ID  │ Username           │ Role     │ Gudang                 │")
	fmt.Println("├─────┼────────────────────┼──────────┼────────────────────────┤")

	for _, u := range users {
		warehouseName := "Semua Gudang"
		if u.WarehouseID != nil {
			warehouse, _ := models.GetWarehouseByID(*u.WarehouseID)
			if warehouse != nil {
				warehouseName = warehouse.Name
			}
		}
		fmt.Printf("│ %-3d │ %-18s │ %-8s │ %-22s │\n", u.ID, u.Username, u.Role, warehouseName)
	}
	fmt.Println("└─────┴────────────────────┴──────────┴────────────────────────┘")
}

func registerUser() {
	fmt.Println("\n═══ TAMBAH USER BARU ═══")

	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	fmt.Print("Role (admin/user): ")
	role, _ := reader.ReadString('\n')
	role = strings.TrimSpace(role)

	if role != "admin" && role != "user" {
		fmt.Println("❌ Role harus 'admin' atau 'user'!")
		return
	}

	var warehouseID *int
	if role == "user" {
		// Tampilkan daftar gudang
		warehouses, err := models.GetAllWarehouses()
		if err != nil {
			fmt.Println("❌ Error:", err)
			return
		}

		fmt.Println("\nDaftar Gudang:")
		for _, w := range warehouses {
			fmt.Printf("  %d. %s\n", w.ID, w.Name)
		}

		fmt.Print("Pilih ID Gudang: ")
		var wID int
		fmt.Scanln(&wID)
		reader.ReadString('\n')
		warehouseID = &wID
	}

	user, err := models.Register(username, password, role, warehouseID)
	if err != nil {
		fmt.Println("❌", err)
		return
	}

	fmt.Printf("✅ User '%s' berhasil ditambahkan dengan ID: %d\n", user.Username, user.ID)
}

func deleteUser() {
	listUsers()

	fmt.Print("\nMasukkan ID user yang akan dihapus (0 untuk batal): ")
	var id int
	fmt.Scanln(&id)
	reader.ReadString('\n')

	if id == 0 {
		return
	}

	// Jangan izinkan hapus diri sendiri
	if models.CurrentUser != nil && models.CurrentUser.ID == id {
		fmt.Println("❌ Tidak bisa menghapus akun sendiri!")
		return
	}

	err := models.DeleteUser(id)
	if err != nil {
		fmt.Println("❌", err)
		return
	}

	fmt.Println("✅ User berhasil dihapus!")
}
