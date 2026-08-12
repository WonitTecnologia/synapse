package synapse

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// ─── Shared WebSocket helpers ─────────────────────────────────────────────────
//
// Used by the domain streams (monitor, chat stream). They encapsulate the
// common connection behaviour: URL conversion, Bearer authentication, the
// Options.Host override (Host header + TLS ServerName) and keepalive.

// buildWSURL converts the API base URL into a WebSocket URL for the given
// path, carrying the session used to resume the server-side delivery queue.
func (c *httpClient) buildWSURL(path, session string) (string, error) {
	u, err := url.Parse(c.baseURL)
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
	u.Path = strings.TrimRight(u.Path, "/") + path
	q := u.Query()
	q.Set("session", session)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// dialWS opens the WebSocket connection with the Bearer token and the
// Options.Host override (Host header and TLS ServerName on wss).
func (c *httpClient) dialWS(ctx context.Context, wsURL string) (*websocket.Conn, error) {
	header := http.Header{}
	header.Set("User-Agent", "synapse-sdk")
	if c.token != "" {
		header.Set("Authorization", "Bearer "+c.token)
	}
	// The Host header key is honoured by the dialer, enabling the internal
	// IP + virtual host connection mode (see Options.Host).
	if c.host != "" {
		header.Set("Host", c.host)
	}

	dialer := websocket.Dialer{HandshakeTimeout: wsHandshakeTimeout}
	if c.host != "" && strings.HasPrefix(wsURL, "wss://") {
		dialer.TLSClientConfig = &tls.Config{ServerName: c.host}
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

// prepareWSConn sets the initial read deadline and the ping handler that
// answers server pings with pongs while extending the read deadline.
func prepareWSConn(conn *websocket.Conn) {
	_ = conn.SetReadDeadline(time.Now().Add(wsReadWait))
	conn.SetPingHandler(func(appData string) error {
		_ = conn.SetReadDeadline(time.Now().Add(wsReadWait))
		err := conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(wsWriteWait))
		if err == websocket.ErrCloseSent {
			return nil
		}
		return err
	})
}
