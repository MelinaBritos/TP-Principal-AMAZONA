package validaciones

import (
	"errors"
	"fmt"

	"github.com/MelinaBritos/TP-Principal-AMAZONA/Proveedor/modelosProveedor"
	"github.com/MelinaBritos/TP-Principal-AMAZONA/baseDeDatos"
)

func ValidarRepuesto(repuesto modelosProveedor.Repuesto) error {

	if !existeCatalogo(repuesto.Id_catalogo) {
		return fmt.Errorf("no existe el catalogo %d", repuesto.Id_catalogo)
	}

	if err := validarNombre(repuesto.Nombre); err != nil {
		return err
	}

	if err := validarStock(repuesto.Stock); err != nil {
		return err
	}

	if err := validarStockMinimo(repuesto.Stock_minimo); err != nil {
		return err
	}

	if err := validarCantidadAComprar(repuesto.Cantidad_a_comprar); err != nil {
		return err
	}

	if err := validarCosto(repuesto.Costo); err != nil {
		return err
	}

	if err := validarDescripcion(repuesto.Descripcion); err != nil {
		return err
	}

	return nil

}

func existeCatalogo(id_catalogo int) bool {

	var catalogo modelosProveedor.Catalogo
	baseDeDatos.DB.Where("id = ?", id_catalogo).First(&catalogo)

	return catalogo.ID != 0
}

func validarStock(stock int) error {

	if stock < 0 {
		return errors.New("el stock no puede ser negativo")
	}

	return nil
}

func validarStockMinimo(stockMinimo int) error {

	if stockMinimo < 0 {
		return errors.New("el stock minimo no puede ser negativo")
	}

	return nil
}

func validarCantidadAComprar(cantidadAComprar int) error {

	if cantidadAComprar < 0 {
		return errors.New("la cantidad a comprar no puede ser negativa")
	}

	return nil
}

func validarCosto(costo float32) error {

	if costo <= 0 {
		return errors.New("el costo no puede ser negativo ni cero")
	}

	return nil
}

func validarDescripcion(descripcion string) error {

	if len(descripcion) > 100 {
		return errors.New("descripcion demasiado larga")
	}

	return nil
}
