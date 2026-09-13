package inventory

import (
	"errors"
	"time"
)

// Entidades principales
type Product struct {
	ID          int
	Name        string
	Description string
	CreatedAt   time.Time
}

type Batch struct {
	ID             int
	ProductID      int
	LotNumber      string
	ExpirationDate time.Time
	Quantity       int
	CreatedAt      time.Time
}

type StockMovement struct {
	ID           int
	BatchID      int
	MovementType string // "entrada" o "salida"
	Quantity     int
	CreatedAt    time.Time
}

// Servicio con reglas de negocio
type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Crear producto
func (s *Service) CreateProduct(name, description string) error {
	if name == "" {
		return errors.New("el nombre del producto no puede estar vacío")
	}
	return s.repo.CreateProduct(name, description)
}

// Registrar lote
func (s *Service) RegisterBatch(productID int, lotNumber string, expirationDate time.Time, quantity int) error {
	if quantity <= 0 {
		return errors.New("la cantidad debe ser mayor a cero")
	}
	if expirationDate.Before(time.Now()) {
		return errors.New("la fecha de vencimiento no puede ser pasada")
	}
	return s.repo.CreateBatch(productID, lotNumber, expirationDate, quantity)
}

// Movimiento de stock
func (s *Service) RegisterMovement(batchID int, movementType string, quantity int) error {
	if quantity <= 0 {
		return errors.New("la cantidad debe ser mayor a cero")
	}
	if movementType != "entrada" && movementType != "salida" {
		return errors.New("tipo de movimiento inválido")
	}
	return s.repo.CreateMovement(batchID, movementType, quantity)
}
