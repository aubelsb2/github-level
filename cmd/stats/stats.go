package main

import (
	"context"
	"encoding/json"
	"github.com/arran4/github-level"
	"os"
)

func main() {
	ctx := context.Background()

	stats := github_level.GetStats(ctx)

	j := json.NewEncoder(os.Stdout)
	j.SetIndent("", "  ")
	_ = j.Encode(stats)
}
