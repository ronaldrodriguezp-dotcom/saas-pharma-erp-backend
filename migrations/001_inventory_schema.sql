-- Módulo de Inventario: Esquema Transaccional y Constraints de Invariantes (PostgreSQL)

CREATE TABLE IF NOT EXISTS inventory_lots (
    tenant_id VARCHAR(64) NOT NULL,
    lot_id VARCHAR(64) NOT NULL,
    presentation_id VARCHAR(64) NOT NULL,
    lot_number VARCHAR(100) NOT NULL,
    expiration_date TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(32) NOT NULL,
    received_at TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY (tenant_id, lot_id)
);

CREATE INDEX IF NOT EXISTS idx_inventory_lots_tenant_product ON inventory_lots(tenant_id, presentation_id);
CREATE INDEX IF NOT EXISTS idx_inventory_lots_expiration ON inventory_lots(tenant_id, expiration_date);

CREATE TABLE IF NOT EXISTS inventory_balances (
    tenant_id VARCHAR(64) NOT NULL,
    balance_id VARCHAR(64) NOT NULL,
    lot_id VARCHAR(64) NOT NULL,
    location_id VARCHAR(64) NOT NULL,
    on_hand BIGINT NOT NULL DEFAULT 0,
    reserved BIGINT NOT NULL DEFAULT 0,
    version BIGINT NOT NULL DEFAULT 1,
    PRIMARY KEY (tenant_id, balance_id),
    CONSTRAINT uk_tenant_lot_location UNIQUE (tenant_id, lot_id, location_id),
    -- Invariantes defensivas a nivel de Base de Datos (INV-001, INV-002, INV-003)
    CONSTRAINT chk_on_hand_non_negative CHECK (on_hand >= 0),
    CONSTRAINT chk_reserved_non_negative CHECK (reserved >= 0),
    CONSTRAINT chk_reserved_lte_on_hand CHECK (reserved <= on_hand)
);

CREATE INDEX IF NOT EXISTS idx_inventory_balances_lookup ON inventory_balances(tenant_id, location_id, lot_id);

CREATE TABLE IF NOT EXISTS inventory_movements (
    tenant_id VARCHAR(64) NOT NULL,
    movement_id VARCHAR(64) NOT NULL,
    lot_id VARCHAR(64) NOT NULL,
    location_id VARCHAR(64) NOT NULL,
    movement_type VARCHAR(64) NOT NULL,
    quantity_change BIGINT NOT NULL,
    previous_on_hand BIGINT NOT NULL,
    resulting_on_hand BIGINT NOT NULL,
    reason TEXT,
    reference_doc VARCHAR(128),
    actor_id VARCHAR(64) NOT NULL,
    correlation_id VARCHAR(64) NOT NULL,
    occurred_at TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY (tenant_id, movement_id),
    CONSTRAINT chk_movement_resulting_non_negative CHECK (resulting_on_hand >= 0)
);

CREATE INDEX IF NOT EXISTS idx_inventory_movements_lot ON inventory_movements(tenant_id, lot_id, occurred_at);

CREATE TABLE IF NOT EXISTS inventory_reservations (
    tenant_id VARCHAR(64) NOT NULL,
    reservation_id VARCHAR(64) NOT NULL,
    lot_id VARCHAR(64) NOT NULL,
    location_id VARCHAR(64) NOT NULL,
    quantity BIGINT NOT NULL,
    status VARCHAR(32) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY (tenant_id, reservation_id),
    CONSTRAINT chk_reservation_qty_positive CHECK (quantity > 0)
);

CREATE INDEX IF NOT EXISTS idx_inventory_reservations_status ON inventory_reservations(tenant_id, status, expires_at);
