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

	getApprovals := func(number int) int {
		count := 0
		reviews, _, _ := client.PullRequests.ListReviews(ctx, "dotnet", "runtime", number, &github.ListOptions{})
		for _, review := range reviews {
			if review.State != nil && *review.State == "APPROVED" {
				count++
			}
		}

		return count
	}

	twoOrMore := 0
	oneOrLess := 0
	page := 0
	for twoOrMore+oneOrLess < 100 {
		options := github.PullRequestListOptions{
			State:     "all",
			Direction: "desc",
		}
		options.Page = page

		prs, _, err := client.PullRequests.List(ctx, "dotnet", "runtime", &options)
		if err != nil {
			panic(err)
		}

		for _, pr := range prs {
			if pr.MergedAt != nil {
				count := getApprovals(*pr.Number)
				fmt.Printf("%d %s\n", count, *pr.HTMLURL)
				if count >= 2 {
					twoOrMore++
				} else {
					oneOrLess++
				}
			}
		}

		page++
	}

	fmt.Printf("2 or more approvals: %d\n", twoOrMore)
	fmt.Printf("1 or less approvals: %d\n", oneOrLess)
	fmt.Printf("Total PRs: %d\n", twoOrMore+oneOrLess)
}
