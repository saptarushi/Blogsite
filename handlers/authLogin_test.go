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

func TestLogin(t *testing.T) {
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
			name: "Incorrect password",
			payload: map[string]string{
				"username": "LoginTestUser",
				"password": "WrongPassword!23",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid credentials",
		},
		{
			name: "Successful login",
			payload: map[string]string{
				"username": "LoginTestUser1",
				"password": "HashedPassword!23", // Use actual password in real application
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "\"token\":\"",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Convert the payload to JSON format
			reqBody, err := json.Marshal(tc.payload)
			if err != nil {
				t.Fatalf("Error marshalling payload: %v", err)
			}

			// Create the HTTP request for login
			req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}
			rr := httptest.NewRecorder()

			// Call the login handler
			handler := http.HandlerFunc(Login)
			handler.ServeHTTP(rr, req)

			// Validate the HTTP status code
			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
			}

			// Validate the response body (for the token, just check if "token" is in the response)
			if tc.expectedStatus == http.StatusOK {
				if !bytes.Contains(rr.Body.Bytes(), []byte(tc.expectedBody)) {
					t.Errorf("handler returned unexpected body: got %v want %v",
						rr.Body.String(), tc.expectedBody)
				}
			} else {
				if !bytes.Contains(rr.Body.Bytes(), []byte(tc.expectedBody)) {
					t.Errorf("handler returned unexpected body: got %v want %v",
						rr.Body.String(), tc.expectedBody)
				}
			}
		})
	}

	// Rollback the transaction to ensure isolation of test data
	tx.Rollback()
}
