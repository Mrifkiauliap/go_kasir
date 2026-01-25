package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// Produk merepresentasikan produk dalam sistem kasir
type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

// Kategori merepresentasikan kategori produk
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// In-memory storage untuk produk dan kategori
var produk = []Produk{
	{ID: 1, Nama: "Indomie Godog", Harga: 3500, Stok: 10},
	{ID: 2, Nama: "Vit 1000ml", Harga: 3000, Stok: 40},
	{ID: 3, Nama: "kecap", Harga: 12000, Stok: 20},
}
var categories = []Category{
	{ID: 1, Name: "Makanan", Description: "Berbagai jenis makanan"},
	{ID: 2, Name: "Minuman", Description: "Berbagai jenis minuman"},
	{ID: 3, Name: "Obat", Description: "Berbagai jenis obat-obatan"},
}

// Mutex dipake untuk mengontrol akses ke data
var produkMutex = &sync.RWMutex{}
var categoryMutex = &sync.RWMutex{}

// Handler untuk /produk (Get All & Create)
func handleProdukCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		produkMutex.RLock()
		defer produkMutex.RUnlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(produk)

	case "POST":
		produkMutex.Lock()
		defer produkMutex.Unlock()

		var produkBaru Produk
		err := json.NewDecoder(r.Body).Decode(&produkBaru)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		produkBaru.ID = len(produk) + 1
		produk = append(produk, produkBaru)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) // 201
		json.NewEncoder(w).Encode(produkBaru)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handler untuk /produk/{id} (Get by ID, Update, Delete)
func handleProdukItem(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		produkMutex.RLock()
		defer produkMutex.RUnlock()

		for _, p := range produk {
			if p.ID == id {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(p)
				return
			}
		}
		http.Error(w, "Produk belum ada", http.StatusNotFound)

	case "PUT":
		produkMutex.Lock()
		defer produkMutex.Unlock()

		var updateProduk Produk
		err := json.NewDecoder(r.Body).Decode(&updateProduk)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		for i := range produk {
			if produk[i].ID == id {
				updateProduk.ID = id
				produk[i] = updateProduk
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(updateProduk)
				return
			}
		}
		http.Error(w, "Produk belum ada", http.StatusNotFound)

	case "DELETE":
		produkMutex.Lock()
		defer produkMutex.Unlock()

		for i, p := range produk {
			if p.ID == id {
				// Bikin slice baru dengan data sebelum dan sesudah index
				produk = append(produk[:i], produk[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		http.Error(w, "Produk belum ada", http.StatusNotFound)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handler untuk /kategori (Get All & Create)
func handleCategoryCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		categoryMutex.RLock()
		defer categoryMutex.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(categories)

	case "POST":
		categoryMutex.Lock()
		defer categoryMutex.Unlock()

		var newCategory Category
		if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		newCategory.ID = len(categories) + 1
		categories = append(categories, newCategory)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newCategory)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handler untuk /kategori/{id} (Get by ID, Update, Delete)
func handleCategoryItem(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Kategori ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		categoryMutex.RLock()
		defer categoryMutex.RUnlock()

		for _, c := range categories {
			if c.ID == id {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(c)
				return
			}
		}
		http.Error(w, "Kategori not found", http.StatusNotFound)

	case "PUT":
		categoryMutex.Lock()
		defer categoryMutex.Unlock()

		var updatedCategory Category
		if err := json.NewDecoder(r.Body).Decode(&updatedCategory); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		for i := range categories {
			if categories[i].ID == id {
				updatedCategory.ID = id
				categories[i] = updatedCategory
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(updatedCategory)
				return
			}
		}
		http.Error(w, "Kategori not found", http.StatusNotFound)

	case "DELETE":
		categoryMutex.Lock()
		defer categoryMutex.Unlock()

		for i, c := range categories {
			if c.ID == id {
				categories = append(categories[:i], categories[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		http.Error(w, "Kategori not found", http.StatusNotFound)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Main function to start the server
func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	// Produk Handlers
	http.HandleFunc("/produk", handleProdukCollection)
	http.HandleFunc("/produk/", handleProdukItem)

	// Kategori Handlers
	http.HandleFunc("/categories", handleCategoryCollection)
	http.HandleFunc("/categories/", handleCategoryItem)

	fmt.Println("Server running di http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}
