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
// query (by swarm ID or by conversation), subagent trace and stop/pause/resume
// control. Use a TENANT token: swarms are scoped
// to the caller tenant, and there is at most one active swarm per conversation.
type SwarmCase interface {
	// GetSwarm returns the full state of a swarm — tasks, counters and the
	// transition timeline (History) — for rendering the swarm panel.
	GetSwarm(ctx context.Context, swarmID string) (*SwarmState, error)

	// GetSwarmByConversation returns the state of the conversation's ACTIVE
	// swarm — same payload as GetSwarm — for rehydrating the swarm panel after
	// a page reload, when only the conversation UUID is known. Returns a 404
	// *APIError when the conversation has no active swarm.
	GetSwarmByConversation(ctx context.Context, conversationUUID string) (*SwarmState, error)

	// GetSwarmTrace returns the execution trace of one task's subagent, kept
	// server-side in Redis for 24h: the ordered entries the subagent produced
	// while working on the task, for post-mortem inspection and for
	// rehydrating the task's live view after a page reload.
	GetSwarmTrace(ctx context.Context, swarmID, taskID string) (*SwarmTrace, error)

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

func (s *swarmClient) GetSwarmByConversation(ctx context.Context, conversationUUID string) (*SwarmState, error) {
	var out SwarmState
	if err := s.http.get(ctx, fmt.Sprintf(pathSystemAgentSwarmByConversation, conversationUUID), nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/swarm.GetSwarmByConversation: %w", err)
	}
	return &out, nil
}

func (s *swarmClient) GetSwarmTrace(ctx context.Context, swarmID, taskID string) (*SwarmTrace, error) {
	var out SwarmTrace
	if err := s.http.get(ctx, fmt.Sprintf(pathSystemAgentSwarmTaskTrace, swarmID, taskID), nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/swarm.GetSwarmTrace: %w", err)
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
