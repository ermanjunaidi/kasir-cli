import { useQuery } from "@tanstack/react-query";

const API_URL = "/api"; // Proxy will handle this

// Types
export interface Product {
    ID: number;
    Name: string;
    Stock: number;
    PurchasePrice: number;
    SellingPrice: number;
}

export interface Transaction {
    ID: number;
    Total: number;
    Profit: number;
    Payment: number;
    Change: number;
    CreatedAt: string;
}

export interface User {
    ID: number;
    Username: string;
    Role: string;
    WarehouseID?: number;
}

export interface Warehouse {
    ID: number;
    Name: string;
    Address: string;
}

export interface Report {
    date: string;
    summary: {
        total_sales: number;
        total_profit: number;
        transaction_count: number;
    };
    transactions: Transaction[];
}

// Helpers
const getAuthHeaders = (): Record<string, string> => {
    const token = localStorage.getItem("auth_token"); // Basic base64
    return token ? { Authorization: `Basic ${token}` } : {};
};

// API Functions
export const api = {
    login: async (username: string, password: string) => {
        // We test credentials by encoding them and calling a guarded endpoint or login endpoint
        const token = btoa(`${username}:${password}`);
        const res = await fetch(`${API_URL}/login`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ username, password }),
        });
        if (!res.ok) throw new Error("Login failed");
        localStorage.setItem("auth_token", token);
        return res.json();
    },

    getProducts: async (): Promise<Product[]> => {
        const res = await fetch(`${API_URL}/products`, {
            headers: getAuthHeaders(),
        });
        if (!res.ok) throw new Error("Failed to fetch products");
        return res.json();
    },

    getTransactions: async (): Promise<Transaction[]> => {
        const res = await fetch(`${API_URL}/transactions`, {
            headers: getAuthHeaders(),
        });
        if (!res.ok) throw new Error("Failed to fetch transactions");
        return res.json();
    },

    createTransaction: async (items: { product_id: number; quantity: number }[], payment: number) => {
        const headers = { "Content-Type": "application/json", ...getAuthHeaders() };
        const res = await fetch(`${API_URL}/transactions`, {
            method: "POST",
            headers: headers,
            body: JSON.stringify({ items, payment }),
        });
        if (!res.ok) {
            const err = await res.text();
            throw new Error(err);
        }
        return res.json();
    },

    // Products
    createProduct: async (data: Omit<Product, "ID">) => {
        const res = await fetch(`${API_URL}/products`, { method: "POST", headers: { "Content-Type": "application/json", ...getAuthHeaders() }, body: JSON.stringify(data) });
        if (!res.ok) throw new Error("Failed to create product");
        return res.json();
    },
    deleteProduct: async (id: number) => {
        const res = await fetch(`${API_URL}/products`, { method: "DELETE", headers: { "Content-Type": "application/json", ...getAuthHeaders() }, body: JSON.stringify({ id }) });
        if (!res.ok) throw new Error("Failed to delete product");
        return res.json();
    },

    // Users
    getUsers: async (): Promise<User[]> => {
        const res = await fetch(`${API_URL}/users`, { headers: getAuthHeaders() });
        if (!res.ok) throw new Error("Failed to fetch users");
        return res.json();
    },
    createUser: async (data: Omit<User, "ID"> & { password?: string }) => { // Assuming password is sent for creation
        const res = await fetch(`${API_URL}/users`, { method: "POST", headers: { "Content-Type": "application/json", ...getAuthHeaders() }, body: JSON.stringify(data) });
        if (!res.ok) throw new Error("Failed to create user");
        return res.json();
    },
    deleteUser: async (id: number) => {
        const res = await fetch(`${API_URL}/users`, { method: "DELETE", headers: { "Content-Type": "application/json", ...getAuthHeaders() }, body: JSON.stringify({ id }) });
        if (!res.ok) throw new Error("Failed to delete user");
        return res.json();
    },

    // Warehouses
    getWarehouses: async (): Promise<Warehouse[]> => {
        const res = await fetch(`${API_URL}/warehouses`, { headers: getAuthHeaders() });
        if (!res.ok) throw new Error("Failed to fetch warehouses");
        return res.json();
    },
    createWarehouse: async (data: Omit<Warehouse, "ID">) => {
        const res = await fetch(`${API_URL}/warehouses`, { method: "POST", headers: { "Content-Type": "application/json", ...getAuthHeaders() }, body: JSON.stringify(data) });
        if (!res.ok) throw new Error("Failed to create warehouse");
        return res.json();
    },
    deleteWarehouse: async (id: number) => {
        const res = await fetch(`${API_URL}/warehouses`, { method: "DELETE", headers: { "Content-Type": "application/json", ...getAuthHeaders() }, body: JSON.stringify({ id }) });
        if (!res.ok) throw new Error("Failed to delete warehouse");
        return res.json();
    },

    // Reports
    getReport: async (dateStr: string): Promise<Report> => {
        const res = await fetch(`${API_URL}/reports?date=${dateStr}`, { headers: getAuthHeaders() });
        if (!res.ok) throw new Error("Failed to fetch report");
        return res.json();
    },
};

// Hooks
export const useProducts = () => useQuery({ queryKey: ["products"], queryFn: api.getProducts });
export const useTransactions = () => useQuery({ queryKey: ["transactions"], queryFn: api.getTransactions });
export const useUsers = () => useQuery({ queryKey: ["users"], queryFn: api.getUsers });
export const useWarehouses = () => useQuery({ queryKey: ["warehouses"], queryFn: api.getWarehouses });
