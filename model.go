package synapse

// ─── Auth ────────────────────────────────────────────────────────────────────

// LoginRequest holds the credentials for the login endpoint.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse is returned on successful login or healthcheck.
type LoginResponse struct {
	Token         string          `json:"token"`
	Expire        string          `json:"expire"`
	SystemTimeUTC string          `json:"system_time_utc"`
	User          UserResponseDto `json:"user"`
}

// OTPRequest requests a one-time password to be sent to the given email.
type OTPRequest struct {
	Email string `json:"email"`
}

// OTPResetPasswordRequest resets a user password using an OTP code.
type OTPResetPasswordRequest struct {
	Email    string `json:"email"`
	OTP      string `json:"otp"`
	Password string `json:"password"`
}

// ApiTokenCreateRequest is the body for creating an API token.
type ApiTokenCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	// ExpireAt is optional; if omitted the server sets a default expiry.
	ExpireAt string `json:"expire_at,omitempty"`
}

// ApiTokenCreateResponse describes a created or listed API token.
type ApiTokenCreateResponse struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Token       string `json:"token"`
	TenantUUID  string `json:"tenant_uuid"`
	UserUUID    string `json:"user_uuid"`
	ExpireAt    string `json:"expire_at"`
	CreatedAt   string `json:"create_at"`
}

// ListApiTokensResponse is the paginated response for the API token list.
type ListApiTokensResponse struct {
	UserUUID  string                   `json:"user_uuid"`
	Page      int                      `json:"page"`
	Size      int                      `json:"size"`
	APITokens []ApiTokenCreateResponse `json:"api_tokens"`
}

// ─── User ────────────────────────────────────────────────────────────────────

// UserRole represents the access level of a user within the system.
type UserRole string

const (
	UserRoleSystemAdmin UserRole = "SYSTEM_ADMIN"
	UserRoleTenantAdmin UserRole = "TENANT_ADMIN"
	UserRoleTenantUser  UserRole = "TENANT_USER"
)

// UserResponseDto is the user shape returned by the API.
type UserResponseDto struct {
	UUID       string   `json:"uuid"`
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Role       UserRole `json:"role"`
	TenantUUID string   `json:"tenant_uuid"`
	Live       bool     `json:"live"`
	CreatedAt  string   `json:"create_at"`
	UpdatedAt  string   `json:"update_at"`
}

// CreateUserRequestDto is the body for creating a user.
type CreateUserRequestDto struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Role     UserRole `json:"role"`
}

// UpdateUserRequestDto is the body for partially updating a user.
type UpdateUserRequestDto struct {
	Name     string   `json:"name,omitempty"`
	Email    string   `json:"email,omitempty"`
	Password string   `json:"password,omitempty"`
	Role     UserRole `json:"role,omitempty"`
}

// ListUsersParams holds query parameters for the user list endpoint.
type ListUsersParams struct {
	Page int
	Size int
	// TenantIdentifier filters by tenant UUID or document (SystemAdmin only).
	TenantIdentifier string
}

// ─── Tenant ──────────────────────────────────────────────────────────────────

// TenantResponseDto is the tenant shape returned by the API.
type TenantResponseDto struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Document  string `json:"document"`
	Live      bool   `json:"live"`
	CreatedAt string `json:"createAt"`
	UpdatedAt string `json:"updateAt"`
}

// TenantsResponseDto is the paginated response for the tenant list.
type TenantsResponseDto struct {
	Page    int                 `json:"page"`
	Size    int                 `json:"size"`
	Tenants []TenantResponseDto `json:"tenants"`
}

// CreateTenantRequestDto is the body for creating a tenant.
type CreateTenantRequestDto struct {
	ID       string `json:"uuid"`
	Name     string `json:"name"`
	Document string `json:"document"`
}

// UpdateTenantRequestDto is the body for partially updating a tenant.
type UpdateTenantRequestDto struct {
	Name     string `json:"name,omitempty"`
	Document string `json:"document,omitempty"`
	// Live uses a pointer so false can be distinguished from omitted.
	Live *bool `json:"live,omitempty"`
}

// ─── Provider ────────────────────────────────────────────────────────────────

