package github

import (
	"bufio"
	"fmt"
	"gofr.dev/pkg/gofr"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

func getGitHubToken() string {
	fmt.Println("generating response....")
	// Check for the token in the environment (GITHUB_TOKEN)
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		return token
	}

	// Prompt the user for the token if not found
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your GitHub token: ")
	token, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading token: %v", err)
	}

	return strings.TrimSpace(token)
}

func RaiseIssue(ctx *gofr.Context) (interface{}, error) {
	// Get GitHub token
	token := getGitHubToken()

	// Repository details
	owner := "Zop-Stars"
	repo := "gofr"

	reader := bufio.NewReader(os.Stdin)
	// Issue details
	fmt.Println("Please enter your Issue Title: ")
	title, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading title: %v", err)
	}
	fmt.Println("Please enter your Issue Description: ")
	body, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading description: %v", err)
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Create the issue
	issueRequest := &github.IssueRequest{
		Title: &title,
		Body:  &body,
	}

	issue, _, err := client.Issues.Create(ctx, owner, repo, issueRequest)
	if err != nil {
		log.Fatalf("Error creating issue: %v", err)
	}

	// Print the created issue URL
	fmt.Printf("Issue created successfully: %s\n", issue.GetHTMLURL())

	return nil, nil
}

func GetIssue(ctx *gofr.Context) (interface{}, error) {
	// Get GitHub token
	token := getGitHubToken()

	// Repository details
	owner := "gofr-dev"
	repo := "gofr"

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	// List all issues
	opts := &github.IssueListByRepoOptions{
		State:       "open", // "open", "closed", or "all"
		ListOptions: github.ListOptions{PerPage: 10},
	}

	for {
		issues, resp, err := client.Issues.ListByRepo(ctx, owner, repo, opts)
		if err != nil {
			log.Fatalf("Error fetching issues: %v", err)
		}

		// Print the issues
		for _, issue := range issues {
			fmt.Printf("#%d: %s (link - %s)\n", *issue.Number, *issue.Title, *issue.URL)
		}

		// Check if there are more pages
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return nil, nil
}

func GetRelease(ctx *gofr.Context) (interface{}, error) {
	// Get GitHub token
	token := getGitHubToken()

	// Repository details
	owner := "gofr-dev"
	repo := "gofr"

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	// List all releases
	opts := &github.ListOptions{
		Page:    1,
		PerPage: 10,
	}

	for {
		releases, resp, err := client.Repositories.ListReleases(ctx, owner, repo, opts)
		if err != nil {
			log.Fatalf("Error fetching releases: %v", err)
		}

		// Print the releases
		for _, release := range releases {
			fmt.Printf("%s: \n%s \n(link - %s)\n", *release.Name, *release.Body, *release.URL)
		}

		// Check if there are more pages
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return nil, nil
}

func GetReleaseLatest(ctx *gofr.Context) (interface{}, error) {
	// Get GitHub token
	token := getGitHubToken()

	// Repository details
	owner := "gofr-dev"
	repo := "gofr"

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	release, _, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		log.Fatalf("Error fetching releases: %v", err)
	}

	// Print the release

	fmt.Printf("%s: \n%s \n(link - %s)\n", *release.Name, *release.Body, *release.URL)

	return nil, nil
}
