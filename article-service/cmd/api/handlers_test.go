package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"websynergy/article-service/services"

	"github.com/gorilla/mux"
)

var testDB *sql.DB
var testApp *AppConfig

// This gets run before your actual test functions do.
func init() {
	var err error
	testDB, err = sql.Open("mysql", "root:somepass@tcp(localhost:3306)/test_db")
	if err != nil {
		fmt.Printf("test db init failed: %s", err)
	}
	testDB.Exec(`IF EXISTS(SELECT *
		FROM   articles)
		DROP TABLE articles;`)

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

	testApp = &AppConfig{dbClient: testDB, articleService: services.NewArticleService(testDB)}
}

func TestAppConfig_GetArticle(t *testing.T) {
	res1 := httptest.NewRecorder()
	req1, err := http.NewRequest("GET", "/articles", nil)
	if err != nil {
		t.Fatal(err)
	}

	res2 := httptest.NewRecorder()
	req1 = mux.SetURLVars(req1, map[string]string{
		"id": "1",
	})
	req2, err := http.NewRequest("GET", "/articles", nil)
	if err != nil {
		t.Fatal(err)
	}

	req2 = mux.SetURLVars(req2, map[string]string{
		"id": "12",
	})
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		req        *http.Request
		res        *httptest.ResponseRecorder
		name       string
		app        *AppConfig
		args       args
		wantStatus any
		want       any
	}{
		// TODO: Add test cases.
		{req1, res1, "pass_test", testApp, args{w: res1, r: req1}, http.StatusOK,
			`{"status":200,"message":"Success","data":{"id":"1","title":"Title 1","content":"lorem ipsum some content etcetc","author":"author 1"}}`,
		},
		{req2, res2, "error_test", testApp, args{w: res2, r: req2}, http.StatusNotFound,
			`{"status":404,"message":"Error: sql: no rows in result set","data":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.app.GetArticle(tt.args.w, tt.args.r)
			if status := tt.res.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantStatus)
			}
			if !reflect.DeepEqual(tt.res.Body.String(), tt.want) {
				t.Errorf("GetArticle() = %v, want %v", tt.res.Body.String(), tt.want)
			}
		})
	}
}

func TestAppConfig_GetArticles(t *testing.T) {
	res := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/articles", nil)
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name       string
		app        *AppConfig
		args       args
		wantStatus any
		want       any
	}{
		{"pass_test", testApp, args{w: res, r: req}, http.StatusOK,
			`{"status":200,"message":"Success","data":[{"id":"1","title":"Title 1","content":"lorem ipsum some content etcetc","author":"author 1"},{"id":"2","title":"Title 2","content":"lorem ipsum some content etcetc","author":"author 2"},{"id":"3","title":"Title 3","content":"lorem ipsum some content etcetc","author":"author 3"}]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.app.GetArticles(tt.args.w, tt.args.r)
			if status := res.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantStatus)
			}
		})
	}
}

func TestAppConfig_CreateArticle(t *testing.T) {
	res := httptest.NewRecorder()
	var jsonStr = []byte(`{
		"title": "Title 1",
		"content": "lorem ipsum some content etcetc",
		"author": "author 1"
	}`)
	req, err := http.NewRequest("POST", "/articles", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name       string
		app        *AppConfig
		args       args
		wantStatus any
	}{
		{"pass_test", testApp, args{w: res, r: req}, http.StatusCreated},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.app.CreateArticle(tt.args.w, tt.args.r)
			if status := res.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantStatus)
			}
		})
	}
}
