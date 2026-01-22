package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Category struct {
	ID          int    `json:"id"`
	Nama        string `json:"nama"`
	Description string `json:"description"`
}

var categories = []Category{
	{ID: 1, Nama: "Elektronik", Description: "Barang-barang elektronik"},
	{ID: 2, Nama: "Fesyen", Description: "Pakaian dan aksesoris"},
	{ID: 3, Nama: "Makanan", Description: "Makanan ringan dan berat"},
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/categories", categoriesHandler)
	http.HandleFunc("/categories/", categoryIDHandler)

	fmt.Println("Server is running at localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "OK",
		"message": "API Running",
	})
}

func categoriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		listCategories(w, r)
	} else if r.Method == "POST" {
		createCategory(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func categoryIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/categories" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Method == "GET" {
		getCategoryByID(w, r)
	} else if r.Method == "PUT" {
		updateCategory(w, r)
	} else if r.Method == "DELETE" {
		deleteCategory(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func listCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func createCategory(w http.ResponseWriter, r *http.Request) {
	var categoryBaru Category
	err := json.NewDecoder(r.Body).Decode(&categoryBaru)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	categoryBaru.ID = len(categories) + 1
	categories = append(categories, categoryBaru)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(categoryBaru)
}

func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	for _, c := range categories {
		if c.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(c)
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

func updateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	var updateData Category
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	for i := range categories {
		if categories[i].ID == id {
			updateData.ID = id
			categories[i] = updateData

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateData)
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	for i, c := range categories {
		if c.ID == id {
			categories = append(categories[:i], categories[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "sukses delete category",
			})
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}
