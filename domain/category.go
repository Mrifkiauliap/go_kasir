package domain

// Category merepresentasikan produk dalam sistem kasir
type Category struct {
	ID        int      `json:"id"`
	Nama      string   `json:"nama"`
	Deskripsi string   `json:"deskripsi"`
	Products  []Produk `json:"products,omitempty"`
}
