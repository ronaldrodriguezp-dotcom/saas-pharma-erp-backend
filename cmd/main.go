package main

import (
	"fmt"
	"fortiasaass-kio-speed/internal/inventory"
)

func main() {
	inv := inventory.NewInventory(5) // stock inicial = 5
	success := inv.TryReserve(3)
	fmt.Println("Reserva de 3 exitosa:", success)
	fmt.Println("Stock disponible:", inv.Available())
}
