package main

import (
	"context"
	"log"
	"postgres/app"
	"postgres/website"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	dbpool, err := pgxpool.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	websiteRepository := website.NewPostgreSQLPGXRepository(dbpool)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	app.RunRepositoryDemo(ctx, websiteRepository)
}
