package services

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"websynergy/article-service/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

var testDB *sql.DB

// This gets run before your actual test functions do.
func init() {
	var err error
	testDB, err = sql.Open("mysql", "root:somepass@tcp(localhost:3306)/test_db")
	if err != nil {
		fmt.Printf("test db init failed: %s", err)
	}
	testDB.Exec(`IF EXISTS(SELECT *
		FROM   dbo.Scores)
		DROP TABLE dbo.Scores;`)
	
	testDB.Exec(`CREATE TABLE articles(
		ID   VARCHAR (40) NOT NULL,
		TITLE VARCHAR (20) NOT NULL,
		CONTENT  VARCHAR (100),
		AUTHOR   VARCHAR (20),       
		PRIMARY KEY (ID)
	 );`)

	testDB.Exec(`INSERT INTO articles VALUES 
		(1, "Title 1", "lorem ipsum some content etcetc", "author 1"),
		(2, "Title 2", "lorem ipsum some content etcetc", "author 2"),
		(3, "Title 3", "lorem ipsum some content etcetc", "author 3");`)
}

func TestNewArticleService(t *testing.T) {
	type args struct {
		dbClient *sql.DB
	}
	tests := []struct {
		name string
		args args
		want *ArticleService
	}{
		// TODO: Add test cases.
		{"pass_test", args{dbClient: testDB}, NewArticleService(testDB)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewArticleService(tt.args.dbClient); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewArticleService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArticleService_GetArticleById(t *testing.T) {
	type args struct {
		ctx context.Context
		ID  string
	}
	tests := []struct {
		name    string
		s       *ArticleService
		args    args
		want    *models.Article
		wantErr bool
	}{
		// TODO: Add test cases.
		{"pass_test", NewArticleService(testDB), args{ctx: context.Background(), ID: "1"}, &models.Article{
			ID:      "1",
			Title:   "Title 1",
			Content: "lorem ipsum some content etcetc",
			Author:  "author 1"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetArticleById(tt.args.ctx, tt.args.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleService.GetArticleById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArticleService.GetArticleById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArticleService_GetArticles(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		s       *ArticleService
		args    args
		want    []models.Article
		wantErr bool
	}{
		// TODO: Add test cases.
		{"pass_test", NewArticleService(testDB), args{ctx: context.Background()}, []models.Article{
			{
				ID:      "1",
				Title:   "Title 1",
				Content: "lorem ipsum some content etcetc",
				Author:  "author 1",
			},
			{
				ID:      "2",
				Title:   "Title 2",
				Content: "lorem ipsum some content etcetc",
				Author:  "author 2",
			},
			{
				ID:      "3",
				Title:   "Title 3",
				Content: "lorem ipsum some content etcetc",
				Author:  "author 3",
			},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetArticles(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleService.GetArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, e := range tt.want {
				if !slices.Contains(got, e) {
					t.Errorf("ArticleService.GetArticles() response: %v, does not contain %v", got, e)
				}
			}
		})
	}
}

func TestArticleService_CreateArticle(t *testing.T) {
	type args struct {
		ctx context.Context
		a   models.Article
	}
	tests := []struct {
		name    string
		s       *ArticleService
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"pass_test", NewArticleService(testDB), args{ctx: context.Background(), a: models.Article{
			Title:   "Title 4",
			Content: "lorem ipsum some content etcetc",
			Author:  "author 4",
		},}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.CreateArticle(tt.args.ctx, tt.args.a)
			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleService.CreateArticle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if _, err := uuid.Parse(got.Id); err != nil {
				t.Errorf("ArticleService.CreateArticle() = %v, want uuid formatted string", got)
			}
		})
	}
}
