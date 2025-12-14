package handlers

import (
	"fmt"
	"kasir/models"
	"time"
)

// ReportMenu menampilkan menu laporan
func ReportMenu() {
	for {
		fmt.Println("\n╔══════════════════════════════════════╗")
		fmt.Println("║         LAPORAN PENJUALAN            ║")
		fmt.Println("╠══════════════════════════════════════╣")
		fmt.Println("║  1. Laporan Hari Ini                 ║")
		fmt.Println("║  2. Laporan Tanggal Tertentu         ║")
		fmt.Println("║  0. Kembali ke Menu Utama            ║")
		fmt.Println("╚══════════════════════════════════════╝")
		fmt.Print("Pilihan: ")

		choice := readInput()
		switch choice {
		case "1":
			showDailyReport(time.Now())
		case "2":
			selectDateReport()
		case "0":
			return
		default:
			fmt.Println("❌ Pilihan tidak valid!")
		}
	}
}

func selectDateReport() {
	fmt.Print("\nMasukkan tanggal (format: DD-MM-YYYY): ")
	dateStr := readInput()

	date, err := time.Parse("02-01-2006", dateStr)
	if err != nil {
		fmt.Println("❌ Format tanggal tidak valid! Gunakan DD-MM-YYYY")
		return
	}

	showDailyReport(date)
}

func showDailyReport(date time.Time) {
	transactions, err := models.GetTransactionsByDate(date)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	total, profit, count, err := models.GetDailyTotal(date)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	// Header laporan
	warehouseInfo := "Semua Gudang"
	if models.CurrentUser != nil && !models.CurrentUser.IsAdmin() && models.CurrentUser.WarehouseID != nil {
		warehouse, _ := models.GetWarehouseByID(*models.CurrentUser.WarehouseID)
		if warehouse != nil {
			warehouseInfo = warehouse.Name
		}
	}

	fmt.Println("\n╔════════════════════════════════════════════════════════════════╗")
	fmt.Printf("║          LAPORAN PENJUALAN: %s                      ║\n", date.Format("02-01-2006"))
	fmt.Println("╠════════════════════════════════════════════════════════════════╣")
	fmt.Printf("║  Gudang          : %-42s ║\n", warehouseInfo)
	fmt.Printf("║  Jumlah Transaksi: %-3d                                         ║\n", count)
	fmt.Printf("║  Total Penjualan : %-20s                     ║\n", formatRupiah(total))
	fmt.Printf("║  Total Profit    : %-20s                     ║\n", formatRupiah(profit))
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")

	if len(transactions) == 0 {
		fmt.Println("\n⚠️  Tidak ada transaksi pada tanggal ini.")
		return
	}

	fmt.Println("\n┌────────────┬──────────┬───────────────┬───────────────┬───────────────┐")
	fmt.Println("│ ID Trans   │ Waktu    │ Total         │ Profit        │ Kembalian     │")
	fmt.Println("├────────────┼──────────┼───────────────┼───────────────┼───────────────┤")

	for _, t := range transactions {
		fmt.Printf("│ TRX-%06d │ %s │ %13s │ %13s │ %13s │\n",
			t.ID,
			t.CreatedAt.Format("15:04"),
			formatRupiah(t.Total),
			formatRupiah(t.Profit),
			formatRupiah(t.Change))
	}
	fmt.Println("└────────────┴──────────┴───────────────┴───────────────┴───────────────┘")

	// Detail per transaksi
	fmt.Println("\n───────────────────────────────────────────────────────────────────")
	fmt.Println("                    DETAIL TRANSAKSI                               ")
	fmt.Println("───────────────────────────────────────────────────────────────────")

	for _, t := range transactions {
		fmt.Printf("\n📋 TRX-%06d (%s)\n", t.ID, t.CreatedAt.Format("15:04:05"))
		for _, item := range t.Items {
			profitItem := (item.SellingPrice - item.PurchasePrice) * float64(item.Quantity)
			fmt.Printf("   • %-20s x%d = %s (profit: %s)\n",
				truncate(item.ProductName, 20),
				item.Quantity,
				formatRupiah(item.Subtotal),
				formatRupiah(profitItem))
		}
	}

	fmt.Println("\n───────────────────────────────────────────────────────────────────")
	fmt.Print("Tekan Enter untuk melanjutkan...")
	readInput()
}
