package github_level

import (
	"context"
	"github.com/google/go-github/github"
	"log"
	"os"
)

type Stats struct {
	PublicRepos int
	Licenses    map[string]int
}

func GetStats(ctx context.Context) (stats *Stats) {
	stats = &Stats{
		Licenses: map[string]int{},
	}
	client := github.NewClient(nil)

	repositories, _, err := client.Repositories.List(ctx, os.ExpandEnv("${GITHUB_REPOSITORY_OWNER}"), &github.RepositoryListOptions{
		Visibility: "public",
	})

	if err != nil {
		log.Panicf("Error: %v", err)
	}

	for _, repository := range repositories {
		if PB(repository.Private) {
			log.Printf("Ignoring privat repo: %s", PS(repository.Name))
			continue
		}
		stats.PublicRepos++
		stats.Licenses[PS(repository.LicenseTemplate)]++
	}
	return
}

func PB(private *bool) bool {
	return private != nil && *private
}

func PS(name *string) string {
	if name != nil {
		return *name
	}
	return ""
}
