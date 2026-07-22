package synapse

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ─── Interface ────────────────────────────────────────────────────────────────

// OpenRouterCase provides operations for the OpenRouter integration.
type OpenRouterCase interface {
	// Sync creates or returns an existing OpenRouter workspace for the authenticated tenant.
	// SystemAdmin callers may target a specific tenant by setting req.TenantUUID.
	Sync(ctx context.Context, req OpenRouterSyncRequest) (*OpenRouterSyncResponse, error)

	// Desync removes the OpenRouter workspace for the authenticated tenant.
	// SystemAdmin callers may target a specific tenant by setting req.TenantUUID.
	Desync(ctx context.Context, req OpenRouterSyncRequest) (*OpenRouterDesyncResponse, error)

	// ListModels returns all models available on OpenRouter.
	// Set freeOnly=true to filter for free-tier models only.
	ListModels(ctx context.Context, freeOnly bool) (*OpenRouterListModelsResponse, error)

	// ListEmbeddingModels fetches live embedding models from the OpenRouter API,
	// filtered by output modality. Each entry includes vector_size (0 = unknown) and
	// context_length so callers can configure Qdrant collections and chunk sizes correctly.
	// Set freeOnly=true to filter for free-tier models only.
	ListEmbeddingModels(ctx context.Context, freeOnly bool) (*OpenRouterListEmbeddingModelsResponse, error)

	// GetMonthlyAnalytics returns tokens_total and total_usage per model/day of a
	// month for the workspace linked to the authenticated tenant.
	GetMonthlyAnalytics(ctx context.Context, params OpenRouterMonthlyAnalyticsParams) (*OpenRouterMonthlyAnalyticsResponse, error)

	// GetAllTenantsMonthlyAnalytics returns the monthly OpenRouter spend (USD),
	// tokens and requests aggregated per tenant, sorted by cost descending, plus
	// the month totals. SYSTEM_ADMIN only — the workspace fan-out runs server-side
	// with the master key.
	GetAllTenantsMonthlyAnalytics(ctx context.Context, params OpenRouterAllTenantsMonthlyParams) (*OpenRouterAllTenantsMonthlyResponse, error)

	// QueryAnalytics runs an analytics query scoped to the tenant's workspace.
	// The workspace filter is applied server-side and must not be present in
	// req.Filters. Metrics/dimensions outside the daily materialized view (only
	// tokens_total, total_usage and the model dimension are MV-compatible) are
	// limited to a 31-day time range, as are the minute and hour granularities.
	QueryAnalytics(ctx context.Context, req OpenRouterAnalyticsQueryRequest) (*OpenRouterAnalyticsResult, error)

	// GetAnalyticsMeta returns the metrics, dimensions, filter operators and
	// granularities available for analytics queries on the tenant's workspace.
	GetAnalyticsMeta(ctx context.Context) (*OpenRouterAnalyticsMeta, error)
}

// OpenRouterMonthlyAnalyticsParams are the optional parameters for GetMonthlyAnalytics.
type OpenRouterMonthlyAnalyticsParams struct {
	// Month in "YYYY-MM" format. Empty = current month.
	Month string
	// Limit caps the number of rows (default 1000, max 1000).
	Limit int
	// TenantUUID targets another tenant's workspace. SYSTEM_ADMIN only; empty uses
	// the caller's own tenant. Used by the Master panel to drill into one company.
	TenantUUID string
}

// OpenRouterAllTenantsMonthlyParams are the optional parameters for GetAllTenantsMonthlyAnalytics.
type OpenRouterAllTenantsMonthlyParams struct {
	// Month in "YYYY-MM" format. Empty = current month.
	Month string
}

// ─── Implementation ───────────────────────────────────────────────────────────

type openrouterClient struct {
	http *httpClient
}

func newOpenRouterClient(hc *httpClient) OpenRouterCase {
	return &openrouterClient{http: hc}
}

