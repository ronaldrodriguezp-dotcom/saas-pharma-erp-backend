package domain

import "errors"

// ErrInsufficientStock se lanza cuando una caja intenta reservar más stock del disponible
var ErrInsufficientStock = errors.New("stock insuficiente para la reserva")
