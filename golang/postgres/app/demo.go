package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"postgres/website"
)

func RunRepositoryDemo(ctx context.Context, websiteRepository website.Repository) {
	fmt.Println("1. Migrate")
	if err := websiteRepository.Migrate(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("2. Create")
	gosamples := website.Website{
		Name: "Go Samples",
		URL:  "https://golangbyexample.com",
		Rank: 1,
	}
	golang := website.Website{
		Name: "Golang",
		URL:  "https://golang.org",
		Rank: 2,
	}
	createdGoSamples, err := websiteRepository.Create(ctx, gosamples)
	if errors.Is(err, website.ErrDuplicate) {
		fmt.Printf("record: %+v already exists", gosamples)
	} else if err != nil {
		log.Fatal(err)
	}
	createdGolang, err := websiteRepository.Create(ctx, golang)
	if errors.Is(err, website.ErrDuplicate) {
		fmt.Printf("record: %+v already exists", golang)
	} else if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n%+v\n", createdGoSamples, createdGolang)

	fmt.Println("3. Get by name")
	gotGoSamples, err := websiteRepository.GetByName(ctx, "GOSAMPLES")
	if errors.Is(err, website.ErrNotExist) {
		fmt.Printf("record GOTSAMLPES does not exist")
	} else if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", gotGoSamples)

	fmt.Println("4. Update")
	updatedItem, err := websiteRepository.Update(ctx, createdGoSamples.ID, website.Website{Name: "Jojo", URL: "https://jojo.com", Rank: 3})
	if errors.Is(err, website.ErrUpdateFailed) {
		fmt.Printf("update failed")
	}
	fmt.Printf("%+v\n", updatedItem)

	fmt.Println("5. Getall")
	all, err := websiteRepository.All(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range all {
		fmt.Printf("%+v\n", item)
	}

	fmt.Println("6. Delete")
	if err := websiteRepository.Delete(ctx, createdGolang.ID); errors.Is(err, website.ErrDeleteFailed) {
		fmt.Printf("delete failed")
	}
}
