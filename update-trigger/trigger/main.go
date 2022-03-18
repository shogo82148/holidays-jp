package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

var githubToken string

type Input struct {
	Owner         string         `json:"owner"`
	Repo          string         `json:"repo"`
	EventType     string         `json:"event_type"`
	ClientPayload map[string]any `json:"client_payload,omitempty"`
}

func handler(ctx context.Context, input *Input) {
	var err error
	if githubToken == "" {
		githubToken, err = loadGitHubToken(ctx)
		if err != nil {
			log.Println("failed to load GITHUB_TOKEN:", err)
			return
		}
	}

	err = triggerWorkflow(ctx, input)
	if err != nil {
		log.Println("failed to trigger the workflow:", err)
		return
	}
}

func loadGitHubToken(ctx context.Context) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", err
	}
	svc := ssm.NewFromConfig(cfg)

	token, err := svc.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(os.Getenv("GITHUB_TOKEN")),
		WithDecryption: true,
	})
	if err != nil {
		return "", err
	}
	return aws.ToString(token.Parameter.Value), nil
}

type DispatchRequest struct {
	EventType     string         `json:"event_type"`
	ClientPayload map[string]any `json:"client_payload,omitempty"`
}

func triggerWorkflow(ctx context.Context, input *Input) error {
	body, err := json.Marshal(&DispatchRequest{
		EventType:     input.EventType,
		ClientPayload: input.ClientPayload,
	})
	if err != nil {
		return err
	}

	// https://docs.github.com/en/rest/reference/repos#create-a-repository-dispatch-event
	u := fmt.Sprintf("%s/repos/%s/%s/dispatches", os.Getenv("GITHUB_API"), input.Owner, input.Repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", "token "+githubToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println(string(data))

	return nil
}

func main() {
	lambda.Start(handler)
}
