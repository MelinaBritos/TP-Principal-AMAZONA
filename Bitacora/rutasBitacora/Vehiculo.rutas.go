package rutasBitacora

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"time"

	"github.com/MelinaBritos/TP-Principal-AMAZONA/Bitacora/modelosBitacora"
	"github.com/MelinaBritos/TP-Principal-AMAZONA/baseDeDatos"
	"github.com/gorilla/mux"
)

func GetVehiculosHandler(w http.ResponseWriter, r *http.Request) {
	var vehiculos []modelosBitacora.Vehiculo

	baseDeDatos.DB.Find(&vehiculos)
	json.NewEncoder(w).Encode(&vehiculos)
}

func GetVehiculoHandler(w http.ResponseWriter, r *http.Request) {
	var vehiculo modelosBitacora.Vehiculo
	parametros := mux.Vars(r)

	baseDeDatos.DB.First(&vehiculo, parametros["id"])

	if vehiculo.ID == 0 {
		w.WriteHeader(http.StatusNotFound) // error 404
		w.Write([]byte("Vehiculo no encontrado"))
		return
	}

	json.NewEncoder(w).Encode(&vehiculo)

}

func GetVehiculosDisponiblesHandler(w http.ResponseWriter, r *http.Request) {
	var vehiculos []modelosBitacora.Vehiculo

	baseDeDatos.DB.Find(&vehiculos)

	var vehiculosDisponibles []modelosBitacora.Vehiculo
	for _, vehiculo := range vehiculos {
		if vehiculo.Estado != "NO APTO PARA CIRCULAR" && vehiculo.Estado != "REPARACIÓN" && vehiculo.Estado != "MANTENIMIENTO" {
			vehiculosDisponibles = append(vehiculosDisponibles, vehiculo)
		}
	}

	json.NewEncoder(w).Encode(&vehiculosDisponibles)
}

func GetMarcasVehiculoHandler(w http.ResponseWriter, r *http.Request) {
	marcas := []string{"Fiat", "Renault", "Peugeot", "Citroën", "Volkswagen", "Ford", "Nissan", "Toyota", "Mercedes-Benz"}
	json.NewEncoder(w).Encode(&marcas)
}

func GetModelosVehiculoHandler(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	switch parametros["marca"] {
	case "Fiat":
		fiat := []string{"Fiat 500", "Fiat Argo", "Fiat Cronos", "Fiat Toro", "Fiat Strada"}
		json.NewEncoder(w).Encode(&fiat)
	case "Renault":
		renault := []string{"Renault Clio", "Renault Sandero", "Renault Kwid", "Renault Duster", "Renault Captur", "Renault Kangoo"}
		json.NewEncoder(w).Encode(&renault)
	case "Peugeot":
		peugeot := []string{"Peugeot 208", "Peugeot 2008", "Peugeot 3008", "Peugeot 308", "Peugeot 408"}
		json.NewEncoder(w).Encode(&peugeot)
	case "Citroën":
		citroen := []string{"Citroën C3", "Citroën C4 Cactus", "Citroën C5 Aircross", "Citroën Berlingo"}
		json.NewEncoder(w).Encode(&citroen)
	case "Volkswagen":
		volkswagen := []string{"Volkswagen Golf", "Volkswagen Polo", "Volkswagen T-Cross", "Volkswagen Tiguan", "Volkswagen Amarok", "Volkswagen Vento"}
		json.NewEncoder(w).Encode(&volkswagen)
	case "Ford":
		ford := []string{"Ford Fiesta", "Ford Focus", "Ford Ranger", "Ford Ecosport", "Ford Kuga", "Ford Mustang"}
		json.NewEncoder(w).Encode(&ford)
	case "Nissan":
		nissan := []string{"Nissan March", "Nissan Versa", "Nissan Sentra", "Nissan X-Trail", "Nissan Kicks", "Nissan Frontier"}
		json.NewEncoder(w).Encode(&nissan)
	case "Toyota":
		toyota := []string{"Toyota Corolla", "Toyota Yaris", "Toyota Hilux", "Toyota RAV4", "Toyota SW4", "Toyota Camry"}
		json.NewEncoder(w).Encode(&toyota)
	case "Mercedes-Benz":
		mercedesBenz := []string{"Mercedes-Benz Clase A", "Mercedes-Benz Clase C", "Mercedes-Benz Clase E", "Mercedes-Benz GLC", "Mercedes-Benz GLE", "Mercedes-Benz Sprinter"}
		json.NewEncoder(w).Encode(&mercedesBenz)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Marca invalida"))
		return
	}
}

