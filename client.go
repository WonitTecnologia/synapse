package synapse

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

const defaultTimeout = 30 * time.Second

// httpClient is the internal transport layer shared by all domain clients.
// It handles auth headers, serialisation, error parsing, and multipart uploads.
type httpClient struct {
	token   string
	baseURL string
	http    *http.Client
}

func newHTTPClient(token, baseURL string, timeout time.Duration) *httpClient {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	if timeout == 0 {
		timeout = defaultTimeout
	}
	return &httpClient{
		token:   token,
		baseURL: baseURL,
		http:    &http.Client{Timeout: timeout},
	}
}

// ─── Core request ─────────────────────────────────────────────────────────────

func (c *httpClient) do(
	ctx context.Context,
	method, path string,
	body any,
	queryParams url.Values,
) ([]byte, int, error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("synapse: marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	fullURL := c.baseURL + path
	if len(queryParams) > 0 {
		fullURL += "?" + queryParams.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("synapse: build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("synapse: execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("synapse: read response body: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, resp.StatusCode, parseAPIError(resp.StatusCode, respBody)
	}

	return respBody, resp.StatusCode, nil
}

// ─── JSON helpers ─────────────────────────────────────────────────────────────

func (c *httpClient) get(ctx context.Context, path string, params url.Values, out any) error {
	body, _, err := c.do(ctx, http.MethodGet, path, nil, params)
	if err != nil {
		return err
	}
	if out != nil && len(body) > 0 {
		if err := json.Unmarshal(body, out); err != nil {
			return fmt.Errorf("synapse: unmarshal response: %w", err)
		}
	}
	return nil
}

func (c *httpClient) post(ctx context.Context, path string, payload, out any) error {
	body, _, err := c.do(ctx, http.MethodPost, path, payload, nil)
	if err != nil {
		return err
	}
	if out != nil && len(body) > 0 {
		if err := json.Unmarshal(body, out); err != nil {
			return fmt.Errorf("synapse: unmarshal response: %w", err)
		}
	}
	return nil
}

func (c *httpClient) put(ctx context.Context, path string, payload, out any) error {
	body, _, err := c.do(ctx, http.MethodPut, path, payload, nil)
	if err != nil {
		return err
	}
	if out != nil && len(body) > 0 {
		if err := json.Unmarshal(body, out); err != nil {
			return fmt.Errorf("synapse: unmarshal response: %w", err)
		}
	}
	return nil
}

func (c *httpClient) patch(ctx context.Context, path string, payload, out any) error {
	body, _, err := c.do(ctx, http.MethodPatch, path, payload, nil)
	if err != nil {
		return err
	}
	if out != nil && len(body) > 0 {
		if err := json.Unmarshal(body, out); err != nil {
			return fmt.Errorf("synapse: unmarshal response: %w", err)
		}
	}
	return nil
}

func (c *httpClient) delete(ctx context.Context, path string, params url.Values) error {
	_, _, err := c.do(ctx, http.MethodDelete, path, nil, params)
	return err
}

func (c *httpClient) deleteJSON(ctx context.Context, path string, payload, out any) error {
	body, _, err := c.do(ctx, http.MethodDelete, path, payload, nil)
	if err != nil {
		return err
	}
	if out != nil && len(body) > 0 {
		if err := json.Unmarshal(body, out); err != nil {
			return fmt.Errorf("synapse: unmarshal response: %w", err)
		}
	}
	return nil
}

// ─── Multipart helper ─────────────────────────────────────────────────────────

// postMultipart sends a multipart/form-data POST request.
// fields contains plain text form fields; fileField/fileName/fileContent describe the file part.
func (c *httpClient) postMultipart(
	ctx context.Context,
	path string,
	fields map[string]string,
	fileField, fileName string,
	fileContent []byte,
	out any,
) error {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	for key, val := range fields {
		if val == "" {
			continue
		}
		if err := w.WriteField(key, val); err != nil {
			return fmt.Errorf("synapse: write multipart field %q: %w", key, err)
		}
	}

	part, err := w.CreateFormFile(fileField, fileName)
	if err != nil {
		return fmt.Errorf("synapse: create form file: %w", err)
	}
	if _, err = part.Write(fileContent); err != nil {
		return fmt.Errorf("synapse: write file content: %w", err)
	}
	w.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, &buf)
	if err != nil {
		return fmt.Errorf("synapse: build multipart request: %w", err)
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Accept", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("synapse: execute multipart request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("synapse: read multipart response: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return parseAPIError(resp.StatusCode, respBody)
	}

	if out != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, out); err != nil {
			return fmt.Errorf("synapse: unmarshal multipart response: %w", err)
		}
	}
	return nil
}
