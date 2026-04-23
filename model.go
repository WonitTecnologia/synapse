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
