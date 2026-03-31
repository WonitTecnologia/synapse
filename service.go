package synapse

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ─── Interface ────────────────────────────────────────────────────────────────

// ServiceCase provides read operations for catalog integration services.
type ServiceCase interface {
	// Get fetches a single integration service by its UUID.
	Get(ctx context.Context, id string) (*ServiceResponseDto, error)

	// List returns a paginated list of all available integration services.
	List(ctx context.Context, page, pageSize int) (*ServicesResponseDto, error)
}

// ─── Implementation ───────────────────────────────────────────────────────────

type serviceClient struct {
	http *httpClient
}

func newServiceClient(hc *httpClient) ServiceCase {
	return &serviceClient{http: hc}
}

func (s *serviceClient) Get(ctx context.Context, id string) (*ServiceResponseDto, error) {
	q := url.Values{}
	q.Set("id", id)

	var out ServiceResponseDto
	if err := s.http.get(ctx, pathService, q, &out); err != nil {
		return nil, fmt.Errorf("synapse/service.Get: %w", err)
	}
	return &out, nil
}

func (s *serviceClient) List(ctx context.Context, page, pageSize int) (*ServicesResponseDto, error) {
	q := url.Values{}
	if page > 0 {
		q.Set("page", strconv.Itoa(page))
	}
	if pageSize > 0 {
		q.Set("pageSize", strconv.Itoa(pageSize))
	}

	var out ServicesResponseDto
	if err := s.http.get(ctx, pathServiceList, q, &out); err != nil {
		return nil, fmt.Errorf("synapse/service.List: %w", err)
	}
	return &out, nil
}
