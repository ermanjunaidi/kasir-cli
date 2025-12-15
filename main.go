package main

import (
	"bufio"
	"flag"
	"fmt"
	"kasir/api"
	"kasir/config"
	"kasir/handlers"
	"kasir/models"
	"os"
	"strings"
)

func main() {
	// Parse flags
	apiMode := flag.Bool("api", false, "Run in API mode")
	port := flag.String("port", "8080", "Port for API server")
	flag.Parse()

	// Banner
	printBanner()

	// Initialize database
	fmt.Println("🔌 Menghubungkan ke database...")
	err := config.InitDB()
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		fmt.Println("\n💡 Pastikan:")
		fmt.Println("   1. PostgreSQL sudah berjalan")
		fmt.Println("   2. Database 'kasir' sudah dibuat")
		fmt.Println("   3. Jalankan migrations/init.sql")
		fmt.Println("\n   Atau set environment variables:")
		fmt.Println("   DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME")
		os.Exit(1)
	}
	defer config.CloseDB()
	fmt.Println("✅ Koneksi database berhasil!")

	// Check if API mode
	if *apiMode {
		api.StartServer(*port)
		return
	}

	// CLI Mode below...
	// Login loop
	reader := bufio.NewReader(os.Stdin)
	for {
		if !handlers.LoginMenu() {
			fmt.Print("\nCoba lagi? (y/n): ")
			input, _ := reader.ReadString('\n')
			if strings.ToLower(strings.TrimSpace(input)) != "y" {
				fmt.Println("\n👋 Sampai jumpa!")
				return
			}
			continue
		}
		break
	}

	// Main loop berdasarkan role
	for {
		if models.CurrentUser.IsAdmin() {
			printAdminMenu()
		} else {
			printUserMenu()
		}

		fmt.Print("Pilihan: ")
		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		if models.CurrentUser.IsAdmin() {
			switch choice {
			case "1":
				handlers.TransactionMenu()
			case "2":
				handlers.ProductMenu()
			case "3":
				handlers.ReportMenu()
			case "4":
				handlers.UserMenu()
			case "5":
				handlers.WarehouseMenu()
			case "6":
				handlers.ChangePassword()
			case "0":
				logout()
				return
			default:
				fmt.Println("❌ Pilihan tidak valid!")
			}
		} else {
			switch choice {
			case "1":
				handlers.TransactionMenu()
			case "2":
				listProductsOnly()
			case "3":
				handlers.ReportMenu()
			case "4":
				handlers.ChangePassword()
			case "0":
				logout()
				return
			default:
				fmt.Println("❌ Pilihan tidak valid!")
			}
		}
	}
}

func printBanner() {
	fmt.Println()
	fmt.Println("╔═════════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                             ║")
	fmt.Println("║   ██╗  ██╗ █████╗ ███████╗██╗██████╗                        ║")
	fmt.Println("║   ██║ ██╔╝██╔══██╗██╔════╝██║██╔══██╗                       ║")
	fmt.Println("║   █████╔╝ ███████║███████╗██║██████╔╝                       ║")
	fmt.Println("║   ██╔═██╗ ██╔══██║╚════██║██║██╔══██╗                       ║")
	fmt.Println("║   ██║  ██╗██║  ██║███████║██║██║  ██║                       ║")
	fmt.Println("║   ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝╚═╝╚═╝  ╚═╝                       ║")
	fmt.Println("║                                                             ║")
	fmt.Println("║            SISTEM KASIR CLI - GOLANG + POSTGRESQL           ║")
	fmt.Println("║                        Version 2.0                          ║")
	fmt.Println("║                                                             ║")
	fmt.Println("╚═════════════════════════════════════════════════════════════╝")
	fmt.Println()
}

func printAdminMenu() {
	fmt.Println("\n╔══════════════════════════════════════╗")
	fmt.Println("║        MENU UTAMA (ADMIN)            ║")
	fmt.Println("╠══════════════════════════════════════╣")
	fmt.Println("║  1. 🛒 Transaksi Baru                ║")
	fmt.Println("║  2. 📦 Manajemen Produk              ║")
	fmt.Println("║  3. 📊 Laporan Penjualan             ║")
	fmt.Println("║  4. 👥 Manajemen User                ║")
	fmt.Println("║  5. 🏭 Manajemen Gudang              ║")
	fmt.Println("║  6. 🔑 Ubah Password                 ║")
	fmt.Println("║  0. 🚪 Logout                        ║")
	fmt.Println("╚══════════════════════════════════════╝")
}

func printUserMenu() {
	warehouseName := "Semua Gudang"
	if models.CurrentUser.WarehouseID != nil {
		warehouse, _ := models.GetWarehouseByID(*models.CurrentUser.WarehouseID)
		if warehouse != nil {
			warehouseName = warehouse.Name
		}
	}

	fmt.Printf("\n╔══════════════════════════════════════╗\n")
	fmt.Printf("║          MENU UTAMA (USER)           ║\n")
	fmt.Printf("╠══════════════════════════════════════╣\n")
	fmt.Printf("║  Gudang: %-27s ║\n", warehouseName)
	fmt.Printf("╠══════════════════════════════════════╣\n")
	fmt.Println("║  1. 🛒 Transaksi Baru                ║")
	fmt.Println("║  2. 📦 Lihat Produk                  ║")
	fmt.Println("║  3. 📊 Laporan Penjualan             ║")
	fmt.Println("║  4. 🔑 Ubah Password                 ║")
	fmt.Println("║  0. 🚪 Logout                        ║")
	fmt.Println("╚══════════════════════════════════════╝")
}

func logout() {
	username := models.CurrentUser.Username
	models.Logout()
	fmt.Printf("\n👋 Sampai jumpa, %s!\n", username)
}

func listProductsOnly() {
	products, err := models.GetAllProducts()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	if len(products) == 0 {
		fmt.Println("\n⚠️  Belum ada produk.")
		return
	}

	fmt.Println("\n┌─────┬────────────────────────┬───────────────┬────────┐")
	fmt.Println("│ ID  │ Nama Produk            │ Harga         │ Stok   │")
	fmt.Println("├─────┼────────────────────────┼───────────────┼────────┤")
	for _, p := range products {
		fmt.Printf("│ %-3d │ %-22s │ %13s │ %6d │\n", p.ID, truncate(p.Name, 22), formatRupiahMain(p.SellingPrice), p.Stock)
	}
	fmt.Println("└─────┴────────────────────────┴───────────────┴────────┘")
}

func formatRupiahMain(amount float64) string {
	intAmount := int64(amount)
	str := fmt.Sprintf("%d", intAmount)

	n := len(str)
	if n <= 3 {
		return "Rp " + str
	}

	var result strings.Builder
	remainder := n % 3
	if remainder > 0 {
		result.WriteString(str[:remainder])
		if n > remainder {
			result.WriteString(".")
		}
	}

	for i := remainder; i < n; i += 3 {
		result.WriteString(str[i : i+3])
		if i+3 < n {
			result.WriteString(".")
		}
	}

	return "Rp " + result.String()
}

func truncate(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}
