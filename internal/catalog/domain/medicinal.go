package domain

import (
"errors"
"time"
)

// PharmaceuticalForm Catálogo de Formas Farmacéuticas (ej. Comprimido, Jarabe, Inyectable)
type PharmaceuticalForm struct {
FormID      string
Name        string
Description string
}

func NewPharmaceuticalForm(formID, name string) (PharmaceuticalForm, error) {
if formID == "" || name == "" {
return PharmaceuticalForm{}, errors.New("el ID y nombre de la forma farmacéutica son obligatorios")
}
return PharmaceuticalForm{
FormID: formID,
Name:   name,
}, nil
}

// DispensingCondition Condición de Venta y Dispensación (versión con vigencia)
type DispensingConditionType string

const (
CondDirectSale      DispensingConditionType = "DIRECT_SALE"       // Venta directa / OTC
CondMedicalPrescr   DispensingConditionType = "MEDICAL_PRESCRIPTION" // Receta médica simple
CondRetainedPrescr  DispensingConditionType = "RETAINED_PRESCRIPTION" // Receta retenida (psicotrópicos/estupefacientes)
CondChequePrescr    DispensingConditionType = "CHEQUE_PRESCRIPTION"  // Receta cheque
)

type DispensingCondition struct {
Type      DispensingConditionType
ValidFrom time.Time
ValidTo   time.Time
}

func NewDispensingCondition(cType DispensingConditionType, validFrom, validTo time.Time) (DispensingCondition, error) {
if validTo.Before(validFrom) {
return DispensingCondition{}, errors.New("la fecha de término de vigencia no puede ser anterior al inicio")
}
return DispensingCondition{
Type:      cType,
ValidFrom: validFrom,
ValidTo:   validTo,
}, nil
}

// HealthRegistration Entidad Versionable del Registro Sanitario (CAT-006)
type HealthRegistration struct {
RegistrationNumber string    // Ej: F-12345 / ISP-C-6789
Holder             string    // Titular del Registro Sanitario (MAH)
Status             string    // Vigente, Suspendido, Caducado
IssuedDate         time.Time
ExpiryDate         time.Time
PreviousVersion    *HealthRegistration // Historial para cumplir CAT-006
}

// MedicinalProductProfile Perfil regulatorio exclusivo para productos de tipo MEDICINAL (CAT-004)
type MedicinalProductProfile struct {
TenantID             string
ProductID            string // Vinculación estricta al Product padre de tipo MEDICINAL
FormID               string // Referencia a PharmaceuticalForm
HealthRegistration   HealthRegistration
DispensingCondition  DispensingCondition
IsControlled         bool // Estupefaciente o Psicotrópico sujeto a ley de control de drogas
}

func NewMedicinalProductProfile(
tenantID string,
productID string,
formID string,
registration HealthRegistration,
dispensing DispensingCondition,
isControlled bool,
) (MedicinalProductProfile, error) {
if tenantID == "" || productID == "" || formID == "" {
return MedicinalProductProfile{}, errors.New("parámetros obligatorios vacíos para crear el perfil medicinal")
}

return MedicinalProductProfile{
TenantID:            tenantID,
ProductID:           productID,
FormID:              formID,
HealthRegistration:  registration,
DispensingCondition: dispensing,
IsControlled:        isControlled,
}, nil
}