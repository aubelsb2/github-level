package main

import (
	"bytes"
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/arran4/github-level"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	ctx := context.Background()

	stats := github_level.GetStats(ctx)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: os.Getenv("GITHUB_TOKEN"),
			TokenType:   "bearer",
		},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	githubUser := os.Getenv("GITHUB_REPOSITORY_OWNER")

	user, _, err := client.Users.Get(ctx, githubUser)
	if err != nil {
		log.Panicf("Error: %v", err)
	}

	masterReadmeContents, _, err := client.Repositories.GetReadme(ctx, githubUser, "github-level", &github.RepositoryContentGetOptions{})
	if err != nil {
		log.Panicf("Readme get fail: %v", err)
	}
	if masterReadmeContents.Content == nil {
		log.Panicf("Readme was nil: %v", err)
	}
	c, err := masterReadmeContents.GetContent()
	if err != nil {
		log.Panicf("Error %v", err)
	}
	presha := sha1.Sum([]byte(c))
	c = ReplaceContent(stats, c)
	postsha := sha1.Sum([]byte(c))
	email := user.GetEmail()
	if len(email) == 0 {
		email = fmt.Sprintf("%s@github.com", user.GetName())
	}
	if presha != postsha {
		_, _, err = client.Repositories.CreateFile(ctx, githubUser, "github-level", masterReadmeContents.GetPath(), &github.RepositoryContentFileOptions{
			Message:   github.String(fmt.Sprintf("Github Level Update: Now %v!", stats.V1().Calculate())),
			Content:   []byte(c),
			SHA:       github.String(masterReadmeContents.GetSHA()),
			Branch:    github.String("main"),
			Committer: &github.CommitAuthor{Name: github.String("Automated " + github_level.PS(user.Name)), Email: &email},
		})
		if err != nil {
			log.Printf("Presha %v postsha %v", presha, postsha)
			log.Printf("Master read me: %v %v %v %v", masterReadmeContents.GetPath(), masterReadmeContents.GetSHA(), masterReadmeContents.GetType(), MustStr(masterReadmeContents.GetContent()))
			log.Panicf("Error creating/updating readme: %v", err)
		}
	} else {
		log.Printf("Presha %v matches post sha %v skipping", presha, postsha)
	}
	if stats.SelfNamedRepo {
		log.Printf("Running in self made profile.")
		RunInSelfNamedRepo(ctx, client, stats, user, githubUser)
	}

}

func MustStr(content string, err error) string {
	return content
}

func RunInSelfNamedRepo(ctx context.Context, client *github.Client, stats *github_level.Stats, user *github.User, githubUser string) {
	userGht := os.Getenv("USER_GITHUB_TOKEN")
	if len(userGht) > 0 {
		log.Printf("Found provided user github token using")
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{
				AccessToken: userGht,
				TokenType:   "bearer",
			},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	} else {
		log.Print("### NOTICE ### If you want this to update your self named repo you will need to create a secret: USER_GITHUB_TOKEN with a github token.")
		return
	}

	masterReadmeContents, _, err := client.Repositories.GetReadme(ctx, githubUser, githubUser, &github.RepositoryContentGetOptions{})
	if err != nil {
		log.Panicf("Readme user get fail: %v", err)
	}
	if masterReadmeContents.Content == nil {
		log.Panicf("Readme user was nil: %v", err)
	}
	c, err := masterReadmeContents.GetContent()
	if err != nil {
		log.Panicf("Error %v", err)
	}
	presha := sha1.Sum([]byte(c))
	c = ReplaceContent(stats, c)
	postsha := sha1.Sum([]byte(c))
	branch := fmt.Sprintf("githublevel-%s", time.Now().Format("D20060102T1504"))
	if presha != postsha {
		reposit, _, err := client.Repositories.Get(ctx, githubUser, githubUser)
		if err != nil {
			log.Panicf("Error getting repo: %v", err)
		}
		mainHeadRef, _, err := client.Git.GetRef(ctx, githubUser, githubUser, "heads/"+reposit.GetDefaultBranch())
		if err != nil {
			log.Panicf("Error getting default branch ref: %v", err)
		}
		_, _, err = client.Git.CreateRef(ctx, githubUser, githubUser, &github.Reference{
			Ref:    github.String("refs/heads/" + branch),
			Object: mainHeadRef.Object,
		})
		if err != nil {
			log.Panicf("Error user creating/updating readme: %v", err)
		}
		_, _, err = client.Repositories.CreateFile(ctx, githubUser, githubUser, masterReadmeContents.GetPath(), &github.RepositoryContentFileOptions{
			Message:   github.String(fmt.Sprintf("Github Level Update: Now %v!", stats.V1().Calculate())),
			Content:   []byte(c),
			SHA:       github.String(masterReadmeContents.GetSHA()),
			Branch:    github.String(branch),
			Committer: &github.CommitAuthor{Name: github.String("Automated " + github_level.PS(user.Name)), Email: user.Email},
		})
		if err != nil {
			log.Panicf("Error user creating/updating readme: %v", err)
		}
		_, _, err = client.PullRequests.Create(ctx, githubUser, githubUser, &github.NewPullRequest{
			Title: github.String(fmt.Sprintf("Github Level (V1): %v", stats.V1().Calculate())),
			Head:  github.String(branch),
			Base:  github.String(reposit.GetDefaultBranch()),
			Body:  github.String(fmt.Sprintf("Please accept: Github Level (V1): %v", stats.V1().Calculate())),
		})
		if err != nil {
			log.Panicf("Error creating pr: %v", err)
		}

	} else {
		log.Printf("Presha %v matches post sha %v skipping", presha, postsha)
	}

}

func ReplaceContent(stats *github_level.Stats, c string) string {
	mrc := bytes.SplitAfter([]byte(c), []byte("\n"))

	newlnchar := "\n"

	shieldLines := make([]string, 0, 1)
	for _, l := range []github_level.Level{
		stats.V1(),
	} {
		shieldLines = append(shieldLines, l.Shield()+newlnchar)
	}

	nrc := make([]string, 0, len(mrc)+2)
	for _, lineB := range mrc {
		line := string(lineB)
		if strings.Contains(line, "id=\"githubLevelId\"") {
			nrc = append(nrc, shieldLines...)
			shieldLines = make([]string, 0, 0)
		} else {
			nrc = append(nrc, line)
		}
	}
	nrc = append(nrc, shieldLines...)
	return strings.Join(nrc, "")
}
