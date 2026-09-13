package ports

import (
"context"
"fortiasaass-kio-speed/internal/catalog/domain"
)

// ProductRepository Puerto para la persistencia del agregado Product y sus presentaciones
type ProductRepository interface {
SaveProduct(ctx context.Context, product domain.Product) error
FindProductByID(ctx context.Context, tenantID string, productID string) (domain.Product, error)

SavePresentation(ctx context.Context, presentation domain.ProductPresentation) error
FindPresentationsByProductID(ctx context.Context, tenantID string, productID string) ([]domain.ProductPresentation, error)
FindPresentationByGTIN(ctx context.Context, tenantID string, gtin string) (domain.ProductPresentation, error)
}

// RegulatoryRepository Puerto para datos regulatorios, perfiles medicinales y registros ISP
type RegulatoryRepository interface {
SaveMedicinalProfile(ctx context.Context, profile domain.MedicinalProductProfile) error
FindMedicinalProfileByProductID(ctx context.Context, tenantID string, productID string) (domain.MedicinalProductProfile, error)

SaveActiveIngredient(ctx context.Context, ingredient domain.ActiveIngredient) error
FindActiveIngredientByID(ctx context.Context, ingredientID string) (domain.ActiveIngredient, error)
}