package handlers

import (
	"Blogsite/config"
	middleware "Blogsite/middlewares"
	"Blogsite/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// CreateBlog handler
func CreateBlog(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uint)
	log.Printf("CreateBlog: user ID from context: %v", userID)

	var blog models.Blog
	if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	blog.UserID = userID

	if err := config.DB.Create(&blog).Error; err != nil {
		log.Printf("Error creating blog: %v", err)
		http.Error(w, "Failed to create blog", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blog)
}

// GetUserBlogs handler (All user tasks)
func GetUserBlogs(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uint)
	log.Printf("GetUserBlogs: user ID from context: %v", userID)

	var blogs []models.Blog
	if err := config.DB.Where("user_id = ?", userID).Find(&blogs).Error; err != nil {
		http.Error(w, "Failed to retrieve blogs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogs)
}

func GetAllBlogs(w http.ResponseWriter, r *http.Request) {
	var blogs []models.Blog
	if err := config.DB.Find(&blogs).Error; err != nil {
		http.Error(w, "Failed to retrieve blogs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogs)
}

// UpdateBlog handler
func UpdateBlog(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uint)
	log.Printf("UpdateBlog: user ID from context: %v", userID)

	vars := mux.Vars(r)
	blogID := vars["id"]

	var blog models.Blog
	if err := config.DB.First(&blog, blogID).Error; err != nil {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}

	if blog.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := config.DB.Save(&blog).Error; err != nil {
		http.Error(w, "Failed to update blog", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blog)
}

// DeleteBlog handler
func DeleteBlog(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uint)
	log.Printf("DeleteBlog: user ID from context: %v", userID)

	vars := mux.Vars(r)
	blogID := vars["id"]

	var blog models.Blog
	if err := config.DB.First(&blog, blogID).Error; err != nil {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}

	if blog.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := config.DB.Delete(&blog).Error; err != nil {
		http.Error(w, "Failed to delete blog", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Blog successfully deleted"}
	json.NewEncoder(w).Encode(response)
}

func GetBlogById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var blog models.Blog
	if err := config.DB.First(&blog, id).Error; err != nil {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blog)
}