// ProviderResponseDto describes a catalog integration provider (e.g. Google, OpenAI).
type ProviderResponseDto struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	LogoURL   string `json:"logo_url"`
	Website   string `json:"website"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ProvidersResponseDto is the paginated response for the provider list.
type ProvidersResponseDto struct {
	Page      int                   `json:"page"`
	Size      int                   `json:"size"`
	Providers []ProviderResponseDto `json:"providers"`
}

// ─── Service ─────────────────────────────────────────────────────────────────

// ServiceResponseDto describes a catalog integration service.
type ServiceResponseDto struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Slug             string `json:"slug"`
	ProviderID       string `json:"provider_id"`
	BaseURL          string `json:"base_url"`
	AuthType         string `json:"auth_type"`
	DocumentationURL string `json:"documentation_url"`
}

// ServicesResponseDto is the paginated response for the service list.
type ServicesResponseDto struct {
	Page     int                  `json:"page"`
	Size     int                  `json:"size"`
	Services []ServiceResponseDto `json:"services"`
}

// ─── OpenAI ──────────────────────────────────────────────────────────────────

// OpenAIModel is the model identifier accepted by the OpenAI integration.
type OpenAIModel string

const (
	OpenAIModelGPT4oMini  OpenAIModel = "gpt-4o-mini"
	OpenAIModelGPT4o      OpenAIModel = "gpt-4o"
	OpenAIModelGPT4_1     OpenAIModel = "gpt-4.1"
	OpenAIModelGPT4_1Mini OpenAIModel = "gpt-4.1-mini"
	OpenAIModelO4Mini     OpenAIModel = "o4-mini"
)

// OpenAICredentialsDTO holds the API key for the OpenAI integration.
type OpenAICredentialsDTO struct {
	Token string `json:"token"`
}

// OpenAISettingsDTO configures model and sampling behaviour.
type OpenAISettingsDTO struct {
	Model OpenAIModel `json:"model"`
	// Temperature ranges 0–2; omitted when zero.
	Temperature float64 `json:"temperature,omitempty"`
}

// ConfigureOpenAIRequest is the body for creating/updating the OpenAI integration.
type ConfigureOpenAIRequest struct {
	Credentials OpenAICredentialsDTO `json:"credentials"`
	Settings    OpenAISettingsDTO    `json:"settings"`
	IsActive    bool                 `json:"is_active,omitempty"`
}

// ChatCompletionRequest is the body for a chat completion call.
type ChatCompletionRequest struct {
	Prompt string `json:"prompt"`
}

// ChatCompletionResponse is the response from a chat completion call.
type ChatCompletionResponse struct {
	Response string `json:"response"`
}

// AnalyzeImageResponse is returned by both OpenAI and Google Vision image endpoints.
type AnalyzeImageResponse struct {
	Response string `json:"response"`
}

// TranscribeAudioRequest groups all parameters for the audio transcription endpoint.
type TranscribeAudioRequest struct {
	// FileName is used as the multipart filename (e.g. "recording.mp3").
	FileName string
	// Content is the raw audio file bytes.
	Content []byte
	// Model selects the transcription engine: whisper-1 | gpt-4o-transcribe | gpt-4o-mini-transcribe.
	Model string
	// Language hints the spoken language (e.g. "pt", "en", "es").
	Language string
	// Prompt is optional auxiliary text to improve transcription accuracy.
	Prompt string
}

// TranscribeAudioFromURLRequest groups all parameters for the audio transcription
// endpoint when the audio file is fetched from a download URL by the server.
type TranscribeAudioFromURLRequest struct {
	// FileURL is the download link the server uses to fetch the audio file.
	FileURL string `json:"file_url"`
	// Model selects the transcription engine: whisper-1 | gpt-4o-transcribe | gpt-4o-mini-transcribe.
	Model string `json:"model,omitempty"`
	// Language hints the spoken language (e.g. "pt", "en", "es").
	Language string `json:"language,omitempty"`
	// Prompt is optional auxiliary text to improve transcription accuracy.
	Prompt string `json:"prompt,omitempty"`
}

// AnalyzeImageFromURLRequest is the body for analysing an image that the server
// fetches from a download URL (OpenAI image analysis).
type AnalyzeImageFromURLRequest struct {
	FileURL string `json:"file_url"`
	Prompt  string `json:"prompt"`
}

// TranscriptionResponse is the response from the audio transcription endpoint.
type TranscriptionResponse struct {
	Response string `json:"response"`
}

// ─── Google ──────────────────────────────────────────────────────────────────

// VisionAICredentialsDTO holds the API key for the Google Vision integration.
type VisionAICredentialsDTO struct {
	Token string `json:"token"`
}

// ConfigureGoogleRequest is the body for creating/updating the Google Vision integration.
type ConfigureGoogleRequest struct {
	Credentials VisionAICredentialsDTO `json:"credentials"`
	IsActive    bool                   `json:"is_active,omitempty"`
}

// VisionOCRFromURLRequest is the body for the Google Vision OCR endpoint when the
// image is fetched by the server from a download URL.
type VisionOCRFromURLRequest struct {
	FileURL string `json:"file_url"`
}

// ─── Chatvolt ────────────────────────────────────────────────────────────────

// ChatvoltCredentialsDTO holds the API key for the Chatvolt integration.
type ChatvoltCredentialsDTO struct {
	Token string `json:"token"`
}

// ConfigureChatvoltRequest is the body for creating/updating the Chatvolt integration.
type ConfigureChatvoltRequest struct {
	Credentials ChatvoltCredentialsDTO `json:"credentials"`
	IsActive    bool                   `json:"is_active,omitempty"`
}

// ChatvoltContact carries optional contact metadata attached to a Chatvolt query.
type ChatvoltContact struct {
	UserID      string `json:"userId,omitempty"`
	Email       string `json:"email,omitempty"`
	FirstName   string `json:"firstName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

