package main

import (
	"context"
	"fmt"
	"go-backend/handlers"
	"go-backend/middleware"
	"go-backend/models"
	"go-backend/repository"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
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
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Veritabanına bağlanılamadı", "error", err)
		os.Exit(1)
	}
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		slog.Error("Migrasyon hatası", "error", err)
		os.Exit(1)
	}
	mux := http.NewServeMux()
	repo := repository.NewSQLUserRepository(db)
	v := validator.New()
	handler := &handlers.UserHandler{Repo: repo, Validate: v}
	mux.Handle("GET /user/{id}", middleware.AuthMiddleware(http.HandlerFunc(handler.GetUser)))
	mux.Handle("POST /user/update-password", middleware.AuthMiddleware(http.HandlerFunc(handler.UpdatePassword)))

	mux.HandleFunc("POST /user", handler.CreateUser)
	mux.HandleFunc("POST /login", handler.Login)

	slog.Info("Sunucu hazırlanıyor", "port", 8080, "env", "production")
	finalHandler := middleware.CORSMiddleware(loggerMiddleware(mux))
	server := &http.Server{
		Addr:    ":8080",
		Handler: finalHandler,
	}
	go func() {
		slog.Info("Sunucu 8080 portunda başlıyor...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Sunucu başlatılamadı", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	slog.Info("Sunucu kapatılıyor...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Sunucu zorla kapatıldı", "error", err)
	}

	slog.Info("Sunucu güvenli bir şekilde durduruldu.")
}
