package modelosPaquete

import (
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type Estado string

var estadosValidos = []Estado{
	SIN_ASIGNAR,
	ASIGNADO,
	EN_VIAJE,
	ENTREGADO,
	NO_ENTREGADO,
}

const (
	SIN_ASIGNAR  Estado = "SIN ASIGNAR"
	ASIGNADO     Estado = "ASIGNADO"
	EN_VIAJE     Estado = "EN VIAJE"
	ENTREGADO    Estado = "ENTREGADO"
	NO_ENTREGADO Estado = "NO ENTREGADO"
)

type Paquete struct {
	gorm.Model

	Id_viaje           int
	Estado             Estado  `gorm:"not null"`
	Peso_kg            float32 `gorm:"not null"`
	Nombre_cliente     string  `gorm:"not null"`
	Tamaño_mts_cubicos float32 `gorm:"not null"`
	Localidad          string  `gorm:"not null"`
	Dir_entrega        string  `gorm:"not null"`
}

func ObtenerEstadosValidos() []Estado {
	return estadosValidos
}

func (p *Paquete) GetIDAsString() string {
	return strconv.Itoa(int(p.ID))
}

func ParseEstado(s string) (Estado, error) {
	switch s {
	case "SIN ASIGNAR":
		return SIN_ASIGNAR, nil
	case "ASIGNADO":
		return ASIGNADO, nil
	case "EN VIAJE":
		return EN_VIAJE, nil
	case "ENTREGADO":
		return ENTREGADO, nil
	case "NO ENTREGADO":
		return NO_ENTREGADO, nil
	default:
		return SIN_ASIGNAR, fmt.Errorf("estado no válido: %s", s)
	}
}
