package api

import (
	"encoding/base64"
	"encoding/json"
	"kasir/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// authMiddleware memverifikasi Basic Auth
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Basic" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(parts[1])
		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		_, err := models.Authenticate(pair[0], pair[1])
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// Helper to get user from request (re-checking basic auth)
func getUserFromRequest(r *http.Request) (*models.User, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, http.ErrNoCookie
	}
	parts := strings.SplitN(authHeader, " ", 2)
	payload, _ := base64.StdEncoding.DecodeString(parts[1])
	pair := strings.SplitN(string(payload), ":", 2)
	return models.Authenticate(pair[0], pair[1])
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	user, err := models.Authenticate(creds.Username, creds.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"user":    user.Username,
		"role":    user.Role,
	})
}

func handleProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		pageStr := r.URL.Query().Get("page")
		limitStr := r.URL.Query().Get("limit")
		search := r.URL.Query().Get("search")

		page := 1
		limit := 10

		if pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}
		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
				limit = l
			}
		}

		products, total, err := models.GetProducts(page, limit, search, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"data": products,
			"meta": map[string]interface{}{
				"current_page": page,
				"limit":        limit,
				"total_items":  total,
				"total_pages":  (total + limit - 1) / limit,
			},
		}

		json.NewEncoder(w).Encode(response)

	case http.MethodPost:
		var req struct {
			Name          string  `json:"name"`
			PurchasePrice float64 `json:"purchase_price"`
			SellingPrice  float64 `json:"selling_price"`
			Stock         int     `json:"stock"`
			WarehouseID   int     `json:"warehouse_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		p, err := models.CreateProduct(req.Name, req.PurchasePrice, req.SellingPrice, req.Stock, req.WarehouseID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(p)

	case http.MethodPut:
		var req struct {
			ID            int     `json:"id"`
			Name          string  `json:"name"`
			PurchasePrice float64 `json:"purchase_price"`
			SellingPrice  float64 `json:"selling_price"`
			Stock         int     `json:"stock"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		err := models.UpdateProduct(req.ID, req.Name, req.PurchasePrice, req.SellingPrice, req.Stock)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Product updated"})

	case http.MethodDelete:
		// Delete using JSON body to be consistent with other endpoints
		var req struct {
			ID int `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		err := models.DeleteProduct(req.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		users, err := models.GetAllUsers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(users)
	case http.MethodPost:
		var req struct {
			Username    string `json:"username"`
			Password    string `json:"password"`
			Role        string `json:"role"`
			WarehouseID *int   `json:"warehouse_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		u, err := models.Register(req.Username, req.Password, req.Role, req.WarehouseID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(u)
	case http.MethodDelete:
		var req struct {
			ID int `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		err := models.DeleteUser(req.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "User deleted"})
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleWarehouses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		data, err := models.GetAllWarehouses()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(data)
	case http.MethodPost:
		var req struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		wObj, err := models.CreateWarehouse(req.Name, req.Address)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(wObj)
	case http.MethodDelete:
		var req struct {
			ID int `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		err := models.DeleteWarehouse(req.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Warehouse deleted"})
	case http.MethodPut:
		var req struct {
			ID      int    `json:"id"`
			Name    string `json:"name"`
			Address string `json:"address"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		err := models.UpdateWarehouse(req.ID, req.Name, req.Address)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Warehouse updated"})
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleReports(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		dateStr = time.Now().Format("02-01-2006")
	}
	date, err := time.Parse("02-01-2006", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format DD-MM-YYYY", http.StatusBadRequest)
		return
	}

	transactions, err := models.GetTransactionsByDate(user, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	total, profit, count, _ := models.GetDailyTotal(user, date)

	resp := map[string]interface{}{
		"date": dateStr,
		"summary": map[string]interface{}{
			"total_sales":       total,
			"total_profit":      profit,
			"transaction_count": count,
		},
		"transactions": transactions,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleTransactions(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodGet {
		// List transactions
		transactions, err := models.GetTransactionsByDate(user, time.Now())
		// Note: Real API would accept date param, defaulting to today for simplicity
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(transactions)
		return
	}

	if r.Method == http.MethodPost {
		// Create transaction
		var req struct {
			Items []struct {
				ProductID int `json:"product_id"`
				Quantity  int `json:"quantity"`
			} `json:"items"`
			Payment float64 `json:"payment"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Build items
		var cart []models.CartItem
		for _, itemReq := range req.Items {
			product, err := models.GetProductByID(itemReq.ProductID)
			if err != nil {
				http.Error(w, "Product not found: "+string(rune(itemReq.ProductID)), http.StatusBadRequest)
				return
			}
			cart = append(cart, models.CartItem{
				Product:  product,
				Quantity: itemReq.Quantity,
			})
		}

		trx, err := models.CreateTransaction(user, cart, req.Payment)
		if err != nil {
			http.Error(w, "Transaction failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(trx)
		return
	}
}
