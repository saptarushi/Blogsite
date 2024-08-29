package handlers

import (
	"Blogsite/config"
	"Blogsite/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestRegister(t *testing.T) {
	// Ensure the database is initialized
	if config.DB == nil {
		dsn := "host=localhost user=postgres password=Postgresql@1234 dbname=blogsite_db port=5432 sslmode=disable"
		var err error
		config.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			t.Fatalf("Failed to connect to database: %v", err)
		}
	}

	// Start a transaction and defer a rollback
	tx := config.DB.Begin()
	defer tx.Rollback()

	// Create a user for testing login
	user := models.User{
		Username: "LoginTestUser",
		Email:    "logintestuser@example.com",
		Password: "HashedPassword!23", // Replace with actual hashing in production
	}

	// Ensure the user is created before running the tests
	if err := tx.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Define the test cases
	tests := []struct {
		name           string
		payload        map[string]string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Invalid request payload",
			payload: map[string]string{
				"username": "username",
				"email":    "invalid-email",
				"password": "Passw0rd!",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid email format",
		},
		{
			name: "Weak password",
			payload: map[string]string{
				"username": "WeakPasswordUser",
				"email":    "weakpassword@example.com",
				"password": "weak",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Password must be at least 8 characters long",
		},
		{
			name: "Successful registration",
			payload: map[string]string{
				"username": "RegisterSuccessUser",
				"email":    "registersuccessuser@example.com",
				"password": "StrongPassw0rd!",
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   "{\"message\":\"User successfully registered\"}",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Clean up any existing users with the same email addresses before running the test
			config.DB.Exec("DELETE FROM users WHERE email IN ('weakpassword@example.com', 'registersuccessuser@example.com')")

			// Convert the payload to JSON format
			reqBody, err := json.Marshal(tc.payload)
			if err != nil {
				t.Fatalf("Error marshalling payload: %v", err)
			}

			// Create the HTTP request for registration
			req, err := http.NewRequest("POST", "/api/register", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}
			rr := httptest.NewRecorder()

			// Call the registration handler
			handler := http.HandlerFunc(Register)
			handler.ServeHTTP(rr, req)

			// Validate the HTTP status code
			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
			}

			// Validate the response body
			if !bytes.Contains(rr.Body.Bytes(), []byte(tc.expectedBody)) {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tc.expectedBody)
			}
		})
	}

	// Rollback the transaction to ensure isolation of test data
	tx.Rollback()
}
