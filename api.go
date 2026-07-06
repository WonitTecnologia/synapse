package synapse

// defaultBaseURL is the production Synapse API base URL.
// Override it via Options.BaseURL when instantiating the client.
const defaultBaseURL = "https://synapse.wonit.net.br"

// ─── Status paths ────────────────────────────────────────────────────────────

const (
	pathStatus = "/api/status"
)

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
	pathGoogleVisionConfig     = "/api/catalog/application/google/vision/config"
	pathGoogleVisionOCR        = "/api/catalog/application/google/vision/image/ocr"
	pathGoogleVisionOCRFromURL = "/api/catalog/application/google/vision/image/ocr/url"
)

// ─── OpenAI paths ─────────────────────────────────────────────────────────────

const (
	pathOpenAIConfig              = "/api/catalog/application/openai/config"
	pathOpenAIChat                = "/api/catalog/application/openai/chat"
	pathOpenAIImageAnalyze        = "/api/catalog/application/openai/image/analyze"
	pathOpenAIImageAnalyzeFromURL = "/api/catalog/application/openai/image/analyze/url"
	pathOpenAITranscribe          = "/api/catalog/application/openai/audio/transcribe"
	pathOpenAITranscribeFromURL   = "/api/catalog/application/openai/audio/transcribe/url"
)

// ─── Chatvolt paths ───────────────────────────────────────────────────────────

const (
	pathChatvoltConfig = "/api/catalog/application/chatvolt/config"
	pathChatvoltQuery  = "/api/catalog/application/chatvolt/query"
	pathChatvoltAgents = "/api/catalog/application/chatvolt/agents"
)

// ─── OpenRouter paths ─────────────────────────────────────────────────────────

const (
	pathOpenRouterSincronismo      = "/api/catalog/application/openrouter/sincronismo"
	pathOpenRouterModels           = "/api/catalog/application/openrouter/models"
	pathOpenRouterEmbeddingModels  = "/api/catalog/application/openrouter/models/embedding"
	pathOpenRouterAnalyticsMonthly = "/api/catalog/application/openrouter/analytics/monthly"
	pathOpenRouterAnalyticsQuery   = "/api/catalog/application/openrouter/analytics/query"
	pathOpenRouterAnalyticsMeta    = "/api/catalog/application/openrouter/analytics/meta"
)

// ─── Knowledge – Collection paths ─────────────────────────────────────────────

const (
	pathKnowledgeCollectionCreate = "/api/knowledge/collection/create"
	pathKnowledgeCollection       = "/api/knowledge/collection"
	pathKnowledgeCollectionList   = "/api/knowledge/collection/list"
)

// ─── Knowledge – Document paths ───────────────────────────────────────────────

const (
	pathKnowledgeDocumentUpload   = "/api/knowledge/application/document/upload"
	pathKnowledgeDocumentEstimate = "/api/knowledge/application/document/estimate"
	pathKnowledgeDocument         = "/api/knowledge/application/document"
	pathKnowledgeDocumentChunks   = "/api/knowledge/application/document/%s/chunks"
)

// ─── Agent paths ──────────────────────────────────────────────────────────────

const (
	pathAgent             = "/api/agent/domain/agent"
	pathAgentChat         = "/api/agent/application/agent/chat"
	pathAgentConversation = "/api/agent/application/agent/conversation"
)

// ─── Agent Prompt paths ───────────────────────────────────────────────────────

const (
	// pathAgentPrompt is formatted with the agent UUID: fmt.Sprintf(pathAgentPrompt, agentUUID)
	pathAgentPrompt       = "/api/agent/domain/agent/%s/prompt"
	pathAgentActivePrompt = "/api/agent/domain/agent/%s/active-prompt"

	pathAgentLogs      = "/api/agent/application/agent/%s/logs"
	pathAgentLogsStats = "/api/agent/application/agent/%s/logs/stats"
	pathAgentThoughts  = "/api/agent/application/agent/%s/thoughts"

	pathWorkspaceCredits  = "/api/catalog/application/openrouter/credits"
	pathWorkspaceActivity = "/api/catalog/application/openrouter/activity"
	pathWorkspaceKey      = "/api/catalog/application/openrouter/key"
)

// ─── MCP Integration paths ────────────────────────────────────────────────────
const (
	pathMcpIntegrations     = "/api/mcp/integrations"
	pathMcpIntegration      = "/api/mcp/integrations/%s"
	pathMcpToggle           = "/api/mcp/integrations/%s/toggle"
	pathMcpIntegrationTools = "/api/mcp/integrations/%s/tools"
)

// ─── External API tool paths ──────────────────────────────────────────────────
const (
	pathExternalApis      = "/api/external-apis"
	pathExternalApi       = "/api/external-apis/%s"
	pathExternalApiToggle = "/api/external-apis/%s/toggle"
)

// ─── Monitor (WebSocket) paths ────────────────────────────────────────────────
const (
	pathMonitorLogsWS = "/api/websocket/application/monitor/logs"
)
