package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
	Name     string
	Email    string
	Tier     string
	Endpoint string
}

var users = []User{
	{Name: "Alice", Email: "alice@example.com", Tier: "premium", Endpoint: "/premium"},
	{Name: "Bob", Email: "bob@example.com", Tier: "basic", Endpoint: "/basic"},
	{Name: "Charlie", Email: "charlie@example.com", Tier: "free", Endpoint: "/free"},
}

func main() {
	// Create a new router
	router := mux.NewRouter()

	// Define handlers for each tier
	router.HandleFunc("/premium", PremiumHandler)
	router.HandleFunc("/basic", BasicHandler)
	router.HandleFunc("/free", FreeHandler)

	// Add middleware to restrict access based on user tier
	router.Use(AuthMiddleware)

	// Start the server on port 8080
	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", router)
}

func PremiumHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Premium endpoint!")
}

func BasicHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Basic endpoint!")
}

func FreeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Free endpoint!")
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the user is authorized to access the endpoint
		userTier := r.Header.Get("X-User-Tier")
		for _, user := range users {
			if user.Email == r.Header.Get("X-User-Email") && user.Tier == userTier {
				// User is authorized, call the next handler
				next.ServeHTTP(w, r)
				return
			}
		}

		// User is not authorized, return a 401 Unauthorized error
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