func (o *openrouterClient) Sync(ctx context.Context, req OpenRouterSyncRequest) (*OpenRouterSyncResponse, error) {
	var out OpenRouterSyncResponse
	if err := o.http.post(ctx, pathOpenRouterSincronismo, req, &out); err != nil {
		return nil, fmt.Errorf("synapse/openrouter.Sync: %w", err)
	}
	return &out, nil
}

func (o *openrouterClient) Desync(ctx context.Context, req OpenRouterSyncRequest) (*OpenRouterDesyncResponse, error) {
	var out OpenRouterDesyncResponse
	if err := o.http.deleteJSON(ctx, pathOpenRouterSincronismo, req, &out); err != nil {
		return nil, fmt.Errorf("synapse/openrouter.Desync: %w", err)
	}
	return &out, nil
}

func (o *openrouterClient) ListModels(ctx context.Context, freeOnly bool) (*OpenRouterListModelsResponse, error) {
	q := url.Values{}
	if freeOnly {
		q.Set("free_only", "true")
	}
	var out OpenRouterListModelsResponse
	if err := o.http.get(ctx, pathOpenRouterModels, q, &out); err != nil {
		return nil, fmt.Errorf("synapse/openrouter.ListModels: %w", err)
	}
	return &out, nil
}

func (o *openrouterClient) ListEmbeddingModels(ctx context.Context, freeOnly bool) (*OpenRouterListEmbeddingModelsResponse, error) {
	q := url.Values{}
	if freeOnly {
		q.Set("free_only", "true")
	}
	var out OpenRouterListEmbeddingModelsResponse
	if err := o.http.get(ctx, pathOpenRouterEmbeddingModels, q, &out); err != nil {
		return nil, fmt.Errorf("synapse/openrouter.ListEmbeddingModels: %w", err)
	}
	return &out, nil
}

func (o *openrouterClient) GetMonthlyAnalytics(ctx context.Context, params OpenRouterMonthlyAnalyticsParams) (*OpenRouterMonthlyAnalyticsResponse, error) {
	q := url.Values{}
	if params.Month != "" {
		q.Set("month", params.Month)
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.TenantUUID != "" {
		q.Set("tenant_uuid", params.TenantUUID)
	}
	var out OpenRouterMonthlyAnalyticsResponse
	if err := o.http.get(ctx, pathOpenRouterAnalyticsMonthly, q, &out); err != nil {
		return nil, fmt.Errorf("synapse/openrouter.GetMonthlyAnalytics: %w", err)
	}
	return &out, nil
}

func (o *openrouterClient) GetAllTenantsMonthlyAnalytics(ctx context.Context, params OpenRouterAllTenantsMonthlyParams) (*OpenRouterAllTenantsMonthlyResponse, error) {
	q := url.Values{}
	if params.Month != "" {
		q.Set("month", params.Month)
	}
	var out OpenRouterAllTenantsMonthlyResponse
	if err := o.http.get(ctx, pathOpenRouterAnalyticsMonthlyAll, q, &out); err != nil {
		return nil, fmt.Errorf("synapse/openrouter.GetAllTenantsMonthlyAnalytics: %w", err)
	}
	return &out, nil
}

func (o *openrouterClient) QueryAnalytics(ctx context.Context, req OpenRouterAnalyticsQueryRequest) (*OpenRouterAnalyticsResult, error) {
	var out OpenRouterAnalyticsResult
	if err := o.http.post(ctx, pathOpenRouterAnalyticsQuery, req, &out); err != nil {
		return nil, fmt.Errorf("synapse/openrouter.QueryAnalytics: %w", err)
	}
	return &out, nil
}

func (o *openrouterClient) GetAnalyticsMeta(ctx context.Context) (*OpenRouterAnalyticsMeta, error) {
	var out OpenRouterAnalyticsMeta
	if err := o.http.get(ctx, pathOpenRouterAnalyticsMeta, nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/openrouter.GetAnalyticsMeta: %w", err)
	}
	return &out, nil
}
