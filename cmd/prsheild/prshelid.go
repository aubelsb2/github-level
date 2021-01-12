package main

import (
	"context"
	github_level "github.com/arran4/github-level"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	ctx := context.Background()

	stats := github_level.GetStats(ctx)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	githubUser := os.Getenv("GITHUB_REPOSITORY_OWNER")

	user, _, err := client.Users.Get(ctx, githubUser)
	if err != nil {
		log.Panicf("Error: %v", err)
	}

	masterReadmeContents, _, _, err := client.Repositories.GetContents(ctx, githubUser, "github-level", "README.md", &github.RepositoryContentGetOptions{})
	if err != nil {
		log.Panicf("Readme get fail: %v", err)
	}
	if masterReadmeContents.Content == nil {
		log.Panicf("Readme was nil: %v", err)
	}
	c := ReplaceContent(stats, *masterReadmeContents.Content)
	_, _, err = client.Repositories.CreateFile(ctx, githubUser, "github-level", "README.md", &github.RepositoryContentFileOptions{
		Message:   github.String("Version Update!"),
		Content:   []byte(c),
		SHA:       masterReadmeContents.SHA,
		Branch:    github.String("main"),
		Committer: &github.CommitAuthor{Name: github.String("Automated " + github_level.PS(user.Name)), Email: user.Email},
	})
	if err != nil {
		log.Panicf("Error creating/updating readme: %v", err)
	}

	if stats.SelfNamedRepo {
		//selfNamedRepo, _, err := client.Repositories.Get(ctx, githubUser, githubUser)
		//if err != nil {
		//	log.Panicf("Error: %v", err)
		//}

	}

}

func ReplaceContent(stats *github_level.Stats, c string) string {
	shieldLines := make([]string, 0, 1)
	for _, l := range []github_level.Level{
		stats.V1(),
	} {
		shieldLines = append(shieldLines, l.Shield())
	}

	mrc := regexp.MustCompile("\r\n|\r|\n").Split(c, -1)
	nrc := make([]string, 0, len(mrc)+2)
	for _, line := range mrc {
		if strings.Contains(line, "id=\"githubLevelId\"") {
			nrc = append(nrc, shieldLines...)
			shieldLines = make([]string, 0, 0)
		} else {
			nrc = append(nrc, line)
		}
	}
	nrc = append(nrc, shieldLines...)
	return strings.Join(nrc, "\n")
}
