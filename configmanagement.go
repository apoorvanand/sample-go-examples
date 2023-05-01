package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type Config struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Store struct {
	Items map[string]string
}

func (s *Store) Get(key string) (string, bool) {
	value, exists := s.Items[key]
	return value, exists
}

func (s *Store) Set(key, value string) {
	s.Items[key] = value
}

func main() {
	store := Store{Items: make(map[string]string)}
	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			key := r.URL.Query().Get("key")
			value, exists := store.Get(key)
			if !exists {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			config := Config{Key: key, Value: value}
			json.NewEncoder(w).Encode(config)
		case http.MethodPost:
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			var config Config
			err = json.Unmarshal(body, &config)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			store.Set(config.Key, config.Value)
			w.WriteHeader(http.StatusCreated)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}
