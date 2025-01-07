package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/google/go-github/v68/github"
	"golang.org/x/oauth2"
)

func main() {
	// Get the GitHub token from the environment variable
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN is not set.")
	}
	// Get the owner, repo, and PR number from the environment variables
	owner := os.Getenv("GITHUB_REPOSITORY_OWNER")
	repo := os.Getenv("GITHUB_REPOSITORY_NAME")
	prNumber := os.Getenv("GITHUB_PR_NUMBER")

    // Create a new GitHub client
    ctx := context.Background()
    ts := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: token},
    )
    tc := oauth2.NewClient(ctx, ts)
    client := github.NewClient(tc)

	// Poll until all checks are finished
	timeout := time.After(10 * time.Minute)
	tick := time.Tick(10 * time.Second)

	for {
		select {
		case <-timeout:
			log.Fatal("Timed out waiting for checks to finish.")
		case <-tick:
			// Retrieve the check runs for the PR
			checks, _, err := client.Checks.ListCheckRunsForRef(ctx, owner, repo, prNumber, &github.ListCheckRunsOptions{})
			if err != nil {
				log.Fatalf("Error getting check runs: %v", err)
			}

			// Check if all checks are complete and successful
			allChecksPassed := true
			for _, check := range checks.CheckRuns {
				if check.Status != nil && *check.Status != "completed" {
					allChecksPassed = false
					break
				}
				if check.Conclusion != nil && *check.Conclusion != "success" {
					allChecksPassed = false
					break
				}
			}

			// If all checks are complete and successful, exit successfully
			if allChecksPassed {
				log.Println("All checks have passed!")
				return
			}
		}
	}
}
