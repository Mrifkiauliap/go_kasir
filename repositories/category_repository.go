package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/domain"
)

// CategoryRepository defines the interface for category data operations.
type CategoryRepository interface {
	GetAll() ([]domain.Category, error)
	GetByID(id int) (domain.Category, error)
	Create(category domain.Category) (domain.Category, error)
	Update(id int, category domain.Category) (domain.Category, error)
	Delete(id int) error
}

type postgresCategoryRepository struct {
	db *sql.DB
}

// NewPostgresCategoryRepository creates a new instance of CategoryRepository.
func NewPostgresCategoryRepository(db *sql.DB) CategoryRepository {
	return &postgresCategoryRepository{
		db: db,
	}
}

func (r *postgresCategoryRepository) GetAll() ([]domain.Category, error) {
	rows, err := r.db.Query("SELECT id, nama, deskripsi FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Nama, &c.Deskripsi); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *postgresCategoryRepository) GetByID(id int) (domain.Category, error) {
	var c domain.Category
	err := r.db.QueryRow("SELECT id, nama, deskripsi FROM categories WHERE id=$1", id).Scan(&c.ID, &c.Nama, &c.Deskripsi)
	if err != nil {
		if err == sql.ErrNoRows {
			return c, fmt.Errorf("category dengan id %d tidak ditemukan", id)
		}
		return c, err
	}
	return c, nil
}

func (r *postgresCategoryRepository) Create(category domain.Category) (domain.Category, error) {
	err := r.db.QueryRow(
		"INSERT INTO categories (nama, deskripsi) VALUES ($1, $2) RETURNING id",
		category.Nama, category.Deskripsi).Scan(&category.ID)
	if err != nil {
		return category, err
	}
	return category, nil
}

func (r *postgresCategoryRepository) Update(id int, category domain.Category) (domain.Category, error) {
	_, err := r.db.Exec(
		"UPDATE categories SET nama=$1, deskripsi=$2 WHERE id=$3",
		category.Nama, category.Deskripsi, id)
	if err != nil {
		return category, err
	}
	category.ID = id
	return category, nil
}

func (r *postgresCategoryRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM categories WHERE id=$1", id)
	return err
}
