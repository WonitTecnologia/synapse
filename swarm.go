package synapse

import (
	"context"
	"fmt"
)

// ─── Interface ────────────────────────────────────────────────────────────────

// SwarmCase queries and controls swarms of the Builder system agent.
//
// A swarm is a set of ephemeral subagents the orchestrator opens via the
// sistema_iniciar_swarm tool to decompose a large task: each task runs in
// parallel as its own dispatch job, with zeroed conversation history and no
// write tools. When every task reaches a terminal status, the reports are
// consolidated back into the main conversation as a synthetic turn.
//
// Live progress is delivered as AgentEvent events of category
// EventCategorySwarm (see Monitor); this case covers the REST side: state
// query and stop/pause/resume control. Use a TENANT token: swarms are scoped
// to the caller tenant, and there is at most one active swarm per conversation.
type SwarmCase interface {
	// GetSwarm returns the full state of a swarm — tasks, counters and the
	// transition timeline (History) — for rendering the swarm panel.
	GetSwarm(ctx context.Context, swarmID string) (*SwarmState, error)

	// StopSwarm cancels the whole swarm: every child job is cancelled (in-flight
	// LLM calls die immediately) and the consolidation runs with whatever
	// reports exist. Returns the updated swarm state.
	StopSwarm(ctx context.Context, swarmID string) (*SwarmState, error)

	// PauseSwarm freezes the swarm: running tasks are cancelled (no token spend)
	// and queued tasks wait for a resume. Returns the updated swarm state.
	PauseSwarm(ctx context.Context, swarmID string) (*SwarmState, error)

	// ResumeSwarm clears the pause flag and re-enqueues the queued/interrupted
	// tasks as new jobs. Returns the updated swarm state.
	ResumeSwarm(ctx context.Context, swarmID string) (*SwarmState, error)
}

// ─── Implementation ───────────────────────────────────────────────────────────

type swarmClient struct {
	http *httpClient
}

func newSwarmClient(hc *httpClient) SwarmCase {
	return &swarmClient{http: hc}
}

func (s *swarmClient) GetSwarm(ctx context.Context, swarmID string) (*SwarmState, error) {
	var out SwarmState
	if err := s.http.get(ctx, fmt.Sprintf(pathSystemAgentSwarm, swarmID), nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/swarm.GetSwarm: %w", err)
	}
	return &out, nil
}

func (s *swarmClient) StopSwarm(ctx context.Context, swarmID string) (*SwarmState, error) {
	var out SwarmState
	if err := s.http.post(ctx, fmt.Sprintf(pathSystemAgentSwarmStop, swarmID), nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/swarm.StopSwarm: %w", err)
	}
	return &out, nil
}

func (s *swarmClient) PauseSwarm(ctx context.Context, swarmID string) (*SwarmState, error) {
	var out SwarmState
	if err := s.http.post(ctx, fmt.Sprintf(pathSystemAgentSwarmPause, swarmID), nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/swarm.PauseSwarm: %w", err)
	}
	return &out, nil
}

func (s *swarmClient) ResumeSwarm(ctx context.Context, swarmID string) (*SwarmState, error) {
	var out SwarmState
	if err := s.http.post(ctx, fmt.Sprintf(pathSystemAgentSwarmResume, swarmID), nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/swarm.ResumeSwarm: %w", err)
	}
	return &out, nil
}
