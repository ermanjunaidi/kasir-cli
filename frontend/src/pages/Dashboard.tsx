import { useState } from "react";
import { useProducts, useUsers, useWarehouses, api, type Product } from "../lib/api";
import { Button } from "../components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "../components/ui/card";
import { Input } from "../components/ui/input";
import { LogOut, ShoppingCart, Trash2, Plus, Users, Package, Home, FileText } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { useMutation, useQueryClient, useQuery } from "@tanstack/react-query";

interface CartItem {
  product: Product;
  quantity: number;
}

export default function Dashboard() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [activeTab, setActiveTab] = useState("pos"); // pos, products, users, warehouses, reports

  const handleLogout = () => {
    localStorage.removeItem("auth_token");
    navigate("/");
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 flex">
        {/* Sidebar */}
        <div className="w-64 bg-white dark:bg-gray-800 border-r p-4 flex flex-col">
            <h1 className="text-2xl font-bold mb-8 text-center text-primary">Kasir Admin</h1>
            <nav className="space-y-2 flex-1">
                <Button variant={activeTab === "pos" ? "default" : "ghost"} className="w-full justify-start" onClick={() => setActiveTab("pos")}>
                    <ShoppingCart className="mr-2 h-4 w-4" /> POS
                </Button>
                <Button variant={activeTab === "products" ? "default" : "ghost"} className="w-full justify-start" onClick={() => setActiveTab("products")}>
                    <Package className="mr-2 h-4 w-4" /> Produk
                </Button>
                <Button variant={activeTab === "users" ? "default" : "ghost"} className="w-full justify-start" onClick={() => setActiveTab("users")}>
                    <Users className="mr-2 h-4 w-4" /> Users
                </Button>
                <Button variant={activeTab === "warehouses" ? "default" : "ghost"} className="w-full justify-start" onClick={() => setActiveTab("warehouses")}>
                    <Home className="mr-2 h-4 w-4" /> Gudang
                </Button>
                <Button variant={activeTab === "reports" ? "default" : "ghost"} className="w-full justify-start" onClick={() => setActiveTab("reports")}>
                    <FileText className="mr-2 h-4 w-4" /> Laporan
                </Button>
            </nav>
            <Button variant="outline" className="mt-auto w-full" onClick={handleLogout}>
                <LogOut className="mr-2 h-4 w-4" /> Logout
            </Button>
        </div>

        {/* Content */}
        <main className="flex-1 p-8 overflow-auto">
            {activeTab === "pos" && <POSContent />}
            {activeTab === "products" && <ProductsContent />}
            {activeTab === "users" && <UsersContent />}
            {activeTab === "warehouses" && <WarehousesContent />}
            {activeTab === "reports" && <ReportsContent />}
        </main>
    </div>
  );
}

