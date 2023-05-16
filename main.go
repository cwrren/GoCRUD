package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var users []User

func createUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate ID
	newUser.ID = len(users) + 1

	// Store user
	users = append(users, newUser)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idParam, ok := vars["id"]
	if !ok {
		http.Error(w, "User ID is missing", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var foundUser User
	for _, user := range users {
		if user.ID == id {
			foundUser = user
			break
		}
	}

	if foundUser.ID == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(foundUser)
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(users)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idParam := vars["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var updatedUser User
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var foundUser *User
	for i := range users {
		if users[i].ID == id {
			foundUser = &users[i]
			break
		}
	}

	if foundUser == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	foundUser.Name = updatedUser.Name
	foundUser.Email = updatedUser.Email
	foundUser.Password = updatedUser.Password

	json.NewEncoder(w).Encode(foundUser)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idParam := vars["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var foundUser *User
	for i := range users {
		if users[i].ID == id {
			foundUser = &users[i]
			users = append(users[:i], users[i+1:]...)
			break
		}
	}

	if foundUser == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func triggerGET() {
	// url endpoint
	url := "http://localhost:8080/user/all"

	// Send the GET request
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error sending GET request: %s\n", err.Error())
		return
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err.Error())
		return
	}

	// Print the response body
	fmt.Printf("Response body: %s\n", body)
}

func triggerPOST() {
	// url endpoint
	url := "https://example.com/endpoint"

	// Create a JSON payload to send in the request body
	payload := []byte(`{"key": "value"}`)

	// Send the POST request
	response, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Printf("Error sending POST request: %s\n", err.Error())
		return
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err.Error())
		return
	}

	// Print the response body
	fmt.Printf("Response body: %s\n", body)
}

func main() {
	triggerGET()
	// Create a new router
	router := mux.NewRouter()

	// CRUD endpoints
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users", getAllUsers).Methods("GET")
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	// Start the server
	log.Println("Server started on port 8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
