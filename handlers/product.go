package handlers

import (
	"fmt"
	"kasir/models"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// ProductMenu menampilkan menu manajemen produk
func ProductMenu() {
	for {
		// Menu berbeda untuk admin dan user
		if models.CurrentUser != nil && models.CurrentUser.IsAdmin() {
			fmt.Println("\n╔══════════════════════════════════════╗")
			fmt.Println("║        MANAJEMEN PRODUK (ADMIN)      ║")
			fmt.Println("╠══════════════════════════════════════╣")
			fmt.Println("║  1. Lihat Produk (Semua Gudang)      ║")
			fmt.Println("║  2. Lihat Stok per Gudang            ║")
			fmt.Println("║  3. Tambah Produk                    ║")
			fmt.Println("║  4. Edit Produk                      ║")
			fmt.Println("║  5. Hapus Produk                     ║")
			fmt.Println("║  6. Export ke Excel                  ║")
			fmt.Println("║  7. Import dari Excel                ║")
			fmt.Println("║  0. Kembali ke Menu Utama            ║")
			fmt.Println("╚══════════════════════════════════════╝")
		} else {
			fmt.Println("\n╔══════════════════════════════════════╗")
			fmt.Println("║        MANAJEMEN PRODUK              ║")
			fmt.Println("╠══════════════════════════════════════╣")
			fmt.Println("║  1. Lihat Daftar Produk              ║")
			fmt.Println("║  2. Tambah Produk                    ║")
			fmt.Println("║  3. Edit Produk                      ║")
			fmt.Println("║  4. Hapus Produk                     ║")
			fmt.Println("║  0. Kembali ke Menu Utama            ║")
			fmt.Println("╚══════════════════════════════════════╝")
		}
		fmt.Print("Pilihan: ")

		choice := readInput()

		if models.CurrentUser != nil && models.CurrentUser.IsAdmin() {
			switch choice {
			case "1":
				ListAllProducts()
			case "2":
				listStockByWarehouse()
			case "3":
				addProduct()
			case "4":
				editProduct()
			case "5":
				deleteProduct()
			case "6":
				exportToExcel()
			case "7":
				importFromExcel()
			case "0":
				return
			default:
				fmt.Println("❌ Pilihan tidak valid!")
			}
		} else {
			switch choice {
			case "1":
				ListProducts()
			case "2":
				addProduct()
			case "3":
				editProduct()
			case "4":
				deleteProduct()
			case "0":
				return
			default:
				fmt.Println("❌ Pilihan tidak valid!")
			}
		}
	}
}

func readInput() string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// ListAllProducts menampilkan semua produk dengan pagination (admin)
func ListAllProducts() {
	page := 1
	limit := 10
	search := ""

	for {
		products, total, err := models.GetProducts(page, limit, search, nil)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		totalPages := (total + limit - 1) / limit
		if totalPages == 0 {
			totalPages = 1
		}

		fmt.Println("\n╔═════════════════════════════════════════════════════════════════════════════════════════════╗")
		fmt.Println("║                                  DAFTAR SEMUA PRODUK                                    ║")
		fmt.Println("╚═════════════════════════════════════════════════════════════════════════════════════════════╝")
		fmt.Printf("🔍 Search: %-20s  📄 Page: %d/%d  📦 Total: %d\n",
			func() string {
				if search == "" {
					return "(none)"
				} else {
					return search
				}
			}(),
			page, totalPages, total)

		fmt.Println("┌─────┬──────────────────────┬─────────────┬─────────────┬──────┬─────────────────────────┐")
		fmt.Println("│ ID  │ Nama Produk          │ Hrg Beli    │ Hrg Jual    │ Stok │ Gudang                  │")
		fmt.Println("├─────┼──────────────────────┼─────────────┼─────────────┼──────┼─────────────────────────┤")

		if len(products) == 0 {
			fmt.Println("│                        T I D A K   A D A   D A T A                              │")
		}

		for _, p := range products {
			warehouseName := "?"
			w, _ := models.GetWarehouseByID(p.WarehouseID)
			if w != nil {
				warehouseName = w.Name
			}

			fmt.Printf("│ %-3d │ %-20s │ %11s │ %11s │ %4d │ %-23s │\n",
				p.ID, truncate(p.Name, 20), formatRupiah(p.PurchasePrice), formatRupiah(p.SellingPrice), p.Stock, truncate(warehouseName, 23))
		}
		fmt.Println("└─────┴──────────────────────┴─────────────┴─────────────┴──────┴─────────────────────────┘")

		fmt.Println("\n[n] Next  [p] Prev  [s] Search  [q] Back")
		fmt.Print("Pilihan: ")
		input := readInput()

		switch strings.ToLower(input) {
		case "n":
			if page < totalPages {
				page++
			} else {
				fmt.Println("⚠️  Sudah di halaman terakhir")
			}
		case "p":
			if page > 1 {
				page--
			} else {
				fmt.Println("⚠️  Sudah di halaman pertama")
			}
		case "s":
			fmt.Print("Masukkan kata kunci: ")
			search = readInput()
			page = 1 // Reset ke halaman 1 saat search baru
		case "q":
			return
		default:
			// check if number
			if pNum, err := strconv.Atoi(input); err == nil && len(input) > 0 {
				if pNum >= 1 && pNum <= totalPages {
					page = pNum
				}
			}
		}
	}
}

// listStockByWarehouse menampilkan ringkasan stok per gudang
func listStockByWarehouse() {
	warehouses, _ := models.GetAllWarehouses()

	fmt.Println("\n╔══════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    RINGKASAN STOK PER GUDANG                     ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════╝")

	fmt.Println("\n┌────────────────────────────┬───────────────┬───────────────┬───────────────────┐")
	fmt.Println("│ Gudang                     │ Jml Produk    │ Total Stok    │ Nilai Stok        │")
	fmt.Println("├────────────────────────────┼───────────────┼───────────────┼───────────────────┤")

	var grandTotalProducts, grandTotalStock int
	var grandTotalValue float64

	for _, w := range warehouses {
		products, _ := models.GetProductsByWarehouse(w.ID)
		var totalStock int
		var totalValue float64

		for _, p := range products {
			totalStock += p.Stock
			totalValue += p.SellingPrice * float64(p.Stock)
		}

		grandTotalProducts += len(products)
		grandTotalStock += totalStock
		grandTotalValue += totalValue

		fmt.Printf("│ %-26s │ %13d │ %13d │ %17s │\n",
			truncate(w.Name, 26), len(products), totalStock, formatRupiah(totalValue))
	}

	fmt.Println("├────────────────────────────┼───────────────┼───────────────┼───────────────────┤")
	fmt.Printf("│ %-26s │ %13d │ %13d │ %17s │\n",
		"TOTAL", grandTotalProducts, grandTotalStock, formatRupiah(grandTotalValue))
	fmt.Println("└────────────────────────────┴───────────────┴───────────────┴───────────────────┘")
}

func ListProducts() {
	page := 1
	limit := 10
	search := ""

	// Ensure we have user warehouse ID if not admin
	var warehouseID *int
	if models.CurrentUser != nil && !models.CurrentUser.IsAdmin() {
		warehouseID = models.CurrentUser.WarehouseID
	}

	for {
		products, total, err := models.GetProducts(page, limit, search, warehouseID)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		totalPages := (total + limit - 1) / limit
		if totalPages == 0 {
			totalPages = 1
		}

		fmt.Println("\n╔══════════════════════════════════════════════════════════╗")
		fmt.Println("║                   DAFTAR PRODUK                          ║")
		fmt.Println("╚══════════════════════════════════════════════════════════╝")
		fmt.Printf("🔍 Search: %-15s  📄 Page: %d/%d  📦 Total: %d\n",
			func() string {
				if search == "" {
					return "(none)"
				} else {
					return search
				}
			}(),
			page, totalPages, total)

		fmt.Println("┌─────┬────────────────────────┬───────────────┬────────┐")
		fmt.Println("│ ID  │ Nama Produk            │ Harga         │ Stok   │")
		fmt.Println("├─────┼────────────────────────┼───────────────┼────────┤")

		if len(products) == 0 {
			fmt.Println("│               T I D A K   A D A   D A T A             │")
		}

		for _, p := range products {
			fmt.Printf("│ %-3d │ %-22s │ %13s │ %6d │\n", p.ID, truncate(p.Name, 22), formatRupiah(p.SellingPrice), p.Stock)
		}
		fmt.Println("└─────┴────────────────────────┴───────────────┴────────┘")

		fmt.Println("\n[n] Next  [p] Prev  [s] Search  [q] Back")
		fmt.Print("Pilihan: ")
		input := readInput()

		switch strings.ToLower(input) {
		case "n":
			if page < totalPages {
				page++
			} else {
				fmt.Println("⚠️  Sudah di halaman terakhir")
			}
		case "p":
			if page > 1 {
				page--
			} else {
				fmt.Println("⚠️  Sudah di halaman pertama")
			}
		case "s":
			fmt.Print("Masukkan kata kunci: ")
			search = readInput()
			page = 1 // Reset ke halaman 1 saat search baru
		case "q":
			return
		default:
			// check if number
			if pNum, err := strconv.Atoi(input); err == nil && len(input) > 0 {
				if pNum >= 1 && pNum <= totalPages {
					page = pNum
				}
			}
		}
	}
}

func truncate(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}

func addProduct() {
	fmt.Println("\n═══ TAMBAH PRODUK BARU ═══")

	fmt.Print("Nama Produk: ")
	name := readInput()
	if name == "" {
		fmt.Println("❌ Nama produk tidak boleh kosong!")
		return
	}

	fmt.Print("Harga Beli: ")
	purchasePriceStr := readInput()
	purchasePrice, err := strconv.ParseFloat(purchasePriceStr, 64)
	if err != nil || purchasePrice < 0 {
		fmt.Println("❌ Harga beli tidak valid!")
		return
	}

	fmt.Print("Harga Jual: ")
	sellingPriceStr := readInput()
	sellingPrice, err := strconv.ParseFloat(sellingPriceStr, 64)
	if err != nil || sellingPrice < 0 {
		fmt.Println("❌ Harga jual tidak valid!")
		return
	}

	if sellingPrice < purchasePrice {
		fmt.Println("⚠️  Peringatan: Harga jual lebih rendah dari harga beli!")
	}

	fmt.Print("Stok: ")
	stockStr := readInput()
	stock, err := strconv.Atoi(stockStr)
	if err != nil || stock < 0 {
		fmt.Println("❌ Stok tidak valid!")
		return
	}

	// Tentukan warehouse
	var warehouseID int
	if models.CurrentUser != nil && !models.CurrentUser.IsAdmin() && models.CurrentUser.WarehouseID != nil {
		warehouseID = *models.CurrentUser.WarehouseID
	} else {
		// Admin pilih warehouse
		warehouses, err := models.GetAllWarehouses()
		if err != nil || len(warehouses) == 0 {
			fmt.Println("❌ Tidak ada gudang tersedia!")
			return
		}

		fmt.Println("\nPilih Gudang:")
		for _, w := range warehouses {
			fmt.Printf("  %d. %s\n", w.ID, w.Name)
		}
		fmt.Print("ID Gudang: ")
		fmt.Scanln(&warehouseID)
		reader.ReadString('\n')
	}

	product, err := models.CreateProduct(name, purchasePrice, sellingPrice, stock, warehouseID)
	if err != nil {
		fmt.Printf("❌ Gagal menambah produk: %v\n", err)
		return
	}

	fmt.Printf("✅ Produk '%s' berhasil ditambahkan dengan ID: %d\n", product.Name, product.ID)
}

func editProduct() {
	ListProducts()

	fmt.Print("\nMasukkan ID produk yang akan diedit: ")
	idStr := readInput()
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("❌ ID tidak valid!")
		return
	}

	product, err := models.GetProductByID(id)
	if err != nil {
		fmt.Println("❌ Produk tidak ditemukan!")
		return
	}

	// Cek akses untuk user biasa
	if models.CurrentUser != nil && !models.CurrentUser.IsAdmin() {
		if models.CurrentUser.WarehouseID != nil && product.WarehouseID != *models.CurrentUser.WarehouseID {
			fmt.Println("❌ Anda tidak memiliki akses ke produk ini!")
			return
		}
	}

	fmt.Printf("\n═══ EDIT PRODUK: %s ═══\n", product.Name)
	fmt.Println("(Tekan Enter untuk tidak mengubah)")

	fmt.Printf("Nama [%s]: ", product.Name)
	name := readInput()
	if name == "" {
		name = product.Name
	}

	fmt.Printf("Harga Beli [%s]: ", formatRupiah(product.PurchasePrice))
	purchasePriceStr := readInput()
	purchasePrice := product.PurchasePrice
	if purchasePriceStr != "" {
		purchasePrice, err = strconv.ParseFloat(purchasePriceStr, 64)
		if err != nil || purchasePrice < 0 {
			fmt.Println("❌ Harga beli tidak valid!")
			return
		}
	}

	fmt.Printf("Harga Jual [%s]: ", formatRupiah(product.SellingPrice))
	sellingPriceStr := readInput()
	sellingPrice := product.SellingPrice
	if sellingPriceStr != "" {
		sellingPrice, err = strconv.ParseFloat(sellingPriceStr, 64)
		if err != nil || sellingPrice < 0 {
			fmt.Println("❌ Harga jual tidak valid!")
			return
		}
	}

	fmt.Printf("Stok [%d]: ", product.Stock)
	stockStr := readInput()
	stock := product.Stock
	if stockStr != "" {
		stock, err = strconv.Atoi(stockStr)
		if err != nil || stock < 0 {
			fmt.Println("❌ Stok tidak valid!")
			return
		}
	}

	err = models.UpdateProduct(id, name, purchasePrice, sellingPrice, stock)
	if err != nil {
		fmt.Printf("❌ Gagal mengupdate produk: %v\n", err)
		return
	}

	fmt.Println("✅ Produk berhasil diupdate!")
}

