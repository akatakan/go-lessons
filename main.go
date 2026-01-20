package main

import (
	"fmt"
	"go-backend/handlers"
	"go-backend/middleware"
	"go-backend/models"
	"go-backend/repository"
	"log/slog"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("İstek Geldi: %s %s\n", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, pass, dbname, port)
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&models.User{})

	mux := http.NewServeMux()
	repo := repository.NewSQLUserRepository(db)
	handler := &handlers.UserHandler{Repo: repo}
	mux.Handle("GET /user/{id}", middleware.AuthMiddleware(http.HandlerFunc(handler.GetUser)))
	mux.Handle("POST /user/update-password", middleware.AuthMiddleware(http.HandlerFunc(handler.UpdatePassword)))

	mux.HandleFunc("POST /user", handler.CreateUser)
	mux.HandleFunc("POST /login", handler.Login)
	http.HandleFunc("GET /user/{id}", handler.GetUser)
	http.HandleFunc("POST /user", handler.CreateUser)

	slog.Info("Sunucu hazırlanıyor", "port", 8080, "env", "production")
	finalHandler := middleware.CORSMiddleware(loggerMiddleware(mux))
	http.ListenAndServe(":8080", finalHandler)
}
