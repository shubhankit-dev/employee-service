package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// User represents a user record
type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Department string `json:"department"`
	Course     string `json:"course"`
	Mobile     string `json:"mobile"`
	City       string `json:"city"`
}

// DB instance
var db *sql.DB

func main() {
	var err error
	// Connect to MySQL database
	db, err = sql.Open("mysql", "kalpit:password@tcp(192.168.100.4:3306)/demo")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/adduser", addUserHandler)
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if user.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Insert user into the database
	_, err := db.Exec("INSERT INTO users (name,city,course,department,mobile) VALUES (?,?,?,?,?)", user.Name, user.City, user.Course, user.Department, user.Mobile)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		return
	}
	fmt.Println("User added successfully")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User added successfully"))
}
