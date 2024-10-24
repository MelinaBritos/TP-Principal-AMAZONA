package modelosBitacora

import (
	"gorm.io/gorm"
)

type Ticket struct {
	gorm.Model

	Username              string `gorm:"not null"`
	MotivoIngreso         string `gorm:"not null"`
	FechaCreacion         string `gorm:"not null"`
	Estado                string `gorm:"not null"`
	Tipo                  string `gorm:"not null"`
	Matricula             string `gorm:"not null"`
	CostoTotal            float32
	DescripcionReparacion string
	FechaFinalizacion     string
	Repuestos             []RepuestoUtilizado `gorm:"foreignKey:IDTicket;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
