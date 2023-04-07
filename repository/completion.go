package repository

import (
	"context"
	"github.com/piroyoung/go-aoai"
)

type CompletionRepository struct {
	client *aoai.AzureOpenAI
}

func NewCompletionRepository(client *aoai.AzureOpenAI) *CompletionRepository {
	return &CompletionRepository{client: client}
}

func (r *CompletionRepository) Complete(ctx context.Context, prompt string, m int) (string, error) {
	request := aoai.CompletionRequest{
		Prompts:   []string{prompt},
		MaxTokens: m,
	}

	response, err := r.client.Completion(ctx, request)
	if err != nil {
		return "", err
	}
	if len(response.Choices) == 0 {
		return "", nil
	} else {
		m := response.Choices[0].Text
		return cleanText(m), nil
	}
}

func (r *CompletionRepository) Stream(ctx context.Context, prompt string, m int, consumer func(t string) error) error {
	request := aoai.CompletionRequest{
		Prompts:   []string{prompt},
		MaxTokens: m,
		Stream:    true,
	}

	return r.client.CompletionStream(ctx, request, func(chunk aoai.CompletionResponse) error {
		if len(chunk.Choices) == 0 {
			return nil
		} else {
			m := chunk.Choices[0].Text
			return consumer(cleanText(m))
		}
	})
}

// replace all punctuation and special characters into space
func cleanText(text string) string {
	//text = strings.ReplaceAll(text, "\n", " ")
	//text = strings.ReplaceAll(text, "\r", "_")
	//text = strings.ReplaceAll(text, "\t", "_")
	//re := regexp.MustCompile(`[\p{P}\p{S}]+`)
	//text = re.ReplaceAllString(text, " ")
	//text = strings.ReplaceAll(text, "  ", " ")
	//text = strings.TrimSpace(text)
	return text
}
