package handlers

import (
	"fmt"
	"kasir/models"
	"strings"
)

// WarehouseMenu menampilkan menu manajemen gudang (admin only)
func WarehouseMenu() {
	for {
		fmt.Println("\n╔══════════════════════════════════════╗")
		fmt.Println("║        MANAJEMEN GUDANG              ║")
		fmt.Println("╠══════════════════════════════════════╣")
		fmt.Println("║  1. Lihat Daftar Gudang              ║")
		fmt.Println("║  2. Tambah Gudang Baru               ║")
		fmt.Println("║  3. Edit Gudang                      ║")
		fmt.Println("║  4. Hapus Gudang                     ║")
		fmt.Println("║  0. Kembali ke Menu Utama            ║")
		fmt.Println("╚══════════════════════════════════════╝")

		fmt.Print("Pilihan: ")
		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		switch choice {
		case "1":
			listWarehouses()
		case "2":
			createWarehouse()
		case "3":
			editWarehouse()
		case "4":
			deleteWarehouse()
		case "0":
			return
		default:
			fmt.Println("❌ Pilihan tidak valid!")
		}
	}
}

func listWarehouses() {
	warehouses, err := models.GetAllWarehouses()
	if err != nil {
		fmt.Println("❌ Error:", err)
		return
	}

	fmt.Println("\n┌─────┬────────────────────────┬────────────────────────────────────┐")
	fmt.Println("│ ID  │ Nama Gudang            │ Alamat                             │")
	fmt.Println("├─────┼────────────────────────┼────────────────────────────────────┤")

	for _, w := range warehouses {
		address := w.Address
		if len(address) > 34 {
			address = address[:31] + "..."
		}
		fmt.Printf("│ %-3d │ %-22s │ %-34s │\n", w.ID, w.Name, address)
	}
	fmt.Println("└─────┴────────────────────────┴────────────────────────────────────┘")
}

func createWarehouse() {
	fmt.Println("\n═══ TAMBAH GUDANG BARU ═══")

	fmt.Print("Nama Gudang: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Alamat: ")
	address, _ := reader.ReadString('\n')
	address = strings.TrimSpace(address)

	warehouse, err := models.CreateWarehouse(name, address)
	if err != nil {
		fmt.Println("❌ Gagal menambah gudang:", err)
		return
	}

	fmt.Printf("✅ Gudang '%s' berhasil ditambahkan dengan ID: %d\n", warehouse.Name, warehouse.ID)
}

func editWarehouse() {
	listWarehouses()

	fmt.Print("\nMasukkan ID gudang yang akan diedit (0 untuk batal): ")
	var id int
	fmt.Scanln(&id)
	reader.ReadString('\n')

	if id == 0 {
		return
	}

	warehouse, err := models.GetWarehouseByID(id)
	if err != nil {
		fmt.Println("❌ Gudang tidak ditemukan!")
		return
	}

	fmt.Printf("\n═══ EDIT GUDANG: %s ═══\n", warehouse.Name)

	fmt.Printf("Nama Baru (kosongkan jika tidak diubah) [%s]: ", warehouse.Name)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if name == "" {
		name = warehouse.Name
	}

	fmt.Printf("Alamat Baru (kosongkan jika tidak diubah) [%s]: ", warehouse.Address)
	address, _ := reader.ReadString('\n')
	address = strings.TrimSpace(address)
	if address == "" {
		address = warehouse.Address
	}

	err = models.UpdateWarehouse(id, name, address)
	if err != nil {
		fmt.Println("❌ Gagal mengupdate gudang:", err)
		return
	}

	fmt.Println("✅ Gudang berhasil diupdate!")
}

func deleteWarehouse() {
	listWarehouses()

	fmt.Print("\nMasukkan ID gudang yang akan dihapus (0 untuk batal): ")
	var id int
	fmt.Scanln(&id)
	reader.ReadString('\n')

	if id == 0 {
		return
	}

	// Cek apakah ada user yang terkait
	users, _ := models.GetAllUsers()
	var linkedUsers []string
	for _, u := range users {
		if u.WarehouseID != nil && *u.WarehouseID == id {
			linkedUsers = append(linkedUsers, u.Username)
		}
	}

	// Cek apakah ada produk yang terkait
	products, _ := models.GetProductsByWarehouse(id)

	if len(linkedUsers) > 0 || len(products) > 0 {
		fmt.Println("\n❌ Tidak dapat menghapus gudang karena masih ada data terkait:")
		if len(linkedUsers) > 0 {
			fmt.Printf("   • %d user: %s\n", len(linkedUsers), strings.Join(linkedUsers, ", "))
		}
		if len(products) > 0 {
			fmt.Printf("   • %d produk\n", len(products))
		}
		fmt.Println("\n💡 Hapus atau pindahkan data tersebut terlebih dahulu.")
		return
	}

	fmt.Print("⚠️  Yakin ingin menghapus gudang ini? (y/n): ")
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm != "y" {
		fmt.Println("Batal menghapus.")
		return
	}

	err := models.DeleteWarehouse(id)
	if err != nil {
		fmt.Println("❌ Gagal menghapus gudang:", err)
		return
	}

	fmt.Println("✅ Gudang berhasil dihapus!")
}
