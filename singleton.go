package synapse

import "time"

// ─── Options ──────────────────────────────────────────────────────────────────

// Options configures the Synapse client at creation time.
// All fields are optional; zero values fall back to library defaults.
type Options struct {
	// BaseURL overrides the default API base URL (https://synapse.wonit.net.br).
	// Useful for staging environments or self-hosted deployments.
	BaseURL string

	// Timeout sets the maximum duration for each HTTP request.
	// Defaults to 30 seconds when zero.
	Timeout time.Duration
}

// ─── Client ───────────────────────────────────────────────────────────────────

// Client is the top-level Synapse SDK entry point.
// Each field exposes a domain-specific interface covering a group of API endpoints.
//
// Usage:
//
//	client, err := synapse.NewClient("your-token", nil)
//
//	// Auth
//	resp, err := client.Auth.Login(ctx, synapse.LoginRequest{...})
//
//	// Tenant
//	tenant, err := client.Tenant.Get(ctx, "uuid", "")
//
//	// Google Vision
//	err = client.Google.Configure(ctx, synapse.ConfigureGoogleRequest{...})
//	ocr, err := client.Google.VisionOCR(ctx, "photo.jpg", imageBytes)
//
//	// OpenAI
//	reply, err := client.OpenAI.Chat(ctx, synapse.ChatCompletionRequest{Prompt: "Hello!"})
//
//	// Chatvolt
//	answer, err := client.Chatvolt.Query(ctx, synapse.ChatvoltAgentQueryRequest{...})
type Client struct {
	// Auth covers login, logout, OTP, password reset, and API token management.
	Auth AuthCase

	// User covers user CRUD operations.
	User UserCase

	// Tenant covers tenant CRUD operations.
	Tenant TenantCase

	// Provider covers read operations for catalog integration providers.
	Provider ProviderCase

	// Service covers read operations for catalog integration services.
	Service ServiceCase

	// Google covers the Google Vision AI integration (configure + OCR).
	Google GoogleCase

	// OpenAI covers the OpenAI integration (configure, chat, image analysis, transcription).
	OpenAI OpenAICase

	// Chatvolt covers the Chatvolt agent integration (configure + query).
	Chatvolt ChatvoltCase
}

// NewClient creates and returns a fully initialised Synapse Client.
//
// token is required and is sent as a Bearer token on every request.
// Pass an *Options to override the base URL or request timeout; pass nil for defaults.
//
//	client, err := synapse.NewClient("sk-...", &synapse.Options{
//	    BaseURL: "https://staging.synapse.example.com",
//	    Timeout: 15 * time.Second,
//	})
func NewClient(token string, opts *Options) (*Client, error) {
	if token == "" {
		return nil, ErrInvalidToken
	}

	var baseURL string
	var timeout time.Duration

	if opts != nil {
		baseURL = opts.BaseURL
		timeout = opts.Timeout
	}

	hc := newHTTPClient(token, baseURL, timeout)

	return &Client{
		Auth:     newAuthClient(hc),
		User:     newUserClient(hc),
		Tenant:   newTenantClient(hc),
		Provider: newProviderClient(hc),
		Service:  newServiceClient(hc),
		Google:   newGoogleClient(hc),
		OpenAI:   newOpenAIClient(hc),
		Chatvolt: newChatvoltClient(hc),
	}, nil
}
