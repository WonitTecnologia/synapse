package synapse

import (
	"context"
	"fmt"
)

// ─── Interface ────────────────────────────────────────────────────────────────

// ChatvoltCase provides operations for the Chatvolt agent integration.
type ChatvoltCase interface {
	// Configure creates or updates the Chatvolt integration credentials
	// for the authenticated tenant.
	Configure(ctx context.Context, req ConfigureChatvoltRequest) error

	// Query sends a message to a Chatvolt agent and returns its response.
	// Set req.ConversationID to continue an existing conversation thread.
	Query(ctx context.Context, req ChatvoltAgentQueryRequest) (*ChatvoltAgentQueryResponse, error)
}

// ─── Implementation ───────────────────────────────────────────────────────────

type chatvoltClient struct {
	http *httpClient
}

func newChatvoltClient(hc *httpClient) ChatvoltCase {
	return &chatvoltClient{http: hc}
}

func (c *chatvoltClient) Configure(ctx context.Context, req ConfigureChatvoltRequest) error {
	if err := c.http.post(ctx, pathChatvoltConfig, req, nil); err != nil {
		return fmt.Errorf("synapse/chatvolt.Configure: %w", err)
	}
	return nil
}

func (c *chatvoltClient) Query(ctx context.Context, req ChatvoltAgentQueryRequest) (*ChatvoltAgentQueryResponse, error) {
	var out ChatvoltAgentQueryResponse
	if err := c.http.post(ctx, pathChatvoltQuery, req, &out); err != nil {
		return nil, fmt.Errorf("synapse/chatvolt.Query: %w", err)
	}
	return &out, nil
}
