package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", loginHandle).Methods("POST")
	http.ListenAndServe(":9090", r)

}

func loginHandle(w http.ResponseWriter, r *http.Request) {
	log.Print("Login Successful")

}
