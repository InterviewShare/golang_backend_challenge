package main

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"websynergy/article-service/models"
)

func TestAppConfig_readJSON(t *testing.T) {
	res := httptest.NewRecorder()
	var jsonStr = []byte(`{
		"title": "Title 1",
		"content": "etcetc",
		"author": "author 1"
	}`)
	req, err := http.NewRequest("POST", "/articles", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	var i models.Article
	type args struct {
		w    http.ResponseWriter
		r    *http.Request
		data any
	}
	tests := []struct {
		name    string
		app     *AppConfig
		args    args
		wantErr bool
	}{
		{"test", testApp, args{w: res, r: req, data: &i}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.app.readJSON(tt.args.w, tt.args.r, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("AppConfig.readJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if i.Title != "Title 1" || i.Content != "etcetc" || i.Author != "author 1" {
				t.Errorf("readJSON() = %v, want %v", i, jsonStr)
			}
		})
	}
}

func TestAppConfig_writeJSON(t *testing.T) {
	res := httptest.NewRecorder()
	responseBody := jsonResponse{
		Status:   http.StatusAccepted,
		Message: "check out this message",
		Data:    "some data is here",
	}
	type args struct {
		w       http.ResponseWriter
		status  int
		data    any
		headers []http.Header
	}
	tests := []struct {
		name    string
		app     *AppConfig
		args    args
		wantErr bool
	}{
		{"test", testApp, args{w: res, status: http.StatusAccepted, data: responseBody, headers: []http.Header{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.app.writeJSON(tt.args.w, tt.args.status, tt.args.data, tt.args.headers...); (err != nil) != tt.wantErr {
				t.Errorf("AppConfig.writeJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAppConfig_errorJSON(t *testing.T) {
	res := httptest.NewRecorder()
	type args struct {
		w      http.ResponseWriter
		err    error
		status []int
	}
	tests := []struct {
		name    string
		app     *AppConfig
		args    args
		wantErr bool
	}{
		{"test", testApp, args{w: res, err: errors.New("some error"), status: []int{http.StatusBadRequest}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.app.errorJSON(tt.args.w, tt.args.err, tt.args.status...); (err != nil) != tt.wantErr {
				t.Errorf("AppConfig.errorJSON() error = %v,wantErr %v", err, tt.wantErr)
			}
		})
	}
}
