package domain

import (
"errors"
"fmt"
)

// GlobalPlatformMaster Representa el repositorio de datos maestros curados a nivel global/plataforma
type GlobalPlatformMaster struct {
ActiveIngredients   map[string]ActiveIngredient
PharmaceuticalForms map[string]PharmaceuticalForm
}

func NewGlobalPlatformMaster() *GlobalPlatformMaster {
return &GlobalPlatformMaster{
ActiveIngredients:   make(map[string]ActiveIngredient),
PharmaceuticalForms: make(map[string]PharmaceuticalForm),
}
}

// TenantCatalogContext Representa el contexto comercial y operativo específico del Tenant
type TenantCatalogContext struct {
TenantID     string
LocalSKUs    map[string]Product
Presentations map[string]ProductPresentation
}

// CatalogInvariantValidator Validador integral para asegurar el cumplimiento de las reglas CAT-001 a CAT-010
type CatalogInvariantValidator struct{}

func (v CatalogInvariantValidator) ValidateProductCreation(product Product) error {
// CAT-001: ProductID inmutable y obligatorio
if product.ProductID == "" {
return errors.New("violación de CAT-001: ProductID es obligatorio e inmutable")
}
return nil
}

func (v CatalogInvariantValidator) ValidatePresentation(pres ProductPresentation, product Product) error {
// CAT-002: GTIN jamás es la PK del producto (validamos que tenga ID propio)
if pres.PresentationID == "" {
return errors.New("violación de CAT-002: la presentación requiere un PresentationID propio, el GTIN no es la PK")
}
// CAT-004: Asegurar vinculación correcta con el producto padre
if pres.ProductID != product.ProductID {
return fmt.Errorf("violación de vinculación: la presentación pertenece a un ProductID diferente (%s != %s)", pres.ProductID, product.ProductID)
}
return nil
}

func (v CatalogInvariantValidator) ValidateMedicinalProfile(profile MedicinalProductProfile, product Product) error {
// CAT-004: Un MedicinalProductProfile debe asociarse exclusivamente a un Product de tipo MEDICINAL
if product.Type != TypeMedicinal {
return errors.New("violación de CAT-004: se intentó asociar un perfil medicinal a un producto que no es de tipo MEDICINAL")
}
if profile.ProductID != product.ProductID {
return errors.New("violación de integridad: el perfil medicinal no coincide con el ProductID del padre")
}
return nil
}

func (v CatalogInvariantValidator) ValidateIngredientAssociation(ingredient ProductIngredient, globalMaster *GlobalPlatformMaster) error {
// CAT-005: Un ProductIngredient debe referenciar un ActiveIngredient válido del maestro global
if _, exists := globalMaster.ActiveIngredients[ingredient.IngredientID]; !exists {
return fmt.Errorf("violación de CAT-005: el principio activo ID '%s' no existe en el maestro global de plataforma", ingredient.IngredientID)
}
return nil
}