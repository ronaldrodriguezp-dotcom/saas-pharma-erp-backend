package inventory

import "sync"

// Inventory representa el stock de un producto
type Inventory struct {
	onHand   int
	reserved int
	mu       sync.Mutex
}

// NewInventory inicializa el inventario con un stock inicial
func NewInventory(initial int) *Inventory {
	return &Inventory{onHand: initial}
}

// TryReserve intenta reservar una cantidad del stock disponible
func (i *Inventory) TryReserve(qty int) bool {
	i.mu.Lock()
	defer i.mu.Unlock()
	if qty <= (i.onHand - i.reserved) {
		i.reserved += qty
		return true
	}
	return false
}

// Available devuelve el stock disponible para vender
func (i *Inventory) Available() int {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.onHand - i.reserved
}
