package github_level

import "math"

type Level interface {
	Calculate() int
	Version() int
	//Badge() image.Image
}

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
		float64(l.UserCreatedDate.Unix() / 60 / 60 / 24 / 30 / 12),
		float64(l.FirstPublicRepoDate.Unix() / 60 / 60 / 24 / 30 / 12),
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
		x := float64(l.SumLastUpdatedDays) / float64(l.TotalReposWithUpdated)
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

//func (GithubLevelV1) Badge() image.Image {
//
//}
