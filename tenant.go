package synapse

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ─── Interface ────────────────────────────────────────────────────────────────

// TenantCase provides CRUD operations for tenants within the Synapse API.
type TenantCase interface {
	// Get fetches a tenant by UUID or document (CNPJ/CPF).
	// At least one of uuid or document must be non-empty.
	Get(ctx context.Context, uuid, document string) (*TenantResponseDto, error)

	// List returns a paginated list of all tenants.
	List(ctx context.Context, page, pageSize int) (*TenantsResponseDto, error)

	// Create registers a new tenant.
	Create(ctx context.Context, req CreateTenantRequestDto) (*TenantResponseDto, error)

	// Update partially updates the tenant identified by uuid.
	Update(ctx context.Context, uuid string, req UpdateTenantRequestDto) (*TenantResponseDto, error)

	// Delete permanently removes a tenant by UUID or document.
	// At least one of uuid or document must be non-empty.
	Delete(ctx context.Context, uuid, document string) error
}

// ─── Implementation ───────────────────────────────────────────────────────────

type tenantClient struct {
	http *httpClient
}

func newTenantClient(hc *httpClient) TenantCase {
	return &tenantClient{http: hc}
}

func (t *tenantClient) Get(ctx context.Context, uuid, document string) (*TenantResponseDto, error) {
	q := url.Values{}
	if uuid != "" {
		q.Set("uuid", uuid)
	}
	if document != "" {
		q.Set("document", document)
	}

	var out TenantResponseDto
	if err := t.http.get(ctx, pathTenant, q, &out); err != nil {
		return nil, fmt.Errorf("synapse/tenant.Get: %w", err)
	}
	return &out, nil
}

func (t *tenantClient) List(ctx context.Context, page, pageSize int) (*TenantsResponseDto, error) {
	q := url.Values{}
	if page > 0 {
		q.Set("page", strconv.Itoa(page))
	}
	if pageSize > 0 {
		q.Set("pageSize", strconv.Itoa(pageSize))
	}

	var out TenantsResponseDto
	if err := t.http.get(ctx, pathTenantList, q, &out); err != nil {
		return nil, fmt.Errorf("synapse/tenant.List: %w", err)
	}
	return &out, nil
}

func (t *tenantClient) Create(ctx context.Context, req CreateTenantRequestDto) (*TenantResponseDto, error) {
	var out TenantResponseDto
	if err := t.http.post(ctx, pathTenantCreate, req, &out); err != nil {
		return nil, fmt.Errorf("synapse/tenant.Create: %w", err)
	}
	return &out, nil
}

func (t *tenantClient) Update(ctx context.Context, uuid string, req UpdateTenantRequestDto) (*TenantResponseDto, error) {
	var out TenantResponseDto
	if err := t.http.patch(ctx, pathTenant+"/"+uuid, req, &out); err != nil {
		return nil, fmt.Errorf("synapse/tenant.Update: %w", err)
	}
	return &out, nil
}

func (t *tenantClient) Delete(ctx context.Context, uuid, document string) error {
	q := url.Values{}
	if uuid != "" {
		q.Set("uuid", uuid)
	}
	if document != "" {
		q.Set("document", document)
	}

	if err := t.http.delete(ctx, pathTenant, q); err != nil {
		return fmt.Errorf("synapse/tenant.Delete: %w", err)
	}
	return nil
}
