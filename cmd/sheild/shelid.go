package main

import (
	"context"
	"github.com/arran4/github-level"
	"log"
)

func main() {
	ctx := context.Background()

	stats := github_level.GetStats(ctx)

	for _, l := range []github_level.Level{
		stats.V1(),
	} {
		log.Printf("%s", l.Shield())
	}

}
