package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var client *mongo.Client
var jwtSecret = []byte("supersecretkey") // Use environment variables in production

type User struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"password" bson:"password"` // Store hashed password
}

// JWT Claims struct
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Hash a user's password before storing it
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Compare hashed password with the provided password during login
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// Extract the token from the "Bearer <token>" format
		tokenString := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// User registration handler
func registerUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

// User login handler and JWT token generation
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

	// Create JWT token upon successful login
	expirationTime := time.Now().Add(time.Hour * 1) // Token expires in 1 hour
	claims := &Claims{
		Username: foundUser.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Return the token in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// Bill struct to represent a bill document
type Bill struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string             `json:"name,omitempty" bson:"name,omitempty"`
	Amount float64            `json:"amount,omitempty" bson:"amount,omitempty"`
}

// Fetch all bills from MongoDB (protected route)
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

// Create a new bill and insert it into MongoDB (protected route)
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

	// Public routes
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/register", registerUser).Methods("POST")
	r.HandleFunc("/login", loginUser).Methods("POST")

	// Protected routes (require JWT token)
	r.Handle("/bills", jwtMiddleware(http.HandlerFunc(getBills))).Methods("GET")
	r.Handle("/bills", jwtMiddleware(http.HandlerFunc(createBill))).Methods("POST")

	// Enable CORS for the server
	corsHandler := cors.AllowAll().Handler(r)

	// Start the server with CORS-enabled routes
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
