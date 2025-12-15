package api

import (
	"fmt"
	"net/http"
)

// StartServer menjalankan web server
func StartServer(port string) {
	mux := http.NewServeMux()

	// Register handlers
	mux.HandleFunc("/api/login", handleLogin)
	mux.HandleFunc("/api/products", authMiddleware(handleProducts))
	mux.HandleFunc("/api/transactions", authMiddleware(handleTransactions))
	mux.HandleFunc("/api/users", authMiddleware(handleUsers))
	mux.HandleFunc("/api/warehouses", authMiddleware(handleWarehouses))
	mux.HandleFunc("/api/reports", authMiddleware(handleReports))

	fmt.Printf("🚀 Server berjalan di port %s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		fmt.Printf("❌ Failed to start server: %v\n", err)
	}
}
