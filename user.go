package synapse

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ─── Interface ────────────────────────────────────────────────────────────────

// UserCase provides CRUD operations for users within the Synapse API.
type UserCase interface {
	// Get returns a single user by UUID or email.
	// Pass the authenticated user's own identifier to fetch their profile.
	Get(ctx context.Context, identifier string) (*UserResponseDto, error)

	// List returns a paginated list of users.
	// SystemAdmin can filter by tenant using params.TenantIdentifier.
	List(ctx context.Context, params ListUsersParams) ([]UserResponseDto, error)

	// Create registers a new user under the given tenant identifier (UUID or document).
	Create(ctx context.Context, tenantIdentifier string, req CreateUserRequestDto) (*UserResponseDto, error)

	// Update partially updates a user identified by UUID or email.
	Update(ctx context.Context, identifier string, req UpdateUserRequestDto) (*UserResponseDto, error)

	// Delete permanently removes a user identified by UUID or email.
	Delete(ctx context.Context, identifier string) error
}

// ─── Implementation ───────────────────────────────────────────────────────────

type userClient struct {
	http *httpClient
}

func newUserClient(hc *httpClient) UserCase {
	return &userClient{http: hc}
}

func (u *userClient) Get(ctx context.Context, identifier string) (*UserResponseDto, error) {
	var out UserResponseDto
	if err := u.http.get(ctx, pathUser+identifier, nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/user.Get: %w", err)
	}
	return &out, nil
}

func (u *userClient) List(ctx context.Context, params ListUsersParams) ([]UserResponseDto, error) {
	q := url.Values{}
	if params.Page > 0 {
		q.Set("page", strconv.Itoa(params.Page))
	}
	if params.Size > 0 {
		q.Set("size", strconv.Itoa(params.Size))
	}
	if params.TenantIdentifier != "" {
		q.Set("tenant_identifier", params.TenantIdentifier)
	}

	var out []UserResponseDto
	if err := u.http.get(ctx, pathUserList, q, &out); err != nil {
		return nil, fmt.Errorf("synapse/user.List: %w", err)
	}
	return out, nil
}

func (u *userClient) Create(ctx context.Context, tenantIdentifier string, req CreateUserRequestDto) (*UserResponseDto, error) {
	var out UserResponseDto
	if err := u.http.post(ctx, pathUser+tenantIdentifier, req, &out); err != nil {
		return nil, fmt.Errorf("synapse/user.Create: %w", err)
	}
	return &out, nil
}

func (u *userClient) Update(ctx context.Context, identifier string, req UpdateUserRequestDto) (*UserResponseDto, error) {
	var out UserResponseDto
	if err := u.http.patch(ctx, pathUser+identifier, req, &out); err != nil {
		return nil, fmt.Errorf("synapse/user.Update: %w", err)
	}
	return &out, nil
}

func (u *userClient) Delete(ctx context.Context, identifier string) error {
	if err := u.http.delete(ctx, pathUser+identifier, nil); err != nil {
		return fmt.Errorf("synapse/user.Delete: %w", err)
	}
	return nil
}
