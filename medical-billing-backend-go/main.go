package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var client *mongo.Client

type User struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"password" bson:"password"` // Store hashed password
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	collection := client.Database("medical-billing").Collection("users")
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collection := client.Database("medical-billing").Collection("users")
	var foundUser User
	err = collection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&foundUser)
	if err != nil || !checkPasswordHash(user.Password, foundUser.Password) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Successful login, set session or JWT here (optional)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(foundUser)
}

// Bill struct to represent a bill document
type Bill struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string             `json:"name,omitempty" bson:"name,omitempty"`
	Amount float64            `json:"amount,omitempty" bson:"amount,omitempty"`
}

// Handler to fetch all bills from the MongoDB database
func getBills(w http.ResponseWriter, r *http.Request) {
	var bills []Bill
	collection := client.Database("medical-billing").Collection("bills")
	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var bill Bill
		err := cur.Decode(&bill)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bills = append(bills, bill)
	}

	if err := cur.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bills)
}

// Handler to create a new bill and insert it into the MongoDB database
func createBill(w http.ResponseWriter, r *http.Request) {
	var bill Bill
	err := json.NewDecoder(r.Body).Decode(&bill)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collection := client.Database("medical-billing").Collection("bills")
	result, err := collection.InsertOne(context.TODO(), bill)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		InsertedID primitive.ObjectID `json:"inserted_id"`
	}{InsertedID: result.InsertedID.(primitive.ObjectID)}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Home page handler
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Medical Billing Portal!")
}

func main() {
	// MongoDB connection setup
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	// Router setup
	r := mux.NewRouter()

	// Route definitions
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/bills", createBill).Methods("POST")
	r.HandleFunc("/bills", getBills).Methods("GET") // Route to fetch bills
	r.HandleFunc("/register", registerUser).Methods("POST")
	r.HandleFunc("/login", loginUser).Methods("POST")

	// Enable CORS for the server
	corsHandler := cors.AllowAll().Handler(r)

	// Start the server with CORS-enabled routes
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
