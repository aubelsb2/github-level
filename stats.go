package github_level

import (
	"context"
	"github.com/google/go-github/github"
	"log"
	"os"
	"time"
)

type Stats struct {
	SumLastUpdatedDays    int64
	TotalReposWithUpdated int
	PublicRepos           int
	ForkedRepos           int
	ArchivedRepos         int
	LicenseSophistication int
	Licenses              map[string]int
	FirstPublicRepoDate   time.Time
	LargestForkCount      int
	LargestStargazerCount int
	LargestWatcherCount   int
}

func GetStats(ctx context.Context) (stats *Stats) {
	stats = &Stats{
		Licenses:            map[string]int{},
		FirstPublicRepoDate: time.Now(), // Born yesterday - https://open.spotify.com/track/22pzlAb4SynBW4aO0HCwo1?si=VvB1IcC3Q_6PwS6-pYR9tA
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
		if stats.Licenses[PS(repository.LicenseTemplate)] == 1 {
			if v, ok := LicenseSophistication[PS(repository.LicenseTemplate)]; ok {
				stats.LicenseSophistication += v
			}
		}
		if PB(repository.Fork) {
			stats.ForkedRepos++
		}
		if PB(repository.Archived) {
			stats.ArchivedRepos++
		}
		if repository.CreatedAt != nil && repository.CreatedAt.Before(stats.FirstPublicRepoDate) {
			stats.FirstPublicRepoDate = repository.CreatedAt.Time
		}
		if repository.ForksCount != nil && PI(repository.ForksCount) > stats.LargestForkCount {
			stats.LargestForkCount = *repository.ForksCount
		}
		if repository.StargazersCount != nil && PI(repository.StargazersCount) > stats.LargestStargazerCount {
			stats.LargestStargazerCount = *repository.StargazersCount
		}
		if repository.WatchersCount != nil && PI(repository.WatchersCount) > stats.LargestWatcherCount {
			stats.LargestWatcherCount = *repository.WatchersCount
		}
		if repository.UpdatedAt != nil && !PB(repository.Archived) {
			stats.SumLastUpdatedDays += repository.UpdatedAt.Time.Unix() / 60 / 60 / 24
			stats.TotalReposWithUpdated++
		}
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
func PI(name *int) int {
	if name != nil {
		return *name
	}
	return 0
}
