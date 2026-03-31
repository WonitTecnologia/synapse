package synapse

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// ─── Sentinel Errors ─────────────────────────────────────────────────────────

var (
	ErrInvalidToken             = errors.New("synapse: token cannot be empty")
	ErrUnauthorized             = errors.New("synapse: unauthorized — invalid or missing token")
	ErrForbidden                = errors.New("synapse: forbidden — insufficient permissions")
	ErrNotFound                 = errors.New("synapse: resource not found")
	ErrConflict                 = errors.New("synapse: conflict — resource already exists or duplicate")
	ErrBadRequest               = errors.New("synapse: bad request — invalid parameters")
	ErrInternalServer           = errors.New("synapse: internal server error")
	ErrBadGateway               = errors.New("synapse: bad gateway — upstream provider error")
	ErrIntegrationNotConfigured = errors.New("synapse: integration not configured for this tenant")
)

// ─── API Error ───────────────────────────────────────────────────────────────

// RestErrCause holds a field-level validation detail returned by the API.
type RestErrCause struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// APIError represents an error response from the Synapse API.
// It carries the full response payload and wraps a typed sentinel error
// so callers can use errors.Is for control flow.
//
// Example:
//
//	_, err := client.Auth.Login(ctx, req)
//	if errors.Is(err, synapse.ErrUnauthorized) { ... }
//
//	if apiErr, ok := synapse.AsAPIError(err); ok {
//	    fmt.Println(apiErr.TraceID, apiErr.Causes)
//	}
type APIError struct {
	// StatusCode is the HTTP status returned by the API.
	StatusCode int `json:"-"`

	// Code mirrors the numeric code in the response body.
	Code int `json:"code"`

	// Err is the short error label from the response body.
	Err string `json:"error"`

	// Message is the human-readable description from the response body.
	Message string `json:"message"`

	// TraceID is the trace identifier returned by the API for debugging.
	TraceID string `json:"trace_id"`

	// Causes holds field-level validation errors when applicable.
	Causes []RestErrCause `json:"causes"`
}

func (e *APIError) Error() string {
	if e.TraceID != "" {
		return fmt.Sprintf("synapse [%d] %s: %s (trace_id=%s)", e.StatusCode, e.Err, e.Message, e.TraceID)
	}
	return fmt.Sprintf("synapse [%d] %s: %s", e.StatusCode, e.Err, e.Message)
}

// Unwrap maps the HTTP status code to the appropriate sentinel error,
// enabling errors.Is checks against the package-level sentinels.
func (e *APIError) Unwrap() error {
	switch e.StatusCode {
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusForbidden:
		return ErrForbidden
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusConflict:
		return ErrConflict
	case http.StatusBadRequest:
		return ErrBadRequest
	case http.StatusInternalServerError:
		return ErrInternalServer
	case http.StatusBadGateway:
		return ErrBadGateway
	default:
		return nil
	}
}

// AsAPIError unwraps err into *APIError and reports whether it succeeded.
// Use this to inspect the full API error payload (causes, trace_id, etc.).
func AsAPIError(err error) (*APIError, bool) {
	var apiErr *APIError
	ok := errors.As(err, &apiErr)
	return apiErr, ok
}

// parseAPIError deserialises the response body into an *APIError.
func parseAPIError(statusCode int, body []byte) *APIError {
	apiErr := &APIError{StatusCode: statusCode}
	_ = json.Unmarshal(body, apiErr)
	if apiErr.Err == "" {
		apiErr.Err = http.StatusText(statusCode)
	}
	return apiErr
}
