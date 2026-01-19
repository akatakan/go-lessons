package handlers

import (
	"encoding/json"
	"go-backend/models"
	"go-backend/repository"
	"net/http"
	"strconv"
)

type UserHandler struct {
	Repo repository.IUserRepository
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Geçersiz ID", http.StatusBadRequest)
		return
	}
	user, ok := h.Repo.GetUserByID(id)
	if !ok {
		http.Error(w, "Kullanıcı bulunamadı", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Geçersiz veri", http.StatusBadRequest)
		return
	}
	h.Repo.Save(user)
	w.WriteHeader(http.StatusCreated)
}
