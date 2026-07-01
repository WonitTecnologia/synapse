package synapse

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// ─── Interface ────────────────────────────────────────────────────────────────

// MonitorCase provides access to the real-time agent event WebSocket.
//
// The stream is receive-only: the server pushes agent execution events (chat,
// tool_call/MCP, RAG, errors, file processing) and the SDK automatically
// acknowledges each delivery. A master (SYSTEM_ADMIN) token receives events
// from every tenant; a tenant token receives only its own tenant's events.
type MonitorCase interface {
	// StreamLogs opens the monitor WebSocket and keeps it alive, reconnecting
	// automatically with the same session until ctx is cancelled or the
	// returned stream is closed. Events are delivered on EventStream.Events().
	StreamLogs(ctx context.Context, opts *StreamLogsOptions) (*EventStream, error)
}

// ─── Options and DTOs ─────────────────────────────────────────────────────────

// StreamLogsOptions tunes the monitor stream. All fields are optional.
type StreamLogsOptions struct {
	// Session (UUID) resumes the same server-side delivery queue across
	// reconnections. Defaults to a random UUID kept for the stream lifetime.
	Session string

	// Buffer is the capacity of the events channel (default 256).
	Buffer int

	// OnError, when set, is invoked with connection/decoding errors. The
	// stream keeps reconnecting regardless; this is for observability only.
	OnError func(error)
}

// Agent event categories delivered by the monitor stream.
const (
	EventCategoryChat        = "chat"
	EventCategoryToolCall    = "tool_call"
	EventCategoryRAG         = "rag"
	EventCategoryError       = "error"
	EventCategoryFileProcess = "file_process"
)

// AgentEventRagChunk is one retrieved RAG chunk with its source file and
// similarity percentage.
type AgentEventRagChunk struct {
	Filename   string  `json:"filename"`
	ChunkIndex int     `json:"chunk_index"`
	Score      float32 `json:"score"`
	ScorePct   float32 `json:"score_pct"`
	Text       string  `json:"text,omitempty"`
}

// AgentEventRag aggregates the RAG stage outcome of a turn.
type AgentEventRag struct {
	ChunksFound int                  `json:"chunks_found"`
	Error       string               `json:"error,omitempty"`
	Chunks      []AgentEventRagChunk `json:"chunks,omitempty"`
}

// AgentEventTokens aggregates the token usage of a turn.
type AgentEventTokens struct {
	Prompt     int `json:"prompt"`
	Completion int `json:"completion"`
	Total      int `json:"total"`
	Embedding  int `json:"embedding"`
}

// AgentEvent is a real-time agent execution event. Unlike the persisted agent
// log, tool parameters, tool results and API responses are NOT truncated.
type AgentEvent struct {
	UUID                   string            `json:"uuid"`
	TenantUUID             string            `json:"tenant_uuid"`
	AgentUUID              string            `json:"agent_uuid"`
	AgentName              string            `json:"agent_name"`
	ConversationUUID       *string           `json:"conversation_uuid,omitempty"`
	ConversationExternalID *string           `json:"conversation_external_id,omitempty"`
	Level                  string            `json:"level"`
	Category               string            `json:"category"`
	Summary                string            `json:"summary"`
	Detail                 map[string]any    `json:"detail,omitempty"`
	ToolName               *string           `json:"tool_name,omitempty"`
	ToolParams             map[string]any    `json:"tool_params,omitempty"`
	ToolSuccess            *bool             `json:"tool_success,omitempty"`
	ToolResult             string            `json:"tool_result,omitempty"`
	APIResponse            string            `json:"api_response,omitempty"`
	Rag                    *AgentEventRag    `json:"rag,omitempty"`
	DurationMs             *int              `json:"duration_ms,omitempty"`
	Model                  *string           `json:"model,omitempty"`
	Tokens                 *AgentEventTokens `json:"tokens,omitempty"`
	CreatedAt              time.Time         `json:"created_at"`
}

// wsEnvelope is the server → client delivery frame; every envelope must be
// acknowledged with wsAck for its UUID (handled automatically by the SDK).
type wsEnvelope struct {
	Type   string     `json:"type"`
	UUID   string     `json:"uuid"`
	SentAt time.Time  `json:"sent_at"`
	Event  AgentEvent `json:"event"`
}

type wsAck struct {
	Type string `json:"type"`
	UUID string `json:"uuid"`
}

// ─── EventStream ──────────────────────────────────────────────────────────────

// EventStream is a live subscription to the monitor WebSocket. Consume events
// from Events(); the channel is closed after Close() or context cancellation.
type EventStream struct {
	events  chan AgentEvent
	session string
	cancel  context.CancelFunc
}

// Events returns the channel where incoming agent events are delivered.
func (s *EventStream) Events() <-chan AgentEvent { return s.events }

// Session returns the session UUID used to resume the delivery queue.
func (s *EventStream) Session() string { return s.session }

// Close stops the stream and closes the events channel.
func (s *EventStream) Close() { s.cancel() }

// ─── Implementation ───────────────────────────────────────────────────────────

const (
	wsHandshakeTimeout = 15 * time.Second
	wsReadWait         = 120 * time.Second
	wsWriteWait        = 5 * time.Second
	wsBackoffMin       = time.Second
	wsBackoffMax       = 30 * time.Second
	wsDefaultBuffer    = 256
)

