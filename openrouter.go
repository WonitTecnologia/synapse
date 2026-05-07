package synapse

import (
	"context"
	"fmt"
	"net/url"
)

// ─── Interface ────────────────────────────────────────────────────────────────

// OpenRouterCase provides operations for the OpenRouter integration.
type OpenRouterCase interface {
	// Sync creates or returns an existing OpenRouter workspace for the authenticated tenant.
	// SystemAdmin callers may target a specific tenant by setting req.TenantUUID.
	Sync(ctx context.Context, req OpenRouterSyncRequest) (*OpenRouterSyncResponse, error)

	// Desync removes the OpenRouter workspace for the authenticated tenant.
	// SystemAdmin callers may target a specific tenant by setting req.TenantUUID.
	Desync(ctx context.Context, req OpenRouterSyncRequest) (*OpenRouterDesyncResponse, error)

	// ListModels returns all models available on OpenRouter.
	// Set freeOnly=true to filter for free-tier models only.
	ListModels(ctx context.Context, freeOnly bool) (*OpenRouterListModelsResponse, error)

	// ListEmbeddingModels returns the catalogue of embedding models available on OpenRouter,
	// including the vector_size required to configure a collection before indexing documents.
	// Set freeOnly=true to filter for free-tier models only.
	ListEmbeddingModels(ctx context.Context, freeOnly bool) (*OpenRouterListEmbeddingModelsResponse, error)
}

// ─── Implementation ───────────────────────────────────────────────────────────

type openrouterClient struct {
	http *httpClient
}

func newOpenRouterClient(hc *httpClient) OpenRouterCase {
	return &openrouterClient{http: hc}
}

func (o *openrouterClient) Sync(ctx context.Context, req OpenRouterSyncRequest) (*OpenRouterSyncResponse, error) {
	var out OpenRouterSyncResponse
	if err := o.http.post(ctx, pathOpenRouterSincronismo, req, &out); err != nil {
		return nil, fmt.Errorf("synapse/openrouter.Sync: %w", err)
	}
	return &out, nil
}

func (o *openrouterClient) Desync(ctx context.Context, req OpenRouterSyncRequest) (*OpenRouterDesyncResponse, error) {
	var out OpenRouterDesyncResponse
	if err := o.http.deleteJSON(ctx, pathOpenRouterSincronismo, req, &out); err != nil {
		return nil, fmt.Errorf("synapse/openrouter.Desync: %w", err)
	}
	return &out, nil
}

func (o *openrouterClient) ListModels(ctx context.Context, freeOnly bool) (*OpenRouterListModelsResponse, error) {
	q := url.Values{}
	if freeOnly {
		q.Set("free_only", "true")
	}
	var out OpenRouterListModelsResponse
	if err := o.http.get(ctx, pathOpenRouterModels, q, &out); err != nil {
		return nil, fmt.Errorf("synapse/openrouter.ListModels: %w", err)
	}
	return &out, nil
}

func (o *openrouterClient) ListEmbeddingModels(ctx context.Context, freeOnly bool) (*OpenRouterListEmbeddingModelsResponse, error) {
	q := url.Values{}
	if freeOnly {
		q.Set("free_only", "true")
	}
	var out OpenRouterListEmbeddingModelsResponse
	if err := o.http.get(ctx, pathOpenRouterEmbeddingModels, q, &out); err != nil {
		return nil, fmt.Errorf("synapse/openrouter.ListEmbeddingModels: %w", err)
	}
	return &out, nil
}
