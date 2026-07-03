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

	// Estimate returns a cost estimate for embedding the file with the given model.
	// Pre-processes text, counts pages/images, and looks up model pricing via OpenRouter.
	Estimate(ctx context.Context, req EstimateDocumentRequest) (*EstimateDocumentResponse, error)

	// ListChunks returns a paginated list of the vectorized text chunks stored in
	// Qdrant for the given document, ordered by chunk_index (reading order).
	// Use params.Size to control how many chunks per page (max 100, default 20).
	ListChunks(ctx context.Context, documentUUID string, params ListChunksParams) (*ListChunksResponse, error)
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

func (d *documentClient) ListChunks(ctx context.Context, documentUUID string, params ListChunksParams) (*ListChunksResponse, error) {
	q := url.Values{}
	if params.Page > 0 {
		q.Set("page", strconv.Itoa(params.Page))
	}
	if params.Size > 0 {
		q.Set("size", strconv.Itoa(params.Size))
	}

	path := fmt.Sprintf(pathKnowledgeDocumentChunks, documentUUID)
	var out ListChunksResponse
	if err := d.http.get(ctx, path, q, &out); err != nil {
		return nil, fmt.Errorf("synapse/document.ListChunks: %w", err)
	}
	return &out, nil
}

func (d *documentClient) Estimate(ctx context.Context, req EstimateDocumentRequest) (*EstimateDocumentResponse, error) {
	fields := map[string]string{
		"embed_model": req.EmbedModel,
	}
	if req.ChunkSize > 0 {
		fields["chunk_size"] = strconv.Itoa(req.ChunkSize)
	}
	if req.Overlap > 0 {
		fields["chunk_overlap"] = strconv.Itoa(req.Overlap)
	}

	var out EstimateDocumentResponse
	err := d.http.postMultipart(
		ctx,
		pathKnowledgeDocumentEstimate,
		fields,
		"file",
		req.FileName,
		req.Content,
		&out,
	)
	if err != nil {
		return nil, fmt.Errorf("synapse/document.Estimate: %w", err)
	}
	return &out, nil
}
