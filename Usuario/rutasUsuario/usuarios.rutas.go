package rutasUsuario

import (
	"encoding/json"
	"errors"

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

	if StatusInternalServerError(w, err, "Error en la solicitud") {return} 
			
	prettyJSON, err := json.MarshalIndent(usuarios, "", "  ")
		
	if StatusInternalServerError(w, err, "Error interno en el servidor"){return}

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
	if StatusInternalServerError(w, err, "Se ha producido un error en el servidor") {return} 

	
	prettyJSON, err := json.MarshalIndent(usuario, "", "  ")

	if StatusInternalServerError(w, err, "Error al formatear los json"){return}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(prettyJSON)
	
}

func GetUsuariosByRolHandler(w http.ResponseWriter, r *http.Request) {

	var usuarios []Usuario
	parametros := mux.Vars(r)
	rol := parametros["rol"]

	err := baseDeDatos.DB.Where("rol = ?", rol).Find(&usuarios).Error

	if StatusInternalServerError(w, err, "Solicitud invalida") {return} 
	if len(usuarios) == 0 {w.WriteHeader(http.StatusNoContent)} 

	
	prettyJSON, err := json.MarshalIndent(usuarios, "", "  ")

	if StatusInternalServerError(w, err, "Error interno del servidor"){return}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(prettyJSON)
	
}

func EditarUsuario(w http.ResponseWriter, r *http.Request) {
	var usuario Usuario
	err := json.NewDecoder(r.Body).Decode(&usuario)

	params := mux.Vars(r)
	username := params["username"]

	if BadRequestError(w, err, "JSON inválido"){return}
	

	if NoExisteNingunCampo(usuario) {
		BadRequestError(w, errors.New(""), "No existe ningun campo")
		return
	}

	errors := verificarAtributos(usuario, SOFT)

	if len(errors) != 0 {
		BadRequestError(w, errors[0], "Atributos invalidos")
		return
	}

	err = baseDeDatos.DB.Model(&usuario).Where("username = ?", username).Updates(&usuario).Error

	if StatusInternalServerError(w, err, "Hubo un problema de actualizacion"){return}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Actualizacion exitosa!"))

}

func CrearUsuario(w http.ResponseWriter, r *http.Request) {

	var usuario Usuario
	err := json.NewDecoder(r.Body).Decode(&usuario)

	if BadRequestError(w, err,  "JSON inválido"){return}
	errors := verificarAtributos(usuario, HARD)

	
	for _, err := range errors {
		if err != nil {
			BadRequestError(w, errors[0], "se ha ingresado algun dato invalido")
			return
		}
	}
		
	usuario = DefinirUsername(usuario)
	usuario.Clave, err = Encriptar(usuario.Clave)

	if StatusInternalServerError(w, err, "error al encriptar la contraseña"){return}
	

	err = baseDeDatos.DB.Model(&usuario).Create(&usuario).Error

	if StatusInternalServerError(w, err, "error al crear el usuario"){return}

	prettyJSON, err := json.MarshalIndent(usuario, "", "  ")

	if StatusInternalServerError(w, err, "error al parsear el json"){return}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(prettyJSON))
}

func CrearUsuarios(w http.ResponseWriter, r *http.Request){

	var usuarios []Usuario

	if err := json.NewDecoder(r.Body).Decode(&usuarios); err != nil {
		http.Error(w, "Error al decodificar los usuarios: "+err.Error(), http.StatusBadRequest)
		return
	}

	for _, usuario := range usuarios {
		if err := verificarAtributos(usuario, HARD); err != nil {
			http.Error(w, "usuario inválido", http.StatusBadRequest)
			return
		}
	}

	tx := baseDeDatos.DB.Begin()
	for _, usuario := range usuarios {

		usuarioCreado := tx.Create(&usuario)

		err := usuarioCreado.Error
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error al crear los usuarios: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tx.Commit()
	w.WriteHeader(http.StatusOK)

}

func EliminarUsuario(w http.ResponseWriter, r *http.Request) {
	var usuario Usuario

	params := mux.Vars(r)
	username := params["username"]

	err := baseDeDatos.DB.Where("username = ?", username).Unscoped().Delete(&usuario).Error
	if StatusInternalServerError(w, err, "Hubo un problema de eliminacion") {return};

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Eliminacion exitosa!"))

}

func Loguearse(w http.ResponseWriter, r *http.Request) {

	var usuario Usuario
	var credencial Credencial
	var err error

	err = json.NewDecoder(r.Body).Decode(&credencial)

	
	if BadRequestError(w, err, "json invalido") {return};

	err = baseDeDatos.DB.Model(&usuario).Where("username = ?", credencial.Username).First(&usuario).Error

	if StatusNotFoundError(w, err, "usuario no encontrado") {return};
	
	err = Equals(credencial.Password, usuario.Clave)
	if StatusUnauthorizedError(w,err, "la contraseña es incorrecta"){return};
	
	prettyJSON, err := json.MarshalIndent(usuario, "", "  ")
	if StatusInternalServerError(w, err, "error al decodificar el usuario") {return};

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(prettyJSON)

}


