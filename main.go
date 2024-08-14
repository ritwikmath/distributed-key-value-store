package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Payload struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type ResponseJson struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

var MyMap map[string]interface{} = make(map[string]interface{})

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Ok")
}

func PutKey(w http.ResponseWriter, r *http.Request) {
	var data Payload
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if MyMap == nil {
		http.Error(w, "Map is not initialized", http.StatusForbidden)
		return
	}
	MyMap[data.Key] = data.Value
	rsp := ResponseJson{
		Status:  200,
		Message: "Key stored",
		Data:    nil,
	}
	stringRsp, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, "Key is not present", http.StatusExpectationFailed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(stringRsp)
}

func GetKey(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	value, ok := MyMap[params["key"]]
	if !ok {
		http.Error(w, "Key is not present", http.StatusNotFound)
		return
	}
	rsp := ResponseJson{
		Status:  200,
		Message: "Key fetched",
		Data:    value,
	}
	stringRsp, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, "Key is not present", http.StatusExpectationFailed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(stringRsp)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", HealthCheck).Methods("GET")
	router.HandleFunc("/put", PutKey).Methods("POST")
	router.HandleFunc("/get/{key}", GetKey).Methods("GET")

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}

	log.Fatal(server.ListenAndServe())
}
