package domain

import (
"errors"
"regexp"
)

var (
ErrInvalidEmail = errors.New("formato de email inválido")
ErrInvalidPhone = errors.New("formato de teléfono inválido")
ErrEmptyAddress = errors.New("los campos obligatorios de la dirección no pueden estar vacíos")
)

type Contact struct {
Name  string
Email string
Phone string
}

func NewContact(name, email, phone string) (Contact, error) {
if email != "" {
re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
if !re.MatchString(email) {
return Contact{}, ErrInvalidEmail
}
}
return Contact{Name: name, Email: email, Phone: phone}, nil
}

type Address struct {
Street     string
City       string
Region     string
Country    string
PostalCode string
}

func NewAddress(street, city, region, country string) (Address, error) {
if street == "" || city == "" || country == "" {
return Address{}, ErrEmptyAddress
}
return Address{
Street:  street,
City:    city,
Region:  region,
Country: country,
}, nil
}