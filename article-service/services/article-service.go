package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"websynergy/article-service/models"

	"github.com/google/uuid"
)

type ArticleService struct {
	dbClient *sql.DB
}

func NewArticleService(dbClient *sql.DB) *ArticleService {
	return &ArticleService{
		dbClient: dbClient,
	}
}

func (s *ArticleService) GetArticleById(ctx context.Context, ID string) (*models.Article, error) {
	query := "SELECT * FROM articles WHERE id = ?"
	ctx, cancelfunc := context.WithTimeout(ctx, 5*time.Second)
    defer cancelfunc()

	var article models.Article
	err := s.dbClient.QueryRowContext(ctx, query, ID).Scan(&article.ID, &article.Title, &article.Content, &article.Author)
	if err != nil {
		fmt.Printf("ERROR SELECT QUERY - %s", err)
		return nil, err
	}
	
	return &article, nil
}

func (s *ArticleService) GetArticles(ctx context.Context) ([]models.Article, error) {
	query := "SELECT * FROM articles"
	ctx, cancelfunc := context.WithTimeout(ctx, 5*time.Second)
    defer cancelfunc()

	rows, err := s.dbClient.QueryContext(ctx, query)
	if err != nil {
		fmt.Printf("ERROR SELECT QUERY - %s", err)
		return nil, err
	}
	var articleList []models.Article
	for rows.Next() {
		var article models.Article
		err = rows.Scan(&article.ID, &article.Title, &article.Content, &article.Author)
		if err != nil {
			fmt.Printf("ERROR QUERY SCAN - %s", err)
			return nil, err
		}
		articleList = append(articleList, article)
	}
	return articleList, nil
}

func (s *ArticleService) CreateArticle(ctx context.Context, a models.Article) (*struct{Id string}, error) {

	uid := uuid.New()
	query := "INSERT INTO articles(id, title, content, author) VALUES(?,?,?,?)"
    ctx, cancelfunc := context.WithTimeout(ctx, 5*time.Second)
    defer cancelfunc()
    stmt, err := s.dbClient.PrepareContext(ctx, query)
    if err != nil {
        fmt.Printf("Error %s when preparing SQL statement", err)
        return nil, err
    }
    defer stmt.Close()
    res, err := stmt.ExecContext(ctx, uid.String(), a.Title, a.Content, a.Author)
    if err != nil {
        fmt.Printf("Error %s when inserting row into products table", err)
        return nil, err
    }

    fmt.Printf("article created: %v", res)
	data := &struct{Id string}{Id: uid.String()}
    return data, nil
}

// func (s *ArticleService) DeleteArticleById(ctx context.Context, ID string) error {
// 	query := "DELETE FROM articles WHERE id = ?"
// 	ctx, cancelfunc := context.WithTimeout(ctx, 5*time.Second)
//     defer cancelfunc()

// 	var article models.Article
// 	err := s.dbClient.QueryRowContext(ctx, query, ID).Scan(&article.ID, &article.Title, &article.Content, &article.Author)
// 	if err != nil {
// 		fmt.Printf("ERROR SELECT QUERY - %s", err)
// 		return nil, err
// 	}
	
// 	return &article, nil
// }