package main

import (
	"encoding/json"
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type User struct {
	Name string `gorethink:"name"`
}

var session *r.Session

// Inicializamos la conexion a la DB
func init() {
	var err error
	session, err = r.Connect(r.ConnectOpts{
		Address:  "127.0.0.1:28015",
		Database: "GolangDB",
	})
	if err != nil {
		printStr(err)
		return
	} else {
		printStr("conexion establecida a la DB")
	}
}

// Metodo para hacer impresiones de pantalla
func printStr(v interface{}) {
	fmt.Println(v)
}

// Metodo para convertir a JSON
func printObj(v interface{}) []byte {
	vBytes, _ := json.Marshal(v)
	return vBytes
}

// Metodo para suscribirse a una tabla y escuchar los nuevos eventos
func Suscribe() {
	result, err := r.Table("users").Changes().Run(session)
	if err != nil {
		printStr(err)
	}
	printStr("*** Escuchando: ***")
	var rs interface{}
	for result.Next(&rs) {
		printStr("*** Nuevo ingreso: ***")
		data_JSON := printObj(rs)
		printStr(string(data_JSON))
	}
}

// Metodo para seleccionar todos los documentos de una tabla
func Select(table string) {
	result, err := r.Table(table).Run(session)
	if err != nil {
		printStr(err)
	}
	var users []interface{}
	err = result.All(&users)
	printStr("*** resultado: ***")
	if err != nil {
		fmt.Print("Error: %s", err)
	}
	for id := range users {
		printStr(users[id])
	}
}

func InsertUser(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		name := req.FormValue("name")
		user := &User{name}
		_, err := r.Table("users").Insert(user).RunWrite(session)
		if err != nil {
			fmt.Print("Error: %s", err)
		} else {
			io.WriteString(res, "Insertado con exito")
			Suscribe()
		}
	}
}

func main() {
	// Select("users")
	r := mux.NewRouter()
	r.HandleFunc("/Users", InsertUser)
	http.ListenAndServe(":5000", r)

}
