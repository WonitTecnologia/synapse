# Synapse Go SDK

Client oficial em Go para a API [Synapse](https://synapse.wonit.net.br).

---

## Instalação

```bash
go get github.com/WonitTecnologia/synapse
```

---

## Início rápido

```go
import "github.com/WonitTecnologia/synapse"

client, err := synapse.NewClient("seu-token", nil)
if err != nil {
    log.Fatal(err)
}
```

---

## Configuração

`NewClient` aceita um token obrigatório e um `*Options` opcional.

```go
client, err := synapse.NewClient("seu-token", &synapse.Options{
    BaseURL: "https://staging.synapse.example.com", // padrão: https://synapse.wonit.net.br
    Timeout: 15 * time.Second,                      // padrão: 30s
})
```

| Campo     | Tipo            | Descrição                                              |
|-----------|-----------------|--------------------------------------------------------|
| `BaseURL` | `string`        | Sobrescreve a URL base da API (staging, self-hosted…) |
| `Host`    | `string`        | Sobrepõe o header `Host` em toda requisição (HTTP e WebSocket). Para conexão interna por IP — ver abaixo |
| `Timeout` | `time.Duration` | Timeout por requisição HTTP. Padrão: `30s`            |

### Conexão interna (IP direto + Host header)

Para conectar pela rede interna sem passar pela URL pública (mesmo padrão do cliente
CSA): aponte a `BaseURL` para o IP e informe o DNS do virtual host em `Host`.

```go
client, err := synapse.NewClient("seu-token", &synapse.Options{
    BaseURL: "http://172.16.50.41",        // IP interno (porta opcional)
    Host:    "synapse-dev.wonit.cloud",    // DNS enviado no header Host
})
```

- Vale para **todas** as requisições HTTP (inclusive upload multipart) **e** para o
  WebSocket de monitoramento (no `wss` o `Host` também é usado como SNI/ServerName).
- Sem `Host`, tudo segue normalmente pela `BaseURL` informada (modo público).

---

## Domínios

| Campo                | Interface         | Endpoints cobertos                               |
|----------------------|-------------------|--------------------------------------------------|
| `client.Auth`        | `AuthCase`        | Login, logout, OTP, reset de senha, API tokens   |
| `client.User`        | `UserCase`        | CRUD de usuários                                 |
| `client.Tenant`      | `TenantCase`      | CRUD de tenants                                  |
| `client.Provider`    | `ProviderCase`    | Listagem de providers do catálogo                |
| `client.Service`     | `ServiceCase`     | Listagem de services do catálogo                 |
| `client.Google`      | `GoogleCase`      | Integração Google Vision AI (OCR)                |
| `client.OpenAI`      | `OpenAICase`      | Chat, análise de imagem, transcrição de áudio    |
| `client.Chatvolt`    | `ChatvoltCase`    | Query a agentes Chatvolt                         |
| `client.OpenRouter`  | `OpenRouterCase`  | Workspace OpenRouter (sync, modelos, analytics)  |
| `client.Collection`  | `CollectionCase`  | Coleções vetoriais (Qdrant) da base de conhecimento |
| `client.Document`    | `DocumentCase`    | Upload e vetorização de documentos               |
| `client.Agent`       | `AgentCase`       | CRUD de agentes de IA + chat (com RAG)           |
| `client.Mcp`         | `McpCase`         | Integrações MCP (Model Context Protocol)         |
| `client.ExternalApi` | `ExternalApiCase` | APIs externas (HTTP cruas) como tools do agente  |
| `client.Monitor`     | `MonitorCase`     | **WebSocket de monitoramento** — stream de eventos do agente em tempo real ([ver seção](#websocket-de-monitoramento-monitor)) |

---

## Auth

### Login

```go
resp, err := client.Auth.Login(ctx, synapse.LoginRequest{
    Email:    "usuario@empresa.com",
    Password: "senha123",
})

fmt.Println(resp.Token)
fmt.Println(resp.User.Name)
```

### Healthcheck (validar token)

```go
resp, err := client.Auth.Healthcheck(ctx)
fmt.Println(resp.User.Email)
```

### Logout

```go
err := client.Auth.Logout(ctx, "token-a-revogar")
```

### Solicitar OTP

```go
err := client.Auth.RequestOTP(ctx, synapse.OTPRequest{
    Email: "usuario@empresa.com",
})
```

### Resetar senha com OTP

```go
err := client.Auth.ResetPassword(ctx, synapse.OTPResetPasswordRequest{
    Email:    "usuario@empresa.com",
    OTP:      "123456",
    Password: "novaSenha@123",
})
```

### API Tokens

```go
// Criar
token, err := client.Auth.CreateAPIToken(ctx, synapse.ApiTokenCreateRequest{
    Name:        "Integração CI",
    Description: "Token para pipeline de deploy",
    ExpireAt:    "2026-12-31T00:00:00Z", // opcional
})
fmt.Println(token.Token)

// Listar
list, err := client.Auth.ListAPITokens(ctx, 1, 20)
for _, t := range list.APITokens {
    fmt.Println(t.Name, t.ExpireAt)
}
```

---

## User

```go
// Buscar por UUID ou email
user, err := client.User.Get(ctx, "uuid-ou-email")

// Listar
users, err := client.User.List(ctx, synapse.ListUsersParams{
    Page:             1,
    Size:             20,
    TenantIdentifier: "uuid-do-tenant", // apenas SystemAdmin
})

// Criar  (identifier = UUID ou documento do tenant)
user, err := client.User.Create(ctx, "uuid-do-tenant", synapse.CreateUserRequestDto{
    Name:     "João Silva",
    Email:    "joao@empresa.com",
    Password: "senha@123",
    Role:     synapse.UserRoleTenantUser,
})

// Atualizar
user, err := client.User.Update(ctx, "uuid-do-usuario", synapse.UpdateUserRequestDto{
    Name: "João Silva Jr.",
})

// Deletar
err := client.User.Delete(ctx, "uuid-ou-email")
```

### Roles disponíveis

| Constante                      | Valor          |
|--------------------------------|----------------|
| `synapse.UserRoleSystemAdmin`  | `SYSTEM_ADMIN` |
| `synapse.UserRoleTenantAdmin`  | `TENANT_ADMIN` |
| `synapse.UserRoleTenantUser`   | `TENANT_USER`  |

---

## Tenant

```go
// Buscar por UUID ou documento (CNPJ/CPF)
tenant, err := client.Tenant.Get(ctx, "uuid", "")
tenant, err := client.Tenant.Get(ctx, "", "12345678000199")

// Listar
tenants, err := client.Tenant.List(ctx, 1, 10)

// Criar
tenant, err := client.Tenant.Create(ctx, synapse.CreateTenantRequestDto{
    Name:     "Empresa XPTO",
    Document: "12345678000199",
})

// Atualizar
live := true
tenant, err := client.Tenant.Update(ctx, "uuid-do-tenant", synapse.UpdateTenantRequestDto{
    Name: "Empresa XPTO Ltda.",
    Live: &live,
})

// Deletar
err := client.Tenant.Delete(ctx, "uuid-do-tenant", "")
```

---

## Provider & Service (Catálogo)

```go
// Providers
provider, err := client.Provider.Get(ctx, "uuid-do-provider")
providers, err := client.Provider.List(ctx, 1, 20)

// Services
service, err := client.Service.Get(ctx, "uuid-do-service")
services, err := client.Service.List(ctx, 1, 20)
```

---

## Google Vision AI

### Configurar integração

```go
err := client.Google.Configure(ctx, synapse.ConfigureGoogleRequest{
    Credentials: synapse.VisionAICredentialsDTO{Token: "sua-api-key-google"},
    IsActive:    true,
})
```

### OCR (extração de texto)

```go
imageBytes, _ := os.ReadFile("nota_fiscal.jpg")

result, err := client.Google.VisionOCR(ctx, "nota_fiscal.jpg", imageBytes)
fmt.Println(result.Response)
```

Formatos aceitos: `png`, `jpg`, `jpeg`, `webp`.

---

## OpenAI

### Configurar integração

```go
err := client.OpenAI.Configure(ctx, synapse.ConfigureOpenAIRequest{
    Credentials: synapse.OpenAICredentialsDTO{Token: "sk-..."},
    Settings: synapse.OpenAISettingsDTO{
        Model:       synapse.OpenAIModelGPT4o,
        Temperature: 0.7,
    },
    IsActive: true,
})
```

### Modelos disponíveis

| Constante                       | Modelo          |
|---------------------------------|-----------------|
| `synapse.OpenAIModelGPT4oMini`  | `gpt-4o-mini`   |
| `synapse.OpenAIModelGPT4o`      | `gpt-4o`        |
| `synapse.OpenAIModelGPT4_1`     | `gpt-4.1`       |
| `synapse.OpenAIModelGPT4_1Mini` | `gpt-4.1-mini`  |
| `synapse.OpenAIModelO4Mini`     | `o4-mini`       |

### Chat Completion

```go
reply, err := client.OpenAI.Chat(ctx, synapse.ChatCompletionRequest{
    Prompt: "Resuma este contrato em 3 pontos.",
})
fmt.Println(reply.Response)
```

### Análise de imagem

```go
imageBytes, _ := os.ReadFile("diagrama.png")

result, err := client.OpenAI.AnalyzeImage(ctx, "diagrama.png", imageBytes, "Descreva o que está nesta imagem.")
fmt.Println(result.Response)
```

### Transcrição de áudio

```go
audioBytes, _ := os.ReadFile("reuniao.mp3")

result, err := client.OpenAI.TranscribeAudio(ctx, synapse.TranscribeAudioRequest{
    FileName: "reuniao.mp3",
    Content:  audioBytes,
    Model:    "whisper-1",    // whisper-1 | gpt-4o-transcribe | gpt-4o-mini-transcribe
    Language: "pt",           // opcional
    Prompt:   "",             // opcional: contexto para melhorar a transcrição
})
fmt.Println(result.Response)
```

---

## Chatvolt

### Configurar integração

```go
err := client.Chatvolt.Configure(ctx, synapse.ConfigureChatvoltRequest{
    Credentials: synapse.ChatvoltCredentialsDTO{Token: "seu-token-chatvolt"},
    IsActive:    true,
})
```

### Query a agente

```go
// Nova conversa
resp, err := client.Chatvolt.Query(ctx, synapse.ChatvoltAgentQueryRequest{
    AgentID: "id-do-agente",
    Query:   "Qual o status do meu pedido #1234?",
})
fmt.Println(resp.Answer)
fmt.Println(resp.ConversationID) // guarde para continuar o contexto

// Continuar conversa existente
resp, err = client.Chatvolt.Query(ctx, synapse.ChatvoltAgentQueryRequest{
    AgentID:        "id-do-agente",
    ConversationID: resp.ConversationID,
    Query:          "E qual a previsão de entrega?",
})
```

### Enviar dados de contato

```go
resp, err := client.Chatvolt.Query(ctx, synapse.ChatvoltAgentQueryRequest{
    AgentID: "id-do-agente",
    Query:   "Preciso de suporte.",
    Contact: &synapse.ChatvoltContact{
        Email:     "cliente@email.com",
        FirstName: "Maria",
        LastName:  "Souza",
    },
})
```

---

## WebSocket de monitoramento (Monitor)

Stream em tempo real dos eventos de execução dos agentes de IA — chat, tool calls
(MCP e APIs externas), RAG, erros e processamento de arquivos. Canal
**somente-recebimento**: o SDK confirma cada entrega automaticamente (protocolo de
ACK interno) e reconecta sozinho; você só consome o canal de eventos.

- Token **master (SYSTEM_ADMIN)** → recebe eventos de **todos os tenants**.
- Token de **tenant** → recebe apenas os eventos do próprio tenant.
- Diferente do log persistido, os eventos chegam **sem truncamento** (parâmetros,
  retorno de API e resultado de tools íntegros).

### Uso básico

```go
stream, err := client.Monitor.StreamLogs(ctx, nil)
if err != nil {
    log.Fatal(err)
}
defer stream.Close()

for evt := range stream.Events() {
    fmt.Printf("[%s] %s — %s\n", evt.Category, evt.AgentName, evt.Summary)
}
// o canal fecha quando ctx é cancelado ou stream.Close() é chamado
```

### Opções (`StreamLogsOptions`)

```go
stream, err := client.Monitor.StreamLogs(ctx, &synapse.StreamLogsOptions{
    Session: "3f2a...-uuid",                 // retoma a fila do servidor após reconexão
    Buffer:  512,                            // capacidade do canal (padrão: 256)
    OnConnect: func(session string) {        // handshake ok (conexão E reconexões)
        log.Println("WS conectado, session:", session)
    },
    OnError: func(err error) {               // erros de conexão (o stream segue tentando)
        log.Println("WS erro:", err)
    },
})
```

| Campo | Tipo | Descrição |
|---|---|---|
| `Session` | `string` | UUID de sessão. Reconexões com a mesma session **retomam a fila de entrega** pendente no servidor. Padrão: UUID aleatório mantido pela vida do stream |
| `Buffer` | `int` | Capacidade do canal de eventos. Padrão: `256` |
| `OnConnect` | `func(session string)` | Disparado a cada handshake bem-sucedido — conexão inicial **e** cada reconexão automática |
| `OnError` | `func(error)` | Erros de conexão/handshake. Apenas observabilidade: a reconexão é automática (backoff 1s → 30s) |

### `EventStream`

| Método | Descrição |
|---|---|
| `Events() <-chan AgentEvent` | Canal dos eventos recebidos (fecha no encerramento) |
| `Session() string` | UUID da sessão em uso (útil para logar/persistir) |
| `Close()` | Encerra o stream e fecha o canal |

### `AgentEvent`

| Campo | Tipo | Descrição |
|---|---|---|
| `UUID` | `string` | Identidade do evento (use para dedup, se necessário) |
| `TenantUUID` | `string` | Tenant dono do evento |
| `AgentUUID` / `AgentName` | `string` | Agente que gerou o evento |
| `ConversationUUID` | `*string` | Conversa interna |
| `ConversationExternalID` | `*string` | ID externo (ex.: protocolo) |
| `Level` | `string` | `info` \| `warn` \| `error` |
| `Category` | `string` | `EventCategoryChat` \| `EventCategoryToolCall` \| `EventCategoryRAG` \| `EventCategoryError` \| `EventCategoryFileProcess` |
| `Summary` | `string` | Resumo humano do evento |
| `Detail` | `map[string]any` | Detalhe por categoria (chat: `user_msg`/`response` íntegros, `reasoning`…) |
| `ToolName` | `*string` | (tool_call) nome da ferramenta |
| `ToolParams` | `map[string]any` | (tool_call) parâmetros **sem truncar** |
| `ToolSuccess` | `*bool` | (tool_call) sucesso |
| `ToolResult` | `string` | (tool_call) resultado **sem truncar** |
| `APIResponse` | `string` | (tool_call de API externa) corpo da resposta **sem truncar** |
| `Rag` | `*AgentEventRag` | (rag) `ChunksFound`, `Error`, `Chunks[]` |
| `DurationMs` | `*int` | Duração da operação |
| `Model` | `*string` | Modelo que atendeu o turno |
| `Tokens` | `*AgentEventTokens` | `Prompt`, `Completion`, `Total`, `Embedding` |
| `CreatedAt` | `time.Time` | UTC |

`AgentEventRagChunk` (cada trecho recuperado pelo RAG): `Filename`, `ChunkIndex`,
`Score`, `ScorePct` (percentual de similaridade) e `Text` (conteúdo completo).

### Semântica de reconexão e entrega

- **Reconexão automática** com backoff exponencial (1s dobrando até 30s), mantendo a
  mesma `Session` — o servidor retoma a fila pendente de onde parou.
- **Entrega confirmada** (*at-least-once*): o servidor reenvia envelopes não
  confirmados; em cenários raros de ACK perdido um evento pode chegar duplicado —
  dedup pelo `evt.UUID` se isso importar para o consumidor.
- O servidor pode fechar o socket com código `1012` (refresh forçado); o SDK trata
  como queda normal e reconecta.
- Conexão interna: as opções `BaseURL` (IP) + `Host` do `NewClient` valem também para
  o WebSocket (ver [Conexão interna](#conexão-interna-ip-direto--host-header)).

> Documentação do lado servidor (rota, escopos, protocolo de ACK, fila Redis e
> política de reentrega): `documentacao/websocket/README.md` no repositório
> `synapse-api`.

### Agente — Logs de execução (REST)

A API REST `client.Agent.ListLogs` / `client.Agent.LogsStats` retorna o histórico
persistido dos eventos do agente. A partir da **v0.0.41**, o `AgentLogItem` foi
pareado com o `AgentEvent` do WebSocket — os mesmos campos estruturados (`tool_result`,
`api_response`, `rag`, `tokens`) agora estão disponíveis **também na resposta REST**,
com os dados íntegros (sem truncamento, que continua sendo regra apenas do `detail`
para compatibilidade).

```go
resp, err := client.Agent.ListLogs(ctx, "agent-uuid", synapse.ListAgentLogsParams{
    ConversationUUID: "conv-uuid",
    Page:             1,
    Size:             20,
})
for _, log := range resp.Logs {
    fmt.Println(log.Category, log.Summary)
    fmt.Println("tokens:", log.Tokens.Prompt, "/", log.Tokens.Completion)
    fmt.Println("resultado tool:", log.ToolResult)
    fmt.Println("api response:", log.APIResponse)
    fmt.Println("rag chunks:", log.Rag)
}
```

#### `AgentLogItem` (v0.0.41+)

| Campo | Tipo | Descrição |
|---|---|---|
| `UUID` | `string` | Identidade do log |
| `TenantUUID` | `string` | Tenant dono do evento |
| `AgentUUID` / `AgentName` | `string` | Agente que gerou o evento |
| `Level` | `string` | `info` \| `warn` \| `error` |
| `Category` | `string` | `chat` \| `tool_call` \| `rag` \| `error` \| `file_process` |
| `Summary` | `string` | Resumo humano (mesmo formato do WS) |
| `Detail` | `any` | Detalhe por categoria (previews — o truncamento é regra só deste campo) |
| `Reasoning` | `*string` | Texto de raciocínio estendido do modelo (extraído do `detail`) |
| `ToolName` | `*string` | (tool_call) nome da ferramenta |
| `ToolParams` | `any` | (tool_call) parâmetros |
| `ToolSuccess` | `*bool` | (tool_call) sucesso |
| `ToolSummary` | `*string` | (tool_call) resumo truncado (legado) |
| `ToolResult` | `string` | (tool_call) resultado **íntegro** — igual ao WS |
| `APIResponse` | `string` | (tool_call) corpo da resposta da API externa **íntegro** — igual ao WS |
| `Rag` | `any` | (rag) `{chunks_found, error?, chunks[]}` com textos **completos** — igual ao WS |
| `DurationMs` | `*int` | Duração da operação |
| `Model` | `*string` | Modelo que atendeu o turno |
| `TokensUsed` | `int` | Total de tokens |
| `PromptTokens` | `int` | Tokens de entrada (prompt) |
| `CompletionTokens` | `int` | Tokens de saída (completion) |
| `EmbeddingTokens` | `int` | Tokens de embedding (RAG) |
| `Tokens` | `*AgentEventTokens` | Objeto agregado `{prompt, completion, total, embedding}` — igual ao WS |
| `CreatedAt` | `string` | Timestamp ISO 8601 |

> **Histórico (REST) vs tempo-real (WS):** ambos agora têm a mesma estrutura de
> campos (`ToolResult`, `APIResponse`, `Rag`, `Tokens`). A diferença é que o WS
> entrega os eventos no momento em que ocorrem (push), enquanto o REST é o
> histórico persistido (pull). Use o WS para monitoramento em tempo real e o REST
> para auditoria/filtro/relatórios.

#### `AgentLogStats`

```go
stats, err := client.Agent.LogsStats(ctx, "agent-uuid", synapse.ListAgentLogsParams{
    ExternalID: "protocolo-123",
})
fmt.Println("chamadas:", stats.TotalCalls)
fmt.Println("tokens (prompt):", stats.TotalPromptTokens)
fmt.Println("tokens (completion):", stats.TotalCompletionTokens)
```

| Campo | Tipo | Descrição |
|---|---|---|
| `TotalCalls` | `int64` | Total de turnos do agente |
| `TotalErrors` | `int64` | Erros |
| `TotalTokens` | `int64` | Soma de todos os tokens |
| `TotalPromptTokens` | `int64` | Tokens de entrada |
| `TotalCompletionTokens` | `int64` | Tokens de saída |
| `TotalEmbeddingTokens` | `int64` | Tokens de embedding |
| `AvgDurationMs` | `float64` | Duração média dos turnos |
| `ByModel` | `[]AgentLogModelStat` | Agregado por modelo |
| `ByConversation` | `[]AgentLogConvStat` | Agregado por conversa |

---

## Tratamento de erros

### Verificar tipo de erro com sentinels

Use `errors.Is` para identificar a categoria do erro sem precisar inspecionar o payload:

```go
resp, err := client.Auth.Login(ctx, req)
if err != nil {
    switch {
    case errors.Is(err, synapse.ErrUnauthorized):
        // 401 — credenciais inválidas ou token expirado
    case errors.Is(err, synapse.ErrForbidden):
        // 403 — sem permissão para este recurso
    case errors.Is(err, synapse.ErrNotFound):
        // 404 — recurso não encontrado
    case errors.Is(err, synapse.ErrConflict):
        // 409 — duplicidade (ex: email já em uso)
    case errors.Is(err, synapse.ErrBadRequest):
        // 400 — parâmetros inválidos
    case errors.Is(err, synapse.ErrBadGateway):
        // 502 — erro no provedor externo (OpenAI, Google…)
    case errors.Is(err, synapse.ErrInternalServer):
        // 500 — erro interno da API
    default:
        // erro de rede, timeout, etc.
    }
}
```

### Inspecionar payload completo da API

Use `synapse.AsAPIError` quando precisar dos detalhes (trace_id, causes):

```go
if apiErr, ok := synapse.AsAPIError(err); ok {
    fmt.Println("HTTP status:", apiErr.StatusCode)
    fmt.Println("Trace ID:",    apiErr.TraceID)
    fmt.Println("Mensagem:",    apiErr.Message)

    for _, cause := range apiErr.Causes {
        fmt.Printf("  campo %q: %s\n", cause.Field, cause.Message)
    }
}
```

### Sentinels disponíveis

| Sentinel                               | HTTP | Quando ocorre                             |
|----------------------------------------|------|-------------------------------------------|
| `synapse.ErrInvalidToken`              | —    | Token vazio ao criar o client             |
| `synapse.ErrUnauthorized`              | 401  | Token inválido ou expirado                |
| `synapse.ErrForbidden`                 | 403  | Sem permissão para o recurso              |
| `synapse.ErrNotFound`                  | 404  | Recurso não encontrado                    |
| `synapse.ErrConflict`                  | 409  | Duplicidade (email, documento, token…)    |
| `synapse.ErrBadRequest`                | 400  | Parâmetros inválidos                      |
| `synapse.ErrInternalServer`            | 500  | Erro interno da API                       |
| `synapse.ErrBadGateway`                | 502  | Falha no provedor externo                 |
| `synapse.ErrIntegrationNotConfigured`  | 409  | Integração não configurada para o tenant  |

---

## Exemplo completo

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/WonitTecnologia/synapse"
)

func main() {
    client, err := synapse.NewClient(os.Getenv("SYNAPSE_TOKEN"), &synapse.Options{
        Timeout: 20 * time.Second,
    })
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Validar token
    me, err := client.Auth.Healthcheck(ctx)
    if err != nil {
        if errors.Is(err, synapse.ErrUnauthorized) {
            log.Fatal("token inválido ou expirado")
        }
        log.Fatal(err)
    }
    fmt.Printf("Logado como: %s (%s)\n", me.User.Name, me.User.Role)

    // Transcrever um áudio
    audio, _ := os.ReadFile("reuniao.mp3")
    transcription, err := client.OpenAI.TranscribeAudio(ctx, synapse.TranscribeAudioRequest{
        FileName: "reuniao.mp3",
        Content:  audio,
        Model:    "whisper-1",
        Language: "pt",
    })
    if err != nil {
        if apiErr, ok := synapse.AsAPIError(err); ok {
            log.Fatalf("erro da API [trace=%s]: %s", apiErr.TraceID, apiErr.Message)
        }
        log.Fatal(err)
    }
    fmt.Println("Transcrição:", transcription.Response)
}
```