func deleteProduct() {
	ListProducts()

	fmt.Print("\nMasukkan ID produk yang akan dihapus: ")
	idStr := readInput()
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("❌ ID tidak valid!")
		return
	}

	product, err := models.GetProductByID(id)
	if err != nil {
		fmt.Println("❌ Produk tidak ditemukan!")
		return
	}

	// Cek akses untuk user biasa
	if models.CurrentUser != nil && !models.CurrentUser.IsAdmin() {
		if models.CurrentUser.WarehouseID != nil && product.WarehouseID != *models.CurrentUser.WarehouseID {
			fmt.Println("❌ Anda tidak memiliki akses ke produk ini!")
			return
		}
	}

	fmt.Printf("⚠️  Yakin ingin menghapus '%s'? (y/n): ", product.Name)
	confirm := readInput()
	if strings.ToLower(confirm) != "y" {
		fmt.Println("Batal menghapus.")
		return
	}

	err = models.DeleteProduct(id)
	if err != nil {
		fmt.Printf("❌ Gagal menghapus produk: %v\n", err)
		return
	}

	fmt.Println("✅ Produk berhasil dihapus!")
}

// exportToExcel mengexport data produk ke file Excel
func exportToExcel() {
	fmt.Println("\n═══ EXPORT DATA PRODUK KE EXCEL ═══")

	// Pilih gudang atau semua
	warehouses, _ := models.GetAllWarehouses()
	fmt.Println("\nPilih Gudang:")
	fmt.Println("  0. Semua Gudang")
	for _, w := range warehouses {
		fmt.Printf("  %d. %s\n", w.ID, w.Name)
	}
	fmt.Print("Pilihan: ")
	var warehouseID int
	fmt.Scanln(&warehouseID)
	reader.ReadString('\n')

	// Buat file Excel
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Produk"
	f.SetSheetName("Sheet1", sheetName)

	// Header style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"4472C4"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	// Set headers
	headers := []string{"ID", "Nama Produk", "Harga Beli", "Harga Jual", "Stok", "Gudang ID", "Nama Gudang"}
	for i, h := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, h)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Set column widths
	f.SetColWidth(sheetName, "A", "A", 8)
	f.SetColWidth(sheetName, "B", "B", 25)
	f.SetColWidth(sheetName, "C", "D", 15)
	f.SetColWidth(sheetName, "E", "E", 10)
	f.SetColWidth(sheetName, "F", "F", 12)
	f.SetColWidth(sheetName, "G", "G", 20)

	// Get products
	var products []models.Product
	if warehouseID == 0 {
		// Semua gudang
		for _, w := range warehouses {
			prods, _ := models.GetProductsByWarehouse(w.ID)
			products = append(products, prods...)
		}
	} else {
		products, _ = models.GetProductsByWarehouse(warehouseID)
	}

	// Data rows
	for i, p := range products {
		row := i + 2
		warehouseName := ""
		warehouse, _ := models.GetWarehouseByID(p.WarehouseID)
		if warehouse != nil {
			warehouseName = warehouse.Name
		}

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), p.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), p.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), p.PurchasePrice)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), p.SellingPrice)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), p.Stock)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), p.WarehouseID)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), warehouseName)
	}

	// Generate filename
	filename := fmt.Sprintf("produk_export_%s.xlsx", time.Now().Format("20060102_150405"))

	// Get current directory and create exports/excel folder
	cwd, _ := os.Getwd()
	excelDir := filepath.Join(cwd, "exports", "excel")
	if err := os.MkdirAll(excelDir, 0755); err != nil {
		fmt.Printf("❌ Gagal membuat folder: %v\n", err)
		return
	}

	filePath := filepath.Join(excelDir, filename)

	// Save file
	if err := f.SaveAs(filePath); err != nil {
		fmt.Printf("❌ Gagal menyimpan file: %v\n", err)
		return
	}

	fmt.Printf("✅ Berhasil export %d produk ke file:\n", len(products))
	fmt.Printf("   📄 %s\n", filePath)
}

