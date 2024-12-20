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
}

// FirebaseDBClient holds a Firestore client to interact with the Firebase database
type FirebaseDBClient struct {
    DB  *firestore.Client
}

// CreateCustomerHandler handles customer creation requests
type CreateCustomerHandler struct {
    client 	FirebaseAPIClient
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

// NewCreateCustomerHandler creates a new CreateCustomerHandler with a given FirebaseAPIClient
func NewCreateCustomerHandler(client FirebaseAPIClient) *CreateCustomerHandler {
    return &CreateCustomerHandler{client: client}
}

// HandleCreateCustomer handles the HTTP request to create a new customer
func (h *CreateCustomerHandler) HandleCreateCustomer(w http.ResponseWriter, r *http.Request) {
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

func (h *CreateCustomerHandler) CreateCustomer(customer map[string]interface{}) (string, error) {
    ctx := context.Background()

    fmt.Printf("Creating customer in Firebase %v", customer)

    firebaseId, err := h.client.AddCustomerToDatabase(ctx, customer)

    if err != nil {
        return "", err
    }

    return firebaseId, nil
}