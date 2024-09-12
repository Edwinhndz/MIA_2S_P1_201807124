// go run main.go
// npm run dev
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"PROYECTO1.com/mod/filesystem"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Comando struct {
	Comando string `json:"comando"`
}

type User struct {
	Carnet int    `json:"carnet"`
	Nombre string `json:"nombre"`
}

type Respuesta struct {
	ResponseBack string `json:"respuesta"`
	Error        bool   `json:"error"`
}

type allTasks []User

var tasks = allTasks{
	{
		Carnet: 201807124,
		Nombre: "Edwin Eduardo Lopez Hernandez",
	},
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func leerComando(w http.ResponseWriter, r *http.Request) {
	var newComando Comando
	var internalComando Comando
	var newRespuesta Respuesta

	// Leer el cuerpo de la solicitud
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		newRespuesta.ResponseBack = "Inserte un comando válido"
		newRespuesta.Error = true
		json.NewEncoder(w).Encode(newRespuesta)
		return
	}

	// Mostrar lo que se recibió
	fmt.Println("Cuerpo de la solicitud:", string(reqBody))

	// Desempaquetar el JSON externo
	err = json.Unmarshal(reqBody, &newComando)
	if err != nil {
		newRespuesta.ResponseBack = "Error en el formato del comando"
		newRespuesta.Error = true
		json.NewEncoder(w).Encode(newRespuesta)
		return
	}

	// Desempaquetar el JSON interno
	err = json.Unmarshal([]byte(newComando.Comando), &internalComando)
	if err != nil {
		newRespuesta.ResponseBack = "Error en el formato del comando interno"
		newRespuesta.Error = true
		json.NewEncoder(w).Encode(newRespuesta)
		return
	}

	// Procesar el comando
	newRespuesta.ResponseBack = filesystem.DividirComando(internalComando.Comando)
	fmt.Println("Respuesta del comando:", newRespuesta.ResponseBack)

	// Enviar la respuesta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newRespuesta)
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bienvenido a mi servidor")
}

func main() {
	// Rutas
	fmt.Println("MIA Edwin Lopez")
	router := mux.NewRouter().StrictSlash(true)
	// Endpoints
	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/Edwin", getTasks).Methods("GET")
	router.HandleFunc("/command", leerComando).Methods("POST")

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	// Servidor
	fmt.Println("Servidor corriendo en el puerto http://localhost:3000")
	http.ListenAndServe(":3000", handlers.CORS(headers, methods, origins)(router))
}
