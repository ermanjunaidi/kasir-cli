package handlers

import (
	"fmt"
	"kasir/models"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// TransactionMenu menampilkan menu transaksi penjualan
func TransactionMenu() {
	cart := []models.CartItem{}

	for {
		fmt.Println("\n╔══════════════════════════════════════╗")
		fmt.Println("║         TRANSAKSI BARU               ║")
		fmt.Println("╠══════════════════════════════════════╣")
		fmt.Println("║  1. Lihat Daftar Produk              ║")
		fmt.Println("║  2. Tambah ke Keranjang              ║")
		fmt.Println("║  3. Lihat Keranjang                  ║")
		fmt.Println("║  4. Hapus dari Keranjang             ║")
		fmt.Println("║  5. Proses Pembayaran                ║")
		fmt.Println("║  0. Batalkan Transaksi               ║")
		fmt.Println("╚══════════════════════════════════════╝")
		fmt.Print("Pilihan: ")

		choice := readInput()
		switch choice {
		case "1":
			listProducts()
		case "2":
			addToCart(&cart)
		case "3":
			viewCart(cart)
		case "4":
			removeFromCart(&cart)
		case "5":
			if processPayment(cart) {
				return // Transaksi selesai, kembali ke menu utama
			}
		case "0":
			if len(cart) > 0 {
				fmt.Print("⚠️  Keranjang tidak kosong. Yakin batalkan? (y/n): ")
				if strings.ToLower(readInput()) != "y" {
					continue
				}
			}
			return
		default:
			fmt.Println("❌ Pilihan tidak valid!")
		}
	}
}

func addToCart(cart *[]models.CartItem) {
	listProducts()

	fmt.Print("\nMasukkan ID produk: ")
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

	// Cek akses warehouse untuk user biasa
	if models.CurrentUser != nil && !models.CurrentUser.IsAdmin() {
		if models.CurrentUser.WarehouseID != nil && product.WarehouseID != *models.CurrentUser.WarehouseID {
			fmt.Println("❌ Produk tidak tersedia di gudang Anda!")
			return
		}
	}

	fmt.Printf("Produk: %s (Stok: %d, Harga: %s)\n", product.Name, product.Stock, formatRupiah(product.SellingPrice))
	fmt.Print("Jumlah: ")
	qtyStr := readInput()
	qty, err := strconv.Atoi(qtyStr)
	if err != nil || qty <= 0 {
		fmt.Println("❌ Jumlah tidak valid!")
		return
	}

	// Cek stok yang tersedia (dikurangi yang sudah di cart)
	availableStock := product.Stock
	for _, item := range *cart {
		if item.Product.ID == product.ID {
			availableStock -= item.Quantity
		}
	}

	if qty > availableStock {
		fmt.Printf("❌ Stok tidak mencukupi! Tersedia: %d\n", availableStock)
		return
	}

	// Cek apakah produk sudah ada di cart
	found := false
	for i, item := range *cart {
		if item.Product.ID == product.ID {
			(*cart)[i].Quantity += qty
			found = true
			break
		}
	}

	if !found {
		*cart = append(*cart, models.CartItem{
			Product:  product,
			Quantity: qty,
		})
	}

	fmt.Printf("✅ %s x%d ditambahkan ke keranjang\n", product.Name, qty)
}

func viewCart(cart []models.CartItem) {
	if len(cart) == 0 {
		fmt.Println("\n⚠️  Keranjang kosong!")
		return
	}

	var total float64
	fmt.Println("\n┌─────┬────────────────────────┬───────────────┬─────┬───────────────┐")
	fmt.Println("│ No  │ Nama Produk            │ Harga         │ Qty │ Subtotal      │")
	fmt.Println("├─────┼────────────────────────┼───────────────┼─────┼───────────────┤")
	for i, item := range cart {
		subtotal := item.Product.SellingPrice * float64(item.Quantity)
		total += subtotal
		fmt.Printf("│ %-3d │ %-22s │ %13s │ %3d │ %13s │\n",
			i+1, truncate(item.Product.Name, 22), formatRupiah(item.Product.SellingPrice), item.Quantity, formatRupiah(subtotal))
	}
	fmt.Println("├─────┴────────────────────────┴───────────────┴─────┼───────────────┤")
	fmt.Printf("│                                       TOTAL        │ %13s │\n", formatRupiah(total))
	fmt.Println("└────────────────────────────────────────────────────┴───────────────┘")
}

func removeFromCart(cart *[]models.CartItem) {
	if len(*cart) == 0 {
		fmt.Println("\n⚠️  Keranjang kosong!")
		return
	}

	viewCart(*cart)

	fmt.Print("\nMasukkan nomor item yang akan dihapus: ")
	noStr := readInput()
	no, err := strconv.Atoi(noStr)
	if err != nil || no < 1 || no > len(*cart) {
		fmt.Println("❌ Nomor tidak valid!")
		return
	}

	removed := (*cart)[no-1]
	*cart = append((*cart)[:no-1], (*cart)[no:]...)
	fmt.Printf("✅ %s dihapus dari keranjang\n", removed.Product.Name)
}

func processPayment(cart []models.CartItem) bool {
	if len(cart) == 0 {
		fmt.Println("\n⚠️  Keranjang kosong! Tambahkan produk terlebih dahulu.")
		return false
	}

	viewCart(cart)

	// Hitung total
	var total float64
	for _, item := range cart {
		total += item.Product.SellingPrice * float64(item.Quantity)
	}

	fmt.Printf("\nTotal Pembayaran: %s\n", formatRupiah(total))
	fmt.Print("Jumlah Bayar: Rp ")
	paymentStr := readInput()

	// Hapus titik pemisah ribuan jika ada (misal: 50.000 -> 50000)
	paymentStr = strings.ReplaceAll(paymentStr, ".", "")

	payment, err := strconv.ParseFloat(paymentStr, 64)
	if err != nil || payment < 0 {
		fmt.Println("❌ Jumlah pembayaran tidak valid!")
		return false
	}

	if payment < total {
		fmt.Printf("❌ Pembayaran kurang! Kurang %s\n", formatRupiah(total-payment))
		return false
	}

	// Proses transaksi
	transaction, err := models.CreateTransaction(cart, payment)
	if err != nil {
		fmt.Printf("❌ Gagal memproses transaksi: %v\n", err)
		return false
	}

	// Cetak struk
	printReceipt(transaction)
	return true
}

func printReceipt(t *models.Transaction) {
	receipt := generateReceiptText(t)

	// Tampilkan struk di layar
	fmt.Print(receipt)

	// Tanyakan apakah ingin menyimpan/print
	fmt.Print("\nSimpan nota ke file? (y/n): ")
	if strings.ToLower(readInput()) == "y" {
		saveReceiptToFile(t, receipt)
	}

	fmt.Print("\nTekan Enter untuk melanjutkan...")
	readInput()
}

func generateReceiptText(t *models.Transaction) string {
	var sb strings.Builder

	sb.WriteString("\n")
	sb.WriteString("═══════════════════════════════════════════\n")
	sb.WriteString("             STRUK PEMBAYARAN              \n")
	sb.WriteString("═══════════════════════════════════════════\n")
	sb.WriteString(fmt.Sprintf("No. Transaksi: TRX-%06d\n", t.ID))
	sb.WriteString(fmt.Sprintf("Tanggal      : %s\n", t.CreatedAt.Format("02-01-2006 15:04:05")))

	if models.CurrentUser != nil {
		sb.WriteString(fmt.Sprintf("Kasir        : %s\n", models.CurrentUser.Username))
	}

	// Tampilkan gudang
	if t.WarehouseID > 0 {
		warehouse, _ := models.GetWarehouseByID(t.WarehouseID)
		if warehouse != nil {
			sb.WriteString(fmt.Sprintf("Gudang       : %s\n", warehouse.Name))
		}
	}

	sb.WriteString("───────────────────────────────────────────\n")

	for _, item := range t.Items {
		sb.WriteString(fmt.Sprintf("%-20s\n", truncate(item.ProductName, 20)))
		sb.WriteString(fmt.Sprintf("  %d x %s = %s\n", item.Quantity, formatRupiah(item.SellingPrice), formatRupiah(item.Subtotal)))
	}

	sb.WriteString("───────────────────────────────────────────\n")
	sb.WriteString(fmt.Sprintf("TOTAL        : %13s\n", formatRupiah(t.Total)))
	sb.WriteString(fmt.Sprintf("BAYAR        : %13s\n", formatRupiah(t.Payment)))
	sb.WriteString(fmt.Sprintf("KEMBALIAN    : %13s\n", formatRupiah(t.Change)))
	sb.WriteString("═══════════════════════════════════════════\n")
	sb.WriteString("    Terima Kasih Atas Kunjungan Anda       \n")
	sb.WriteString("═══════════════════════════════════════════\n")

	return sb.String()
}

func saveReceiptToFile(t *models.Transaction, receipt string) {
	// Get current directory
	cwd, _ := os.Getwd()

	// Create exports/nota folder if not exists
	notaDir := cwd + "/exports/nota"
	if err := os.MkdirAll(notaDir, 0755); err != nil {
		fmt.Printf("❌ Gagal membuat folder: %v\n", err)
		return
	}

	// Generate filename
	filename := fmt.Sprintf("nota_TRX-%06d_%s.txt", t.ID, t.CreatedAt.Format("20060102_150405"))
	filepath := notaDir + "/" + filename

	// Write to file
	err := os.WriteFile(filepath, []byte(receipt), 0644)
	if err != nil {
		fmt.Printf("❌ Gagal menyimpan nota: %v\n", err)
		return
	}

	fmt.Printf("✅ Nota berhasil disimpan ke: %s\n", filepath)

	// Try to print (Linux)
	fmt.Print("Cetak langsung ke printer? (y/n): ")
	if strings.ToLower(readInput()) == "y" {
		printFile(filepath)
	}
}

func printFile(filepath string) {
	// Check if lp command exists (Linux printing)
	_, err := exec.LookPath("lp")
	if err != nil {
		fmt.Println("⚠️  Perintah 'lp' tidak ditemukan. Install CUPS untuk mencetak.")
		fmt.Println("   Atau cetak manual file: " + filepath)
		return
	}

	// Print using lp command
	cmd := exec.Command("lp", filepath)
	err = cmd.Run()
	if err != nil {
		fmt.Printf("❌ Gagal mencetak: %v\n", err)
		fmt.Println("   Cetak manual file: " + filepath)
		return
	}

	fmt.Println("✅ Nota dikirim ke printer!")
}

// getCurrentTime returns formatted current time
func getCurrentTime() string {
	return time.Now().Format("02-01-2006 15:04:05")
}
