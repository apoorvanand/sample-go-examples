package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "sync"
)

type keyValuePair struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}

type keyValueStore struct {
    sync.RWMutex
    items       []*keyValuePair
    lookupTable map[string]int
}

func (kvs *keyValueStore) get(key string) (string, bool) {
    kvs.RLock()
    defer kvs.RUnlock()

    index, ok := kvs.lookupTable[key]
    if !ok {
        return "", false
    }
    return kvs.items[index].Value, true
}

func (kvs *keyValueStore) set(key string, value string) {
    kvs.Lock()
    defer kvs.Unlock()

    index, ok := kvs.lookupTable[key]
    if ok {
        kvs.items[index].Value = value
        return
    }

    kvs.lookupTable[key] = len(kvs.items)
    kvs.items = append(kvs.items, &keyValuePair{
        Key:   key,
        Value: value,
    })
}

func (kvs *keyValueStore) delete(key string) {
    kvs.Lock()
    defer kvs.Unlock()

    index, ok := kvs.lookupTable[key]
    if !ok {
        return
    }

    // Delete the item at the index and shift the remaining items to the left.
    kvs.items = append(kvs.items[:index], kvs.items[index+1:]...)
    delete(kvs.lookupTable, key)

    // Decrement the indexes in the lookup table for items that come after the deleted index.
    for k, v := range kvs.lookupTable {
        if v > index {
            kvs.lookupTable[k] = v - 1
        }
    }
}

func (kvs *keyValueStore) itemsInInsertionOrder() []*keyValuePair {
    kvs.RLock()
    defer kvs.RUnlock()

    return kvs.items
}

var store = keyValueStore{
    items:       make([]*keyValuePair, 0),
    lookupTable: make(map[string]int),
}

func handleGet(w http.ResponseWriter, r *http.Request) {
    key := r.URL.Query().Get("key")
    value, ok := store.get(key)
    if !ok {
        http.NotFound(w, r)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(keyValuePair{
        Key:   key,
        Value: value,
    })
}

func handlePost(w http.ResponseWriter, r *http.Request) {
    var kvp keyValuePair
    if err := json.NewDecoder(r.Body).Decode(&kvp); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    store.set(kvp.Key, kvp.Value)
    w.WriteHeader(http.StatusCreated)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
    key := r.URL.Query().Get("key")
    store.delete(key)
    w.WriteHeader(http.StatusNoContent)
}

func handleList(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(store.itemsInInsertionOrder())
}

func main() {
    //sync 
    http.HandleFunc("/get", handleGet)
    http.HandleFunc("/set", handlePost)
    http.HandleFunc("/delete", handleDelete)
    http.HandleFunc("/list", handleList)

   // Create a new HTTP server and register the handler function
	server := http.Server{
		Addr:    ":8080",
		
	}

	// Start the server
	fmt.Println("Starting server...")
	err := server.ListenAndServe()
	if err != nil {
        log.Fatal(err)
		panic(err)
	}
    }