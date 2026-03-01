package services

import (
	"context"
	"fmt"
	"os"

	"merope/models"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiService struct {
	enabled bool
	client  *genai.Client
	model   *genai.GenerativeModel
}

func NewGeminiService() *GeminiService {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return &GeminiService{enabled: false}
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		fmt.Printf("Failed to create Gemini client: %v\n", err)
		return &GeminiService{enabled: false}
	}

	model := client.GenerativeModel("gemini-flash-latest")

	return &GeminiService{
		enabled: true,
		client:  client,
		model:   model,
	}
}

func (g *GeminiService) AnalyzeAlert(alert *models.Alert) (string, error) {
	if !g.enabled {
		return "", nil
	}

	ctx := context.Background()
	prompt := fmt.Sprintf(`You are a DevOps expert assistant. Analyze this system alert and provide a very brief, actionable recommendation (max 2 sentences).
Alert: %s
Message: %s
Level: %s`, alert.Title, alert.Message, alert.Level)

	resp, err := g.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		if txt, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			return string(txt), nil
		}
	}

	return "", nil
}

func (g *GeminiService) Close() {
	if g.client != nil {
		g.client.Close()
	}
}

func (g *GeminiService) IsEnabled() bool {
	return g.enabled
}
