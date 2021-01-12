package main

import (
	"context"
	"github.com/google/go-github/v33/github"
	"log"
	"os"
)

func main() {
	client := github.NewClient(nil)

	// list all organizations for user "willnorris"
	ctx := context.Background()
	repositories, _, err := client.Repositories.List(ctx, os.ExpandEnv("${GITHUB_REPOSITORY_OWNER}"), nil)

	if err != nil {
		log.Panicf("Error: %v", err)
	}

	for _, repository := range repositories {
		log.Printf("%s", PS(repository.Name))
	}
}

func PS(name *string) string {
	if name != nil {
		return *name
	}
	return ""
}
