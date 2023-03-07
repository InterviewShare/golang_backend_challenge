package main

import (
	"errors"
	"net/http"
	"websynergy/article-service/models"

	"github.com/gorilla/mux"
)

// handles the request to GET an article
func (app *AppConfig) GetArticle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		app.errorJSON(w, errors.New("id is missing in parameters"), http.StatusBadRequest)
		return
	}
	article, err := app.articleService.GetArticleById(ctx, id)
	if err != nil {
		app.errorJSON(w, err, http.StatusNotFound)
		return
	}

	responseBody := jsonResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    *article,
	}
	app.writeJSON(w, http.StatusOK, responseBody)
}

// handles the request to GET article list
func (app *AppConfig) GetArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := app.articleService.GetArticles(r.Context())
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	responseBody := jsonResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    articles,
	}
	app.writeJSON(w, http.StatusOK, responseBody)
}

// handles the request to create an article
func (app *AppConfig) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var itemToCreate models.Article
	err := app.readJSON(w, r, &itemToCreate)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	data, err := app.articleService.CreateArticle(r.Context(), itemToCreate)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	responseBody := jsonResponse{
		Status:  http.StatusCreated,
		Message: "Success",
		Data:    *data,
	}
	app.writeJSON(w, http.StatusCreated, responseBody)
}