// ChatvoltAgentQueryRequest is the body for querying a Chatvolt agent.
type ChatvoltAgentQueryRequest struct {
	AgentID string `json:"agentId"`
	Query   string `json:"query"`
	// ConversationID keeps context across turns; omit to start a new conversation.
	ConversationID string `json:"conversationId,omitempty"`
	// Streaming should generally be false for synchronous responses.
	Streaming bool             `json:"streaming"`
	Contact   *ChatvoltContact `json:"contact,omitempty"`
}

// ChatvoltAgentVisibility indicates whether a Chatvolt agent is publicly accessible.
type ChatvoltAgentVisibility string

const (
	ChatvoltAgentVisibilityPublic  ChatvoltAgentVisibility = "public"
	ChatvoltAgentVisibilityPrivate ChatvoltAgentVisibility = "private"
)

// ChatvoltAgentItem describes a single Chatvolt agent returned by the catalog.
type ChatvoltAgentItem struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	IconUrl     *string                 `json:"iconUrl,omitempty"`
	Description string                  `json:"description"`
	ModelName   string                  `json:"modelName"`
	Visibility  ChatvoltAgentVisibility `json:"visibility"`
}

// ListChatvoltAgentsResponse is the response for the Chatvolt agents list endpoint.
type ListChatvoltAgentsResponse struct {
	Agents []ChatvoltAgentItem `json:"agents"`
}

// ChatvoltAgentQueryResponse is the response from a Chatvolt agent query.
type ChatvoltAgentQueryResponse struct {
	Answer         string      `json:"answer"`
	ConversationID string      `json:"conversationId"`
	MessageID      string      `json:"messageId"`
	VisitorID      string      `json:"visitorId"`
	Sources        interface{} `json:"sources"`
	Metadata       interface{} `json:"metadata"`
}

// ─── OpenRouter ───────────────────────────────────────────────────────────────

// OpenRouterSyncRequest is the body for syncing a workspace with OpenRouter.
// TenantUUID is optional and only honoured when the caller is a SystemAdmin.
type OpenRouterSyncRequest struct {
	TenantUUID string `json:"tenant_uuid,omitempty"`
}

// OpenRouterSyncResponse is returned after a successful workspace sync.
type OpenRouterSyncResponse struct {
	TenantUUID    string `json:"tenant_uuid"`
	WorkspaceID   string `json:"workspace_id"`
	KeyHash       string `json:"key_hash"`
	Active        bool   `json:"active"`
	AlreadySynced bool   `json:"already_synced"`
	SyncedAt      string `json:"synced_at"`
}

// OpenRouterDesyncResponse is returned after a successful workspace removal.
type OpenRouterDesyncResponse struct {
	TenantUUID string `json:"tenant_uuid"`
	Removed    bool   `json:"removed"`
}

// OpenRouterModelPricing holds per-token pricing for an OpenRouter model.
type OpenRouterModelPricing struct {
	Prompt     string `json:"prompt"`
	Completion string `json:"completion"`
}

// OpenRouterModelInfo describes a single model available on OpenRouter.
type OpenRouterModelInfo struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description,omitempty"`
	ContextLength int                    `json:"context_length,omitempty"`
	Pricing       OpenRouterModelPricing `json:"pricing"`
	IsFree        bool                   `json:"is_free"`
}

