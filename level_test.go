package github_level

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"testing"
	"time"
)

func TestL1(t *testing.T) {
	for _, each := range []struct {
		Name  string
		Level GithubLevelV1
	}{
		{
			Name: "arran4",
			Level: GithubLevelV1{
				SumLastUpdatedDays:    406628,
				TotalReposWithUpdated: 23,
				PublicRepos:           30,
				ForkedRepos:           17,
				ArchivedRepos:         7,
				LicenseSophistication: 0,
				Licenses: map[string]int{
					"": 30,
				},
				FirstPublicRepoDate:   TimeMustParse(time.Parse(time.RFC3339, "2010-01-31T04:07:52Z")),
				LargestForkCount:      22,
				LargestStargazerCount: 49,
				LargestWatcherCount:   49,
				SelfNamedRepo:         false,
				UserCreatedDate:       TimeMustParse(time.Parse(time.RFC3339, "2009-08-04T03:16:20Z")),
				Followers:             17,
				Following:             34,
				OwnedPrivateRepos:     0,
				PublicGists:           62,
				PRRequests:            3,
				IssuesLogged:          6,
			},
		},
	} {
		t.Run(each.Name, func(t *testing.T) {
			v := each.Level.Calculate()
			if v > 1000 || v < 10 {
				t.Fatalf("Got level 1 level: %v", v)
			}
			t.Logf("Got level 1 level: %v", v)
			i := each.Level.Badge()
			w, err := os.Create(fmt.Sprintf("badge-version-%v.png", each.Level.Version()))
			if err != nil {
				t.Fatalf("OS Create error: %v", err)
			}
			if err := png.Encode(w, i); err != nil {
				t.Fatalf("PNG encode error: %v", err)
			}
			if err := w.Close(); err != nil {
				t.Fatalf("Writer close error: %v", err)
			}
		})
	}
}

func TimeMustParse(parse time.Time, err error) time.Time {
	log.Printf("ERROR: %v", err)
	return parse
}
