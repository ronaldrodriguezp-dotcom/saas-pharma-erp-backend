package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"fortiasaass-kio-speed/internal/purchasing/domain"
	"fortiasaass-kio-speed/internal/purchasing/ports"
)

var (
	ErrGoodsReceiptNotFound = errors.New("recepción de mercancía no encontrada")
)

type goodsReceiptRepository struct {
	db *sql.DB
}

func NewGoodsReceiptRepository(db *sql.DB) ports.GoodsReceiptRepository {
	return &goodsReceiptRepository{db: db}
}

func (r *goodsReceiptRepository) Create(ctx context.Context, gr *domain.GoodsReceipt) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Insertar la cabecera del GoodsReceipt
	qHeader := `
INSERT INTO goods_receipts (
id, tenant_id, purchase_order_id, supplier_id, warehouse_id, 
document_type, document_number, status, received_by, received_at, 
created_at, updated_at, version
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err = tx.ExecContext(ctx, qHeader,
		gr.ID, gr.TenantID, gr.PurchaseOrderID, gr.SupplierID, gr.WarehouseID,
		gr.DocumentType, gr.DocumentNumber, gr.Status, gr.ReceivedBy, gr.ReceivedAt,
		timeNow(), timeNow(), gr.Version,
	)
	if err != nil {
		return fmt.Errorf("error insertando cabecera de recepción: %w", err)
	}

	// 2. Insertar las líneas asociadas
	qLine := `
INSERT INTO goods_receipt_lines (
tenant_id, goods_receipt_id, po_line_id, presentation_id, 
lot_number, expiration_date, received_quantity, accepted_quantity, 
rejected_quantity, rejection_reason
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	for _, line := range gr.Lines {
		_, err = tx.ExecContext(ctx, qLine,
			gr.TenantID, gr.ID, line.POLineID, line.PresentationID,
			line.LotNumber, line.ExpirationDate, line.ReceivedQuantity,
			line.AcceptedQuantity, line.RejectedQuantity, line.RejectionReason,
		)
		if err != nil {
			return fmt.Errorf("error insertando línea de recepción: %w", err)
		}
	}

	return tx.Commit()
}

func (r *goodsReceiptRepository) Update(ctx context.Context, gr *domain.GoodsReceipt) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Actualizar estado y versión con concurrencia optimista
	qUpdate := `
UPDATE goods_receipts 
SET status = $1, updated_at = NOW(), version = version + 1
WHERE id = $2 AND tenant_id = $3 AND version = $4`

	res, err := tx.ExecContext(ctx, qUpdate, gr.Status, gr.ID, gr.TenantID, gr.Version)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("conflicto de concurrencia o recepción no encontrada al actualizar")
	}

	err = tx.Commit()
	if err == nil {
		gr.Version++
	}
	return err
}

func (r *goodsReceiptRepository) GetByID(ctx context.Context, tenantID domain.TenantID, id domain.GoodsReceiptID) (*domain.GoodsReceipt, error) {
	qHeader := `
SELECT purchase_order_id, supplier_id, warehouse_id, document_type, document_number, status, received_by, received_at, version
FROM goods_receipts WHERE tenant_id = $1 AND id = $2`

	gr := &domain.GoodsReceipt{
		ID:       id,
		TenantID: tenantID,
		Lines:    make([]domain.GoodsReceiptLine, 0),
	}

	var poID, suppID, whID, docType, docNum, status, receiver string
	err := r.db.QueryRowContext(ctx, qHeader, tenantID, id).Scan(
		&poID, &suppID, &whID, &docType, &docNum, &status, &receiver, &gr.ReceivedAt, &gr.Version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGoodsReceiptNotFound
		}
		return nil, err
	}

	gr.PurchaseOrderID = domain.PurchaseOrderID(poID)
	gr.SupplierID = domain.SupplierID(suppID)
	gr.WarehouseID = domain.LocationID(whID)
	gr.DocumentType = domain.DocumentType(docType)
	gr.DocumentNumber = docNum
	gr.Status = domain.GoodsReceiptStatus(status)
	gr.ReceivedBy = receiver

	// Cargar líneas
	qLines := `
SELECT id, po_line_id, presentation_id, lot_number, expiration_date, received_quantity, accepted_quantity, rejected_quantity, rejection_reason
FROM goods_receipt_lines WHERE tenant_id = $1 AND goods_receipt_id = $2`

	rows, err := r.db.QueryContext(ctx, qLines, tenantID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var l domain.GoodsReceiptLine
		var lineID string // Cambiado de int a string (el ID en la BD es VARCHAR)

		err := rows.Scan(
			&lineID,
			&l.POLineID,
			&l.PresentationID,
			&l.LotNumber,
			&l.ExpirationDate,
			&l.ReceivedQuantity,
			&l.AcceptedQuantity,
			&l.RejectedQuantity,
			&l.RejectionReason,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando línea de recepción: %w", err)
		}

		l.ID = lineID
		gr.Lines = append(gr.Lines, l)
	}

	return gr, nil
}

func timeNow() interface{} {
	return sql.NullTime{} // Se usa la función nativa NOW() o se pasa desde dominio
}
