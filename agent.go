package synapse

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ─── Interface ────────────────────────────────────────────────────────────────

// AgentCase provides CRUD, chat and conversation operations for AI agents.
type AgentCase interface {
	// Create registers a new AI agent.
	// SYSTEM_ADMIN must set req.TenantUUID; other roles use the tenant from the token.
	Create(ctx context.Context, req CreateAgentRequest) (*AgentResponse, error)

	// Get returns an agent by its UUID.
	// SYSTEM_ADMIN can fetch agents from any tenant.
	Get(ctx context.Context, agentUUID string) (*AgentResponse, error)

	// List returns a paginated list of agents.
	// SYSTEM_ADMIN lists agents from all tenants; other roles are scoped to their own tenant.
	// Use page=0 and size=0 to rely on server defaults.
	List(ctx context.Context, page, size int) (*ListAgentsResponse, error)

	// Update replaces all mutable fields of an agent (PUT semantics).
	// SYSTEM_ADMIN can update agents from any tenant.
	Update(ctx context.Context, agentUUID string, req UpdateAgentRequest) (*AgentResponse, error)

	// Patch updates a single field of an agent (PATCH semantics).
	// Set exactly one field in req; the server rejects requests with 0 or more than 1 field.
	// SYSTEM_ADMIN can patch agents from any tenant.
	Patch(ctx context.Context, agentUUID string, req UpdateAgentRequest) (*AgentResponse, error)

	// Delete permanently removes an agent and all its conversations.
	// SYSTEM_ADMIN can delete agents from any tenant.
	Delete(ctx context.Context, agentUUID string) error

	// Chat sends a message to an agent and returns the response.
	// Omit req.ConversationUUID to start a new conversation; the returned
	// ChatResponse.ConversationUUID can be passed in subsequent calls to maintain context.
	// req.ConversationUUID also accepts any client-defined identifier instead of a UUID
	// (e.g. "2024123ABCabc"): the backend generates and links a UUID to it on the first
	// call, scoped by tenant + agent, and reuses the same conversation on later calls with
	// the same identifier — see ChatRequest.ConversationUUID for details.
	// SYSTEM_ADMIN can chat with agents from any tenant.
	Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)

	// ListConversations returns a paginated list of conversation summaries.
	// SYSTEM_ADMIN lists conversations from all tenants; other roles are scoped to their own tenant.
	// Use params.AgentUUID to filter by a specific agent.
	ListConversations(ctx context.Context, params ListConversationsParams) (*ListConversationsResponse, error)

	// ListPrompts returns a paginated list of versioned prompts saved for an agent.
	// The IsActive field in each response indicates which prompt is currently active.
	ListPrompts(ctx context.Context, agentUUID string, page, size int) (*ListAgentPromptsResponse, error)

	// GetPrompt returns a single versioned prompt by its UUID.
	GetPrompt(ctx context.Context, agentUUID, promptUUID string) (*AgentPromptResponse, error)

	// CreatePrompt saves a new versioned prompt for an agent.
	// Multiple prompts with the same name are allowed; each is versioned by creation date.
	CreatePrompt(ctx context.Context, agentUUID string, req CreateAgentPromptRequest) (*AgentPromptResponse, error)

	// UpdatePrompt updates the name and/or content of a saved prompt.
	// If the prompt is currently active, the agent's prompt content is synced automatically.
	UpdatePrompt(ctx context.Context, agentUUID, promptUUID string, req UpdateAgentPromptRequest) (*AgentPromptResponse, error)

	// DeletePrompt permanently removes a saved prompt.
	// If the deleted prompt was the active one, the agent will have no active prompt (active_prompt_uuid = null).
	DeletePrompt(ctx context.Context, agentUUID, promptUUID string) error

	// ActivatePrompt sets a saved prompt as the active prompt for the agent.
	// The agent's prompt content is updated immediately to reflect the selected version.
	ActivatePrompt(ctx context.Context, agentUUID, promptUUID string) error

	// DeactivatePrompt removes the active prompt link from the agent (active_prompt_uuid = null).
	DeactivatePrompt(ctx context.Context, agentUUID string) error

	// ListLogs returns execution logs for an agent, filterable by conversation UUID or external ID.
	ListLogs(ctx context.Context, agentUUID string, params ListAgentLogsParams) (*ListAgentLogsResponse, error)

	// LogsStats returns aggregated token usage statistics for an agent, grouped by model and
	// conversation. Optionally filter by conversation UUID or external ID.
	LogsStats(ctx context.Context, agentUUID string, params ListAgentLogsParams) (*AgentLogStats, error)

	// SearchThoughts returns the thoughts (working memory) stored by the agent during a
	// conversation. conversationUUID is required; query filters by keyword in content/label
	// (case-insensitive, empty returns all). The thoughts are Redis-backed, per-conversation,
	// and expire after TTL (default 24h of inactivity).
	SearchThoughts(ctx context.Context, agentUUID string, params ThoughtSearchParams) (*ThoughtSearchResponse, error)

	GetCredits(ctx context.Context) (*WorkspaceCredits, error)
	GetActivity(ctx context.Context) ([]ActivityItem, error)
	GetActivityWithDate(ctx context.Context, date string) ([]ActivityItem, error)
	GetKeyInfo(ctx context.Context) (*KeyInfo, error)
}

