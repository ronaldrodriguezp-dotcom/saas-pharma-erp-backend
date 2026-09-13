-- ==============================================================================
-- Módulo: Purchasing (Abastecimiento)
-- Agregado: GoodsReceipt (Recepción de Mercadería)
-- ==============================================================================

CREATE TABLE IF NOT EXISTS goods_receipts (
    id VARCHAR(128) NOT NULL,
    tenant_id VARCHAR(64) NOT NULL,
    purchase_order_id VARCHAR(128) NOT NULL,
    supplier_id VARCHAR(128) NOT NULL,
    warehouse_id VARCHAR(128) NOT NULL,
    document_type VARCHAR(50) NOT NULL,
    document_number VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL,
    received_by VARCHAR(128) NOT NULL,
    received_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    version INT NOT NULL DEFAULT 1,
    
    -- PK compuesta para garantizar el particionamiento lógico del tenant
    PRIMARY KEY (tenant_id, id),
    
    -- Defensas de integridad que reflejan los Enums de Go
    CONSTRAINT chk_gr_status CHECK (status IN ('DRAFT', 'CONFIRMED', 'CANCELLED')),
    CONSTRAINT chk_gr_doc_type CHECK (document_type IN ('INVOICE', 'DELIVERY_NOTE'))
);

-- Líneas de la recepción (Registro físico por Lote y Vencimiento)
CREATE TABLE IF NOT EXISTS goods_receipt_lines (
    id SERIAL,
    tenant_id VARCHAR(64) NOT NULL,
    goods_receipt_id VARCHAR(128) NOT NULL,
    po_line_id VARCHAR(128) NOT NULL,
    presentation_id VARCHAR(128) NOT NULL,
    lot_number VARCHAR(100) NOT NULL DEFAULT '',
    expiration_date DATE NOT NULL,
    received_quantity INT NOT NULL,
    accepted_quantity INT NOT NULL,
    rejected_quantity INT NOT NULL,
    rejection_reason TEXT NOT NULL DEFAULT '',
    
    PRIMARY KEY (tenant_id, id),
    
    -- FK Compuesta: Garantiza que la línea pertenezca a la recepción dentro del mismo tenant
    CONSTRAINT fk_gr_line_receipt FOREIGN KEY (tenant_id, goods_receipt_id) 
        REFERENCES goods_receipts (tenant_id, id) ON DELETE CASCADE,
        
    -- Defensas lógicas cuantitativas a nivel de base de datos
    CONSTRAINT chk_quantities_positive CHECK (received_quantity >= 0 AND accepted_quantity >= 0 AND rejected_quantity >= 0),
    CONSTRAINT chk_quantity_math CHECK (received_quantity = accepted_quantity + rejected_quantity)
);