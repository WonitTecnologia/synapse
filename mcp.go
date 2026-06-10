package synapse

import (
	"context"
	"fmt"
)

// ─── Interface ────────────────────────────────────────────────────────────────

// McpCase provides CRUD operations for MCP (Model Context Protocol) server integrations.
type McpCase interface {
	Create(ctx context.Context, req CreateMcpIntegrationRequest) (*McpIntegrationResponse, error)
	Get(ctx context.Context, uuid string) (*McpIntegrationResponse, error)
	List(ctx context.Context) (*McpIntegrationListResponse, error)
	Update(ctx context.Context, uuid string, req UpdateMcpIntegrationRequest) (*McpIntegrationResponse, error)
	Toggle(ctx context.Context, uuid string, req ToggleMcpIntegrationRequest) (*McpIntegrationResponse, error)
	Delete(ctx context.Context, uuid string) error
	// GetTools lists the tools exposed by the remote MCP server behind this integration.
	GetTools(ctx context.Context, uuid string) (*McpToolsListResponse, error)
}

// ─── Implementation ───────────────────────────────────────────────────────────

type mcpClient struct {
	http *httpClient
}

func newMcpClient(hc *httpClient) McpCase {
	return &mcpClient{http: hc}
}

func (m *mcpClient) Create(ctx context.Context, req CreateMcpIntegrationRequest) (*McpIntegrationResponse, error) {
	var out McpIntegrationResponse
	if err := m.http.post(ctx, pathMcpIntegrations, req, &out); err != nil {
		return nil, fmt.Errorf("synapse/mcp.Create: %w", err)
	}
	return &out, nil
}

func (m *mcpClient) Get(ctx context.Context, uuid string) (*McpIntegrationResponse, error) {
	var out McpIntegrationResponse
	if err := m.http.get(ctx, fmt.Sprintf(pathMcpIntegration, uuid), nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/mcp.Get: %w", err)
	}
	return &out, nil
}

func (m *mcpClient) List(ctx context.Context) (*McpIntegrationListResponse, error) {
	var out McpIntegrationListResponse
	if err := m.http.get(ctx, pathMcpIntegrations, nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/mcp.List: %w", err)
	}
	return &out, nil
}

func (m *mcpClient) Update(ctx context.Context, uuid string, req UpdateMcpIntegrationRequest) (*McpIntegrationResponse, error) {
	var out McpIntegrationResponse
	if err := m.http.put(ctx, fmt.Sprintf(pathMcpIntegration, uuid), req, &out); err != nil {
		return nil, fmt.Errorf("synapse/mcp.Update: %w", err)
	}
	return &out, nil
}

func (m *mcpClient) Toggle(ctx context.Context, uuid string, req ToggleMcpIntegrationRequest) (*McpIntegrationResponse, error) {
	var out McpIntegrationResponse
	if err := m.http.patch(ctx, fmt.Sprintf(pathMcpToggle, uuid), req, &out); err != nil {
		return nil, fmt.Errorf("synapse/mcp.Toggle: %w", err)
	}
	return &out, nil
}

func (m *mcpClient) Delete(ctx context.Context, uuid string) error {
	if err := m.http.delete(ctx, fmt.Sprintf(pathMcpIntegration, uuid), nil); err != nil {
		return fmt.Errorf("synapse/mcp.Delete: %w", err)
	}
	return nil
}

func (m *mcpClient) GetTools(ctx context.Context, uuid string) (*McpToolsListResponse, error) {
	var out McpToolsListResponse
	if err := m.http.get(ctx, fmt.Sprintf(pathMcpIntegrationTools, uuid), nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/mcp.GetTools: %w", err)
	}
	return &out, nil
}
