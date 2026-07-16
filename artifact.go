package synapse

import (
	"context"
	"fmt"
	"net/url"
)

// ─── Interface ────────────────────────────────────────────────────────────────

// ApiArtifactCase manages API Artifacts — per-agent pipelines of chained
// external API tools. Each active artifact becomes ONE tool for the main agent
// (artefato_<name>), executed by an internal agent that resolves the step
// parameters (declarative mappings + conversation) and caches every step's
// return in Redis per conversation. Requires the agent flag
// api_artifacts_enabled; with the flag off, API tools behave exactly as before.
type ApiArtifactCase interface {
	// List returns all artifacts of an agent.
	List(ctx context.Context, agentUUID string) (*ApiArtifactListResponse, error)
	Create(ctx context.Context, agentUUID string, req CreateApiArtifactRequest) (*ApiArtifactResponse, error)
	Update(ctx context.Context, agentUUID, artifactUUID string, req UpdateApiArtifactRequest) (*ApiArtifactResponse, error)
	// Toggle flips the is_active state of an artifact.
	Toggle(ctx context.Context, agentUUID, artifactUUID string) (*ApiArtifactResponse, error)
	Delete(ctx context.Context, agentUUID, artifactUUID string) error
	// Cache returns the cached step results of the agent's artifacts within a
	// conversation (the executor's memory) — audit view for monitoring.
	Cache(ctx context.Context, agentUUID, conversationUUID string) (*ArtifactCacheResponse, error)
}

// ─── Implementation ───────────────────────────────────────────────────────────

type apiArtifactClient struct {
	http *httpClient
}

func newApiArtifactClient(hc *httpClient) ApiArtifactCase {
	return &apiArtifactClient{http: hc}
}

func (a *apiArtifactClient) List(ctx context.Context, agentUUID string) (*ApiArtifactListResponse, error) {
	var out ApiArtifactListResponse
	if err := a.http.get(ctx, fmt.Sprintf(pathAgentArtifacts, agentUUID), nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/artifact.List: %w", err)
	}
	return &out, nil
}

func (a *apiArtifactClient) Create(ctx context.Context, agentUUID string, req CreateApiArtifactRequest) (*ApiArtifactResponse, error) {
	var out ApiArtifactResponse
	if err := a.http.post(ctx, fmt.Sprintf(pathAgentArtifacts, agentUUID), req, &out); err != nil {
		return nil, fmt.Errorf("synapse/artifact.Create: %w", err)
	}
	return &out, nil
}

func (a *apiArtifactClient) Update(ctx context.Context, agentUUID, artifactUUID string, req UpdateApiArtifactRequest) (*ApiArtifactResponse, error) {
	var out ApiArtifactResponse
	if err := a.http.put(ctx, fmt.Sprintf(pathAgentArtifact, agentUUID, artifactUUID), req, &out); err != nil {
		return nil, fmt.Errorf("synapse/artifact.Update: %w", err)
	}
	return &out, nil
}

func (a *apiArtifactClient) Toggle(ctx context.Context, agentUUID, artifactUUID string) (*ApiArtifactResponse, error) {
	var out ApiArtifactResponse
	if err := a.http.patch(ctx, fmt.Sprintf(pathAgentArtifactToggle, agentUUID, artifactUUID), nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/artifact.Toggle: %w", err)
	}
	return &out, nil
}

func (a *apiArtifactClient) Delete(ctx context.Context, agentUUID, artifactUUID string) error {
	if err := a.http.delete(ctx, fmt.Sprintf(pathAgentArtifact, agentUUID, artifactUUID), nil); err != nil {
		return fmt.Errorf("synapse/artifact.Delete: %w", err)
	}
	return nil
}

func (a *apiArtifactClient) Cache(ctx context.Context, agentUUID, conversationUUID string) (*ArtifactCacheResponse, error) {
	params := url.Values{}
	params.Set("conversation_uuid", conversationUUID)
	var out ArtifactCacheResponse
	if err := a.http.get(ctx, fmt.Sprintf(pathAgentArtifactCache, agentUUID), params, &out); err != nil {
		return nil, fmt.Errorf("synapse/artifact.Cache: %w", err)
	}
	return &out, nil
}
