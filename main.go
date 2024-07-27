package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Usuario struct {
	Id_usuario int
	Nombre     string
	Email      string
}

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {
	http.HandleFunc("/", Inicio)
	http.HandleFunc("/crear", PlantillaCrear)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/eliminar", Eliminar)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/actualizar", Actualizar)

	log.Println("Servidor corriendo... localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func Actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id_usuario := r.FormValue("id")
		nombre := r.FormValue("nombre")
		email := r.FormValue("email")

		conexion := ConexionDB()

		actualizar, err := conexion.Prepare(" UPDATE usuarios SET nombre=?, email=? WHERE=id_usuario=?")
		if err != nil {
			panic(err.Error())
		}
		actualizar.Exec(nombre, email, id_usuario)
		http.Redirect(w, r, "/", r.Response.StatusCode)
	}
}
func Editar(w http.ResponseWriter, r *http.Request) {
	id_usuario := r.URL.Query().Get("id_usuario")
	conexion := ConexionDB()
	usuario := Usuario{}

	registro, err := conexion.Query(" SELECT * FROM usuarios WHERE id_usuario=?", id_usuario)
	if err != nil {
		panic(err.Error())
	}
	for registro.Next() {
		var nombre, email string
		err = registro.Scan(&id_usuario, &nombre, &email)
		if err != nil {
			panic(err.Error())
		}
		usuario.Nombre = nombre
		usuario.Email = email
	}
	plantillas.ExecuteTemplate(w, "/editar", usuario)

}

func Eliminar(w http.ResponseWriter, r *http.Request) {
	id_usuario := r.URL.Query().Get("id_usuario")

	conexion := ConexionDB()

	eliminar, err := conexion.Prepare("DELETE FROM usuarios WHERE id_usuario=?)")
	if err != nil {
		panic(err.Error())
	}
	eliminar.Exec(id_usuario)
	http.Redirect(w, r, "/", r.Response.StatusCode)
}

func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		email := r.FormValue("email")

		conexion := ConexionDB()

		registar, err := conexion.Prepare("INSERT INTO usuarios(nombre,email) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		registar.Exec(nombre, email)
		http.Redirect(w, r, "/", r.Response.StatusCode)
	}
}

func Inicio(w http.ResponseWriter, r *http.Request) {
	conexion := ConexionDB()
	usuario := Usuario{}
	arregloUsuario := []Usuario{}

	registar, err := conexion.Query(" SELECT * FROM usuarios")
	if err != nil {
		panic(err.Error())
	}
	for registar.Next() {
		var id_usuario int
		var nombre, email string
		err = registar.Scan(&id_usuario, &nombre, &email)
		if err != nil {
			panic(err.Error())
		}
		usuario.Nombre = nombre
		usuario.Email = email
		arregloUsuario = append(arregloUsuario, usuario)

	}
	println(arregloUsuario)
	plantillas.ExecuteTemplate(w, "inicio", arregloUsuario)
}

func PlantillaCrear(w http.ResponseWriter, r *http.Request) {
	plantillas = template.Must(template.ParseGlob("plantillas/*"))
	plantillas.ExecuteTemplate(w, "crear", nil)
}

func ConexionDB() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contraseña := ""
	Nombre := "aplicacion"

	conexion, err := sql.Open(Driver, Usuario+":"+Contraseña+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}
