package models

type Expiry struct {
	Month string
	Year  string
}

type Card struct {
	Name     string
	Postcode string
	PAN      string
	CVV      uint
	Expiry   Expiry
}
