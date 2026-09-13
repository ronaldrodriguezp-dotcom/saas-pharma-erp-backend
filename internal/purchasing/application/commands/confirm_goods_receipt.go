package commands

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// ConfirmGoodsReceiptCommand define la estructura de entrada esperada
type ConfirmGoodsReceiptCommand struct {
	ReceiptID string `json:"receipt_id"`
	Actor     string `json:"actor"`
}

// ConfirmGoodsReceiptHandler maneja la ejecución transaccional del comando
type ConfirmGoodsReceiptHandler struct {
	db *sql.DB
}

func NewConfirmGoodsReceiptHandler(db *sql.DB) *ConfirmGoodsReceiptHandler {
	return &ConfirmGoodsReceiptHandler{db: db}
}

func (h *ConfirmGoodsReceiptHandler) Handle(ctx context.Context, tenantID string, cmd ConfirmGoodsReceiptCommand) error {
	// 1. Iniciar la transacción para garantizar atomicidad (ACID)
	tx, err := h.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %w", err)
	}
	defer tx.Rollback()

	// 2. Verificar estado actual de la recepción y bloquear la fila (FOR UPDATE)
	var status, warehouseID string
	queryCheck := `SELECT status, warehouse_id FROM goods_receipts WHERE tenant_id = $1 AND id = $2 FOR UPDATE`
	err = tx.QueryRowContext(ctx, queryCheck, tenantID, cmd.ReceiptID).Scan(&status, &warehouseID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("la recepción no existe")
		}
		return fmt.Errorf("error al consultar recepción: %w", err)
	}

	if status != "DRAFT" {
		return fmt.Errorf("transición de estado inválida: la recepción está en estado %s", status)
	}

	// 3. Actualizar la recepción a CONFIRMED
	queryUpdateReceipt := `
		UPDATE goods_receipts 
		SET status = 'CONFIRMED', updated_at = NOW() 
		WHERE tenant_id = $1 AND id = $2`
	_, err = tx.ExecContext(ctx, queryUpdateReceipt, tenantID, cmd.ReceiptID)
	if err != nil {
		return fmt.Errorf("error al actualizar estado de recepción: %w", err)
	}

	// 4. Obtener las líneas de la recepción
	queryLines := `
		SELECT id, product_id, received_quantity, lot_number, expiration_date 
		FROM goods_receipt_lines 
		WHERE tenant_id = $1 AND goods_receipt_id = $2`
	rows, err := tx.QueryContext(ctx, queryLines, tenantID, cmd.ReceiptID)
	if err != nil {
		return fmt.Errorf("error al obtener líneas de recepción: %w", err)
	}
	defer rows.Close()

	// 5. Procesar e impactar el stock por lote
	for rows.Next() {
		var lineID, productID, lotNumber string
		var receivedQty float64
		var expirationDate time.Time

		if err := rows.Scan(&lineID, &productID, &receivedQty, &lotNumber, &expirationDate); err != nil {
			return fmt.Errorf("error escaneando línea para inventario: %w", err)
		}

		stockID := fmt.Sprintf("stk-%s-%s", warehouseID, lineID)

		queryUpsertStock := `
			INSERT INTO inventory_stocks (id, tenant_id, warehouse_id, product_id, lot_number, expiration_date, quantity, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
			ON CONFLICT (tenant_id, warehouse_id, product_id, lot_number)
			DO UPDATE SET 
				quantity = inventory_stocks.quantity + EXCLUDED.quantity,
				updated_at = NOW()`

		_, err = tx.ExecContext(ctx, queryUpsertStock, stockID, tenantID, warehouseID, productID, lotNumber, expirationDate, receivedQty)
		if err != nil {
			return fmt.Errorf("error al actualizar el stock de inventario: %w", err)
		}
	}

	// 6. Confirmar la transacción
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error al confirmar la transacción (commit): %w", err)
	}

	return nil
}
