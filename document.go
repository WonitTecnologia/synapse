package synapse

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ─── Interface ────────────────────────────────────────────────────────────────

// DocumentCase provides operations for uploading and managing vectorized documents.
type DocumentCase interface {
	// Upload sends a document file to the server, which parses it, generates embeddings
	// via OpenRouter, and stores the vectors in Qdrant. Processing happens in the background;
	// the returned DocumentResponse will have Status="processing" immediately.
	Upload(ctx context.Context, req UploadDocumentRequest) (*DocumentResponse, error)

	// Get returns the metadata of a document by its UUID.
	Get(ctx context.Context, documentUUID string) (*DocumentResponse, error)

	// List returns a paginated list of documents belonging to the given collection.
	List(ctx context.Context, params ListDocumentsParams) (*ListDocumentsResponse, error)

	// Delete removes a document and all its associated Qdrant vectors.
	Delete(ctx context.Context, documentUUID string) error
}

// ─── Implementation ───────────────────────────────────────────────────────────

type documentClient struct {
	http *httpClient
}

func newDocumentClient(hc *httpClient) DocumentCase {
	return &documentClient{http: hc}
}

func (d *documentClient) Upload(ctx context.Context, req UploadDocumentRequest) (*DocumentResponse, error) {
	fields := map[string]string{
		"collection_uuid": req.CollectionUUID,
		"embed_model":     req.EmbedModel,
	}
	if req.TenantUUID != nil && *req.TenantUUID != "" {
		fields["tenant_uuid"] = *req.TenantUUID
	}
	if req.ChunkSize > 0 {
		fields["chunk_size"] = strconv.Itoa(req.ChunkSize)
	}
	if req.ChunkOverlap > 0 {
		fields["chunk_overlap"] = strconv.Itoa(req.ChunkOverlap)
	}

	var out DocumentResponse
	err := d.http.postMultipart(
		ctx,
		pathKnowledgeDocumentUpload,
		fields,
		"file",
		req.FileName,
		req.Content,
		&out,
	)
	if err != nil {
		return nil, fmt.Errorf("synapse/document.Upload: %w", err)
	}
	return &out, nil
}

func (d *documentClient) Get(ctx context.Context, documentUUID string) (*DocumentResponse, error) {
	var out DocumentResponse
	if err := d.http.get(ctx, pathKnowledgeDocument+"/"+documentUUID, nil, &out); err != nil {
		return nil, fmt.Errorf("synapse/document.Get: %w", err)
	}
	return &out, nil
}

func (d *documentClient) List(ctx context.Context, params ListDocumentsParams) (*ListDocumentsResponse, error) {
	q := url.Values{}
	if params.CollectionUUID != "" {
		q.Set("collection_uuid", params.CollectionUUID)
	}
	if params.Page > 0 {
		q.Set("page", strconv.Itoa(params.Page))
	}
	if params.Size > 0 {
		q.Set("size", strconv.Itoa(params.Size))
	}

	var out ListDocumentsResponse
	if err := d.http.get(ctx, pathKnowledgeDocument, q, &out); err != nil {
		return nil, fmt.Errorf("synapse/document.List: %w", err)
	}
	return &out, nil
}

func (d *documentClient) Delete(ctx context.Context, documentUUID string) error {
	if err := d.http.delete(ctx, pathKnowledgeDocument+"/"+documentUUID, nil); err != nil {
		return fmt.Errorf("synapse/document.Delete: %w", err)
	}
	return nil
}
