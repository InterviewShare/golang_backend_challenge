package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"websynergy/article-service/services"

	_ "github.com/go-sql-driver/mysql"
)

const port = "8080"

var tryCount int64

type AppConfig struct {
	dbClient       *sql.DB
	articleService *services.ArticleService
}

func main() {
	log.Println("Starting article service")

	dbConn := connectToDB()
	if dbConn == nil {
		log.Panic("Cannot connect to database.")
	}
	defer dbConn.Close()

	app := AppConfig{
		dbClient:       dbConn,
		articleService: services.NewArticleService(dbConn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", port),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Mysql db not ready yet...")
			tryCount++
		} else {
			log.Println("Connected to MySQL!")
			return connection
		}

		if tryCount > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Waiting for 2 seconds")
		time.Sleep(2 * time.Second)
		continue
	}
}