func PostVehiculoHandler(w http.ResponseWriter, r *http.Request) {
	var vehiculos []modelosBitacora.Vehiculo

	if err := json.NewDecoder(r.Body).Decode(&vehiculos); err != nil {
		http.Error(w, "Error al decodificar los vehiculos: "+err.Error(), http.StatusBadRequest)
		return
	}

	for _, vehiculo := range vehiculos {
		if err := validarVehiculo(vehiculo); err != nil {
			http.Error(w, "Vehiculo inválido: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	tx := baseDeDatos.DB.Begin()
	for _, vehiculo := range vehiculos {
		vehiculo.FechaIngreso = time.Now()
		vehiculoCreado := tx.Create(&vehiculo)

		err := vehiculoCreado.Error
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error al crear los vehiculos: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tx.Commit()
	w.WriteHeader(http.StatusOK)
}

func PutVehiculoHandler(w http.ResponseWriter, r *http.Request) {
	var vehiculo modelosBitacora.Vehiculo
	var vehiculoInput modelosBitacora.Vehiculo

	if err := json.NewDecoder(r.Body).Decode(&vehiculoInput); err != nil {
		http.Error(w, "Error al decodificar el vehículo: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := validarVehiculo(vehiculoInput); err != nil {
		http.Error(w, "Vehiculo inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	tx := baseDeDatos.DB.Begin()
	if err := tx.First(&vehiculo, "matricula = ?", vehiculoInput.Matricula).Error; err != nil {
		http.Error(w, "Vehículo no encontrado: "+err.Error(), http.StatusNotFound)
		return
	}

	if err := tx.Model(&vehiculo).Omit("FechaIngreso").Updates(vehiculoInput).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al actualizar el vehículo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tx.Commit()
	w.WriteHeader(http.StatusOK)
}

func DeleteVehiculoHandler(w http.ResponseWriter, r *http.Request) {
	var vehiculo modelosBitacora.Vehiculo
	parametros := mux.Vars(r)

	baseDeDatos.DB.First(&vehiculo, parametros["id"])

	if vehiculo.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Vehiculo no encontrado"))
		return
	}

	baseDeDatos.DB.Unscoped().Delete(&vehiculo)
	w.WriteHeader(http.StatusOK)

}

func validarVehiculo(vehiculo modelosBitacora.Vehiculo) error {

	switch vehiculo.Estado {
	case "NO APTO PARA CIRCULAR", "APTO PARA CIRCULAR":
	default:
		return errors.New("estado invalido")
	}
	switch vehiculo.EstadoVTV {
	case "VENCIDA", "APROBADA", "RECHAZADA":
	default:
		return errors.New("estado de VTV invalido")
	}

	// validar fechaVTV, no vencida
	if vehiculo.FechaVTV.IsZero() {
		return errors.New("la fecha no puede ser nula o vacía")
	}
	Ahora := time.Now()
	UnAñoAtras := Ahora.AddDate(-1, 0, 0)
	if vehiculo.FechaVTV.Before(UnAñoAtras) || vehiculo.FechaVTV.After(Ahora) {
		if vehiculo.Estado == "APROBADA" {
			return errors.New("la vtv esta vencida")
		}
	}

	if vehiculo.Estado == "APTO PARA CIRCULAR" && (vehiculo.EstadoVTV == "RECHAZADA" || vehiculo.EstadoVTV == "VENCIDA") {
		return errors.New("el vehiculo no es apto para circular")
	}

	switch vehiculo.Marca {
	case "Fiat", "Renault", "Peugeot", "Citroën", "Volkswagen", "Ford", "Nissan", "Toyota", "Mercedes-Benz":
	default:
		return errors.New("marca invalida")
	}
	if vehiculo.PesoAdmitido <= 0 {
		return errors.New("peso admitido invalido")
	}
	if vehiculo.VolumenAdmitidoMtsCubicos <= 0 {
		return errors.New("volumen admitido invalido")
	}
	if vehiculo.Año > 2024 || vehiculo.Año < 2000 {
		return errors.New("año invalido")
	}

	matriculaVieja := regexp.MustCompile(`^[A-Z]{3}[0-9]{3}$`)
	matriculaNueva := regexp.MustCompile(`^[A-Z]{2}[0-9]{3}[A-Z]{2}$`)

	if !matriculaVieja.MatchString(vehiculo.Matricula) && !matriculaNueva.MatchString(vehiculo.Matricula) {
		return errors.New("matricula invalida")
	}

	return nil
}
