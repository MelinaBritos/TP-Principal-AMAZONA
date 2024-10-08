package rutasProveedor

import (
	"encoding/json"
	"net/http"

	"github.com/MelinaBritos/TP-Principal-AMAZONA/Proveedor/modelosProveedor"
	"github.com/MelinaBritos/TP-Principal-AMAZONA/baseDeDatos"
	"github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ola mundo"))
}

func GetProveedoresHandler(w http.ResponseWriter, r *http.Request) {
	//aca va la logica para obtener los proveedores
	var proveedores []modelosProveedor.Proveedor
	baseDeDatos.DB.Find(&proveedores)
	json.NewEncoder(w).Encode(&proveedores)

}

func GetProveedorHandler(w http.ResponseWriter, r *http.Request) {
	//aca va la logica para obtener un solo proveedor
	//w.Write([]byte("ola mundo proveedor"))
	var proveedor modelosProveedor.Proveedor
	params := mux.Vars(r)
	idProveedor := params["id_proveedor"]

	//baseDeDatos.DB.First(&proveedor, params["id_proveedor"])
	baseDeDatos.DB.Where("id_proveedor = ?", idProveedor).First(&proveedor)

	if proveedor.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("El proveedor no existe"))
		return
	}
	json.NewEncoder(w).Encode(&proveedor)
}

func PostProveedorHandler(w http.ResponseWriter, r *http.Request) {
	//aca va la logica para agregar un nuevo proveedor
	//w.Write([]byte("ola mundo post proveedor"))
	var proveedor modelosProveedor.Proveedor

	if err := json.NewDecoder(r.Body).Decode(&proveedor); err != nil {
		http.Error(w, "Error al decodificar el proveedor: "+err.Error(), http.StatusBadRequest)
		return
	}

	tx := baseDeDatos.DB.Begin()

	if err := tx.Create(&proveedor); err.Error != nil {
		tx.Rollback()
		http.Error(w, "Error al crear el proveedor: "+err.Error.Error(), http.StatusInternalServerError)
		return
	}

	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&proveedor)
}

func PutProveedorHandler(w http.ResponseWriter, r *http.Request) {

	//aca va la logica para modificar los datos de un proveedor
	w.Write([]byte("ola mundo put proveedor"))
}
