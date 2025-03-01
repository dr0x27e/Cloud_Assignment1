package handlers

import (
	"log"
	"net/http"
)

func EmptyHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Request Path: ", r.URL.Path)
	log.Println("Request Method: ", r.Method)
	log.Println("Path Does Not Exist")
	w.WriteHeader(http.StatusNotFound)
}
