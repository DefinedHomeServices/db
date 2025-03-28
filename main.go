package main

import (
	"db/customers"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting db server at part 8080 exposed on port 8084")
	firebaseCustomersClient := customers.NewFirebaseClient()
    createCustomerHandler := customers.NewCustomerHandler(firebaseCustomersClient)

    // CUSTOMER (Firebase) HANDLERS
    http.HandleFunc("/customer/new", createCustomerHandler.HandleCreateCustomer)
	http.HandleFunc("/customer/get", createCustomerHandler.HandleGetCustomer)
	http.ListenAndServe(":8080", nil)
}