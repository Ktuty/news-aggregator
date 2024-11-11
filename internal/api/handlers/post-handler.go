package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"news/internal/models"
	"strconv"
)

// http метод для получения постов
func (api *API) Posts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	quantity, err := strconv.Atoi(mux.Vars(r)["quantity"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if quantity < 1 {
		quantity = 10
	}

	var posts []models.Post

	posts, err = api.db.Posts(quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err = json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
