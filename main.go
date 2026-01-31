package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Load config from .env file if exists
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}
	}

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	productRepo := repositories.NewPostgresProductRepository(db)
	categoryRepo := repositories.NewPostgresCategoryRepository(db)

	// Initialize services
	productService := services.NewProductService(productRepo, categoryRepo)
	categoryService := services.NewCategoryService(categoryRepo, productRepo)

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Initialize router
	r := mux.NewRouter()

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	}).Methods("GET")

	// Product routes
	productRouter := r.PathPrefix("/produk").Subrouter()
	productRouter.HandleFunc("", productHandler.GetAllProducts).Methods("GET")
	productRouter.HandleFunc("/{id:[0-9]+}", productHandler.GetProductByID).Methods("GET")
	productRouter.HandleFunc("", productHandler.CreateProduct).Methods("POST")
	productRouter.HandleFunc("/{id:[0-9]+}", productHandler.UpdateProduct).Methods("PUT")
	productRouter.HandleFunc("/{id:[0-9]+}", productHandler.DeleteProduct).Methods("DELETE")

	// Category routes
	categoryRouter := r.PathPrefix("/categories").Subrouter()
	categoryRouter.HandleFunc("", categoryHandler.GetAllCategories).Methods("GET")
	categoryRouter.HandleFunc("/{id:[0-9]+}", categoryHandler.GetCategoryByID).Methods("GET")
	categoryRouter.HandleFunc("", categoryHandler.CreateCategory).Methods("POST")
	categoryRouter.HandleFunc("/{id:[0-9]+}", categoryHandler.UpdateCategory).Methods("PUT")
	categoryRouter.HandleFunc("/{id:[0-9]+}", categoryHandler.DeleteCategory).Methods("DELETE")

	// Start server
	port := viper.GetString("SERVER_PORT")
	if port == "" {
		port = "8080" // Default port
	}

	fmt.Printf("Server running at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
