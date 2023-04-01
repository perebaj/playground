package app

import (
	"context"
	"fmt"
	"log"
	"postgres/website"
)

func RunRepositoryDemo(ctx context.Context, websiteRepository website.Repository) {
	fmt.Println("1. Migrate")
	if err := websiteRepository.Migrate(ctx); err != nil {
		log.Fatal(err)
	}
}
