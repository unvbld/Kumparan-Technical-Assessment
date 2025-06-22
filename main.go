package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/unvbld/Kumparan-Technical-Assessment/handler"
	"github.com/unvbld/Kumparan-Technical-Assessment/repository"
	"golang.org/x/time/rate"
)

func rateLimitMiddleware(limiter *rate.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=kumparan sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	if err := db.Ping(); err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
	log.Println("Connected to database successfully")
	repo := &repository.ArticleRepository{DB: db}
	h := &handler.ArticleHandler{Repo: repo}

	limiter := rate.NewLimiter(rate.Every(time.Minute/100), 100)
	r := mux.NewRouter()

	r.Use(corsMiddleware)
	r.Use(rateLimitMiddleware(limiter))

	r.HandleFunc("/articles", h.PostArticle).Methods("POST")
	r.HandleFunc("/articles", h.GetArticles).Methods("GET")
	r.HandleFunc("/articles/{id}", h.GetArticleByID).Methods("GET")
	r.HandleFunc("/articles/{id}", h.DeleteArticle).Methods("DELETE")

	log.Println("Server running at :8080")
	log.Println("Rate limit: 100 requests per minute")
	log.Println("Database pool: 25 max connections, 5 idle connections")
	http.ListenAndServe(":8080", r)
}
