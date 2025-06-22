package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/unvbld/Kumparan-Technical-Assessment/model"
	"github.com/unvbld/Kumparan-Technical-Assessment/repository"
)

type ArticleHandler struct {
	Repo *repository.ArticleRepository
}

func (h *ArticleHandler) PostArticle(w http.ResponseWriter, r *http.Request) {
	var a model.Article
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if a.Title == "" || a.Body == "" || a.Author == "" {
		http.Error(w, "Title, body, and author are required", http.StatusBadRequest)
		return
	}

	if len(a.Title) > 200 {
		http.Error(w, "Title too long (max 200 characters)", http.StatusBadRequest)
		return
	}

	if err := h.Repo.CreateArticle(&a); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Article created successfully"})
}

func (h *ArticleHandler) GetArticles(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	author := r.URL.Query().Get("author")

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	articles, err := h.Repo.GetArticles(query, author, page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}

func (h *ArticleHandler) GetArticleByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	article, err := h.Repo.GetArticleByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if article == nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}

func (h *ArticleHandler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	err = h.Repo.DeleteArticle(id)
	if err != nil {
		if err.Error() == "article with id "+vars["id"]+" not found" {
			http.Error(w, "Article not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Article deleted successfully"})
}
