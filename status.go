package synapse

import (
	"context"
	"fmt"
	"time"
)

// ─── Interface ────────────────────────────────────────────────────────────────

// StatusCase provides access to the service status endpoint.
type StatusCase interface {
	// Get returns build info, pod, dependency health (Postgres/Redis/Qdrant),
	// the validity of the token used on the request and the server-side
	// processing time.
	Get(ctx context.Context) (*ServiceStatus, error)
}

// ─── DTOs ─────────────────────────────────────────────────────────────────────

// StatusCheck is the health of a single dependency.
type StatusCheck struct {
	Status    string `json:"status"` // online | offline | disabled
	LatencyMs int64  `json:"latency_ms"`
	Error     string `json:"error,omitempty"`
}

// StatusToken describes the validity of the token used on the request.
type StatusToken struct {
	ExpiresAt        time.Time `json:"expires_at"`
	ExpiresInSeconds int64     `json:"expires_in_seconds"`
}

// ServiceStatus is the full status payload of the Synapse API.
type ServiceStatus struct {
	Service        string      `json:"service"`
	Version        string      `json:"version"`
	Commit         string      `json:"commit,omitempty"`
	BuildTime      string      `json:"build_time,omitempty"`
	GoVersion      string      `json:"go_version"`
	Pod            string      `json:"pod"`
	UptimeSeconds  int64       `json:"uptime_seconds"`
	SystemTimeUTC  time.Time   `json:"system_time_utc"`
	Database       StatusCheck `json:"database"`
	Redis          StatusCheck `json:"redis"`
	Qdrant         StatusCheck `json:"qdrant"`
	Token          StatusToken `json:"token"`
	ResponseTimeMs int64       `json:"response_time_ms"`
}

// ─── Implementation ───────────────────────────────────────────────────────────

type statusClient struct {
	http *httpClient
}

func newStatusClient(hc *httpClient) StatusCase {
	return &statusClient{http: hc}
}

func (s *statusClient) Get(ctx context.Context) (*ServiceStatus, error) {
	var out ServiceStatus
	if err := s.http.get(ctx, pathStatus, nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/status.Get: %w", err)
	}
	return &out, nil
}
