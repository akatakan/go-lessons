package main

import (
	"fmt"
	"go-backend/handlers"
	"go-backend/models"
	"go-backend/repository"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID           int    `json:"id"`
	KullaniciAdi string `json:"kullanici_adi"`
	Email        string `json:"email"`
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// İşlem başlamadan önce:
		fmt.Printf("İstek Geldi: %s %s\n", r.Method, r.URL.Path)

		// Asıl fonksiyona (Handler) devret:
		next.ServeHTTP(w, r)
	})
}

func main() {
	host := os.Getenv("DB_HOST") // Docker Compose'dan "db" gelecek
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, pass, dbname, port)
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&models.User{})

	repo := repository.NewSQLUserRepository(db)
	handler := &handlers.UserHandler{Repo: repo}
	http.HandleFunc("GET /user/{id}", handler.GetUser)
	http.HandleFunc("POST /user", handler.CreateUser)
	fmt.Println("Sunucu 8080 portunda başlıyor...")
	http.ListenAndServe(":8080", loggerMiddleware(http.DefaultServeMux))
}
