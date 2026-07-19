package synapse

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ─── Interface ────────────────────────────────────────────────────────────────

// SystemAgentCase manages system agents: platform-owned agents (tenant NULL)
// that any tenant can chat with, always billed on the CALLER tenant's
// OpenRouter key. All operations require a SYSTEM_ADMIN token.
type SystemAgentCase interface {
	// Create registers a new system agent (SYSTEM_ADMIN only).
	Create(ctx context.Context, req CreateSystemAgentRequest) (*AgentResponse, error)

	// Get returns a system agent by its UUID, including prompt, enabled tools
	// and configured models.
	Get(ctx context.Context, agentUUID string) (*AgentResponse, error)

	// List returns a paginated list of system agents.
	// Use page=0 and size=0 to rely on server defaults.
	List(ctx context.Context, page, size int) (*ListAgentsResponse, error)

	// Update replaces mutable fields of a system agent (PUT semantics).
	Update(ctx context.Context, agentUUID string, req UpdateAgentRequest) (*AgentResponse, error)

	// Patch updates a single field of a system agent (PATCH semantics) — e.g.
	// switching the LLM model from the Master panel. Set exactly one field in
	// req; the server rejects requests with 0 or more than 1 field.
	Patch(ctx context.Context, agentUUID string, req UpdateAgentRequest) (*AgentResponse, error)

	// Delete permanently removes a system agent and all its conversations.
	Delete(ctx context.Context, agentUUID string) error

	// Chat talks to a system agent through its DEDICATED execution pipeline
	// (text-only, no judge, sistema_* tools). Use a TENANT token: the turn is
	// always billed on the caller tenant's OpenRouter key. ConversationUUID
	// accepts an existing conversation UUID or any client-defined slug (linked
	// per tenant+agent). Attachment is ignored by this pipeline.
	Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
}

// CreateSystemAgentRequest mirrors CreateAgentRequest without TenantUUID:
// system agents belong to the platform, not to a tenant.
type CreateSystemAgentRequest struct {
	Name                string   `json:"name"`
	Description         string   `json:"description,omitempty"`
	Model               string   `json:"model"`
	Prompt              string   `json:"prompt"`
	CollectionUUIDs     []string `json:"collection_uuids,omitempty"`
	MaxContext          int      `json:"max_context,omitempty"`
	Temperature         *float64 `json:"temperature,omitempty"`
	McpEnabled          *bool    `json:"mcp_enabled,omitempty"`
	McpIntegrationUUIDs []string `json:"mcp_integration_uuids,omitempty"`
	McpDisabledTools    []string `json:"mcp_disabled_tools,omitempty"`
	ApiToolsEnabled     *bool    `json:"api_tools_enabled,omitempty"`
	ApiToolUUIDs        []string `json:"api_tool_uuids,omitempty"`
	ThoughtsEnabled     *bool    `json:"thoughts_enabled,omitempty"`
	ApiArtifactsEnabled *bool    `json:"api_artifacts_enabled,omitempty"`
	AcceptFiles         *bool    `json:"accept_files,omitempty"`
	FileModel           string   `json:"file_model,omitempty"`
	TextModel           string   `json:"text_model,omitempty"`
	ImageModel          string   `json:"image_model,omitempty"`
	AudioModel          string   `json:"audio_model,omitempty"`
	TextFallbackModel   string   `json:"text_fallback_model,omitempty"`
	ImageFallbackModel  string   `json:"image_fallback_model,omitempty"`
	AudioFallbackModel  string   `json:"audio_fallback_model,omitempty"`
}

// ─── Implementation ───────────────────────────────────────────────────────────

type systemAgentClient struct {
	http *httpClient
}

func newSystemAgentClient(hc *httpClient) SystemAgentCase {
	return &systemAgentClient{http: hc}
}

func (a *systemAgentClient) Create(ctx context.Context, req CreateSystemAgentRequest) (*AgentResponse, error) {
	var out AgentResponse
	if err := a.http.post(ctx, pathSystemAgent, req, &out); err != nil {
		return nil, fmt.Errorf("synapse/systemagent.Create: %w", err)
	}
	return &out, nil
}

func (a *systemAgentClient) Get(ctx context.Context, agentUUID string) (*AgentResponse, error) {
	var out AgentResponse
	if err := a.http.get(ctx, pathSystemAgent+"/"+agentUUID, nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/systemagent.Get: %w", err)
	}
	return &out, nil
}

func (a *systemAgentClient) List(ctx context.Context, page, size int) (*ListAgentsResponse, error) {
	q := url.Values{}
	if page > 0 {
		q.Set("page", strconv.Itoa(page))
	}
	if size > 0 {
		q.Set("size", strconv.Itoa(size))
	}

	var out ListAgentsResponse
	if err := a.http.get(ctx, pathSystemAgent, q, &out); err != nil {
		return nil, fmt.Errorf("synapse/systemagent.List: %w", err)
	}
	return &out, nil
}

func (a *systemAgentClient) Update(ctx context.Context, agentUUID string, req UpdateAgentRequest) (*AgentResponse, error) {
	var out AgentResponse
	if err := a.http.put(ctx, pathSystemAgent+"/"+agentUUID, req, &out); err != nil {
		return nil, fmt.Errorf("synapse/systemagent.Update: %w", err)
	}
	return &out, nil
}

func (a *systemAgentClient) Patch(ctx context.Context, agentUUID string, req UpdateAgentRequest) (*AgentResponse, error) {
	var out AgentResponse
	if err := a.http.patch(ctx, pathSystemAgent+"/"+agentUUID, req, &out); err != nil {
		return nil, fmt.Errorf("synapse/systemagent.Patch: %w", err)
	}
	return &out, nil
}

func (a *systemAgentClient) Delete(ctx context.Context, agentUUID string) error {
	if err := a.http.delete(ctx, pathSystemAgent+"/"+agentUUID, nil); err != nil {
		return fmt.Errorf("synapse/systemagent.Delete: %w", err)
	}
	return nil
}

func (a *systemAgentClient) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	var out ChatResponse
	if err := a.http.post(ctx, pathSystemAgentChat, req, &out); err != nil {
		return nil, fmt.Errorf("synapse/systemagent.Chat: %w", err)
	}
	return &out, nil
}
