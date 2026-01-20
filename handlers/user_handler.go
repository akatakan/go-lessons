package handlers

import (
	"encoding/json"
	"fmt"
	"go-backend/models"
	"go-backend/repository"
	"go-backend/utils"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	Repo     repository.IUserRepository
	Validate *validator.Validate
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		KullaniciAdi string `json:"kullanici_adi"`
		Password     string `json:"password"`
	}

	json.NewDecoder(r.Body).Decode(&credentials)

	user, err := h.Repo.GetUserByUsername(credentials.KullaniciAdi)
	if err != nil {
		http.Error(w, "Kullanıcı bulunamadı", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPasswordHash(credentials.Password, user.PasswordHash) {
		http.Error(w, "Geçersiz şifre", http.StatusUnauthorized)
		return
	}

	token, _ := utils.GenerateToken(user.ID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
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

	if err := h.Validate.Struct(user); err != nil {
		http.Error(w, "Validasyon hatası: "+err.Error(), http.StatusBadRequest)
		return
	}
	h.Repo.Save(user)
	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EskiSifre string `json:"eski_sifre"`
		YeniSifre string `json:"yeni_sifre"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	// 1. Token'dan gelen kullanıcı ID'sini al (Context üzerinden)
	userID := r.Context().Value("user_id").(int)

	// 2. Kullanıcıyı DB'den bul
	user, _ := h.Repo.GetUserByID(userID)

	// 3. Eski şifre doğru mu?
	if !utils.CheckPasswordHash(req.EskiSifre, user.PasswordHash) {
		http.Error(w, "Mevcut şifreniz hatalı", http.StatusUnauthorized)
		return
	}

	// 4. Yeni şifreyi hashle ve kaydet
	hashedNew, _ := utils.HashPassword(req.YeniSifre)
	user.PasswordHash = hashedNew
	h.Repo.Update(user)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Şifre başarıyla güncellendi.")
}
