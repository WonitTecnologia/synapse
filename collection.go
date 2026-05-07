package synapse

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ─── Interface ────────────────────────────────────────────────────────────────

// CollectionCase provides CRUD operations for Qdrant vector collections.
type CollectionCase interface {
	// Create registers a new vector collection for the authenticated tenant.
	Create(ctx context.Context, req CreateCollectionRequest) (*CollectionResponse, error)

	// Get returns a collection by its UUID.
	Get(ctx context.Context, collectionUUID string) (*CollectionResponse, error)

	// List returns a paginated list of collections for the authenticated tenant.
	// Use page=0 and size=0 to rely on server defaults.
	List(ctx context.Context, page, size int) (*CollectionsResponse, error)

	// Delete permanently removes a collection and its Qdrant data by UUID.
	Delete(ctx context.Context, collectionUUID string) error
}

// ─── Implementation ───────────────────────────────────────────────────────────

type collectionClient struct {
	http *httpClient
}

func newCollectionClient(hc *httpClient) CollectionCase {
	return &collectionClient{http: hc}
}

func (c *collectionClient) Create(ctx context.Context, req CreateCollectionRequest) (*CollectionResponse, error) {
	var out CollectionResponse
	if err := c.http.post(ctx, pathKnowledgeCollectionCreate, req, &out); err != nil {
		return nil, fmt.Errorf("synapse/collection.Create: %w", err)
	}
	return &out, nil
}

func (c *collectionClient) Get(ctx context.Context, collectionUUID string) (*CollectionResponse, error) {
	q := url.Values{}
	q.Set("uuid", collectionUUID)

	var out CollectionResponse
	if err := c.http.get(ctx, pathKnowledgeCollection, q, &out); err != nil {
		return nil, fmt.Errorf("synapse/collection.Get: %w", err)
	}
	return &out, nil
}

func (c *collectionClient) List(ctx context.Context, page, size int) (*CollectionsResponse, error) {
	q := url.Values{}
	if page > 0 {
		q.Set("page", strconv.Itoa(page))
	}
	if size > 0 {
		q.Set("size", strconv.Itoa(size))
	}

	var out CollectionsResponse
	if err := c.http.get(ctx, pathKnowledgeCollectionList, q, &out); err != nil {
		return nil, fmt.Errorf("synapse/collection.List: %w", err)
	}
	return &out, nil
}

func (c *collectionClient) Delete(ctx context.Context, collectionUUID string) error {
	q := url.Values{}
	q.Set("uuid", collectionUUID)

	if err := c.http.delete(ctx, pathKnowledgeCollection, q); err != nil {
		return fmt.Errorf("synapse/collection.Delete: %w", err)
	}
	return nil
}
