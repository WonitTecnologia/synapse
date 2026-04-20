package synapse

// defaultBaseURL is the production Synapse API base URL.
// Override it via Options.BaseURL when instantiating the client.
const defaultBaseURL = "https://synapse.wonit.net.br"

// ─── Auth paths ──────────────────────────────────────────────────────────────

const (
	pathAuthLogin         = "/api/auth/login"
	pathAuthLogout        = "/api/auth/logout/"
	pathAuthHealthcheck   = "/api/auth/healthcheck"
	pathAuthOTP           = "/api/auth/otp"
	pathAuthPasswordReset = "/api/auth/password/reset"
	pathAuthAPIToken      = "/api/auth/api-token"
)

// ─── User paths ──────────────────────────────────────────────────────────────

const (
	pathUserList = "/api/user/list"
	pathUser     = "/api/user/"
)

// ─── Tenant paths ─────────────────────────────────────────────────────────────

const (
	pathTenant       = "/api/tenant"
	pathTenantCreate = "/api/tenant/create"
	pathTenantList   = "/api/tenant/list"
)

// ─── Provider paths ───────────────────────────────────────────────────────────

const (
	pathProvider     = "/api/provider"
	pathProviderList = "/api/provider/list"
)

// ─── Service paths ────────────────────────────────────────────────────────────

const (
	pathService     = "/api/services"
	pathServiceList = "/api/services/list"
)

// ─── Google paths ─────────────────────────────────────────────────────────────

const (
	pathGoogleVisionConfig = "/api/catalog/application/google/vision/config"
	pathGoogleVisionOCR    = "/api/catalog/application/google/vision/image/ocr"
)

// ─── OpenAI paths ─────────────────────────────────────────────────────────────

const (
	pathOpenAIConfig       = "/api/catalog/application/openai/config"
	pathOpenAIChat         = "/api/catalog/application/openai/chat"
	pathOpenAIImageAnalyze = "/api/catalog/application/openai/image/analyze"
	pathOpenAITranscribe   = "/api/catalog/application/openai/audio/transcribe"
)

// ─── Chatvolt paths ───────────────────────────────────────────────────────────

const (
	pathChatvoltConfig = "/api/catalog/application/chatvolt/config"
	pathChatvoltQuery  = "/api/catalog/application/chatvolt/query"
	pathChatvoltAgents = "/api/catalog/application/chatvolt/agents"
)