type monitorClient struct {
	http *httpClient
}

func newMonitorClient(hc *httpClient) MonitorCase {
	return &monitorClient{http: hc}
}

func (m *monitorClient) StreamLogs(ctx context.Context, opts *StreamLogsOptions) (*EventStream, error) {
	if opts == nil {
		opts = &StreamLogsOptions{}
	}

	session := opts.Session
	if session == "" {
		var err error
		session, err = randomUUID()
		if err != nil {
			return nil, fmt.Errorf("synapse/monitor.StreamLogs: generate session: %w", err)
		}
	}

	wsURL, err := m.buildURL(session)
	if err != nil {
		return nil, fmt.Errorf("synapse/monitor.StreamLogs: %w", err)
	}

	buffer := opts.Buffer
	if buffer <= 0 {
		buffer = wsDefaultBuffer
	}

	streamCtx, cancel := context.WithCancel(ctx)
	stream := &EventStream{
		events:  make(chan AgentEvent, buffer),
		session: session,
		cancel:  cancel,
	}

	go m.run(streamCtx, stream, wsURL, opts.OnError)
	return stream, nil
}

// buildURL converts the API base URL into the monitor WebSocket URL.
func (m *monitorClient) buildURL(session string) (string, error) {
	u, err := url.Parse(m.http.baseURL)
	if err != nil {
		return "", fmt.Errorf("parse base URL: %w", err)
	}
	switch u.Scheme {
	case "http":
		u.Scheme = "ws"
	case "https":
		u.Scheme = "wss"
	case "ws", "wss":
		// already a websocket URL
	default:
		return "", fmt.Errorf("unsupported base URL scheme %q", u.Scheme)
	}
	u.Path = strings.TrimRight(u.Path, "/") + pathMonitorLogsWS
	q := u.Query()
	q.Set("session", session)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// run keeps the subscription alive: dial, consume, and reconnect with backoff
// until the context is cancelled. The same session resumes the server queue.
func (m *monitorClient) run(ctx context.Context, stream *EventStream, wsURL string, onError func(error)) {
	defer close(stream.events)

	report := func(err error) {
		if onError != nil && err != nil {
			onError(err)
		}
	}

	backoff := wsBackoffMin
	for {
		if ctx.Err() != nil {
			return
		}

		conn, err := m.dial(ctx, wsURL)
		if err != nil {
			report(fmt.Errorf("synapse/monitor: connect: %w", err))
		} else {
			backoff = wsBackoffMin
			err = m.consume(ctx, conn, stream)
			_ = conn.Close()
			if ctx.Err() != nil {
				return
			}
			report(fmt.Errorf("synapse/monitor: connection lost: %w", err))
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(backoff):
		}
		backoff *= 2
		if backoff > wsBackoffMax {
			backoff = wsBackoffMax
		}
	}
}

func (m *monitorClient) dial(ctx context.Context, wsURL string) (*websocket.Conn, error) {
	header := http.Header{}
	header.Set("User-Agent", "synapse-sdk")
	if m.http.token != "" {
		header.Set("Authorization", "Bearer "+m.http.token)
	}
	// The Host header key is honoured by the dialer, enabling the internal
	// IP + virtual host connection mode (see Options.Host).
	if m.http.host != "" {
		header.Set("Host", m.http.host)
	}

	dialer := websocket.Dialer{HandshakeTimeout: wsHandshakeTimeout}
	if m.http.host != "" && strings.HasPrefix(wsURL, "wss://") {
		dialer.TLSClientConfig = &tls.Config{ServerName: m.http.host}
	}

	conn, resp, err := dialer.DialContext(ctx, wsURL, header)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		if resp != nil {
			return nil, fmt.Errorf("%w (status %d)", err, resp.StatusCode)
		}
		return nil, err
	}
	return conn, nil
}

// consume reads envelopes until the connection drops, acknowledging each one
// and delivering the events to the stream channel.
func (m *monitorClient) consume(ctx context.Context, conn *websocket.Conn, stream *EventStream) error {
	_ = conn.SetReadDeadline(time.Now().Add(wsReadWait))
	conn.SetPingHandler(func(appData string) error {
		_ = conn.SetReadDeadline(time.Now().Add(wsReadWait))
		err := conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(wsWriteWait))
		if err == websocket.ErrCloseSent {
			return nil
		}
		return err
	})

	// Unblock the read loop as soon as the context is cancelled.
	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-ctx.Done():
			_ = conn.Close()
		case <-done:
		}
	}()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		_ = conn.SetReadDeadline(time.Now().Add(wsReadWait))

		var env wsEnvelope
		if err := json.Unmarshal(data, &env); err != nil || env.Type != "event" {
			continue
		}

		// Acknowledge this envelope; without the "ok" the server retries.
		_ = conn.SetWriteDeadline(time.Now().Add(wsWriteWait))
		if err := conn.WriteJSON(wsAck{Type: "ok", UUID: env.UUID}); err != nil {
			return err
		}

		select {
		case stream.events <- env.Event:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// randomUUID generates a version 4 UUID without external dependencies.
func randomUUID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // RFC 4122 variant
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}
