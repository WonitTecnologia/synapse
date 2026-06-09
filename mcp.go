package synapse

import (
	"context"
	"fmt"
)

// McpCase provides operations for managing MCP (Model Context Protocol) server integrations.
// An MCP integration links an AI agent tenant to an external MCP server (e.g. PABX)
// so that the agent can discover and execute tools exposed by that server.
type McpCase interface {
	// Create registers a new MCP server integration for the authenticated tenant.
	Create(ctx context.Context, req CreateMcpIntegrationRequest) (*McpIntegrationResponse, error)

	// Get returns a single MCP integration by its UUID.
	Get(ctx context.Context, uuid string) (*McpIntegrationResponse, error)

	// List returns all MCP integrations for the authenticated tenant.
	List(ctx context.Context) (*McpIntegrationListResponse, error)

	// Update partially updates an MCP integration.
	Update(ctx context.Context, uuid string, req UpdateMcpIntegrationRequest) (*McpIntegrationResponse, error)

	// Toggle activates or deactivates an MCP integration.
	Toggle(ctx context.Context, uuid string, active bool) error

	// Delete permanently removes an MCP integration.
	Delete(ctx context.Context, uuid string) error
}

type mcpClient struct {
	http *httpClient
}

func newMcpClient(hc *httpClient) McpCase {
	return &mcpClient{http: hc}
}

func (c *mcpClient) Create(ctx context.Context, req CreateMcpIntegrationRequest) (*McpIntegrationResponse, error) {
	var out McpIntegrationResponse
	if err := c.http.post(ctx, pathMcp, req, &out); err != nil {
		return nil, fmt.Errorf("synapse/mcp.Create: %w", err)
	}
	return &out, nil
}

func (c *mcpClient) Get(ctx context.Context, uuid string) (*McpIntegrationResponse, error) {
	var out McpIntegrationResponse
	if err := c.http.get(ctx, fmt.Sprintf(pathMcpByUUID, uuid), nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/mcp.Get: %w", err)
	}
	return &out, nil
}

func (c *mcpClient) List(ctx context.Context) (*McpIntegrationListResponse, error) {
	var out McpIntegrationListResponse
	if err := c.http.get(ctx, pathMcp, nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/mcp.List: %w", err)
	}
	return &out, nil
}

func (c *mcpClient) Update(ctx context.Context, uuid string, req UpdateMcpIntegrationRequest) (*McpIntegrationResponse, error) {
	var out McpIntegrationResponse
	if err := c.http.put(ctx, fmt.Sprintf(pathMcpByUUID, uuid), req, &out); err != nil {
		return nil, fmt.Errorf("synapse/mcp.Update: %w", err)
	}
	return &out, nil
}

func (c *mcpClient) Toggle(ctx context.Context, uuid string, active bool) error {
	if err := c.http.patch(ctx, fmt.Sprintf(pathMcpToggle, uuid), ToggleMcpIntegrationRequest{IsActive: active}, nil); err != nil {
		return fmt.Errorf("synapse/mcp.Toggle: %w", err)
	}
	return nil
}

func (c *mcpClient) Delete(ctx context.Context, uuid string) error {
	if err := c.http.delete(ctx, fmt.Sprintf(pathMcpByUUID, uuid), nil); err != nil {
		return fmt.Errorf("synapse/mcp.Delete: %w", err)
	}
	return nil
}
