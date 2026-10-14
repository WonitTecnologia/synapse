package synapse

import (
	"context"
	"fmt"
)

// ─── Interface ────────────────────────────────────────────────────────────────

// ClassificationCase manages the custom efficiency classifications of an agent:
// labels applied at the end of an attendance (e.g. "resolved", "no answer"),
// each with an automatic action (transfer/close). Used together with the agent
// flag efficiency_tagging_required.
type ClassificationCase interface {
	// List returns all efficiency classifications of an agent.
	List(ctx context.Context, agentUUID string) (*EfficiencyClassificationListResponse, error)
	Create(ctx context.Context, agentUUID string, req CreateEfficiencyClassificationRequest) (*EfficiencyClassification, error)
	Update(ctx context.Context, agentUUID, classificationUUID string, req UpdateEfficiencyClassificationRequest) (*EfficiencyClassification, error)
	Delete(ctx context.Context, agentUUID, classificationUUID string) error
}

// ─── Implementation ───────────────────────────────────────────────────────────

type classificationClient struct {
	http *httpClient
}

func newClassificationClient(hc *httpClient) ClassificationCase {
	return &classificationClient{http: hc}
}

func (c *classificationClient) List(ctx context.Context, agentUUID string) (*EfficiencyClassificationListResponse, error) {
	var out EfficiencyClassificationListResponse
	if err := c.http.get(ctx, fmt.Sprintf(pathEfficiencyClassifications, agentUUID), nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/classification.List: %w", err)
	}
	return &out, nil
}

func (c *classificationClient) Create(ctx context.Context, agentUUID string, req CreateEfficiencyClassificationRequest) (*EfficiencyClassification, error) {
	var out EfficiencyClassification
	if err := c.http.post(ctx, fmt.Sprintf(pathEfficiencyClassifications, agentUUID), req, &out); err != nil {
		return nil, fmt.Errorf("synapse/classification.Create: %w", err)
	}
	return &out, nil
}

func (c *classificationClient) Update(ctx context.Context, agentUUID, classificationUUID string, req UpdateEfficiencyClassificationRequest) (*EfficiencyClassification, error) {
	var out EfficiencyClassification
	if err := c.http.put(ctx, fmt.Sprintf(pathEfficiencyClassification, agentUUID, classificationUUID), req, &out); err != nil {
		return nil, fmt.Errorf("synapse/classification.Update: %w", err)
	}
	return &out, nil
}

func (c *classificationClient) Delete(ctx context.Context, agentUUID, classificationUUID string) error {
	if err := c.http.delete(ctx, fmt.Sprintf(pathEfficiencyClassification, agentUUID, classificationUUID), nil); err != nil {
		return fmt.Errorf("synapse/classification.Delete: %w", err)
	}
	return nil
}
