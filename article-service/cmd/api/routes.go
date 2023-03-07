package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// sets up the routes for http server
func (app *AppConfig) routes() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/articles", app.GetArticles).Methods("GET")
	router.HandleFunc("/articles", app.CreateArticle).Methods("POST")
	router.HandleFunc("/articles/{id}", app.GetArticle).Methods("GET")

	return router
}