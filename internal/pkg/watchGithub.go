package pkg

import (
	"context"
	"log"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func SearchGithub(githubToken string) {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	client := github.NewClient(oauth2.NewClient(ctx, ts))

	query := "AKIA filename:.env OR filename:.ini OR filename:.yml OR filename:.yaml OR filename:.json"
	opt := &github.SearchOptions{Sort: "indexed", Order: "desc"}

	results, _, err := client.Search.Code(ctx, query, opt)
	if err != nil {
		log.Fatalf("Error searching GitHub: %v", err)
	}

	for _, file := range results.CodeResults {
		checkFileContent(ctx, client, &file)
	}
}
