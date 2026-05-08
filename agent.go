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
	// SYSTEM_ADMIN can chat with agents from any tenant.
	Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)

	// ListConversations returns a paginated list of conversation summaries.
	// SYSTEM_ADMIN lists conversations from all tenants; other roles are scoped to their own tenant.
	// Use params.AgentUUID to filter by a specific agent.
	ListConversations(ctx context.Context, params ListConversationsParams) (*ListConversationsResponse, error)
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
