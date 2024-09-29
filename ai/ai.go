package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type AI struct {
	prompt string
	client *openai.Client
}

type Response struct {
	Title     string   `json:"title"`
	Artist    string   `json:"artist"`
	Genres    []string `json:"genres"`
	Comments  string   `json:"comments"`
	Decade    string   `json:"decade"`
	Year      int      `json:"year"`
	Confident bool     `json:"confident"`
}

func New(prompt string) (*AI, error) {
	promptBytes, err := os.ReadFile(prompt)
	if err != nil {
		return nil, err
	}
	promptString := string(promptBytes)

	cfg := openai.DefaultConfig("blah")
	cfg.BaseURL = "http://127.0.0.1:1337/v1"

	client := openai.NewClientWithConfig(cfg)

	return &AI{prompt: promptString, client: client}, nil
}

func (a *AI) TryFile(filename string, hint string) (Response, error) {
	fmt.Println(hint)
	req := openai.ChatCompletionRequest{
		Model: "llama3.1-8b-instruct",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: a.prompt,
			},
		},
	}

	req.Messages = append(req.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: fmt.Sprintf("filename: %s\nhint: %s", filename, hint),
	})
	resp, err := a.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return Response{}, err
	}

	ct := resp.Choices[0].Message.Content

	// Remove "<|eot_id|>" from the end of the content if present
	ct = strings.TrimSuffix(ct, "<|eot_id|>")

	var response Response
	err = json.Unmarshal([]byte(ct), &response)
	if err != nil {
		return Response{}, fmt.Errorf("failed to parse JSON response: %w, %s", err, ct)
	}

	return response, nil

}