// ─── Implementation ───────────────────────────────────────────────────────────

type agentClient struct {
	http *httpClient
}

func newAgentClient(hc *httpClient) AgentCase {
	return &agentClient{http: hc}
}

func (a *agentClient) Create(ctx context.Context, req CreateAgentRequest) (*AgentResponse, error) {
	var out AgentResponse
	if err := a.http.post(ctx, pathAgent, req, &out); err != nil {
		return nil, fmt.Errorf("synapse/agent.Create: %w", err)
	}
	return &out, nil
}

func (a *agentClient) Get(ctx context.Context, agentUUID string) (*AgentResponse, error) {
	var out AgentResponse
	if err := a.http.get(ctx, pathAgent+"/"+agentUUID, nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/agent.Get: %w", err)
	}
	return &out, nil
}

func (a *agentClient) List(ctx context.Context, page, size int) (*ListAgentsResponse, error) {
	q := url.Values{}
	if page > 0 {
		q.Set("page", strconv.Itoa(page))
	}
	if size > 0 {
		q.Set("size", strconv.Itoa(size))
	}

	var out ListAgentsResponse
	if err := a.http.get(ctx, pathAgent, q, &out); err != nil {
		return nil, fmt.Errorf("synapse/agent.List: %w", err)
	}
	return &out, nil
}

func (a *agentClient) Update(ctx context.Context, agentUUID string, req UpdateAgentRequest) (*AgentResponse, error) {
	var out AgentResponse
	if err := a.http.put(ctx, pathAgent+"/"+agentUUID, req, &out); err != nil {
		return nil, fmt.Errorf("synapse/agent.Update: %w", err)
	}
	return &out, nil
}

func (a *agentClient) Patch(ctx context.Context, agentUUID string, req UpdateAgentRequest) (*AgentResponse, error) {
	var out AgentResponse
	if err := a.http.patch(ctx, pathAgent+"/"+agentUUID, req, &out); err != nil {
		return nil, fmt.Errorf("synapse/agent.Patch: %w", err)
	}
	return &out, nil
}

func (a *agentClient) Delete(ctx context.Context, agentUUID string) error {
	if err := a.http.delete(ctx, pathAgent+"/"+agentUUID, nil); err != nil {
		return fmt.Errorf("synapse/agent.Delete: %w", err)
	}
	return nil
}

func (a *agentClient) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	var out ChatResponse
	if err := a.http.post(ctx, pathAgentChat, req, &out); err != nil {
		return nil, fmt.Errorf("synapse/agent.Chat: %w", err)
	}
	return &out, nil
}

func (a *agentClient) ListConversations(ctx context.Context, params ListConversationsParams) (*ListConversationsResponse, error) {
	q := url.Values{}
	if params.AgentUUID != "" {
		q.Set("agent_uuid", params.AgentUUID)
	}
	if params.ExternalID != "" {
		q.Set("external_id", params.ExternalID)
	}
	if params.Page > 0 {
		q.Set("page", strconv.Itoa(params.Page))
	}
	if params.Size > 0 {
		q.Set("size", strconv.Itoa(params.Size))
	}

	var out ListConversationsResponse
	if err := a.http.get(ctx, pathAgentConversation, q, &out); err != nil {
		return nil, fmt.Errorf("synapse/agent.ListConversations: %w", err)
	}
	return &out, nil
}

// ─── Prompt methods ───────────────────────────────────────────────────────────

func (a *agentClient) ListPrompts(ctx context.Context, agentUUID string, page, size int) (*ListAgentPromptsResponse, error) {
	q := url.Values{}
	if page > 0 {
		q.Set("page", strconv.Itoa(page))
	}
	if size > 0 {
		q.Set("size", strconv.Itoa(size))
	}
	var out ListAgentPromptsResponse
	path := fmt.Sprintf(pathAgentPrompt, agentUUID)
	if err := a.http.get(ctx, path, q, &out); err != nil {
		return nil, fmt.Errorf("synapse/agent.ListPrompts: %w", err)
	}
	return &out, nil
}

func (a *agentClient) GetPrompt(ctx context.Context, agentUUID, promptUUID string) (*AgentPromptResponse, error) {
	var out AgentPromptResponse
	path := fmt.Sprintf(pathAgentPrompt, agentUUID) + "/" + promptUUID
	if err := a.http.get(ctx, path, nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/agent.GetPrompt: %w", err)
	}
	return &out, nil
}

func (a *agentClient) CreatePrompt(ctx context.Context, agentUUID string, req CreateAgentPromptRequest) (*AgentPromptResponse, error) {
	var out AgentPromptResponse
	path := fmt.Sprintf(pathAgentPrompt, agentUUID)
	if err := a.http.post(ctx, path, req, &out); err != nil {
		return nil, fmt.Errorf("synapse/agent.CreatePrompt: %w", err)
	}
	return &out, nil
}

