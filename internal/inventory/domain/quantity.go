package domain

import (
"errors"
"fmt"
)

// UnitOfMeasure Define las unidades de medida farmacéuticas y logísticas
type UnitOfMeasure string

const (
UnitTablet  UnitOfMeasure = "TABLET"     // Comprimido / Pastilla unitaria (Unidad base habitual)
UnitCapsule UnitOfMeasure = "CAPSULE"    // Cápsula unitaria
UnitBox     UnitOfMeasure = "BOX"        // Caja / Envase primario
UnitBlister UnitOfMeasure = "BLISTER"    // Blíster
UnitMl      UnitOfMeasure = "ML"         // Mililitros (Unidad base para líquidos)
UnitGram    UnitOfMeasure = "GRAM"       // Gramos (Unidad base para semisólidos)
UnitUnit    UnitOfMeasure = "UNIT"       // Unidad genérica
)

// Quantity Representa una cantidad asociada a su unidad de medida de forma segura
type Quantity struct {
BaseUnits int64         // Valor absoluto expresado en la unidad base mínima (ej. número total de comprimidos o mL)
Unit      UnitOfMeasure // Unidad de presentación o registro
Factor    int64         // Factor de conversión hacia la unidad base (ej. 1 para comprimido, 30 si es una caja de 30)
}

// NewQuantity Crea una cantidad validando que no sea negativa
func NewQuantity(value int64, unit UnitOfMeasure, factor int64) (Quantity, error) {
if value < 0 {
return Quantity{}, errors.New("la cantidad no puede ser negativa")
}
if factor <= 0 {
factor = 1
}

return Quantity{
BaseUnits: value * factor,
Unit:      unit,
Factor:    factor,
}, nil
}

// Add Suma cantidades asegurando compatibilidad o conversión a unidad base
func (q Quantity) Add(other Quantity) (Quantity, error) {
return Quantity{
BaseUnits: q.BaseUnits + other.BaseUnits,
Unit:      q.Unit,
Factor:    q.Factor,
}, nil
}

// Subtract Resta cantidades asegurando que no queden valores negativos
func (q Quantity) Subtract(other Quantity) (Quantity, error) {
if q.BaseUnits < other.BaseUnits {
return Quantity{}, errors.New("la operación resultaría en una cantidad negativa en unidades base")
}
return Quantity{
BaseUnits: q.BaseUnits - other.BaseUnits,
Unit:      q.Unit,
Factor:    q.Factor,
}, nil
}

func (q Quantity) String() string {
return fmt.Sprintf("%d (BaseUnits) [%s]", q.BaseUnits, q.Unit)
}
