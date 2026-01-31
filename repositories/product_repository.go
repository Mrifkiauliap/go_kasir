package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/domain"
)

// ProductRepository defines the interface for product data operations.
type ProductRepository interface {
	GetAll() ([]domain.Produk, error)
	GetByID(id int) (domain.Produk, error)
	Create(product domain.Produk) (domain.Produk, error)
	Update(id int, product domain.Produk) (domain.Produk, error)
	Delete(id int) error
	GetByCategoryID(categoryID int) ([]domain.Produk, error)
}

type postgresProductRepository struct {
	db *sql.DB
}

// NewPostgresProductRepository creates a new instance of ProductRepository.
func NewPostgresProductRepository(db *sql.DB) ProductRepository {
	return &postgresProductRepository{
		db: db,
	}
}

func (r *postgresProductRepository) GetAll() ([]domain.Produk, error) {
	query := `
		SELECT p.id, p.nama, p.harga, p.stok, p.category_id, COALESCE(c.nama, '') as category_name
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Produk
	for rows.Next() {
		var p domain.Produk
		if err := rows.Scan(&p.ID, &p.Nama, &p.Harga, &p.Stok, &p.CategoryID, &p.CategoryName); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *postgresProductRepository) GetByID(id int) (domain.Produk, error) {
	query := `
		SELECT p.id, p.nama, p.harga, p.stok, p.category_id, COALESCE(c.nama, '') as category_name
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id=$1`
	var p domain.Produk
	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.Nama, &p.Harga, &p.Stok, &p.CategoryID, &p.CategoryName)
	if err != nil {
		if err == sql.ErrNoRows {
			return p, fmt.Errorf("product dengan id %d tidak ditemukan", id)
		}
		return p, err
	}
	return p, nil
}

func (r *postgresProductRepository) Create(product domain.Produk) (domain.Produk, error) {
	query := "INSERT INTO products (nama, harga, stok, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	err := r.db.QueryRow(
		query,
		product.Nama, product.Harga, product.Stok, product.CategoryID).Scan(&product.ID)
	if err != nil {
		return product, err
	}
	return product, nil
}

func (r *postgresProductRepository) Update(id int, product domain.Produk) (domain.Produk, error) {
	query := "UPDATE products SET nama=$1, harga=$2, stok=$3, category_id=$4 WHERE id=$5"
	_, err := r.db.Exec(
		query,
		product.Nama, product.Harga, product.Stok, product.CategoryID, id)
	if err != nil {
		return product, err
	}
	product.ID = id
	return product, nil
}

func (r *postgresProductRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id=$1", id)
	return err
}

func (r *postgresProductRepository) GetByCategoryID(categoryID int) ([]domain.Produk, error) {
	rows, err := r.db.Query("SELECT id, nama, harga, stok, category_id FROM products WHERE category_id=$1", categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Produk
	for rows.Next() {
		var p domain.Produk
		if err := rows.Scan(&p.ID, &p.Nama, &p.Harga, &p.Stok, &p.CategoryID); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
