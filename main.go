package main

import (
	"log"
	"net/http"

	c "password/controller"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	// Create a new router
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/sighup", c.SignUp).Methods("POST")
	router.HandleFunc("/login", c.Login).Methods("POST")
	router.HandleFunc("/store_password", c.StorePassword).Methods("POST")
	router.HandleFunc("/get_password", c.GetPassword).Methods("GET")
	router.HandleFunc("/verifier", c.PasswordQuality).Methods("POST")

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
