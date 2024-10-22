package rutasUsuario

import (
	"encoding/json"
	"errors"
	"fmt"

	"net/http"

	"github.com/MelinaBritos/TP-Principal-AMAZONA/Usuario/modelosUsuario"
	"github.com/MelinaBritos/TP-Principal-AMAZONA/baseDeDatos"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Usuario = modelosUsuario.Usuario
type Credencial = modelosUsuario.Credencial

func GetUsuariosHandler(w http.ResponseWriter, r *http.Request) {

	var usuarios []Usuario
	err := baseDeDatos.DB.Find(&usuarios).Error

	if StatusInternalServerError(w, err, "Error en la solicitud") {
		return
	}

	prettyJSON, err := json.MarshalIndent(usuarios, "", "  ")

	if StatusInternalServerError(w, err, "Error interno en el servidor") {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(prettyJSON)

}

func GetUsuarioByUsernameHandler(w http.ResponseWriter, r *http.Request) {

	var usuario Usuario
	parametros := mux.Vars(r)
	username := parametros["username"]

	err := baseDeDatos.DB.Where("username = ?", username).First(&usuario).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		StatusNotFoundError(w, err, "Usuario no encontrado")
		return
	}
	if StatusInternalServerError(w, err, "Se ha producido un error en el servidor") {
		return
	}

	prettyJSON, err := json.MarshalIndent(usuario, "", "  ")

	if StatusInternalServerError(w, err, "Error al formatear los json") {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(prettyJSON)

}

func GetUsuariosByRolHandler(w http.ResponseWriter, r *http.Request) {

	var usuarios []Usuario
	parametros := mux.Vars(r)
	rol := parametros["rol"]

	err := baseDeDatos.DB.Where("rol = ?", rol).Find(&usuarios).Error

	if StatusInternalServerError(w, err, "Solicitud invalida") {
		return
	}
	if len(usuarios) == 0 {
		w.WriteHeader(http.StatusNoContent)
	}

	prettyJSON, err := json.MarshalIndent(usuarios, "", "  ")

	if StatusInternalServerError(w, err, "Error interno del servidor") {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(prettyJSON)

}

func EditarUsuario(w http.ResponseWriter, r *http.Request) {

	var usuario Usuario

	err := json.NewDecoder(r.Body).Decode(&usuario)

	params := mux.Vars(r)
	username := params["username"]

	if BadRequestError(w, err, "JSON inválido") {
		return
	}

	if NoExisteNingunCampo(usuario) {
		BadRequestError(w, errors.New(""), "No existe ningun campo")
		return
	}

	errors := verificarAtributos(usuario, SOFT)

	if len(errors) != 0 {
		BadRequestError(w, errors[0], "Atributos invalidos")
		return
	}

	if usuario.Clave != "" {
		usuario.Clave, err = Encriptar(usuario.Clave)
		if StatusInternalServerError(w, err, "error al encriptar la nueva clave") {
			return
		}
	}

	if usuario.Nombre != "" || usuario.Apellido != "" || usuario.Dni != "" {

		var usuarioActual Usuario
		err := baseDeDatos.DB.Where("username = ?", username).First(&usuarioActual).Error
		if StatusInternalServerError(w, err, "error al obtener al usuario") {return}

		if usuario.Nombre != "" {

			usuario = DefinirUsuarioSegunApellido(usuario, usuarioActual)
			
			
		} else {

			usuario = DefinirUsuarioSegunNombreVacio(usuario, usuarioActual)
		}
		usuario = DefinirUsername(usuario)
	}

	err = baseDeDatos.DB.Model(&usuario).Where("username = ?", username).Updates(&usuario).Error

	if StatusInternalServerError(w, err, "Hubo un problema de actualizacion") {return}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Actualizacion exitosa!"))

}



func EditarUsuarios(w http.ResponseWriter, r *http.Request) {

	var usuarios []Usuario

	err := json.NewDecoder(r.Body).Decode(&usuarios)
	if BadRequestError(w, err, "Error al decodificar los usuarios: ") {
		return
	}

	for _, usuario := range usuarios {

		if usuario.Username == "" {
			BadRequestError(w, errors.New(""), "Hay un usuario sin username")
			return
		}
		errors := verificarAtributos(usuario, SOFT)
		if len(errors) != 0 {
			BadRequestError(w, errors[0], "Alguna informacion del usuario es incorrecta")
		}

	}

	tx := baseDeDatos.DB.Begin()
	for _, usuario := range usuarios {

		var actualUsername string

		if usuario.Clave != "" {
			usuario.Clave, err = Encriptar(usuario.Clave)
			if StatusInternalServerError(w, err, "error al encriptar la nueva clave") {
				tx.Rollback()
				return
			}
		}

		if usuario.Nombre != "" || usuario.Apellido != "" || usuario.Dni != "" {

			
			var usuarioActual Usuario
			err := baseDeDatos.DB.Where("username = ?", usuario.Username).First(&usuarioActual).Error

			fmt.Println(usuario.Nombre)
			if StatusInternalServerError(w, err, "error al obtener al usuario") {
				tx.Rollback()
				return
			}
			actualUsername = usuarioActual.Username

			if usuario.Nombre != "" {
				usuario = DefinirUsuarioSegunApellido(usuario, usuarioActual)
				
			} else {
				usuario = DefinirUsuarioSegunNombreVacio(usuario, usuarioActual)
			}
			usuario = DefinirUsername(usuario)
			
		} else{
			actualUsername = usuario.Username
		}

		err = baseDeDatos.DB.Model(&usuario).Where("username = ?", actualUsername).Updates(&usuario).Error

		if StatusInternalServerError(w, err, "Hubo un problema de actualizacion") {
			tx.Rollback()
			return
		}

	}

	tx.Commit()
	w.WriteHeader(http.StatusOK)
}

func CrearUsuario(w http.ResponseWriter, r *http.Request) {

	var usuario Usuario
	err := json.NewDecoder(r.Body).Decode(&usuario)

	if BadRequestError(w, err, "JSON inválido") {
		return
	}
	errors := verificarAtributos(usuario, HARD)

	for _, err := range errors {
		if err != nil {
			BadRequestError(w, errors[0], "se ha ingresado algun dato invalido")
			return
		}
	}

	usuario = DefinirUsername(usuario)
	usuario.Clave, err = Encriptar(usuario.Clave)

	if StatusInternalServerError(w, err, "error al encriptar la contraseña") {
		return
	}

	err = baseDeDatos.DB.Model(&usuario).Create(&usuario).Error

	if StatusInternalServerError(w, err, "error al crear el usuario") {
		return
	}

	prettyJSON, err := json.MarshalIndent(usuario, "", "  ")

	if StatusInternalServerError(w, err, "error al parsear el json") {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(prettyJSON))
}

func CrearUsuarios(w http.ResponseWriter, r *http.Request) {

	var usuarios []Usuario

	err := json.NewDecoder(r.Body).Decode(&usuarios)
	if BadRequestError(w, err, "Error al decodificar los usuarios: ") {
		return
	}

	for _, usuario := range usuarios {

		errors := verificarAtributos(usuario, HARD)
		if len(errors) != 0 {
			BadRequestError(w, errors[0], "Alguna informacion del usuario es incorrecta")
		}

	}

	tx := baseDeDatos.DB.Begin()
	for _, usuario := range usuarios {

		usuario = DefinirUsername(usuario)
		usuario.Clave, err = Encriptar(usuario.Clave)

		StatusInternalServerError(w, err, "no se pudo encriptar la contraseña")

		usuarioCreado := tx.Create(&usuario)

		err := usuarioCreado.Error

		if err != nil {
			tx.Rollback()
			StatusInternalServerError(w, err, "error al insertar usuario")
			return
		}
	}

	tx.Commit()
	w.WriteHeader(http.StatusCreated)

}

func EliminarUsuario(w http.ResponseWriter, r *http.Request) {
	var usuario Usuario

	params := mux.Vars(r)
	username := params["username"]

	err := baseDeDatos.DB.Where("username = ?", username).Unscoped().Delete(&usuario).Error
	if StatusInternalServerError(w, err, "Hubo un problema de eliminacion") {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Eliminacion exitosa!"))

}

func EliminarUsuarios(w http.ResponseWriter, r *http.Request) {

	var credenciales []Credencial

	err := json.NewDecoder(r.Body).Decode(&credenciales)
	if BadRequestError(w, err, "Error al decodificar las credenciales de los usuarios: ") {
		return
	}

	tx := baseDeDatos.DB.Begin()

	for _, credencial := range credenciales {

		if StatusInternalServerError(w, err, "error en la codificacion de la contraseña") {
			return
		}

		err = baseDeDatos.DB.Where(&Usuario{Username: credencial.Username}).Unscoped().Delete(Usuario{}).Error

		if StatusInternalServerError(w, err, "error en la eliminacion de algun usuario") {
			tx.Rollback()
			return
		}
	}

	tx.Commit()
	w.WriteHeader(http.StatusOK)

}

func Loguearse(w http.ResponseWriter, r *http.Request) {

	var usuario Usuario
	var credencial Credencial
	var err error

	err = json.NewDecoder(r.Body).Decode(&credencial)

	if BadRequestError(w, err, "json invalido") {
		return
	}

	err = baseDeDatos.DB.Model(&usuario).Where("username = ?", credencial.Username).First(&usuario).Error

	if StatusNotFoundError(w, err, "usuario no encontrado") {
		return
	}

	err = Equals(credencial.Password, usuario.Clave)
	if StatusUnauthorizedError(w, err, "la contraseña es incorrecta") {
		return
	}

	prettyJSON, err := json.MarshalIndent(usuario, "", "  ")
	if StatusInternalServerError(w, err, "error al decodificar el usuario") {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(prettyJSON)

}
