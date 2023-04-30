package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var m sync.Map

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/keys", getAllKeys).Methods("GET")
	r.HandleFunc("/keys/{key}", getValueByKey).Methods("GET")
	r.HandleFunc("/keys/{key}", setValueByKey).Methods("PUT")
	r.HandleFunc("/keys/{key}", deleteValueByKey).Methods("DELETE")

	http.ListenAndServe(":8080", r)
}

func getAllKeys(w http.ResponseWriter, r *http.Request) {
	keys := []string{}
	m.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})
	json.NewEncoder(w).Encode(keys)
}

func getValueByKey(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	value, ok := m.Load(key)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(value)
}

func setValueByKey(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	var keyValue KeyValue
	err := json.NewDecoder(r.Body).Decode(&keyValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	m.Store(key, keyValue.Value)
	fmt.Fprintf(w, "Key %s set to %s", key, keyValue.Value)
}

func deleteValueByKey(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	m.Delete(key)
	fmt.Fprintf(w, "Key %s deleted", key)
}
