package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func handleRedashNotify(w http.ResponseWriter, r *http.Request) {
	// TODO implement
}

func handleRequests() {
	router := mux.NewRouter().StrictSlack(true)
	router.HandleFunc("/", handleRedashNotify)
}

func main() {
	handleRequests()
}
