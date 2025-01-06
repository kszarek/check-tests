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

    // Create a status for the commit
    createStatus(ctx, client, owner, repo, sha, "pending", "Follow-On Test Status", "The follow-on tests have started")

    // Simulate some work
    log.Println("The follow-on works!!")
		time.Sleep(5 * time.Second)

    // Update the status to success
    createStatus(ctx, client, owner, repo, sha, "success", "Follow-On Test Status", "The follow-on tests completed successfully")
}

func createStatus(ctx context.Context, client *github.Client, owner, repo, sha, state, context, description string) {
    status := &github.RepoStatus{
        State:       github.Ptr(state),
        Context:     github.Ptr(context),
        Description: github.Ptr(description),
    }

		log.Printf("Creating status: state=%s, context=%s, description=%s", *status.State, *status.Context, *status.Description)
		repoStatus, respose, err := client.Repositories.CreateStatus(ctx, owner, repo, sha, status)
		if err != nil {
				log.Fatalf("Error creating status: %v", err)
		}
		log.Printf("Response: %+v", respose)
		log.Printf("Status response: %+v", repoStatus)
		log.Printf("Status created successfully")
}
