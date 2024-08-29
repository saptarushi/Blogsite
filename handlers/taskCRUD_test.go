package handlers

import (
	"Blogsite/config"
	"Blogsite/middlewares"
	"Blogsite/models"
	"Blogsite/utils"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func TestBlogCRUD(t *testing.T) {
	// Initialize the database
	config.InitDB()

	// Start a transaction
	tx := config.DB.Begin()
	defer tx.Rollback()

	// Create a user in the database for testing
	user := models.User{
		Username: "TestUserCRUD1",
		Email:    "testusercrud1@example.com",
		Password: "Hashedpassword$43", // Replace with appropriate hash function in a real app
	}

	if err := tx.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Commit the user creation transaction so the user is available in the DB
	tx.Commit()

	// Re-open a new transaction for the CRUD tests
	tx = config.DB.Begin()
	defer tx.Rollback()

	// Generate a valid JWT token for testing
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	// Create a new blog post
	t.Run("CreateBlog", func(t *testing.T) {
		authMiddleware := middlewares.AuthMiddleware(http.HandlerFunc(CreateBlog))

		blog := models.Blog{
			Title:       "Test Blog CRUD",
			Description: "This is a test blog for CRUD.",
			Completed:   false,
			UserID:      user.ID,
		}

		body, _ := json.Marshal(blog)
		req := httptest.NewRequest("POST", "/api/user/blog", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		ctx := req.Context()
		ctx = context.WithValue(ctx, middlewares.UserIDKey, user.ID)
		req = req.WithContext(ctx)

		authMiddleware.ServeHTTP(w, req)

		if status := w.Code; status != http.StatusOK {
			t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var createdBlog models.Blog
		if err := json.NewDecoder(w.Body).Decode(&createdBlog); err != nil {
			t.Fatalf("Could not decode response: %v", err)
		}

		if createdBlog.Title != blog.Title {
			t.Fatalf("Expected blog title to be %v, got %v", blog.Title, createdBlog.Title)
		}
	})

	// Retrieve the created blog post to get its ID
	var blog models.Blog
	if err := config.DB.Where("title = ?", "Test Blog CRUD").First(&blog).Error; err != nil {
		t.Fatalf("Failed to find created blog: %v", err)
	}

	// Update the created blog post
	t.Run("UpdateBlog", func(t *testing.T) {
		blog.Description = "Updated description for CRUD"
		body, _ := json.Marshal(blog)
		req := httptest.NewRequest("PUT", "/api/user/blog/"+strconv.Itoa(int(blog.ID)), bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(int(blog.ID))})
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		ctx := req.Context()
		ctx = context.WithValue(ctx, middlewares.UserIDKey, blog.UserID)
		req = req.WithContext(ctx)

		authMiddleware := middlewares.AuthMiddleware(http.HandlerFunc(UpdateBlog))
		authMiddleware.ServeHTTP(w, req)

		if status := w.Code; status != http.StatusOK {
			t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var updatedBlog models.Blog
		if err := json.NewDecoder(w.Body).Decode(&updatedBlog); err != nil {
			t.Fatalf("Could not decode response: %v", err)
		}

		if updatedBlog.Description != "Updated description for CRUD" {
			t.Fatalf("Expected blog description to be updated, got %v", updatedBlog.Description)
		}
	})

	// Get the blog by ID
	t.Run("GetBlogById", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/blog/"+strconv.Itoa(int(blog.ID)), nil)
		req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(int(blog.ID))})

		w := httptest.NewRecorder()

		GetBlogById(w, req)

		if status := w.Code; status != http.StatusOK {
			t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var fetchedBlog models.Blog
		if err := json.NewDecoder(w.Body).Decode(&fetchedBlog); err != nil {
			t.Fatalf("Could not decode response: %v", err)
		}

		if fetchedBlog.ID != blog.ID {
			t.Fatalf("Expected blog ID to be %v, got %v", blog.ID, fetchedBlog.ID)
		}
	})

	// Delete the created blog post
	t.Run("DeleteBlog", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/user/blog/"+strconv.Itoa(int(blog.ID)), nil)
		req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(int(blog.ID))})
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		ctx := req.Context()
		ctx = context.WithValue(ctx, middlewares.UserIDKey, blog.UserID)
		req = req.WithContext(ctx)

		authMiddleware := middlewares.AuthMiddleware(http.HandlerFunc(DeleteBlog))
		authMiddleware.ServeHTTP(w, req)

		if status := w.Code; status != http.StatusOK {
			t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// Check if the blog was deleted from the database
		var deletedBlog models.Blog
		err := config.DB.First(&deletedBlog, blog.ID).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			t.Fatalf("Expected blog to be deleted, but found: %v", err)
		}
	})

	// Rollback the transaction after the test
	tx.Rollback()
}
