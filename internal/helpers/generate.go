//go:generate mockgen -package=helpers -source=generate.go -destination=./generate_mock.go Generator
package helpers

import "github.com/google/uuid"

type Generator interface {
	GenerateUUID() (uuid.UUID, error)
	AsString(input uuid.UUID) string
}

type Generate struct {
}

func (g *Generate) GenerateUUID() (uuid.UUID, error) {
	return uuid.NewUUID()
}

func (g *Generate) AsString(input uuid.UUID) string {
	return input.String()
}

func NewGenerator() *Generate {
	return &Generate{}
}
