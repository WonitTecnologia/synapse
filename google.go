package synapse

import (
	"context"
	"fmt"
)

// ─── Interface ────────────────────────────────────────────────────────────────

// GoogleCase provides operations for the Google Vision AI integration.
type GoogleCase interface {
	// Configure creates or updates the Google Vision AI integration credentials
	// for the authenticated tenant.
	Configure(ctx context.Context, req ConfigureGoogleRequest) error

	// VisionOCR extracts text from an image using Google Vision TEXT_DETECTION.
	// fileName is used as the multipart filename (e.g. "photo.jpg").
	// fileContent is the raw image bytes (png, jpg, jpeg, webp).
	VisionOCR(ctx context.Context, fileName string, fileContent []byte) (*AnalyzeImageResponse, error)
}

// ─── Implementation ───────────────────────────────────────────────────────────

type googleClient struct {
	http *httpClient
}

func newGoogleClient(hc *httpClient) GoogleCase {
	return &googleClient{http: hc}
}

func (g *googleClient) Configure(ctx context.Context, req ConfigureGoogleRequest) error {
	if err := g.http.post(ctx, pathGoogleVisionConfig, req, nil); err != nil {
		return fmt.Errorf("synapse/google.Configure: %w", err)
	}
	return nil
}

func (g *googleClient) VisionOCR(ctx context.Context, fileName string, fileContent []byte) (*AnalyzeImageResponse, error) {
	var out AnalyzeImageResponse
	err := g.http.postMultipart(
		ctx,
		pathGoogleVisionOCR,
		nil,
		"file",
		fileName,
		fileContent,
		&out,
	)
	if err != nil {
		return nil, fmt.Errorf("synapse/google.VisionOCR: %w", err)
	}
	return &out, nil
}
