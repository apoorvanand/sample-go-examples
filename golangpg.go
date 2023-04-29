package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
)

type Book struct {
	ID    string `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
	Author string `json:"author,omitempty"`
}

var db *sql.DB 

func main() {
	var err error
	db, err = sql.Open("postgres", "postgresql://username:password@localhost/bookstore?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", createBook).Methods("POST")
	router.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	var books []Book
	rows, err := db.Query("SELECT id, title, author FROM books")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author)
		if err != nil {
			log.Fatal(err)
		}
		books = append(books, book)
	}
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var book Book
	row := db.QueryRow("SELECT id, title, author FROM books WHERE id=$1", params["id"])
	err := row.Scan(&book.ID, &book.Title, &book.Author)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(book)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	_, err := db.Exec("INSERT INTO books(id, title, author) VALUES($1, $2, $3)", book.ID, book.Title, book.Author)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	_, err := db.Exec("UPDATE books SET title=$1, author=$2 WHERE id=$3", book.Title, book.Author, params["id"])
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(book)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	_, err := db.Exec("DELETE FROM books WHERE id=$1", params["id"])
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode("Book deleted successfully")
}
