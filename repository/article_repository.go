package repository

import (
	"database/sql"
	"fmt"

	"github.com/unvbld/Kumparan-Technical-Assessment/model"
)

type ArticleRepository struct {
	DB *sql.DB
}

func (r *ArticleRepository) CreateArticle(a *model.Article) error {
	query := `INSERT INTO articles (title, body, author) VALUES ($1, $2, $3)`
	_, err := r.DB.Exec(query, a.Title, a.Body, a.Author)
	return err
}

type ArticleResponse struct {
	Articles []model.Article `json:"articles"`
	Total    int             `json:"total"`
	Page     int             `json:"page"`
	Limit    int             `json:"limit"`
	HasNext  bool            `json:"has_next"`
}

func (r *ArticleRepository) GetArticles(query, author string, page, limit int) (*ArticleResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	countQuery := "SELECT COUNT(*) FROM articles WHERE 1=1"
	countArgs := []interface{}{}
	countIdx := 1
	if query != "" {
		countQuery += fmt.Sprintf(" AND (to_tsvector('english', title) @@ plainto_tsquery('english', $%d) OR to_tsvector('english', body) @@ plainto_tsquery('english', $%d))", countIdx, countIdx+1)
		countArgs = append(countArgs, query, query)
		countIdx += 2
	}

	if author != "" {
		countQuery += fmt.Sprintf(" AND author ILIKE $%d", countIdx)
		countArgs = append(countArgs, "%"+author+"%")
	}
	var total int
	err := r.DB.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, err
	}

	var articles []model.Article
	baseQuery := "SELECT id, title, body, author, created_at FROM articles WHERE 1=1"
	args := []interface{}{}
	idx := 1
	if query != "" {
		baseQuery += fmt.Sprintf(" AND (to_tsvector('english', title) @@ plainto_tsquery('english', $%d) OR to_tsvector('english', body) @@ plainto_tsquery('english', $%d))", idx, idx+1)
		args = append(args, query, query)
		idx += 2
	}

	if author != "" {
		baseQuery += fmt.Sprintf(" AND author ILIKE $%d", idx)
		args = append(args, "%"+author+"%")
		idx++
	}

	baseQuery += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", idx, idx+1)
	args = append(args, limit, offset)

	rows, err := r.DB.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a model.Article
		if err := rows.Scan(&a.ID, &a.Title, &a.Body, &a.Author, &a.CreatedAt); err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}

	hasNext := (page * limit) < total

	return &ArticleResponse{
		Articles: articles,
		Total:    total,
		Page:     page,
		Limit:    limit,
		HasNext:  hasNext,
	}, nil
}

func (r *ArticleRepository) GetArticleByID(id int) (*model.Article, error) {
	query := "SELECT id, title, body, author, created_at FROM articles WHERE id = $1"
	row := r.DB.QueryRow(query, id)

	var a model.Article
	err := row.Scan(&a.ID, &a.Title, &a.Body, &a.Author, &a.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &a, nil
}

func (r *ArticleRepository) DeleteArticle(id int) error {
	query := "DELETE FROM articles WHERE id = $1"
	result, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("article with id %d not found", id)
	}

	return nil
}
