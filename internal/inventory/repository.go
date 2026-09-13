package inventory

import (
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateProduct(name, description string) error {
	_, err := r.db.Exec("INSERT INTO products (name, description) VALUES ($1, $2)", name, description)
	return err
}