func (a *agentClient) UpdatePrompt(ctx context.Context, agentUUID, promptUUID string, req UpdateAgentPromptRequest) (*AgentPromptResponse, error) {
	var out AgentPromptResponse
	path := fmt.Sprintf(pathAgentPrompt, agentUUID) + "/" + promptUUID
	if err := a.http.put(ctx, path, req, &out); err != nil {
		return nil, fmt.Errorf("synapse/agent.UpdatePrompt: %w", err)
	}
	return &out, nil
}

func (a *agentClient) DeletePrompt(ctx context.Context, agentUUID, promptUUID string) error {
	path := fmt.Sprintf(pathAgentPrompt, agentUUID) + "/" + promptUUID
	if err := a.http.delete(ctx, path, nil); err != nil {
		return fmt.Errorf("synapse/agent.DeletePrompt: %w", err)
	}
	return nil
}

func (a *agentClient) ActivatePrompt(ctx context.Context, agentUUID, promptUUID string) error {
	path := fmt.Sprintf(pathAgentPrompt, agentUUID) + "/" + promptUUID + "/activate"
	if err := a.http.patch(ctx, path, struct{}{}, nil); err != nil {
		return fmt.Errorf("synapse/agent.ActivatePrompt: %w", err)
	}
	return nil
}

func (a *agentClient) DeactivatePrompt(ctx context.Context, agentUUID string) error {
	path := fmt.Sprintf(pathAgentActivePrompt, agentUUID)
	if err := a.http.delete(ctx, path, nil); err != nil {
		return fmt.Errorf("synapse/agent.DeactivatePrompt: %w", err)
	}
	return nil
}

func (a *agentClient) ListLogs(ctx context.Context, agentUUID string, params ListAgentLogsParams) (*ListAgentLogsResponse, error) {
	q := url.Values{}
	if params.ConversationUUID != "" {
		q.Set("conversation_uuid", params.ConversationUUID)
	}
	if params.ExternalID != "" {
		q.Set("external_id", params.ExternalID)
	}
	if params.Page > 0 {
		q.Set("page", strconv.Itoa(params.Page))
	}
	if params.Size > 0 {
		q.Set("size", strconv.Itoa(params.Size))
	}
	path := fmt.Sprintf(pathAgentLogs, agentUUID)
	var out ListAgentLogsResponse
	if err := a.http.get(ctx, path, q, &out); err != nil {
		return nil, fmt.Errorf("synapse/agent.ListLogs: %w", err)
	}
	return &out, nil
}

func (a *agentClient) LogsStats(ctx context.Context, agentUUID string, params ListAgentLogsParams) (*AgentLogStats, error) {
	q := url.Values{}
	if params.ConversationUUID != "" {
		q.Set("conversation_uuid", params.ConversationUUID)
	}
	if params.ExternalID != "" {
		q.Set("external_id", params.ExternalID)
	}
	path := fmt.Sprintf(pathAgentLogsStats, agentUUID)
	var out AgentLogStats
	if err := a.http.get(ctx, path, q, &out); err != nil {
		return nil, fmt.Errorf("synapse/agent.LogsStats: %w", err)
	}
	return &out, nil
}

func (a *agentClient) SearchThoughts(ctx context.Context, agentUUID string, params ThoughtSearchParams) (*ThoughtSearchResponse, error) {
	q := url.Values{}
	q.Set("conversation_uuid", params.ConversationUUID)
	if params.Query != "" {
		q.Set("query", params.Query)
	}
	path := fmt.Sprintf(pathAgentThoughts, agentUUID)
	var out ThoughtSearchResponse
	if err := a.http.get(ctx, path, q, &out); err != nil {
		return nil, fmt.Errorf("synapse/agent.SearchThoughts: %w", err)
	}
	return &out, nil
}

func (a *agentClient) GetCredits(ctx context.Context) (*WorkspaceCredits, error) {
	var out WorkspaceCredits
	if err := a.http.get(ctx, pathWorkspaceCredits, nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/agent.GetCredits: %w", err)
	}
	return &out, nil
}

func (a *agentClient) GetActivity(ctx context.Context) ([]ActivityItem, error) {
	var out []ActivityItem
	if err := a.http.get(ctx, pathWorkspaceActivity, nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/agent.GetActivity: %w", err)
	}
	return out, nil
}

func (a *agentClient) GetActivityWithDate(ctx context.Context, date string) ([]ActivityItem, error) {
	q := url.Values{}
	q.Set("date", date)
	var out []ActivityItem
	if err := a.http.get(ctx, pathWorkspaceActivity, q, &out); err != nil {
		return nil, fmt.Errorf("synapse/agent.GetActivity: %w", err)
	}
	return out, nil
}

func (a *agentClient) GetKeyInfo(ctx context.Context) (*KeyInfo, error) {
	var out KeyInfo
	if err := a.http.get(ctx, pathWorkspaceKey, nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/agent.GetKeyInfo: %w", err)
	}
	return &out, nil
}
