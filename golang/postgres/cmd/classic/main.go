package main

import (
	"context"
	"database/sql"
	"log"
	"postgres/app"
	"postgres/website"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	db, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	websiteRepository := website.NewPostgreSQLClassicRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app.RunRepositoryDemo(ctx, websiteRepository)
}
