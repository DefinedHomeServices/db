package customers

import (
	"context"
	"errors"
	"testing"
	"time"

	app_firebase "defined-home-services/internal/firebase"

	"github.com/stretchr/testify/assert"
)

// MockFirebaseAPIClient is a mock implementation of the FirebaseAPIClient interface
type MockFirebaseAPIClient struct {
    CreateCustomerFunc func(ctx context.Context, customer app_firebase.Customer) error
}

func (m *MockFirebaseAPIClient) CreateCustomer(ctx context.Context, customer app_firebase.Customer) error {
    return m.CreateCustomerFunc(ctx, customer)
}

func TestCreateCustomerHandler_HandleCreateCustomer(t *testing.T) {
    tests := []struct {
        name           string
        mockFunc       func(ctx context.Context, customer app_firebase.Customer) error
        customer       app_firebase.Customer
        expectedError  error
    }{
        {
            name: "successful customer creation",
            mockFunc: func(ctx context.Context, customer app_firebase.Customer) error {
                return nil
            },
            customer: app_firebase.Customer{
                ID:               "1",
                FirstName:        "John",
                LastName:         "Doe",
                Phone:      "123-456-7890",
                Email:            "john.doe@example.com",
                DateCreated:      time.Now(),
                LastUpdated:      time.Now(),
                Location: app_firebase.CustomerLocation{
                    Address: app_firebase.CustomerAddress{
                        Line1: "123 Main St",
                        Line2: "",
                    },
                    City:             "Anytown",
                    State:            "CA",
                    Zip:              "12345",
                },
                OptedInEmails:    true,
                StripeCustomerID: "cus_123456789",
            },
            expectedError: nil,
        },
        {
            name: "error creating customer",
            mockFunc: func(ctx context.Context, customer app_firebase.Customer) error {
                return errors.New("failed to create customer")
            },
            customer: app_firebase.Customer{
                ID:               "2",
                FirstName:        "Jane",
                LastName:         "Doe",
                Phone:            "098-765-4321",
                Email:            "jane.doe@example.com",
                DateCreated:      time.Now(),
                LastUpdated:      time.Now(),
                Location: app_firebase.CustomerLocation{
                    Address: app_firebase.CustomerAddress{
                        Line1: "456 Elm St",
                        Line2: "",
                    },
                    City:             "Othertown",
                    State:            "NY",
                    Zip:              "54321",
                },
                OptedInEmails:    false,
                StripeCustomerID: "cus_987654321",
            },
            expectedError: errors.New("failed to create customer"),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockClient := &MockFirebaseAPIClient{
                CreateCustomerFunc: tt.mockFunc,
            }
            handler := NewCreateCustomerHandler(mockClient)

            err := handler.CreateCustomer(tt.customer)
            if tt.expectedError != nil {
                assert.ErrorContains(t, err, tt.expectedError.Error())
            } else {
                assert.NoError(t, err)
            }
        })
    }
}