function POSContent() {
    const [page, setPage] = useState(1);
    const [search, setSearch] = useState("");
    const { data: productsData } = useProducts({ page, limit: 12, search });
    const products = productsData?.data || [];
    const meta = productsData?.meta;

    const [cart, setCart] = useState<CartItem[]>([]);
    const [payment, setPayment] = useState("");
    const queryClient = useQueryClient();

    const addToCart = (product: Product) => {
        setCart((prev) => {
          const existing = prev.find((item) => item.product.ID === product.ID);
          if (existing) {
            return prev.map((item) =>
              item.product.ID === product.ID
                ? { ...item, quantity: item.quantity + 1 }
                : item
            );
          }
          return [...prev, { product, quantity: 1 }];
        });
    };

    const removeFromCart = (id: number) => {
        setCart((prev) => prev.filter((item) => item.product.ID !== id));
    };

    const total = cart.reduce(
        (acc, item) => acc + item.product.SellingPrice * item.quantity,
        0
    );

    const checkoutMutation = useMutation({
        mutationFn: async () => {
            const items = cart.map(c => ({ product_id: c.product.ID, quantity: c.quantity }));
            return api.createTransaction(items, parseFloat(payment));
        },
        onSuccess: () => {
            alert("Transaksi Berhasil!");
            setCart([]);
            setPayment("");
            queryClient.invalidateQueries({ queryKey: ["products"] });
        },
        onError: (err) => {
            alert("Gagal: " + err.message);
        }
    });

    return (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="md:col-span-2 space-y-4">
                <h2 className="text-xl font-semibold mb-4">Katalog Produk</h2>
                <div className="flex gap-2">
                    <Input placeholder="Cari produk..." value={search} onChange={e => { setSearch(e.target.value); setPage(1); }} />
                </div>
                <div className="grid grid-cols-2 lg:grid-cols-3 gap-4">
                {products?.map((product) => (
                    <Card key={product.ID} className="cursor-pointer hover:shadow-lg transition" onClick={() => addToCart(product)}>
                    <CardHeader className="p-4 pb-2">
                        <CardTitle className="text-lg">{product.Name}</CardTitle>
                    </CardHeader>
                    <CardContent className="p-4 pt-0">
                        <p className="font-bold text-primary">Rp {product.SellingPrice.toLocaleString()}</p>
                        <p className="text-sm text-gray-500">Stok: {product.Stock}</p>
                    </CardContent>
                    </Card>
                ))}
                </div>
                {meta && (
                     <div className="flex justify-between items-center mt-4">
                        <Button variant="outline" onClick={() => setPage(p => Math.max(1, p - 1))} disabled={page === 1}>Prev</Button>
                        <span>Page {meta.current_page} of {meta.total_pages}</span>
                        <Button variant="outline" onClick={() => setPage(p => Math.min(meta.total_pages, p + 1))} disabled={page === meta.total_pages}>Next</Button>
                    </div>
                )}
            </div>
            <div className="md:col-span-1">
                <Card className="sticky top-6">
                    <CardHeader><CardTitle>Keranjang</CardTitle></CardHeader>
                    <CardContent>
                        {cart.length === 0 ? <p className="text-center text-muted-foreground">Kosong</p> : (
                            <div className="space-y-4">
                                {cart.map(item => (
                                    <div key={item.product.ID} className="flex justify-between items-center">
                                        <div>
                                            <p className="font-medium">{item.product.Name}</p>
                                            <p className="text-sm text-muted-foreground">{item.quantity} x {item.product.SellingPrice}</p>
                                        </div>
                                        <Button variant="ghost" size="icon" onClick={() => removeFromCart(item.product.ID)}><Trash2 className="h-4 w-4 text-red-500"/></Button>
                                    </div>
                                ))}
                                <div className="border-t pt-2 font-bold flex justify-between">
                                    <span>Total</span>
                                    <span>Rp {total.toLocaleString()}</span>
                                </div>
                                <Input placeholder="Bayar..." type="number" value={payment} onChange={e => setPayment(e.target.value)} />
                                <Button className="w-full" onClick={() => checkoutMutation.mutate()} disabled={checkoutMutation.isPending}>Bayar</Button>
                            </div>
                        )}
                    </CardContent>
                </Card>
            </div>
        </div>
    )
}

