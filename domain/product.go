package domain

// Produk merepresentasikan produk dalam sistem kasir
type Produk struct {
	ID           int    `json:"id"`
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name,omitempty"`
	Nama         string `json:"nama"`
	Harga        int    `json:"harga"`
	Stok         int    `json:"stok"`
}
