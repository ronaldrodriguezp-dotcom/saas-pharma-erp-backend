package inventory

import (
	"database/sql"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Crear producto
func (r *Repository) CreateProduct(name, description string) error {
	_, err := r.db.Exec("INSERT INTO products (name, description) VALUES ($1, $2)", name, description)
	return err
}

// Crear lote
func (r *Repository) CreateBatch(productID int, lotNumber string, expirationDate time.Time, quantity int) error {
	_, err := r.db.Exec(
		"INSERT INTO batches (product_id, lot_number, expiration_date, quantity) VALUES ($1, $2, $3, $4)",
		productID, lotNumber, expirationDate, quantity,
	)
	return err
}

// Registrar movimiento de stock
func (r *Repository) CreateMovement(batchID int, movementType string, quantity int) error {
	_, err := r.db.Exec(
		"INSERT INTO stock_movements (batch_id, movement_type, quantity) VALUES ($1, $2, $3)",
		batchID, movementType, quantity,
	)
	return err
}

// Obtener stock actual
func (r *Repository) GetStocks() ([]Batch, error) {
    if r.db == nil {
        return nil, fmt.Errorf("DB no inicializada")
    }

    rows, err := r.db.Query("SELECT id, product_id, lot_number, expiration_date, quantity, created_at FROM batches")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    // ...
}

	defer rows.Close()

	var batches []Batch
	for rows.Next() {
		var b Batch
		if err := rows.Scan(&b.ID, &b.ProductID, &b.LotNumber, &b.ExpirationDate, &b.Quantity, &b.CreatedAt); err != nil {
			return nil, err
		}
		batches = append(batches, b)
	}
	return batches, nil
}
