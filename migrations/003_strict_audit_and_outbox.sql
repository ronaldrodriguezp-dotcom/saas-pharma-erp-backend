-- FIX-05: Identidad única de InventoryLot (Evita duplicados concurrentes)
ALTER TABLE inventory_lots 
    DROP CONSTRAINT IF EXISTS uk_tenant_presentation_lot;
ALTER TABLE inventory_lots 
    ADD CONSTRAINT uk_tenant_presentation_lot UNIQUE (tenant_id, presentation_id, lot_number);

-- FIX-04: Identidad y eliminación de redundancias en InventoryBalance
ALTER TABLE inventory_balances 
    DROP CONSTRAINT IF EXISTS uk_tenant_balance;
ALTER TABLE inventory_balances 
    ADD CONSTRAINT uk_tenant_lot_location UNIQUE (tenant_id, lot_id, location_id);

-- FIX-02: Integridad referencial Cross-Tenant (Aislamiento físico)
ALTER TABLE inventory_lots ADD CONSTRAINT uq_tenant_lot_id UNIQUE (tenant_id, lot_id);
ALTER TABLE inventory_balances 
    ADD CONSTRAINT fk_balance_tenant_lot FOREIGN KEY (tenant_id, lot_id) 
    REFERENCES inventory_lots (tenant_id, lot_id) ON DELETE RESTRICT;

-- FIX-01: Restricciones defensivas contra saldos inválidos en PostgreSQL
ALTER TABLE inventory_balances ADD CONSTRAINT chk_on_hand_positive CHECK (on_hand >= 0);
ALTER TABLE inventory_balances ADD CONSTRAINT chk_reserved_positive CHECK (reserved >= 0);
ALTER TABLE inventory_balances ADD CONSTRAINT chk_reserved_lte_on_hand CHECK (reserved <= on_hand);

-- FIX-07: Infraestructura base de Transactional Outbox
CREATE TABLE IF NOT EXISTS transactional_outbox (
    event_id VARCHAR(128) PRIMARY KEY,
    event_type VARCHAR(255) NOT NULL,
    event_version INT NOT NULL DEFAULT 1,
    tenant_id VARCHAR(64) NOT NULL,
    aggregate_id VARCHAR(128) NOT NULL,
    occurred_at TIMESTAMP WITH TIME ZONE NOT NULL,
    correlation_id VARCHAR(128),
    causation_id VARCHAR(128),
    payload JSONB NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'PENDING',
    attempts INT NOT NULL DEFAULT 0,
    next_attempt_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMP WITH TIME ZONE
);
CREATE INDEX IF NOT EXISTS idx_outbox_pending ON transactional_outbox(tenant_id, status) WHERE status = 'PENDING';