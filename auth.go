package synapse

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ─── Interface ────────────────────────────────────────────────────────────────

// AuthCase provides all authentication-related operations against the Synapse API.
type AuthCase interface {
	// Login authenticates a user with email and password and returns an access token.
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)

	// Logout revokes the given access token.
	Logout(ctx context.Context, token string) error

	// Healthcheck validates the current token and returns the logged-in user's data.
	Healthcheck(ctx context.Context) (*LoginResponse, error)

	// RequestOTP sends a one-time password to the given email address.
	RequestOTP(ctx context.Context, req OTPRequest) error

	// ResetPassword validates the OTP and replaces the user's password.
	ResetPassword(ctx context.Context, req OTPResetPasswordRequest) error

	// ListAPITokens returns all API tokens linked to the authenticated user.
	// Use page=0 and pageSize=0 to rely on server defaults.
	ListAPITokens(ctx context.Context, page, pageSize int) (*ListApiTokensResponse, error)

	// CreateAPIToken creates a new API token for the authenticated user.
	CreateAPIToken(ctx context.Context, req ApiTokenCreateRequest) (*ApiTokenCreateResponse, error)
}

// ─── Implementation ───────────────────────────────────────────────────────────

type authClient struct {
	http *httpClient
}

func newAuthClient(hc *httpClient) AuthCase {
	return &authClient{http: hc}
}

func (a *authClient) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	var out LoginResponse
	if err := a.http.post(ctx, pathAuthLogin, req, &out); err != nil {
		return nil, fmt.Errorf("synapse/auth.Login: %w", err)
	}
	return &out, nil
}

func (a *authClient) Logout(ctx context.Context, token string) error {
	if err := a.http.post(ctx, pathAuthLogout+token, nil, nil); err != nil {
		return fmt.Errorf("synapse/auth.Logout: %w", err)
	}
	return nil
}

func (a *authClient) Healthcheck(ctx context.Context) (*LoginResponse, error) {
	var out LoginResponse
	if err := a.http.get(ctx, pathAuthHealthcheck, nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/auth.Healthcheck: %w", err)
	}
	return &out, nil
}

func (a *authClient) RequestOTP(ctx context.Context, req OTPRequest) error {
	if err := a.http.post(ctx, pathAuthOTP, req, nil); err != nil {
		return fmt.Errorf("synapse/auth.RequestOTP: %w", err)
	}
	return nil
}

func (a *authClient) ResetPassword(ctx context.Context, req OTPResetPasswordRequest) error {
	if err := a.http.post(ctx, pathAuthPasswordReset, req, nil); err != nil {
		return fmt.Errorf("synapse/auth.ResetPassword: %w", err)
	}
	return nil
}

func (a *authClient) ListAPITokens(ctx context.Context, page, pageSize int) (*ListApiTokensResponse, error) {
	params := url.Values{}
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	if pageSize > 0 {
		params.Set("pageSize", strconv.Itoa(pageSize))
	}

	var out ListApiTokensResponse
	if err := a.http.get(ctx, pathAuthAPIToken, params, &out); err != nil {
		return nil, fmt.Errorf("synapse/auth.ListAPITokens: %w", err)
	}
	return &out, nil
}

func (a *authClient) CreateAPIToken(ctx context.Context, req ApiTokenCreateRequest) (*ApiTokenCreateResponse, error) {
	var out ApiTokenCreateResponse
	if err := a.http.post(ctx, pathAuthAPIToken, req, &out); err != nil {
		return nil, fmt.Errorf("synapse/auth.CreateAPIToken: %w", err)
	}
	return &out, nil
}
