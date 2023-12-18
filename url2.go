package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// Item represents a simple data structure for our example
type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var (
	items      = make(map[string]Item)
	itemsMutex sync.RWMutex
)

// Handlers for CRUD operations

func getAllItems(w http.ResponseWriter, r *http.Request) {
	itemsMutex.RLock()
	defer itemsMutex.RUnlock()

	// Convert map of items to a slice
	var itemList []Item
	for _, item := range items {
		itemList = append(itemList, item)
	}

	// Convert the slice to JSON and write it as the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(itemList)
}

func getItem(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	itemsMutex.RLock()
	defer itemsMutex.RUnlock()

	item, found := items[id]
	if !found {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	var newItem Item
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	itemsMutex.Lock()
	defer itemsMutex.Unlock()

	items[newItem.ID] = newItem

	w.WriteHeader(http.StatusCreated)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	itemsMutex.Lock()
	defer itemsMutex.Unlock()

	if _, found := items[id]; !found {
		http.NotFound(w, r)
		return
	}

	delete(items, id)

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	// Define API routes and corresponding handlers
	http.HandleFunc("/items", getAllItems)
	http.HandleFunc("/item", getItem)
	http.HandleFunc("/item/create", createItem)
	http.HandleFunc("/item/delete", deleteItem)

	// Start the server on port 8080
	fmt.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
