package items

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
)

// FirebaseDBClient holds a Firestore client to interact with the Firebase database
type FirebaseDBClient struct {
    db *firestore.Client
}

// FirebaseAPIClient defines a method for adding an item to the database
type FirebaseAPIClient interface {
    AddItemToCollection(ctx context.Context, collection string, item interface{}) error
}

// CRUDHandler uses a FirebaseAPIClient to handle item addition
type CRUDHandler struct {
    client FirebaseAPIClient
}

// AddItem adds an item to the specified collection in Firestore
func (c *FirebaseDBClient) AddItemToCollection(ctx context.Context, collection string, item interface{}) error {
    fmt.Println("Adding item to Firebase")
    docRef, _, err := c.db.Collection(collection).Add(ctx, item)
    if err != nil {
        fmt.Println("Error adding item:", err)
        return err
    }
    fmt.Println("Item added with ID:", docRef.ID)
    return nil
}

func (handler *CRUDHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var item map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	collection := item["collection"].(string)
	if collection == "" {
		http.Error(w, "Collection not specified", http.StatusBadRequest)
		return
	}

	value := item["value"]
	if value == nil {
		http.Error(w, "Value not specified", http.StatusBadRequest)
		return
	}

	if err := handler.client.AddItemToCollection(context.Background(), collection, value); err != nil {
		http.Error(w, "Failed to add item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Item added successfully")
}