// OpenRouterListModelsResponse is the paginated list of OpenRouter models.
type OpenRouterListModelsResponse struct {
	Models []OpenRouterModelInfo `json:"models"`
	Total  int                   `json:"total"`
}

// OpenRouterEmbeddingModel describes a single embedding model on OpenRouter.
type OpenRouterEmbeddingModel struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Provider    string `json:"provider"`
	VectorSize  int    `json:"vector_size"`
	Description string `json:"description"`
	PricePer1M  string `json:"price_per_1m_tokens_usd"`
	IsFree      bool   `json:"is_free"`
}

// OpenRouterListEmbeddingModelsResponse is the paginated list of OpenRouter embedding models.
type OpenRouterListEmbeddingModelsResponse struct {
	Models []OpenRouterEmbeddingModel `json:"models"`
	Total  int                        `json:"total"`
}

// ─── Knowledge – Collection ───────────────────────────────────────────────────

// CreateCollectionRequest is the body for creating a Qdrant vector collection.
type CreateCollectionRequest struct {
	// TenantUUID is required when the caller is SYSTEM_ADMIN; ignored for other roles.
	TenantUUID *string `json:"tenant_uuid,omitempty"`
	Name       string  `json:"name"`
	// VectorSize is the dimension of the embedding vectors (e.g. 1536 for text-embedding-ada-002).
	// When zero the server uses its default (1536).
	VectorSize uint64 `json:"vector_size,omitempty"`
	// Distance is the similarity metric: "Cosine" (default), "Euclid", or "Dot".
	Distance string `json:"distance,omitempty"`
}

// CollectionResponse describes a single Qdrant collection registered for a tenant.
type CollectionResponse struct {
	UUID                 string `json:"uuid"`
	TenantUUID           string `json:"tenant_uuid"`
	Name                 string `json:"name"`
	QdrantCollectionName string `json:"qdrant_collection_name"`
	VectorSize           uint64 `json:"vector_size"`
	Distance             string `json:"distance"`
	CreatedAt            string `json:"createAt"`
	UpdatedAt            string `json:"updateAt"`
}

// CollectionsResponse is the paginated list of collections for a tenant.
type CollectionsResponse struct {
	Collections []CollectionResponse `json:"collections"`
	Page        int                  `json:"page"`
	Size        int                  `json:"size"`
}

// ─── Knowledge – Document ─────────────────────────────────────────────────────

// UploadDocumentRequest groups all parameters for the document upload endpoint.
type UploadDocumentRequest struct {
	// TenantUUID is required when the caller is SYSTEM_ADMIN; ignored for other roles.
	TenantUUID *string
	// CollectionUUID is the UUID of the target Qdrant collection.
	CollectionUUID string
	// EmbedModel is the OpenRouter embedding model ID (e.g. "openai/text-embedding-ada-002").
	EmbedModel string
	// ChunkSize is the chunk size in characters (default: 1000 when zero).
	ChunkSize int
	// ChunkOverlap is the overlap between chunks in characters (default: 200 when zero).
	ChunkOverlap int
	// FileName is the multipart filename (e.g. "report.pdf").
	FileName string
	// Content is the raw file bytes (PDF, HTML, DOCX, XLSX, CSV).
	Content []byte
}

