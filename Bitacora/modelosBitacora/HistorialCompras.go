package modelosBitacora

import (
	"github.com/MelinaBritos/TP-Principal-AMAZONA/Proveedor/modelosProveedor"
	"gorm.io/gorm"
)

type HistorialCompras struct {
	gorm.Model

	RepuestoCompradoID int
	RepuestoComprado   modelosProveedor.Repuesto
	Cantidad           int
	Costo              float32
}