function ProductsContent() {
    const [page, setPage] = useState(1);
    const [search, setSearch] = useState("");
    const { data: productsData } = useProducts({ page, limit: 10, search });
    const products = productsData?.data;
    const meta = productsData?.meta;
    const queryClient = useQueryClient();
    const [isCreating, setIsCreating] = useState(false);
    
    // Simple form state
    const [formData, setFormData] = useState({ Name: "", PurchasePrice: 0, SellingPrice: 0, Stock: 0, WarehouseID: 1 });

    const createMutation = useMutation({
        mutationFn: api.createProduct,
        onSuccess: () => {
            setIsCreating(false);
            queryClient.invalidateQueries({ queryKey: ["products"] });
            alert("Produk dibuat");
        }
    });

    const deleteMutation = useMutation({
        mutationFn: api.deleteProduct,
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ["products"] })
    });

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <h2 className="text-2xl font-bold">Manajemen Produk</h2>
                <Button onClick={() => setIsCreating(!isCreating)}><Plus className="mr-2 h-4 w-4" /> Tambah Produk</Button>
            </div>
            
            <div className="flex gap-2 w-full md:w-1/3">
                 <Input placeholder="Cari produk..." value={search} onChange={e => { setSearch(e.target.value); setPage(1); }} />
            </div>

            {isCreating && (
                <Card className="p-4 bg-gray-50 mb-4">
                    <div className="grid grid-cols-2 gap-4">
                        <Input placeholder="Nama Produk" onChange={e => setFormData({...formData, Name: e.target.value})} />
                        <Input placeholder="Stok" type="number" onChange={e => setFormData({...formData, Stock: parseInt(e.target.value)})} />
                        <Input placeholder="Harga Beli" type="number" onChange={e => setFormData({...formData, PurchasePrice: parseFloat(e.target.value)})} />
                        <Input placeholder="Harga Jual" type="number" onChange={e => setFormData({...formData, SellingPrice: parseFloat(e.target.value)})} />
                        <Input placeholder="ID Gudang" type="number" value={formData.WarehouseID} onChange={e => setFormData({...formData, WarehouseID: parseInt(e.target.value)})} />
                        <Button onClick={() => createMutation.mutate(formData)}>Simpan</Button>
                    </div>
                </Card>
            )}

            <div className="bg-white rounded-lg border">
                <table className="w-full text-sm text-left">
                    <thead className="text-xs uppercase bg-gray-50 border-b">
                        <tr>
                            <th className="px-6 py-3">ID</th>
                            <th className="px-6 py-3">Nama</th>
                            <th className="px-6 py-3">Stok</th>
                            <th className="px-6 py-3">Harga Jual</th>
                            <th className="px-6 py-3">Aksi</th>
                        </tr>
                    </thead>
                    <tbody>
                        {products?.map(p => (
                            <tr key={p.ID} className="bg-white border-b hover:bg-gray-50">
                                <td className="px-6 py-4">{p.ID}</td>
                                <td className="px-6 py-4">{p.Name}</td>
                                <td className="px-6 py-4">{p.Stock}</td>
                                <td className="px-6 py-4">{p.SellingPrice}</td>
                                <td className="px-6 py-4">
                                    <Button variant="destructive" size="sm" onClick={() => { if(confirm("Hapus?")) deleteMutation.mutate(p.ID) }}>Hapus</Button>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
            {meta && (
                 <div className="flex justify-between items-center mt-4">
                    <Button variant="outline" onClick={() => setPage(p => Math.max(1, p - 1))} disabled={page === 1}>Prev</Button>
                    <span>Page {meta.current_page} of {meta.total_pages}</span>
                    <Button variant="outline" onClick={() => setPage(p => Math.min(meta.total_pages, p + 1))} disabled={page === meta.total_pages}>Next</Button>
                </div>
            )}
        </div>
    )
}

function UsersContent() {
    const { data: users } = useUsers();
    const queryClient = useQueryClient();
    const [formData, setFormData] = useState({ Username: "", Password: "", Role: "user", WarehouseID: 1 });

    const createMutation = useMutation({
        mutationFn: (data: any) => api.createUser(data), // Cast to any or match type, api expects Omit<User...> & {password}
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["users"] });
            alert("User dibuat");
        }
    });
    
    const deleteMutation = useMutation({
        mutationFn: api.deleteUser,
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ["users"] })
    });

    return (
        <div className="space-y-6">
            <h2 className="text-2xl font-bold">Manajemen User</h2>
             <Card className="p-4 bg-gray-50">
                <div className="flex gap-4">
                    <Input placeholder="Username" onChange={e => setFormData({...formData, Username: e.target.value})} />
                    <Input placeholder="Password" type="password" onChange={e => setFormData({...formData, Password: e.target.value})} />
                    <select className="border rounded p-2" onChange={e => setFormData({...formData, Role: e.target.value})}>
                        <option value="user">User</option>
                        <option value="admin">Admin</option>
                    </select>
                    <Button onClick={() => createMutation.mutate({ ...formData, password: formData.Password })}>Tambah</Button>
                </div>
            </Card>

            <div className="bg-white rounded-lg border">
                <table className="w-full text-sm text-left">
                    <thead className="text-xs uppercase bg-gray-50 border-b">
                        <tr>
                            <th className="px-6 py-3">ID</th>
                            <th className="px-6 py-3">Username</th>
                            <th className="px-6 py-3">Role</th>
                            <th className="px-6 py-3">Aksi</th>
                        </tr>
                    </thead>
                    <tbody>
                        {users?.map(u => (
                            <tr key={u.ID} className="bg-white border-b hover:bg-gray-50">
                                <td className="px-6 py-4">{u.ID}</td>
                                <td className="px-6 py-4">{u.Username}</td>
                                <td className="px-6 py-4">{u.Role}</td>
                                <td className="px-6 py-4">
                                    <Button variant="destructive" size="sm" onClick={() => deleteMutation.mutate(u.ID)}>Hapus</Button>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    )
}

function WarehousesContent() {
    const { data: warehouses } = useWarehouses();
    const queryClient = useQueryClient();
    const [formData, setFormData] = useState({ Name: "", Address: "" });

    const createMutation = useMutation({
        mutationFn: api.createWarehouse,
        onSuccess: () => {
             queryClient.invalidateQueries({ queryKey: ["warehouses"] });
             setFormData({ Name: "", Address: "" });
        }
    });

     const deleteMutation = useMutation({
        mutationFn: api.deleteWarehouse,
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ["warehouses"] })
    });

    return (
        <div className="space-y-6">
            <h2 className="text-2xl font-bold">Manajemen Gudang</h2>
            <Card className="p-4 bg-gray-50">
                <div className="flex gap-4">
                    <Input placeholder="Nama Gudang" value={formData.Name} onChange={e => setFormData({...formData, Name: e.target.value})} />
                    <Input placeholder="Alamat" value={formData.Address} onChange={e => setFormData({...formData, Address: e.target.value})} />
                    <Button onClick={() => createMutation.mutate(formData)}>Tambah</Button>
                </div>
            </Card>

            <div className="bg-white rounded-lg border">
                <table className="w-full text-sm text-left">
                    <thead className="text-xs uppercase bg-gray-50 border-b">
                        <tr>
                            <th className="px-6 py-3">ID</th>
                            <th className="px-6 py-3">Nama</th>
                            <th className="px-6 py-3">Alamat</th>
                            <th className="px-6 py-3">Aksi</th>
                        </tr>
                    </thead>
                    <tbody>
                        {warehouses?.map(w => (
                            <tr key={w.ID} className="bg-white border-b hover:bg-gray-50">
                                <td className="px-6 py-4">{w.ID}</td>
                                <td className="px-6 py-4">{w.Name}</td>
                                <td className="px-6 py-4">{w.Address}</td>
                                <td className="px-6 py-4">
                                    <Button variant="destructive" size="sm" onClick={() => deleteMutation.mutate(w.ID)}>Hapus</Button>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    )
}

function ReportsContent() {
    const [date, setDate] = useState(new Date().toISOString().split('T')[0]); // YYYY-MM-DD
    const { data: report, refetch } = useQuery({
        queryKey: ["report", date],
        queryFn: () => {
             // Convert YYYY-MM-DD to DD-MM-YYYY
             const [y, m, d] = date.split("-");
             return api.getReport(`${d}-${m}-${y}`);
        }
    });
    // Import useQuery here since it is inside function? No, need to import at top. 
    // Wait, useQuery is imported.
    // 'useQuery' ...
    
    // We already imported useMutation, useQueryClient. Need useQuery.
    // Actually imported in api.ts hook, but here we use raw.
    // Let's rely on my import at top.

    return (
        <div className="space-y-6">
             <h2 className="text-2xl font-bold">Laporan Penjualan</h2>
             <div className="flex items-center gap-4">
                 <Input type="date" value={date} onChange={e => setDate(e.target.value)} className="w-48" />
                 <Button onClick={() => refetch()}>Refresh</Button>
             </div>

             {report && (
                 <div className="grid grid-cols-3 gap-6">
                     <Card>
                         <CardHeader><CardTitle>Total Penjualan</CardTitle></CardHeader>
                         <CardContent className="text-2xl font-bold">Rp {report.summary.total_sales.toLocaleString()}</CardContent>
                     </Card>
                     <Card>
                         <CardHeader><CardTitle>Total Profit</CardTitle></CardHeader>
                         <CardContent className="text-2xl font-bold text-green-600">Rp {report.summary.total_profit.toLocaleString()}</CardContent>
                     </Card>
                     <Card>
                         <CardHeader><CardTitle>Jumlah Transaksi</CardTitle></CardHeader>
                         <CardContent className="text-2xl font-bold">{report.summary.transaction_count}</CardContent>
                     </Card>
                 </div>
             )}

             <div className="bg-white rounded-lg border mt-6">
                <table className="w-full text-sm text-left">
                    <thead className="text-xs uppercase bg-gray-50 border-b">
                        <tr>
                            <th className="px-6 py-3">ID</th>
                            <th className="px-6 py-3">Waktu</th>
                            <th className="px-6 py-3">Total</th>
                            <th className="px-6 py-3">Profit</th>
                        </tr>
                    </thead>
                    <tbody>
                        {report?.transactions?.map((t: any) => (
                             <tr key={t.ID} className="bg-white border-b hover:bg-gray-50">
                                <td className="px-6 py-4">TRX-{t.ID}</td>
                                <td className="px-6 py-4">{new Date(t.CreatedAt).toLocaleTimeString()}</td>
                                <td className="px-6 py-4">Rp {t.Total.toLocaleString()}</td>
                                <td className="px-6 py-4">Rp {t.Profit.toLocaleString()}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
             </div>
        </div>
    )
}
