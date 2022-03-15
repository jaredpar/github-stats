package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v43/github"
)

func CreateClient(ctx context.Context) (*github.Client, error) {
	token, err := os.ReadFile(`c:\users\jaredpar\.token`)
	if err != nil {
		return nil, err
	}

	tokenStr := string(token)
	/*
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: tokenStr},
		)
		tc := oauth2.NewClient(ctx, ts)
	*/

	transport := github.BasicAuthTransport{
		Username: "jaredpar",
		Password: strings.TrimSuffix(tokenStr, "\r\n"),
	}
	client := github.NewClient(transport.Client())
	zen, _, _ := client.Zen(ctx)
	fmt.Println(zen)
	return client, nil
}

func main() {
	ctx := context.Background()
	client, err := CreateClient(ctx)
	if err != nil {
		panic(err)
	}

	org := "dotnet"
	repo := "sdk"
	getApprovals := func(number int, headSha string) (totalCount, lastPushCount int) {
		reviews, _, _ := client.PullRequests.ListReviews(ctx, org, repo, number, &github.ListOptions{})
		for _, review := range reviews {
			if review.State != nil && *review.State == "APPROVED" {
				totalCount++
				if review.CommitID != nil && *review.CommitID == headSha {
					lastPushCount++
				}
			}
		}

		return
	}

	twoOrMore := 0
	exactlyOne := 0
	lastPushTwoOrMore := 0
	lastPushExactlyOne := 0
	page := 0
	prCount := 0
	for prCount < 1000 {
		options := github.PullRequestListOptions{
			State:     "all",
			Direction: "desc",
		}
		options.Page = page

		prs, _, err := client.PullRequests.List(ctx, org, repo, &options)
		if err != nil {
			panic(err)
		}

		for _, pr := range prs {
			if pr.MergedAt != nil {
				headSha := pr.Head.SHA
				count, lastPushCount := getApprovals(*pr.Number, *headSha)
				fmt.Printf("%d %d %s\n", count, lastPushCount, *pr.HTMLURL)
				if count >= 2 {
					twoOrMore++
				} else if count == 1 {
					exactlyOne++
				}

				if lastPushCount >= 2 {
					lastPushTwoOrMore++
				} else if lastPushCount == 1 {
					lastPushExactlyOne++
				}

				prCount++
			}
		}

		page++
		fmt.Printf("2+ approvals: %d\n", twoOrMore)
		fmt.Printf("2+ approvals last push: %d\n", lastPushTwoOrMore)
		fmt.Printf("1 approval: %d\n", exactlyOne)
		fmt.Printf("1 approval last push: %d\n", lastPushExactlyOne)
		fmt.Printf("Total PRs: %d\n", prCount)
	}
}
