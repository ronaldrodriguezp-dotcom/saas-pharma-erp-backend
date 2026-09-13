-- ==============================================================================
-- Módulo: Purchasing (Abastecimiento)
-- Agregado: Supplier (Proveedor)
-- ==============================================================================

CREATE TABLE IF NOT EXISTS suppliers (
    id VARCHAR(128) NOT NULL,
    tenant_id VARCHAR(64) NOT NULL,
    legal_name VARCHAR(255) NOT NULL,
    trade_name VARCHAR(255) NOT NULL DEFAULT '',
    tax_id VARCHAR(20) NOT NULL,
    supplier_type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    payment_terms VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    version INT NOT NULL DEFAULT 1,
    
    -- PK compuesta para garantizar el particionamiento lógico del tenant
    PRIMARY KEY (tenant_id, id),
    
    -- Invariante: No puede haber dos proveedores con el mismo RUT en el mismo tenant
    CONSTRAINT uq_tenant_tax_id UNIQUE (tenant_id, tax_id),
    
    -- Defensas de integridad que reflejan los Enums de Go
    CONSTRAINT chk_supplier_status CHECK (status IN ('ACTIVE', 'SUSPENDED', 'BLOCKED', 'INACTIVE')),
    CONSTRAINT chk_supplier_type CHECK (supplier_type IN ('DISTRIBUTOR', 'LABORATORY', 'WHOLESALER', 'GENERAL_SUPPLIER', 'SERVICE_PROVIDER', 'OTHER')),
    CONSTRAINT chk_payment_terms CHECK (payment_terms IN ('CASH', 'NET_30', 'NET_60', 'NET_90'))
);

-- Como Contacts y Addresses son Value Objects (sin ciclo de vida independiente),
-- y para respetar el aislamiento estricto, usamos tablas hijas con FK compuestas.

CREATE TABLE IF NOT EXISTS supplier_addresses (
    id SERIAL,
    tenant_id VARCHAR(64) NOT NULL,
    supplier_id VARCHAR(128) NOT NULL,
    street VARCHAR(255) NOT NULL,
    city VARCHAR(128) NOT NULL,
    region VARCHAR(128) NOT NULL,
    country VARCHAR(128) NOT NULL,
    postal_code VARCHAR(32) NOT NULL DEFAULT '',
    
    PRIMARY KEY (tenant_id, id),
    -- FK Compuesta: Garantiza que la dirección pertenezca al mismo tenant que el proveedor
    CONSTRAINT fk_supplier_address FOREIGN KEY (tenant_id, supplier_id) 
        REFERENCES suppliers (tenant_id, id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS supplier_contacts (
    id SERIAL,
    tenant_id VARCHAR(64) NOT NULL,
    supplier_id VARCHAR(128) NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL DEFAULT '',
    phone VARCHAR(50) NOT NULL DEFAULT '',
    
    PRIMARY KEY (tenant_id, id),
    -- FK Compuesta: Garantiza que el contacto pertenezca al mismo tenant que el proveedor
    CONSTRAINT fk_supplier_contact FOREIGN KEY (tenant_id, supplier_id) 
        REFERENCES suppliers (tenant_id, id) ON DELETE CASCADE
);