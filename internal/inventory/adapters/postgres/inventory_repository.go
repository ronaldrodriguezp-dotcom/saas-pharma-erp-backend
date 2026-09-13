package postgres

import (
"context"
"database/sql"
"fmt"
"time"

"fortiasaass-kio-speed/internal/inventory/domain"
)

type InventoryRepository struct {
db *sql.DB
}

func NewInventoryRepository(db *sql.DB) *InventoryRepository {
return &InventoryRepository{db: db}
}

// ReceiveGoods procesa las líneas aceptadas de una recepción de mercancía,
// creando o incrementando el saldo de los lotes físicos bajo control estricto de tenant.
func (r *InventoryRepository) ReceiveGoods(ctx context.Context, tenantID string, receiptID string, lines []ReceiptLineData) error {
tx, err := r.db.BeginTx(ctx, nil)
if err != nil {
return err
}
defer tx.Rollback()

qUpsert := `
INSERT INTO inventory_lots (
id, tenant_id, presentation_id, warehouse_id, lot_number, 
expiration_date, received_at, available_qty, version
)
VALUES ($1, $2, $3, $4, $5, $6, NOW(), $7, 1)
ON CONFLICT (tenant_id, id) 
DO UPDATE SET 
available_qty = inventory_lots.available_qty + EXCLUDED.available_qty,
version = inventory_lots.version + 1`

for i, line := range lines {
lotID := fmt.Sprintf("lot-%s-%d", receiptID, i)
_, err = tx.ExecContext(ctx, qUpsert,
lotID,
tenantID,
line.PresentationID,
line.LocationID,
line.LotNumber,
line.ExpirationDate,
line.Quantity,
)
if err != nil {
return fmt.Errorf("error persistiendo lote de inventario: %w", err)
}
}

return tx.Commit()
}

// GetAvailableLots recupera todos los lotes con stock > 0 para un producto y bodega específicos dentro del tenant.
func (r *InventoryRepository) GetAvailableLots(ctx context.Context, tenantID, warehouseID, presentationID string) ([]domain.StockLot, error) {
query := `
SELECT id, presentation_id, lot_number, expiration_date, received_at, available_qty
FROM inventory_lots
WHERE tenant_id = $1 AND warehouse_id = $2 AND presentation_id = $3 AND available_qty > 0`

rows, err := r.db.QueryContext(ctx, query, tenantID, warehouseID, presentationID)
if err != nil {
return nil, fmt.Errorf("error consultando lotes disponibles: %w", err)
}
defer rows.Close()

var lots []domain.StockLot
for rows.Next() {
var lot domain.StockLot
if err := rows.Scan(&lot.ID, &lot.PresentationID, &lot.LotNumber, &lot.ExpirationDate, &lot.ReceivedAt, &lot.AvailableQty); err != nil {
return nil, fmt.Errorf("error leyendo fila de lote: %w", err)
}
lots = append(lots, lot)
}

return lots, nil
}

// UpdateLotQuantity actualiza el stock disponible de un lote con aislamiento multitenant.
func (r *InventoryRepository) UpdateLotQuantity(ctx context.Context, tenantID string, lotID string, newQty int) error {
query := `
UPDATE inventory_lots
SET available_qty = $1, version = version + 1
WHERE tenant_id = $2 AND id = $3`

res, err := r.db.ExecContext(ctx, query, newQty, tenantID, lotID)
if err != nil {
return fmt.Errorf("error actualizando cantidad de lote: %w", err)
}

rows, err := res.RowsAffected()
if err != nil {
return err
}
if rows == 0 {
return fmt.Errorf("lote %s no encontrado para el tenant", lotID)
}

return nil
}

type ReceiptLineData struct {
PresentationID string
LotNumber      string
ExpirationDate time.Time
Quantity       int
LocationID     string
}