// importFromExcel mengimport data produk dari file Excel
func importFromExcel() {
	fmt.Println("\n═══ IMPORT DATA PRODUK DARI EXCEL ═══")
	fmt.Println("Format Excel yang diperlukan:")
	fmt.Println("  Kolom A: Nama Produk")
	fmt.Println("  Kolom B: Harga Beli")
	fmt.Println("  Kolom C: Harga Jual")
	fmt.Println("  Kolom D: Stok")
	fmt.Println("  Kolom E: Gudang ID")
	fmt.Println("  (Baris pertama = header, data mulai baris 2)")

	fmt.Print("\nMasukkan path file Excel: ")
	filePath := readInput()

	if filePath == "" {
		fmt.Println("❌ Path file tidak boleh kosong!")
		return
	}

	// Open Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		fmt.Printf("❌ Gagal membuka file: %v\n", err)
		return
	}
	defer f.Close()

	// Get first sheet
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		fmt.Printf("❌ Gagal membaca file: %v\n", err)
		return
	}

	if len(rows) < 2 {
		fmt.Println("❌ File tidak memiliki data (minimal 2 baris: header + data)")
		return
	}

	// Process data (skip header)
	var successCount, failCount int
	for i, row := range rows[1:] {
		if len(row) < 5 {
			fmt.Printf("⚠️  Baris %d: Data tidak lengkap, dilewati\n", i+2)
			failCount++
			continue
		}

		name := row[0]
		purchasePrice, err1 := strconv.ParseFloat(row[1], 64)
		sellingPrice, err2 := strconv.ParseFloat(row[2], 64)
		stock, err3 := strconv.Atoi(row[3])
		warehouseID, err4 := strconv.Atoi(row[4])

		if name == "" || err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			fmt.Printf("⚠️  Baris %d: Data tidak valid, dilewati\n", i+2)
			failCount++
			continue
		}

		_, err := models.CreateProduct(name, purchasePrice, sellingPrice, stock, warehouseID)
		if err != nil {
			fmt.Printf("⚠️  Baris %d: Gagal import '%s': %v\n", i+2, name, err)
			failCount++
			continue
		}

		successCount++
	}

	fmt.Printf("\n✅ Import selesai!\n")
	fmt.Printf("   ✓ Berhasil: %d produk\n", successCount)
	if failCount > 0 {
		fmt.Printf("   ✗ Gagal: %d produk\n", failCount)
	}
}
