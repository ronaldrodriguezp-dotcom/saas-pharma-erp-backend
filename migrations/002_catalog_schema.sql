-- Módulo de Catálogo Maestro: Esquema Transaccional y Constraints Regulatorios (PostgreSQL)

CREATE TABLE IF NOT EXISTS catalog_products (
    tenant_id VARCHAR(64) NOT NULL,
    product_id VARCHAR(64) NOT NULL,
    sku VARCHAR(64) NOT NULL,
    commercial_name VARCHAR(255) NOT NULL,
    brand VARCHAR(128) NOT NULL,
    product_type VARCHAR(32) NOT NULL,
    status VARCHAR(32) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY (tenant_id, product_id),
    CONSTRAINT uk_tenant_sku UNIQUE (tenant_id, sku)
);

CREATE INDEX IF NOT EXISTS idx_catalog_products_type ON catalog_products(tenant_id, product_type);

CREATE TABLE IF NOT EXISTS catalog_presentations (
    tenant_id VARCHAR(64) NOT NULL,
    presentation_id VARCHAR(64) NOT NULL,
    product_id VARCHAR(64) NOT NULL,
    name VARCHAR(255) NOT NULL,
    base_units_quantity BIGINT NOT NULL,
    is_franchisable BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY (tenant_id, presentation_id),
    CONSTRAINT chk_base_units_positive CHECK (base_units_quantity > 0),
    CONSTRAINT fk_presentation_product FOREIGN KEY (tenant_id, product_id) 
        REFERENCES catalog_products(tenant_id, product_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_catalog_presentations_product ON catalog_presentations(tenant_id, product_id);

CREATE TABLE IF NOT EXISTS catalog_presentation_barcodes (
    tenant_id VARCHAR(64) NOT NULL,
    presentation_id VARCHAR(64) NOT NULL,
    code VARCHAR(128) NOT NULL,
    barcode_type VARCHAR(32) NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (tenant_id, presentation_id, code),
    -- CAT-003: Un GTIN activo no puede identificar simultáneamente dos presentaciones activas dentro del mismo tenant
    CONSTRAINT uk_tenant_active_gtin UNIQUE (tenant_id, code),
    CONSTRAINT fk_barcode_presentation FOREIGN KEY (tenant_id, presentation_id) 
        REFERENCES catalog_presentations(tenant_id, presentation_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS global_active_ingredients (
    ingredient_id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    cas_number VARCHAR(64),
    description TEXT
);

CREATE TABLE IF NOT EXISTS catalog_medicinal_profiles (
    tenant_id VARCHAR(64) NOT NULL,
    product_id VARCHAR(64) NOT NULL,
    form_id VARCHAR(64) NOT NULL,
    registration_number VARCHAR(128) NOT NULL,
    holder_name VARCHAR(255) NOT NULL,
    registration_status VARCHAR(64) NOT NULL,
    issued_date TIMESTAMP WITH TIME ZONE NOT NULL,
    expiry_date TIMESTAMP WITH TIME ZONE NOT NULL,
    dispensing_condition VARCHAR(64) NOT NULL,
    is_controlled BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (tenant_id, product_id),
    CONSTRAINT fk_medicinal_product FOREIGN KEY (tenant_id, product_id) 
        REFERENCES catalog_products(tenant_id, product_id) ON DELETE CASCADE
);