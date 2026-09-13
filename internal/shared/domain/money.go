package domain

import "errors"

type Currency string

const (
CLP Currency = "CLP"
USD Currency = "USD"
)

// Constante de escala para costos internos (4 decimales = 10000)
// Ej: $1.33 CLP de costo se guarda como 13333.
const CostScaleFactor int64 = 10000

type Money struct {
Amount   int64 // Representación fixed-point
Currency Currency
}

func NewCostMoney(amount int64, currency Currency) (Money, error) {
if amount < 0 {
return Money{}, errors.New("el monto monetario no puede ser negativo")
}
if currency == "" {
return Money{}, errors.New("la moneda es obligatoria")
}
return Money{Amount: amount, Currency: currency}, nil
}