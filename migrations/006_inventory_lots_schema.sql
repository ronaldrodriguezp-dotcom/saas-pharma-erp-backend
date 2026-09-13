-- ==============================================================================
-- Módulo: Inventory (Inventario y Control de Lotes)
-- ==============================================================================

CREATE TABLE IF NOT EXISTS inventory_lots (
    id VARCHAR(128) NOT NULL,
    tenant_id VARCHAR(64) NOT NULL,
    presentation_id VARCHAR(128) NOT NULL,
    warehouse_id VARCHAR(128) NOT NULL,
    lot_number VARCHAR(100) NOT NULL,
    expiration_date DATE NOT NULL,
    received_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    available_qty INT NOT NULL DEFAULT 0,
    version INT NOT NULL DEFAULT 1,

    PRIMARY KEY (tenant_id, id),
    
    -- Defensas lógicas para evitar cantidades negativas en inventario físico
    CONSTRAINT chk_inventory_qty_positive CHECK (available_qty >= 0)
);

-- Índice optimizado para búsquedas rápidas bajo políticas FEFO (por fecha de expiración) y FIFO (por fecha de recepción)
CREATE INDEX IF NOT EXISTS idx_inventory_lots_fefo 
ON inventory_lots (tenant_id, presentation_id, warehouse_id, expiration_date);