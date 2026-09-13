package domain

import (
"sort"
"time"
)

type StockLot struct {
ID             string
PresentationID string
LotNumber      string
ExpirationDate time.Time
ReceivedAt     time.Time
AvailableQty   int
}

// DischargeStrategy define el contrato para ordenar los lotes según la política del negocio.
type DischargeStrategy interface {
SortLots(lots []StockLot) []StockLot
}

// FefoStrategy: Prioriza los lotes más próximos a vencer (First Expired, First Out).
// Esencial para farmacias, clínicas y productos regulados sanitariamente.
type FefoStrategy struct{}

func (s *FefoStrategy) SortLots(lots []StockLot) []StockLot {
sort.Slice(lots, func(i, j int) bool {
// Si las fechas de expiración son iguales, desempatamos por fecha de recepción (FIFO secundario)
if lots[i].ExpirationDate.Equal(lots[j].ExpirationDate) {
return lots[i].ReceivedAt.Before(lots[j].ReceivedAt)
}
return lots[i].ExpirationDate.Before(lots[j].ExpirationDate)
})
return lots
}

// FifoStrategy: Prioriza los lotes que ingresaron primero (First In, First Out).
// Utilizado para retail general, ferreterías y bienes no perecederos.
type FifoStrategy struct{}

func (s *FifoStrategy) SortLots(lots []StockLot) []StockLot {
sort.Slice(lots, func(i, j int) bool {
return lots[i].ReceivedAt.Before(lots[j].ReceivedAt)
})
return lots
}