package postgres

import (
	"context"
	"database/sql"
	"fmt"
	
	"github.com/lib/pq"
	"github.com/stackus/errors"
	"middleman/support/internal/domain"
)

type KnowledgeArticleCatalogRepository struct {
	tableName string
	db        *sql.DB
}

var _ domain.KnowledgeArticleCatalogRepository = (*KnowledgeArticleCatalogRepository)(nil)

func NewKnowledgeArticleCatalogRepository(tableName string, db *sql.DB) KnowledgeArticleCatalogRepository {
	return KnowledgeArticleCatalogRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r KnowledgeArticleCatalogRepository) Add(ctx context.Context, article *domain.KnowledgeArticleCatalog) error {
	const query = `
		INSERT INTO %s (id, title, categories, public, view_count, average_rating, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	_, err := r.db.ExecContext(ctx, r.table(query),
		article.ID,
		article.Title,
		pq.Array(article.Categories),
		article.Public,
		article.ViewCount,
		article.AverageRating,
		article.CreatedAt,
	)
	
	return err
}

func (r KnowledgeArticleCatalogRepository) Update(ctx context.Context, article *domain.KnowledgeArticleCatalog) error {
	const query = `
		UPDATE %s 
		SET title = $2, categories = $3, public = $4, view_count = $5, average_rating = $6
		WHERE id = $1
	`
	
	_, err := r.db.ExecContext(ctx, r.table(query),
		article.ID,
		article.Title,
		pq.Array(article.Categories),
		article.Public,
		article.ViewCount,
		article.AverageRating,
	)
	
	return err
}

func (r KnowledgeArticleCatalogRepository) Delete(ctx context.Context, articleID string) error {
	const query = `DELETE FROM %s WHERE id = $1`
	
	_, err := r.db.ExecContext(ctx, r.table(query), articleID)
	
	return err
}

func (r KnowledgeArticleCatalogRepository) Find(ctx context.Context, articleID string) (*domain.KnowledgeArticleCatalog, error) {
	const query = `
		SELECT id, title, categories, public, view_count, average_rating, created_at
		FROM %s
		WHERE id = $1
	`
	
	article := &domain.KnowledgeArticleCatalog{}
	var categories pq.StringArray
	
	err := r.db.QueryRowContext(ctx, r.table(query), articleID).Scan(
		&article.ID,
		&article.Title,
		&categories,
		&article.Public,
		&article.ViewCount,
		&article.AverageRating,
		&article.CreatedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(errors.ErrNotFound, "knowledge article not found")
		}
		return nil, err
	}
	
	article.Categories = []string(categories)
	
	return article, nil
}

func (r KnowledgeArticleCatalogRepository) Search(ctx context.Context, searchQuery string, categories []string, publicOnly bool, limit, offset int) ([]*domain.KnowledgeArticleCatalog, error) {
	query := `
		SELECT id, title, categories, public, view_count, average_rating, created_at
		FROM %s
		WHERE 1=1
	`
	
	args := []interface{}{}
	argCount := 0
	
	// Add search query if provided
	if searchQuery != "" {
		argCount++
		query += fmt.Sprintf(" AND title ILIKE $%d", argCount)
		args = append(args, "%"+searchQuery+"%")
	}
	
	// Filter by categories if provided
	if len(categories) > 0 {
		argCount++
		query += fmt.Sprintf(" AND categories && $%d", argCount)
		args = append(args, pq.Array(categories))
	}
	
	// Filter by public status
	if publicOnly {
		argCount++
		query += fmt.Sprintf(" AND public = $%d", argCount)
		args = append(args, true)
	}
	
	// Order by relevance (view count and rating)
	query += " ORDER BY average_rating DESC, view_count DESC"
	
	if limit > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, limit)
	}
	
	if offset > 0 {
		argCount++
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, offset)
	}
	
	rows, err := r.db.QueryContext(ctx, r.table(query), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return r.scanArticles(rows)
}

func (r KnowledgeArticleCatalogRepository) GetPopular(ctx context.Context, limit int) ([]*domain.KnowledgeArticleCatalog, error) {
	query := `
		SELECT id, title, categories, public, view_count, average_rating, created_at
		FROM %s
		WHERE public = true
		ORDER BY view_count DESC, average_rating DESC
		LIMIT $1
	`
	
	rows, err := r.db.QueryContext(ctx, r.table(query), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return r.scanArticles(rows)
}

func (r KnowledgeArticleCatalogRepository) scanArticles(rows *sql.Rows) ([]*domain.KnowledgeArticleCatalog, error) {
	var articles []*domain.KnowledgeArticleCatalog
	
	for rows.Next() {
		article := &domain.KnowledgeArticleCatalog{}
		var categories pq.StringArray
		
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&categories,
			&article.Public,
			&article.ViewCount,
			&article.AverageRating,
			&article.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		article.Categories = []string(categories)
		articles = append(articles, article)
	}
	
	return articles, rows.Err()
}

func (r KnowledgeArticleCatalogRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}