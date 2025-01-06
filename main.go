package main

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/v68/github"
	"golang.org/x/oauth2"
)

func main() {
    // Get environment variables
    owner := os.Getenv("GITHUB_REPOSITORY_OWNER")
    repo := os.Getenv("GITHUB_REPOSITORY_NAME")
    sha := os.Getenv("GITHUB_SHA")
    token := os.Getenv("GITHUB_TOKEN")

    // Create a new GitHub client
    ctx := context.Background()
    ts := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: token},
    )
    tc := oauth2.NewClient(ctx, ts)
    client := github.NewClient(tc)

    // Create a check run
    createCheckRun(ctx, client, owner, repo, sha, "in_progress", "Follow-On Test Status", "The follow-on tests have started")

    // Simulate some work
    log.Println("The follow-on works!!")

    // Update the check run to success
    createCheckRun(ctx, client, owner, repo, sha, "completed", "Follow-On Test Status", "The follow-on tests completed successfully")
}

func createCheckRun(ctx context.Context, client *github.Client, owner, repo, sha, status, name, summary string) {
		checkRun := &github.CreateCheckRunOptions{
			Name:    name,
			HeadSHA: sha,
			Status:  github.Ptr(status),
			Output: &github.CheckRunOutput{
				Title:   github.Ptr(name),
				Summary: github.Ptr(summary),
			},
		}

    if status == "completed" {
        checkRun.Conclusion = github.Ptr("success")
    }

		checkRunResult, resp, err := client.Checks.CreateCheckRun(ctx, owner, repo, *checkRun)
		if err != nil {
			log.Fatalf("Error creating check run: %v", err)
		}
		log.Printf("Check run created: ID=%d, Status=%s, Response Status=%s",
			checkRunResult.GetID(),
			checkRunResult.GetStatus(),
			resp.Status)
}
