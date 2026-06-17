package synapse

import (
	"context"
	"fmt"
)

// ─── Interface ────────────────────────────────────────────────────────────────

// ExternalApiCase provides CRUD operations for external API tools — raw HTTP APIs
// registered by a tenant that an AI agent can invoke as function-calling tools.
type ExternalApiCase interface {
	Create(ctx context.Context, req CreateExternalApiRequest) (*ExternalApiResponse, error)
	Get(ctx context.Context, uuid string) (*ExternalApiResponse, error)
	List(ctx context.Context) (*ExternalApiListResponse, error)
	Update(ctx context.Context, uuid string, req UpdateExternalApiRequest) (*ExternalApiResponse, error)
	// Toggle flips the is_active state of an external API tool.
	Toggle(ctx context.Context, uuid string) (*ExternalApiResponse, error)
	Delete(ctx context.Context, uuid string) error
}

// ─── Implementation ───────────────────────────────────────────────────────────

type externalApiClient struct {
	http *httpClient
}

func newExternalApiClient(hc *httpClient) ExternalApiCase {
	return &externalApiClient{http: hc}
}

func (m *externalApiClient) Create(ctx context.Context, req CreateExternalApiRequest) (*ExternalApiResponse, error) {
	var out ExternalApiResponse
	if err := m.http.post(ctx, pathExternalApis, req, &out); err != nil {
		return nil, fmt.Errorf("synapse/externalapi.Create: %w", err)
	}
	return &out, nil
}

func (m *externalApiClient) Get(ctx context.Context, uuid string) (*ExternalApiResponse, error) {
	var out ExternalApiResponse
	if err := m.http.get(ctx, fmt.Sprintf(pathExternalApi, uuid), nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/externalapi.Get: %w", err)
	}
	return &out, nil
}

func (m *externalApiClient) List(ctx context.Context) (*ExternalApiListResponse, error) {
	var out ExternalApiListResponse
	if err := m.http.get(ctx, pathExternalApis, nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/externalapi.List: %w", err)
	}
	return &out, nil
}

func (m *externalApiClient) Update(ctx context.Context, uuid string, req UpdateExternalApiRequest) (*ExternalApiResponse, error) {
	var out ExternalApiResponse
	if err := m.http.put(ctx, fmt.Sprintf(pathExternalApi, uuid), req, &out); err != nil {
		return nil, fmt.Errorf("synapse/externalapi.Update: %w", err)
	}
	return &out, nil
}

func (m *externalApiClient) Toggle(ctx context.Context, uuid string) (*ExternalApiResponse, error) {
	var out ExternalApiResponse
	if err := m.http.patch(ctx, fmt.Sprintf(pathExternalApiToggle, uuid), nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/externalapi.Toggle: %w", err)
	}
	return &out, nil
}

func (m *externalApiClient) Delete(ctx context.Context, uuid string) error {
	if err := m.http.delete(ctx, fmt.Sprintf(pathExternalApi, uuid), nil); err != nil {
		return fmt.Errorf("synapse/externalapi.Delete: %w", err)
	}
	return nil
}
