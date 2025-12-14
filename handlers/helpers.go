package handlers

import (
	"fmt"
	"strings"
)

// formatRupiah memformat angka menjadi format Rupiah dengan pemisah ribuan
func formatRupiah(amount float64) string {
	// Konversi ke integer untuk menghilangkan desimal
	intAmount := int64(amount)

	// Konversi ke string
	str := fmt.Sprintf("%d", intAmount)

	// Tambahkan pemisah ribuan (titik)
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

// formatNumber memformat angka dengan pemisah ribuan tanpa prefix Rp
func formatNumber(amount float64) string {
	intAmount := int64(amount)
	str := fmt.Sprintf("%d", intAmount)

	n := len(str)
	if n <= 3 {
		return str
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

	return result.String()
}
