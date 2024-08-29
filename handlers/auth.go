package handlers

import (
	"Blogsite/config"
	"Blogsite/models"
	"Blogsite/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Decode the request body into user struct
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate username
	if err := validateUsername(user.Username); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate email format
	if !isValidEmail(user.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Validate password strength
	if err := validatePassword(user.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Create user in the database
	if err := config.DB.Create(&user).Error; err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Send success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := map[string]string{"message": "User successfully registered"}
	json.NewEncoder(w).Encode(response)
}

func validateUsername(username string) error {
	if len(username) < 6 {
		return fmt.Errorf("Username must be at least 6 characters long")
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(username) {
		return fmt.Errorf("Username must contain only alphanumeric characters")
	}

	return nil
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func validatePassword(password string) error {
	var (
		minLength   = 8
		uppercase   = regexp.MustCompile(`[A-Z]`)
		lowercase   = regexp.MustCompile(`[a-z]`)
		number      = regexp.MustCompile(`[0-9]`)
		specialChar = regexp.MustCompile(`[!@#~$%^&*()+|_]{1}`)
	)

	if len(password) < minLength {
		return fmt.Errorf("Password must be at least %d characters long", minLength)
	}
	if !uppercase.MatchString(password) {
		return fmt.Errorf("Password must contain at least one uppercase letter")
	}
	if !lowercase.MatchString(password) {
		return fmt.Errorf("Password must contain at least one lowercase letter")
	}
	if !number.MatchString(password) {
		return fmt.Errorf("Password must contain at least one number")
	}
	if !specialChar.MatchString(password) {
		return fmt.Errorf("Password must contain at least one special character (!@#~$%%^&*()+|_)")
	}
	return nil
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if creds.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	var user models.User

	if err := config.DB.Where("username = ?", creds.Username).First(&user).Error; err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