// DocumentResponse describes a document that has been uploaded and vectorized.
type DocumentResponse struct {
	UUID           string `json:"uuid"`
	TenantUUID     string `json:"tenant_uuid"`
	CollectionUUID string `json:"collection_uuid"`
	Filename       string `json:"filename"`
	FileType       string `json:"file_type"`
	EmbedModel     string `json:"embed_model"`
	ChunkCount     int    `json:"chunk_count"`
	VectorSize     int    `json:"vector_size"`
	TokensUsed     int64  `json:"tokens_used"`
	// Status is "processing", "ready", or "error".
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message,omitempty"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// ListDocumentsResponse is the paginated list of documents in a collection.
type ListDocumentsResponse struct {
	Documents []DocumentResponse `json:"documents"`
	Page      int                `json:"page"`
	Size      int                `json:"size"`
}

// ListDocumentsParams holds query parameters for the document list endpoint.
type ListDocumentsParams struct {
	// CollectionUUID filters documents by collection.
	// Required for TENANT_ADMIN/TENANT_USER; optional for SYSTEM_ADMIN.
	CollectionUUID string
	Page           int
	Size           int
}

// ─── Agent ────────────────────────────────────────────────────────────────────

// CreateAgentRequest is the body for creating an AI agent.
type CreateAgentRequest struct {
	// TenantUUID is required when the caller is SYSTEM_ADMIN; ignored for other roles.
	TenantUUID  *string  `json:"tenant_uuid,omitempty"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Model       string   `json:"model"`
	Prompt      string   `json:"prompt"`
	// CollectionUUID enables RAG — set to the UUID of an indexed collection.
	// The embedding model is resolved automatically from the collection documents.
	CollectionUUID *string `json:"collection_uuid,omitempty"`
	// MaxContext accepted values: 10000, 15000, 20000 (default: 10000).
	MaxContext int `json:"max_context,omitempty"`
	// Temperature controls randomness: 0.0 (deterministic) – 0.7 (max without hallucination).
	Temperature *float64 `json:"temperature,omitempty"`
}

// UpdateAgentRequest is used for both full (PUT) and partial (PATCH) agent updates.
// For PATCH, set exactly one field; for PUT you may set multiple.
type UpdateAgentRequest struct {
	Name           *string  `json:"name,omitempty"`
	Description    *string  `json:"description,omitempty"`
	Model          *string  `json:"model,omitempty"`
	Prompt         *string  `json:"prompt,omitempty"`
	CollectionUUID *string  `json:"collection_uuid,omitempty"`
	MaxContext     *int     `json:"max_context,omitempty"`
	Temperature    *float64 `json:"temperature,omitempty"`
	Active         *bool    `json:"active,omitempty"`
}

// AgentResponse describes an AI agent.
type AgentResponse struct {
	UUID            string  `json:"uuid"`
	TenantUUID      string  `json:"tenant_uuid"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Model           string  `json:"model"`
	Prompt          string  `json:"prompt"`
	CollectionUUID  *string `json:"collection_uuid"`
	QueryEmbedModel string  `json:"query_embed_model,omitempty"`
	MaxContext      int     `json:"max_context"`
	Temperature     float64 `json:"temperature"`
	Active          bool    `json:"active"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

// ListAgentsResponse is the paginated list of agents for a tenant.
type ListAgentsResponse struct {
	Agents []AgentResponse `json:"agents"`
	Page   int             `json:"page"`
	Size   int             `json:"size"`
}

// ChatRequest is the body for sending a message to an AI agent.
type ChatRequest struct {
	AgentUUID string `json:"agent_uuid"`
	Message   string `json:"message"`
	// ConversationUUID continues an existing conversation; omit to start a new one.
	ConversationUUID *string `json:"conversation_uuid,omitempty"`
}

// ChatReference is a document chunk retrieved from Qdrant and used in the RAG context.
type ChatReference struct {
	Filename   string  `json:"filename"`
	ChunkIndex int     `json:"chunk_index"`
	Score      float32 `json:"score"`
	ScorePct   float32 `json:"score_pct"`
	Text       string  `json:"text"`
}

// ChatRagInfo reports what happened during the semantic search step.
type ChatRagInfo struct {
	Enabled     bool   `json:"enabled"`
	ChunksFound int    `json:"chunks_found"`
	Error       string `json:"error,omitempty"`
}

// ChatResponse is the reply from the AI agent.
type ChatResponse struct {
	ConversationUUID string          `json:"conversation_uuid"`
	Message          string          `json:"message"`
	References       []ChatReference `json:"references"`
	Rag              ChatRagInfo     `json:"rag"`
}

// ConversationResponse is a summary of a stored conversation (without full message content).
type ConversationResponse struct {
	UUID         string `json:"uuid"`
	TenantUUID   string `json:"tenant_uuid"`
	AgentUUID    string `json:"agent_uuid"`
	// MessageCount is the total number of messages (user + assistant turns).
	MessageCount int    `json:"message_count"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// ListConversationsResponse is the paginated list of conversation summaries.
type ListConversationsResponse struct {
	Conversations []ConversationResponse `json:"conversations"`
	Page          int                   `json:"page"`
	Size          int                   `json:"size"`
}

// ListConversationsParams holds query parameters for the conversation list endpoint.
type ListConversationsParams struct {
	// AgentUUID filters conversations by a specific agent (optional).
	AgentUUID string
	Page      int
	Size      int
}
