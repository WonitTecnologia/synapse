package synapse

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// DispatchCase provides access to the dispatch queue observability endpoints.
type DispatchCase interface {
	// GetQueueStats returns the state of the jobs and outputs Redis streams,
	// including active consumers (workers) and their idle times.
	GetQueueStats(ctx context.Context) (*QueueStatsResponse, error)
	// ListJobs returns pending and processing jobs. Use params.TenantUUID to
	// filter by tenant; params.Count defaults to 50 (max 200).
	ListJobs(ctx context.Context, params ListQueuedJobsParams) (*PagedJobs, error)
	// DeleteJob removes a job from the stream (XDEL). If the job is in PEL,
	// XACK is also issued. The redisID must be the Redis Stream message ID
	// (e.g. "1783450838712-0").
	DeleteJob(ctx context.Context, redisID string) error
}

type dispatchClient struct {
	http *httpClient
}

func newDispatchClient(hc *httpClient) DispatchCase {
	return &dispatchClient{http: hc}
}

func (d *dispatchClient) GetQueueStats(ctx context.Context) (*QueueStatsResponse, error) {
	var out QueueStatsResponse
	if err := d.http.get(ctx, pathDispatchQueueStats, nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/dispatch.GetQueueStats: %w", err)
	}
	return &out, nil
}

func (d *dispatchClient) ListJobs(ctx context.Context, params ListQueuedJobsParams) (*PagedJobs, error) {
	q := url.Values{}
	if params.TenantUUID != "" {
		q.Set("tenant_uuid", params.TenantUUID)
	}
	if params.Cursor != "" {
		q.Set("cursor", params.Cursor)
	}
	if params.Count > 0 {
		q.Set("count", strconv.FormatInt(params.Count, 10))
	}
	var out PagedJobs
	if err := d.http.get(ctx, pathDispatchQueueJobs, q, &out); err != nil {
		return nil, fmt.Errorf("synapse/dispatch.ListJobs: %w", err)
	}
	return &out, nil
}

func (d *dispatchClient) DeleteJob(ctx context.Context, redisID string) error {
	path := fmt.Sprintf(pathDispatchQueueJobDelete, redisID)
	if err := d.http.delete(ctx, path, nil); err != nil {
		return fmt.Errorf("synapse/dispatch.DeleteJob: %w", err)
	}
	return nil
}
