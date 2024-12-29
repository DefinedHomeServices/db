package customers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
)

// FirebaseAPIClient defines a method for creating a customer in the database
type FirebaseAPIClient interface {
    AddCustomerToDatabase(ctx context.Context, customer map[string]interface{}) (string, error)
    GetCustomer(ctx context.Context, email string) (map[string]interface{}, error)
}

// FirebaseDBClient holds a Firestore client to interact with the Firebase database
type FirebaseDBClient struct {
    DB  *firestore.Client
}

// CustomerHandler handles customer creation requests
type CustomerHandler struct {
    client 	FirebaseAPIClient
}

func (h *FirebaseDBClient) GetCustomer(ctx context.Context, email string) (map[string]interface{}, error) {
    fmt.Println("Getting customer from email:", email)
    docRef := h.DB.Collection("customers").Where("email", "==", email).Documents(ctx)
    docs, err := docRef.GetAll()
    
    if (len(docs ) == 0) {
        return nil, nil
    }

    if err != nil {
        fmt.Println("Error getting customer:", err)
        return nil, err
    }
    customer := docs[0].Data()
    
    return customer, nil
}

func (h *FirebaseDBClient) AddCustomerToDatabase(ctx context.Context, customer map[string]interface{}) (string, error) {
    fmt.Println("Creating customer in Firebase")
    docRef, _, err := h.DB.Collection("customers").Add(ctx, customer)
    if err != nil {
        fmt.Println("Error creating customer:", err)
        return "", err
    }
    fmt.Println("Customer created with ID:", docRef.ID)
    return docRef.ID, nil
}

// NewCustomerHandler creates a new CreateCustomerHandler with a given FirebaseAPIClient
func NewCustomerHandler(client FirebaseAPIClient) *CustomerHandler {
    return &CustomerHandler{client: client}
}

func (h *CustomerHandler) HandleGetCustomer(w http.ResponseWriter, r *http.Request) {
    email := r.URL.Query().Get("email")
    if email == "" {
        http.Error(w, "Missing email parameter", http.StatusBadRequest)
        return
    }
    customer, err := h.client.GetCustomer(r.Context(), email)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to get customer: %v", err), http.StatusInternalServerError)
        return
    }
    if customer == nil {
        http.Error(w, "Customer not found", http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(customer)
    return
}

// HandleCreateCustomer handles the HTTP request to create a new customer
func (h *CustomerHandler) HandleCreateCustomer(w http.ResponseWriter, r *http.Request) {
    var customer map[string]interface{}
    
    err := json.NewDecoder(r.Body).Decode(&customer)
    
    fmt.Printf("Customer Decoded: %v", customer)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    firebaseId, err := h.CreateCustomer(customer)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to create customer: %v", err), http.StatusInternalServerError)
        return
    }

    response := map[string]string{
        "firebase_id": firebaseId,
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(response)    
}

func (h *CustomerHandler) CreateCustomer(customer map[string]interface{}) (string, error) {
    ctx := context.Background()

    fmt.Printf("Creating customer in Firebase %v", customer)

    firebaseId, err := h.client.AddCustomerToDatabase(ctx, customer)

    if err != nil {
        return "", err
    }

    return firebaseId, nil
}