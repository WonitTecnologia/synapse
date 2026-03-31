package synapse

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ─── Interface ────────────────────────────────────────────────────────────────

// ProviderCase provides read operations for catalog integration providers.
type ProviderCase interface {
	// Get fetches a single provider by its UUID.
	Get(ctx context.Context, id string) (*ProviderResponseDto, error)

	// List returns a paginated list of all available providers.
	List(ctx context.Context, page, pageSize int) (*ProvidersResponseDto, error)
}

// ─── Implementation ───────────────────────────────────────────────────────────

type providerClient struct {
	http *httpClient
}

func newProviderClient(hc *httpClient) ProviderCase {
	return &providerClient{http: hc}
}

func (p *providerClient) Get(ctx context.Context, id string) (*ProviderResponseDto, error) {
	q := url.Values{}
	q.Set("id", id)

	var out ProviderResponseDto
	if err := p.http.get(ctx, pathProvider, q, &out); err != nil {
		return nil, fmt.Errorf("synapse/provider.Get: %w", err)
	}
	return &out, nil
}

func (p *providerClient) List(ctx context.Context, page, pageSize int) (*ProvidersResponseDto, error) {
	q := url.Values{}
	if page > 0 {
		q.Set("page", strconv.Itoa(page))
	}
	if pageSize > 0 {
		q.Set("pageSize", strconv.Itoa(pageSize))
	}

	var out ProvidersResponseDto
	if err := p.http.get(ctx, pathProviderList, q, &out); err != nil {
		return nil, fmt.Errorf("synapse/provider.List: %w", err)
	}
	return &out, nil
}
