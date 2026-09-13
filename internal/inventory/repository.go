package inventory

import (
    "context"
    "database/sql"
    "time"
)

type StockItem struct {
    ID             string    `json:"id"`
    TenantID       string    `json:"tenant_id"`
    WarehouseID    string    `json:"warehouse_id"`
    ProductID      string    `json:"product_id"`
    LotNumber      string    `json:"lot_number"`
    ExpirationDate time.Time `json:"expiration_date"`
    Quantity       float64   `json:"quantity"`
    UpdatedAt      time.Time `json:"updated_at"`
}

type Repository struct {
    db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
    return &Repository{db: db}
}

func (r *Repository) GetStocks(ctx context.Context, tenantID string) ([]StockItem, error) {
    query := `
        SELECT id, tenant_id, warehouse_id, product_id, lot_number, expiration_date, quantity, updated_at
        FROM inventory_stocks
        WHERE tenant_id = $1 AND quantity > 0
        ORDER BY expiration_date ASC`

    rows, err := r.db.QueryContext(ctx, query, tenantID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var stocks []StockItem
    for rows.Next() {
        var s StockItem
        if err := rows.Scan(&s.ID, &s.TenantID, &s.WarehouseID, &s.ProductID, &s.LotNumber, &s.ExpirationDate, &s.Quantity, &s.UpdatedAt); err != nil {
            return nil, err
        }
        stocks = append(stocks, s)
    }

    return stocks, nil
}