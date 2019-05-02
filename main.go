package main

import (
	"database/sql"
	"github.com/denysvitali/tsp-node-server/server"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

var db *sql.DB

func main() {
	connStr := "postgres://postgres:postgres_tsp@172.17.0.2/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	var state server.State
	state.Db = db


	http.HandleFunc("/api/v1/upload", state.UploadJson)
	_ = http.ListenAndServe(":12538", nil)
}