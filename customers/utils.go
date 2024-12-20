package customers

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// NewFirebaseClient initializes Firebase and Firestore using a service account JSON file
func NewFirebaseClient() *FirebaseDBClient {
    ctx := context.Background()
    fmt.Println("Starting Firebase initialization")
    
    // Initialize Firebase Admin SDK with credentials
    opt := option.WithCredentialsFile("./service-account.json")
    app, err := firebase.NewApp(ctx, nil, opt)
    
    if err != nil {
        log.Fatalf("Error initializing Firebase: %v", err)
    } else {
        fmt.Println("Firebase initialized successfully")
    }

    client, err := app.Firestore(ctx)
    if err != nil {
        log.Fatalf("Error initializing Firestore: %v", err)
    }

    return &FirebaseDBClient{DB: client}
}