package main

import (
	"encoding/json"
	"gestanestle/aizm-server/internal/db"
	"gestanestle/aizm-server/internal/mqttc"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	d := db.NewCon()
	defer d.Close()

	mqttc.Subscribe()

	r := mux.NewRouter()
	r.HandleFunc("/", Hello).Methods("GET")
	http.ListenAndServe(":4000", r)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	res := Response{
		Status:  "OK",
		Message: "Hello from AIZM Server",
	}
	json.NewEncoder(w).Encode(res)
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
