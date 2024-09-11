package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	d "password/db"
	m "password/models"
	util "password/utils"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func SignUp(w http.ResponseWriter, r *http.Request) {
	var user m.User
	// Parse the JSON request body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the user request using go-playground/validator
	err = validate.Struct(user)
	if err != nil {
		// Extract validation errors and return the first one
		for _, err := range err.(validator.ValidationErrors) {
			http.Error(w, fmt.Sprintf("Validation failed for field '%s': %s", err.Field(), err.Tag()), http.StatusBadRequest)
			return
		}
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error while hashing password", http.StatusInternalServerError)
		return
	}

	// Store user in DB
	db := d.GetDBConn()
	defer db.Close()
	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, string(hashedPassword))
	if err != nil {
		println(err.Error())
		http.Error(w, "Error storing user", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("User signed up successfully!"))
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user m.User
	json.NewDecoder(r.Body).Decode(&user)

	// Get the stored hash from the DB
	var storedHash string
	db := d.GetDBConn()
	defer db.Close()
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", user.Username).Scan(&storedHash)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Write([]byte("Login successful!"))
}

func StorePassword(w http.ResponseWriter, r *http.Request) {
	var passwordData m.Password
	json.NewDecoder(r.Body).Decode(&passwordData)

	// Encrypt the password before storing
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordData.StoredPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error while encrypting password", http.StatusInternalServerError)
		return
	}

	// Store the password in DB
	db := d.GetDBConn()
	defer db.Close()
	_, err = db.Exec("INSERT INTO passwords (username, stored_password) VALUES (?, ?)", passwordData.Username, string(encryptedPassword))
	if err != nil {
		http.Error(w, "Error storing password", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Password stored successfully!"))
}

func GetPassword(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	// Retrieve the encrypted password from DB
	var encryptedPassword string
	db := d.GetDBConn()
	defer db.Close()
	err := db.QueryRow("SELECT stored_password FROM passwords WHERE username = ?", username).Scan(&encryptedPassword)
	if err == sql.ErrNoRows {
		http.Error(w, "No password found for user", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Encrypted password: " + encryptedPassword))
}

func PasswordQuality(w http.ResponseWriter, r *http.Request) {
	var password m.Pass
	json.NewDecoder(r.Body).Decode(&password)

	result := util.EvaluatePassword(password.Password)

	// Set the Content-Type header to indicate we're returning JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode the struct to JSON and write it to the response
	json.NewEncoder(w).Encode(result)
}
