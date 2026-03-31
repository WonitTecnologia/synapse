package synapse

import (
	"context"
	"fmt"
)

// ─── Interface ────────────────────────────────────────────────────────────────

// OpenAICase provides operations for the OpenAI integration.
type OpenAICase interface {
	// Configure creates or updates the OpenAI integration credentials and settings
	// for the authenticated tenant.
	Configure(ctx context.Context, req ConfigureOpenAIRequest) error

	// Chat executes a Chat Completion call using the tenant's configured OpenAI account.
	Chat(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error)

	// AnalyzeImage analyses an image using OpenAI vision.
	// fileName is the multipart filename (e.g. "image.png").
	// fileContent is the raw image bytes (png, jpg, jpeg, webp).
	// prompt describes what analysis to perform.
	AnalyzeImage(ctx context.Context, fileName string, fileContent []byte, prompt string) (*AnalyzeImageResponse, error)

	// TranscribeAudio converts an audio file to text using OpenAI Whisper or GPT-4o.
	// See TranscribeAudioRequest for the full set of options.
	TranscribeAudio(ctx context.Context, req TranscribeAudioRequest) (*TranscriptionResponse, error)
}

// ─── Implementation ───────────────────────────────────────────────────────────

type openaiClient struct {
	http *httpClient
}

func newOpenAIClient(hc *httpClient) OpenAICase {
	return &openaiClient{http: hc}
}

func (o *openaiClient) Configure(ctx context.Context, req ConfigureOpenAIRequest) error {
	if err := o.http.post(ctx, pathOpenAIConfig, req, nil); err != nil {
		return fmt.Errorf("synapse/openai.Configure: %w", err)
	}
	return nil
}

func (o *openaiClient) Chat(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	var out ChatCompletionResponse
	if err := o.http.post(ctx, pathOpenAIChat, req, &out); err != nil {
		return nil, fmt.Errorf("synapse/openai.Chat: %w", err)
	}
	return &out, nil
}

func (o *openaiClient) AnalyzeImage(ctx context.Context, fileName string, fileContent []byte, prompt string) (*AnalyzeImageResponse, error) {
	var out AnalyzeImageResponse
	err := o.http.postMultipart(
		ctx,
		pathOpenAIImageAnalyze,
		map[string]string{"prompt": prompt},
		"file",
		fileName,
		fileContent,
		&out,
	)
	if err != nil {
		return nil, fmt.Errorf("synapse/openai.AnalyzeImage: %w", err)
	}
	return &out, nil
}

func (o *openaiClient) TranscribeAudio(ctx context.Context, req TranscribeAudioRequest) (*TranscriptionResponse, error) {
	fields := map[string]string{
		"model":    req.Model,
		"language": req.Language,
		"prompt":   req.Prompt,
	}

	var out TranscriptionResponse
	err := o.http.postMultipart(
		ctx,
		pathOpenAITranscribe,
		fields,
		"file",
		req.FileName,
		req.Content,
		&out,
	)
	if err != nil {
		return nil, fmt.Errorf("synapse/openai.TranscribeAudio: %w", err)
	}
	return &out, nil
}
