package github_level

import (
	"fmt"
	"github.com/fogleman/gg"
	"image"
	"image/color"
	"math"
	"os"
)

type Level interface {
	Calculate() int
	Version() int
	Shield() string
	//Badge() image.Image
}

var (
	blue = color.RGBA{
		R: 0,
		G: 0,
		B: 255,
		A: 255,
	}
	green = color.RGBA{
		R: 0,
		G: 255,
		B: 0,
		A: 255,
	}
	gold = color.RGBA{
		R: 222,
		G: 200,
		B: 55,
		A: 255,
	}
)

type GithubLevelV1 Stats

func (l GithubLevelV1) Calculate() (r int) {
	for _, v := range []float64{
		float64(l.PublicGists),
		float64(l.PublicRepos),
		float64(l.Following),
		float64(l.Followers),
		float64(l.ArchivedRepos),
		float64(l.LargestForkCount),
		float64(l.LargestStargazerCount),
		float64(l.LargestWatcherCount),
		float64(l.LicenseSophistication),
		float64(l.UserCreatedDate.Unix() / 60 / 60 / 24 / 365),
		float64(l.FirstPublicRepoDate.Unix() / 60 / 60 / 24 / 365),
		float64(len(l.Licenses)),
	} {
		if v > 0 {
			r += int(math.Round(v * math.Log10(v+1)))
		}
	}
	for _, v := range []float64{
		float64(l.PRRequests),
		float64(l.IssuesLogged),
	} {
		if v > 0 {
			r += int(math.Round(v * math.Log2(v+1)))
		}
	}
	if l.TotalReposWithUpdated > 0 {
		x := (float64(l.SumLastUpdatedDays) / float64(l.TotalReposWithUpdated)) / 60 / 60 / 24 / 365
		v := int(math.Round(x * math.Log10(x)))
		if v < r {
			r -= v
		} else {
			r = 1
		}
	}
	if l.SelfNamedRepo {
		r += 1
	}

	return
}

func (GithubLevelV1) Version() int {
	return 1
}

func (l GithubLevelV1) Shield() string {
	v := l.Calculate()
	githubUser := os.Getenv("GITHUB_REPOSITORY_OWNER")
	url := fmt.Sprintf("https://github.com/%s/github-level", githubUser)
	return fmt.Sprintf(`<a id="githubLevelId" href="%s"> <img src="https://img.shields.io/badge/%s%%20version%v-%v-yellowgreen" alt="Github level %v"/></a>`, url, "Github Level",
		l.Version(), v, v)
}

func (GithubLevelV1) Badge() image.Image {
	// Never minds found https://github.com/badges/shields
	dc := gg.NewContext(200, 50)
	drawBadgeBack(dc, 200, 50)
	return dc.Image()
}

func drawBadgeBack(dc *gg.Context, width float64, height float64) {
	dc.Push()
	dc.DrawRoundedRectangle(10, 10, width-20, height-20, 10)
	dc.SetColor(green)
	dc.FillPreserve()
	dc.ClearPath()
	dc.Pop()
}
