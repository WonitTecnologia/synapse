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
| `client.OpenRouter`  | `OpenRouterCase`  | Workspace OpenRouter (sync, modelos, service tiers, analytics) |
| `client.Collection`  | `CollectionCase`  | Coleções vetoriais (Qdrant) da base de conhecimento |
| `client.Document`    | `DocumentCase`    | Upload e vetorização de documentos               |
| `client.Agent`       | `AgentCase`       | CRUD de agentes de IA + chat (com RAG)           |
| `client.SystemAgent` | `SystemAgentCase` | Agentes de sistema (plataforma) + chat dedicado  |
| `client.Swarm`       | `SwarmCase`       | Modo Swarm do Construtor — estado e controle (stop/pause/resume) ([ver seção](#swarm-modo-swarm-do-construtor)) |
| `client.Mcp`         | `McpCase`         | Integrações MCP (Model Context Protocol)         |
| `client.ExternalApi` | `ExternalApiCase` | APIs externas (HTTP cruas) como tools do agente ([ver seção](#external-apis-api-tools)) |
| `client.Monitor`     | `MonitorCase`     | **WebSocket de monitoramento** — stream de eventos do agente em tempo real ([ver seção](#websocket-de-monitoramento-monitor)) |
| `client.ChatStream`  | `ChatStreamCase`  | **WebSocket de chat** — conversa bidirecional com agentes em tempo real ([ver seção](#websocket-de-chat-chatstream)) |

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

## OpenRouter

### Gasto mensal de um agente de sistema por tenant

Cada tipo de agente de sistema (ex.: construtor de agentes) usa uma chave
OpenRouter própria por tenant — provisionada automaticamente no primeiro chat —
isolando o gasto do agente do gasto geral do tenant.
`GetAgentMonthlyAnalytics` expõe esse gasto agregado por tenant (custo em USD,
tokens e requests do mês), ordenado por custo desc, com os totais gerais.

TENANT_ADMIN vê apenas o próprio tenant. SYSTEM_ADMIN pode informar
`TenantUUID` (drill-down) ou omiti-lo para o rollup de todos os tenants com
chave ativa daquele agente.

```go
resp, err := client.OpenRouter.GetAgentMonthlyAnalytics(ctx, synapse.OpenRouterAgentMonthlyParams{
    AgentUUID: "uuid-do-agente-de-sistema",
    Month:     "2026-08", // opcional; vazio = mês corrente
    // TenantUUID: "uuid-do-tenant", // opcional (SYSTEM_ADMIN); vazio = rollup
})
for _, item := range resp.Items {
    fmt.Println(item.TenantName, item.TotalUsage, item.TokensTotal, item.Requests)
}
```

### Service tiers de um modelo

`ListModelTiers` retorna os service tiers (`priority`, `flex`) disponíveis para
um modelo do OpenRouter, com os preços crus por token (USD, strings exatamente
como o OpenRouter devolve) e os provedores que atendem cada tier. Um tier `nil`
significa que o modelo não o oferece.

```go
resp, err := client.OpenRouter.ListModelTiers(ctx, "openai/gpt-5.2")
if resp.Tiers.Priority != nil {
    fmt.Println(resp.Tiers.Priority.PromptPrice, resp.Tiers.Priority.CompletionPrice)
    fmt.Println(resp.Tiers.Priority.Providers)
}
if resp.Tiers.Flex != nil {
    fmt.Println(resp.Tiers.Flex.PromptPrice, resp.Tiers.Flex.CompletionPrice)
    fmt.Println(resp.Tiers.Flex.Providers)
}
```

| Campo (`OpenRouterModelTierInfo`) | Tipo | Descrição |
|---|---|---|
| `PromptPrice` | `string` | Preço do token de prompt em USD (string crua do OpenRouter) |
| `CompletionPrice` | `string` | Preço do token de completion em USD (string crua do OpenRouter) |
| `Providers` | `[]string` | Provedores que atendem o tier (ex.: `openai/fast`) |

### Parâmetros suportados por modelo

Cada `OpenRouterModelInfo` devolvido por `ListModels` traz o campo
`SupportedParameters` (`supported_parameters`): a lista de parâmetros que o
modelo aceita nas chamadas (ex.: `temperature`, `top_p`, `max_tokens`,
`tools`). É a fonte que o servidor usa para calcular o `temperature_enabled`
do agente — um modelo sem `temperature` na lista faz a temperature
configurada ser ignorada nas chamadas (ver
[`temperature_enabled`](#temperature_enabled-somente-leitura)).

---

## Agent

### `temperature_enabled` (somente leitura)

O `AgentResponse` traz o campo `TemperatureEnabled` (`temperature_enabled`),
gerenciado pelo sistema e calculado do catálogo OpenRouter: `true` por
padrão; `false` quando o modelo configurado no agente não aceita o parâmetro
`temperature` — nesse caso a temperature configurada **não é enviada** nas
chamadas ao modelo. Vale para agentes comuns e de sistema. Requests de
create/update não enviam o campo (o backend o ignora/recalcula); para saber
de antemão se um modelo aceita temperature, consulte `SupportedParameters`
em [`OpenRouter.ListModels`](#parâmetros-suportados-por-modelo).

### Listar conversas

`client.Agent.ListConversations` retorna a lista paginada de conversas
(`ConversationResponse`): `UUID`, `TenantUUID`, `AgentUUID`, `ExternalID`,
`MessageCount`, `CreatedAt`, `UpdatedAt` — e, em conversas atendidas por
agentes de sistema/construtor (ex.: Builder Agent):

| Campo | Tipo | Descrição |
|---|---|---|
| `Title` | `string` | Título da conversa gerado por IA (vazio para agentes comuns) |
| `Focus` | `*ConversationFocus` | O que o agente de sistema está fazendo na conversa agora (`nil` para agentes comuns): `Action` (`analisando`/`criando`/`modificando`), `AgentName`, `TargetAgentUUID` (agente sendo criado/modificado, quando aplicável) e `UpdatedAt` |

### Transferência entre agentes de IA

Um agente pode transferir o atendimento para outro agente de IA do mesmo tenant.
A configuração é feita pelo campo `TransferAgentUUIDs` (`transfer_agent_uuids`),
presente em `CreateAgentRequest`, `UpdateAgentRequest` e `AgentResponse` — a lista
de UUIDs dos agentes para os quais ele pode transferir a conversa:

```go
resp, err := client.Agent.Create(ctx, synapse.CreateAgentRequest{
    Name:               "Atendente N1",
    Model:              "openai/gpt-4o-mini",
    Prompt:             "...",
    TransferAgentUUIDs: []string{"uuid-do-agente-n2"},
})
```

No update (`UpdateAgentRequest`), o campo é `*[]string` e segue a semântica dos
demais campos de lista: `nil` = sem alteração, `[]string{}` = remove todos,
`["uuid1"]` = substitui a lista inteira.

Comportamento de cascata: quando um agente é removido, ele some automaticamente
das listas de transferência dos demais agentes — não é preciso atualizar cada
agente manualmente.

No chat, quando ocorre uma transferência, o `ChatResponse` traz o campo `Transfer`
(`*AgentTransferInfo`, `nil` quando não houve transferência):

```go
chat, err := client.Agent.Chat(ctx, synapse.ChatRequest{ /* ... */ })
if chat.Transfer != nil {
    fmt.Printf("transferido para %s (%s): %s\n",
        chat.Transfer.TargetAgentName,
        chat.Transfer.TargetAgentUUID,
        chat.Transfer.Summary,
    )
}
```

| Campo             | Tipo     | Descrição                                    |
|-------------------|----------|----------------------------------------------|
| `TargetAgentUUID` | `string` | UUID do agente que recebeu a conversa        |
| `TargetAgentName` | `string` | Nome do agente que recebeu a conversa        |
| `Summary`         | `string` | Resumo do atendimento repassado ao agente alvo |

---

## System Agent

Agentes de sistema são agentes da plataforma (sem tenant dono) usados em fluxos
internos — como o construtor de agentes. O chat roda num pipeline dedicado
(sem juiz, tools `sistema_*`) e a cobrança cai na chave OpenRouter própria do
tipo de agente dentro do workspace do tenant chamador (provisionada
automaticamente no primeiro chat), isolada do gasto geral do tenant — ver
[Gasto mensal de um agente de sistema por tenant](#gasto-mensal-de-um-agente-de-sistema-por-tenant).

No chat com agentes de sistema, o `ChatResponse` pode trazer dois campos extras
gerados por IA no turno: `Title` (título da conversa, presente no turno em que
foi gerado) e `Suggestions` (sugestões de próxima mensagem do usuário). Agentes
comuns nunca preenchem esses campos. No transporte WebSocket eles chegam no
`ChatStreamMessage` recebido (ver [WebSocket de chat](#websocket-de-chat-chatstream)).

### Cancelar job de chat

No transporte WebSocket o chat responde `202` com um `job_id` e a resposta chega
depois pelo stream. `CancelChat` cancela esse job: se ainda estiver na fila, ele
é descartado; se já estiver em execução, o turno é abortado no worker. `job_id`
desconhecido **não** é erro — a flag de cancelamento vale até o job aparecer ou
expirar (15 min).

```go
err := client.SystemAgent.CancelChat(ctx, synapse.CancelChatRequest{
    JobID: "uuid-do-job",
})
```

---

## Swarm (Modo Swarm do Construtor)

O Modo Swarm permite ao Agente Construtor (system agent) **orquestrar a si
mesmo**: diante de uma tarefa grande e decomponível (ex.: "crie 3 agentes de
vendas", "audite os prompts de todos os agentes"), ele abre um **swarm** — até
8 subagentes efêmeros do próprio Construtor rodando **em paralelo**, cada um
com contexto zerado (somente o prompt da subtarefa + system prompt base) e
**sem tools de escrita** (planos/apply ficam bloqueados no papel subagente — só
leitura, análise e rascunho em texto). Quando todas as tarefas terminam, os
relatórios voltam agregados como um turno sintético na conversa principal e o
orquestrador consolida o resultado para o usuário.

- Há no máximo **um swarm ativo por conversa** — abrir outro retorna erro
  orientando concluir/parar o atual.
- **Parar/pausar param os serviços de verdade**: o cancelamento propaga para
  cada job filho e derruba as chamadas LLM em voo (sem queimar token parado).
- Toda transição (queued→running→done, pausa, stop, stall detectado) é
  registrada na **timeline** (`History`) e publicada como evento em tempo real
  (ver [Eventos do swarm](#eventos-do-swarm-em-tempo-real)).

### Modos de permissão e modo swarm no chat

O `ChatRequest` do system agent aceita dois campos por conversa — enviados a
cada mensagem e persistidos pelo backend (sobrevivem a reload da página).
Agentes comuns (não-system) ignoram os dois campos.

| Campo | Tipo | Valores | Descrição |
|---|---|---|---|
| `PermissionMode` | `string` | `always_ask` (default) \| `ask_when_needed` \| `never_ask` | Como o apply de planos é autorizado nesta conversa |
| `SwarmMode` | `string` | `on` \| `off` (default) | `on` instrui o Construtor a decompor a tarefa e abrir um swarm |

| Modo de permissão | Comportamento |
|---|---|
| `PermissionModeAlwaysAsk` (`always_ask`) | Padrão, comportamento atual: todo apply de plano só ocorre após autorização explícita do usuário |
| `PermissionModeAskWhenNeeded` (`ask_when_needed`) | Planos **somente-criação** (classificação determinística no backend) são aplicados automaticamente; planos com update/delete ainda perguntam |
| `PermissionModeNeverAsk` (`never_ask`) | "Yolo": todo plano finalizado é aplicado **na mesma cadeia**, sem mensagem de autorização. Cada auto-apply é auditado nos logs do agente |

Os dois modos também trafegam no [WebSocket de chat](#websocket-de-chat-chatstream):
`ChatStreamMessage.PermissionMode`/`SwarmMode` são propagados pelo
`ChatStream.Send` e serializados no frame de envio como `permission_mode`/
`swarm_mode` — mesma semântica do REST (a cada mensagem, persistidos por
conversa, ignorados por agentes não-system).

### Tipos

`SwarmState` — estado completo do swarm, mantido pelo backend no Redis
(TTL de 24h renovado a cada acesso):

| Campo | Tipo | Descrição |
|---|---|---|
| `SwarmID` | `string` | Identificador do swarm |
| `ConversationUUID` | `string` | Conversa principal (do Construtor) dona do swarm |
| `ParentJobID` | `string` | Job de dispatch do turno orquestrador que abriu o swarm |
| `Status` | `SwarmStatus` | Status do swarm (mesma escala das tarefas) |
| `PermissionMode` | `string` | Modo de permissão vigente quando o swarm foi aberto |
| `Tasks` | `[]SwarmTask` | Uma entrada por tarefa/subagente |
| `History` | `[]SwarmHistoryEntry` | Timeline de transições (ring buffer, máx. 50) — trilha de auditoria |
| `Counters` | `SwarmCounters` | Progresso agregado: `Total`, `Queued`, `Running`, `Done`, `Failed`, `Paused`, `Stopped` |

`SwarmTask` — uma tarefa executada por um subagente efêmero:

| Campo | Tipo | Descrição |
|---|---|---|
| `TaskID` | `string` | Identificador da tarefa dentro do swarm |
| `Title` | `string` | Título curto da subtarefa |
| `Role` | `string` | Papel pedido pelo orquestrador (`deep` = modelo cheio; vazio = tier econômico do modelo do Construtor) |
| `Model` | `string` | Modelo LLM efetivamente atribuído ao subagente |
| `Status` | `SwarmStatus` | `queued` \| `running` \| `done` \| `failed` \| `paused` \| `stopped` \| `stalled` (constantes `SwarmStatus*`) |
| `JobID` | `string` | Job de dispatch executando a tarefa (correlaciona com os logs do agente) |
| `Report` | `string` | Relatório final do subagente (agregado na consolidação) |
| `QueuedAt` / `StartedAt` / `FinishedAt` | `string` | Marcos temporais — `queued_at → started_at` é o tempo real de fila |
| `LastProgressAt` | `string` | Sinal de vida, renovado a cada tool call do subagente (base do detector de travamento) |
| `Error` | `string` | Erro instrutivo (e relatório parcial) de uma tarefa falha |

`SwarmHistoryEntry` — uma entrada da timeline: `Ts` (timestamp), `TaskID`
(vazio em eventos do swarm inteiro), `From`/`To` (`SwarmStatus`) e `Detail`.

### Métodos (`client.Swarm`)

| Método | Endpoint | Descrição |
|---|---|---|
| `GetSwarm(ctx, swarmID)` | `GET /api/agent/application/system-agent/swarm/:id` | Estado completo do swarm (tarefas, contadores e timeline) |
| `StopSwarm(ctx, swarmID)` | `POST .../swarm/:id/stop` | Cancela **todos** os jobs filhos (LLM em voo morre na hora); a consolidação roda com os relatórios existentes |
| `PauseSwarm(ctx, swarmID)` | `POST .../swarm/:id/pause` | Congela o swarm: tarefas running são canceladas (sem gasto de token) e queued aguardam retomada |
| `ResumeSwarm(ctx, swarmID)` | `POST .../swarm/:id/resume` | Limpa a pausa e re-enfileira as tarefas pendentes como jobs novos |

Todos retornam o `*SwarmState` atualizado. Use token de **tenant**: swarms são
escopados pelo tenant chamador.

### Eventos do swarm em tempo real

Cada transição de tarefa (e do swarm inteiro) é publicada no
[WebSocket de monitoramento](#websocket-de-monitoramento-monitor) como um
`AgentEvent` de categoria `EventCategorySwarm` (`"swarm"`). O campo `Detail`
carrega um `SwarmEventPayload`:

| Campo | Tipo | Descrição |
|---|---|---|
| `SwarmID` | `string` | Swarm de origem |
| `TaskID` | `string` | Tarefa (vazio em eventos do swarm inteiro, ex.: pausa geral) |
| `Title` / `Role` / `Model` | `string` | Descrição da tarefa |
| `Status` | `SwarmStatus` | Novo status da tarefa/swarm |
| `Progress` | `int` | Progresso da tarefa, quando informado |
| `Error` | `string` | Erro, em transições para `failed` |
| `ConversationUUID` / `TenantUUID` | `string` | Correlação com a conversa e o tenant |

### Exemplo completo

```go
// 1. Chat com o Construtor ativando swarm e modo de permissão da conversa.
resp, err := client.SystemAgent.Chat(ctx, synapse.ChatRequest{
	AgentUUID:        "uuid-do-construtor",
	Message:          "Crie 3 agentes de vendas (N1, N2 e retenção).",
	ConversationUUID: &convUUID,
	PermissionMode:   synapse.PermissionModeAskWhenNeeded,
	SwarmMode:        synapse.SwarmModeOn,
})
if err != nil {
	log.Fatal(err)
}

// 2. Acompanhamento ao vivo: eventos de categoria swarm no Monitor.
stream, _ := client.Monitor.StreamLogs(ctx, nil)
defer stream.Close()
go func() {
	for evt := range stream.Events() {
		if evt.Category != synapse.EventCategorySwarm {
			continue
		}
		var p synapse.SwarmEventPayload
		b, _ := json.Marshal(evt.Detail)
		if json.Unmarshal(b, &p) == nil {
			fmt.Printf("[swarm %s] %s %s → %s\n", p.SwarmID, p.TaskID, p.Title, p.Status)
		}
	}
}()

// 3. Estado pontual + timeline (painel de acompanhamento).
state, err := client.Swarm.GetSwarm(ctx, swarmID)
if err != nil {
	log.Fatal(err)
}
fmt.Printf("progresso: %d/%d concluídas\n", state.Counters.Done, state.Counters.Total)
for _, t := range state.Tasks {
	fmt.Printf("- %s [%s] %s\n", t.TaskID, t.Status, t.Title)
}
for _, h := range state.History {
	fmt.Printf("  %s %s: %s → %s %s\n", h.Ts, h.TaskID, h.From, h.To, h.Detail)
}

// 4. Controle: pausar (congela sem gasto), retomar ou parar tudo.
state, err = client.Swarm.PauseSwarm(ctx, swarmID)
state, err = client.Swarm.ResumeSwarm(ctx, swarmID)
state, err = client.Swarm.StopSwarm(ctx, swarmID)
```

---

## External APIs (API tools)

APIs externas são endpoints HTTP crus cadastrados pelo tenant que o agente de IA
invoca como tools de function calling (`client.ExternalApi` — CRUD completo:
`Create`, `Get`, `List`, `Update`, `Toggle`, `Delete`). Cada parâmetro declarado
em `Parameters` (`ExternalApiParamDef`) vira uma propriedade do JSON Schema da
tool e é interpolado como `{{nome}}` na URL, nos headers ou no corpo, conforme o
`Location`.

### Tipo de corpo (`body_type`)

O campo `BodyType` (`body_type`) — presente em `CreateExternalApiRequest`,
`UpdateExternalApiRequest` e `ExternalApiResponse` — define como o executor
monta o corpo da requisição: `none` | `json` | `multipart` | `form_urlencoded` |
`raw`. Vazio = `json` (retrocompatível).

| Valor | Comportamento do executor |
|---|---|
| `json` (default) | `BodyTemplate` interpolado, enviado com `Content-Type: application/json` |
| `raw` | `BodyTemplate` interpolado enviado como string crua; o `Content-Type` vem dos headers configurados (não é forçado json) |
| `multipart` | Corpo `multipart/form-data` **combinável**: parâmetros com `Location: "form"` viram campos texto (com `Type: "object"`/`"array"` são serializados como JSON no campo — ex.: `metadata={"a":1}`); parâmetros com `Location: "file"` / `Type: "file"` viram campos arquivo |
| `form_urlencoded` | Parâmetros com `Location: "form"` viram corpo `application/x-www-form-urlencoded` (`object`/`array` também viram JSON no valor) |
| `none` | Sem corpo |

Parâmetros `query`, `path` e `header` combinam livremente com qualquer
`body_type` — o tipo de corpo só define como o **corpo** é montado.

### Parâmetros de arquivo (`type: "file"`)

`Type` aceita `string|number|integer|boolean|array|object|file` e `Location`
aceita `query|path|header|body|form|file`. No schema de function calling, um
parâmetro `file` aparece ao LLM como `string` (URL), com descrição orientando
que o valor deve ser a URL de um arquivo — por exemplo, a URL de mídia de um
anexo recebido na conversa. O executor baixa essa URL e a anexa como binário no
campo de arquivo do multipart (também aceita valores `data:...;base64,...`).

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
| `Category` | `string` | `EventCategoryChat` \| `EventCategoryToolCall` \| `EventCategoryRAG` \| `EventCategoryError` \| `EventCategoryFileProcess` \| `EventCategorySwarm` |
| `Summary` | `string` | Resumo humano do evento |
| `Detail` | `map[string]any` | Detalhe por categoria (chat: `user_msg`/`response` íntegros, `reasoning`…; swarm: `SwarmEventPayload` — ver [Swarm](#swarm-modo-swarm-do-construtor)) |
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

## WebSocket de chat (ChatStream)

Conversa **bidirecional** em tempo real com agentes de IA na mesma conexão:
você envia mensagens com `Send` e recebe, em canais separados, as confirmações
de envio (`accepted`/`error`), as respostas do agente, os chunks incrementais
da resposta (streaming efêmero, sem ACK) e os eventos de execução
das conversas da sessão. O SDK confirma automaticamente os envelopes recebidos
(protocolo de ACK interno, igual ao monitor) e reconecta sozinho mantendo a
mesma `Session`.

- Token **master (SYSTEM_ADMIN)** → conversa com **qualquer agente**.
- Token de **tenant** → conversa apenas com agentes do **próprio tenant**.

### Uso básico

```go
stream, err := client.ChatStream.Stream(ctx, nil)
if err != nil {
    log.Fatal(err)
}
defer stream.Close()

// enviar mensagem (UUID gerado pelo SDK se vazio)
err = stream.Send(synapse.ChatStreamMessage{
    AgentUUID: "agent-uuid",
    Message:   "Olá!",
})
if errors.Is(err, synapse.ErrStreamNotConnected) {
    // desconectado/reconectando: reenvie a mensagem (não há fila no cliente)
}

// consumir confirmações, respostas e eventos
for {
    select {
    case ack := <-stream.Accepted():
        if ack.Error != "" {
            log.Printf("mensagem %s rejeitada: %s", ack.UUID, ack.Error)
        } else {
            log.Printf("mensagem %s aceita (%s), job %s", ack.UUID, ack.Status, ack.JobID)
        }
    case msg := <-stream.Messages():
        fmt.Printf("[%s] %s\n", msg.AgentName, msg.Message)
    case evt := <-stream.Events():
        fmt.Printf("evento %s — %s\n", evt.Category, evt.Summary)
    case <-ctx.Done():
        return
    }
}
// os canais fecham quando ctx é cancelado ou stream.Close() é chamado
```

### Opções (`ChatStreamOptions`)

Mesma semântica do monitor (`StreamLogsOptions`):

| Campo | Tipo | Descrição |
|---|---|---|
| `Session` | `string` | UUID de sessão. Reconexões com a mesma session **retomam a fila de entrega** pendente no servidor. Padrão: UUID aleatório mantido pela vida do stream |
| `Buffer` | `int` | Capacidade dos canais (`Accepted`/`Messages`/`Chunks`/`Events`). Padrão: `256` |
| `OnConnect` | `func(session string)` | Disparado a cada handshake bem-sucedido — conexão inicial **e** cada reconexão automática |
| `OnError` | `func(error)` | Erros de conexão/handshake. Apenas observabilidade: a reconexão é automática (backoff 1s → 30s) |

### `ChatStream`

| Método | Descrição |
|---|---|
| `Send(msg ChatStreamMessage) error` | Envia mensagem ao agente. Gera UUID v4 se `msg.UUID` vazio. Thread-safe. **Sem fila no cliente**: retorna `ErrStreamNotConnected` imediatamente se desconectado — o chamador reenvia |
| `Accepted() <-chan ChatStreamAccepted` | Confirmações de envio (`Status` = `queued`/`throttled`, ou `Error` preenchido) |
| `Messages() <-chan ChatStreamMessage` | Respostas do agente (ACK automático). Erro terminal do job chega aqui com `Error` preenchido e `Message` vazia |
| `Chunks() <-chan ChatStreamChunk` | Chunks incrementais da resposta (streaming). **Efêmeros, sem ACK**: a resposta definitiva continua chegando em `Messages()` — chunks após ela devem ser ignorados. `Kind` distingue `content` (texto da resposta), `reasoning` (raciocínio do modelo) e `boundary` (fecha o parcial como mensagem intermediária; o próximo `content` abre bolha nova). `Reset=true` no 1º chunk de uma nova tentativa de modelo (fallback): zere o texto e o raciocínio parciais acumulados antes de aplicar o `Delta` |
| `Events() <-chan AgentEvent` | Eventos de execução das conversas da sessão (mesmo tipo do monitor, ACK automático) |
| `Session() string` | UUID da sessão em uso |
| `Close()` | Encerra o stream e fecha os canais |

### `ChatStreamMessage`

| Campo | Tipo | Descrição |
|---|---|---|
| `UUID` | `string` | UUID do frame (envio: gerado pelo SDK se vazio; use para correlacionar com `Accepted`) |
| `JobID` | `string` | (recebimento) job que processou a mensagem |
| `AgentUUID` | `string` | Agente destino (envio) / origem (recebimento) |
| `AgentName` | `string` | (recebimento) nome do agente |
| `ConversationUUID` | `string` | Conversa (opcional no envio; preenchido nas respostas) |
| `Message` | `string` | Texto da mensagem |
| `Context` | `string` | (envio, opcional) contexto adicional |
| `Attachment` | `*ChatStreamAttachment` | (opcional) anexo: `URL`, `Type` (`image`/`audio`/`document`), `MimeType`, `FileName` |
| `PermissionMode` | `string` | (envio, opcional) modo de permissão da conversa em system agents (constantes `PermissionMode*` — ver [Modos de permissão e modo swarm no chat](#modos-de-permissão-e-modo-swarm-no-chat)); serializado no frame como `permission_mode` |
| `SwarmMode` | `string` | (envio, opcional) modo swarm da conversa em system agents (constantes `SwarmMode*` — ver [Modos de permissão e modo swarm no chat](#modos-de-permissão-e-modo-swarm-no-chat)); serializado no frame como `swarm_mode` |
| `Title` | `string` | (recebimento) título da conversa gerado por IA — presente no turno em que foi gerado. Hoje só agentes de sistema/construtor (ex.: Builder Agent) produzem |
| `Suggestions` | `[]string` | (recebimento) sugestões geradas por IA para a próxima mensagem do usuário — presentes no turno em que foram geradas. Hoje só agentes de sistema/construtor (ex.: Builder Agent) produzem |
| `Error` | `string` | (recebimento) erro terminal do job (ex.: agente não encontrado) — quando preenchido, `Message` vem vazia |

### `ChatStreamChunk`

| Campo | Tipo | Descrição |
|---|---|---|
| `JobID` | `string` | Job que está gerando a resposta |
| `ConversationUUID` | `string` | Conversa da resposta |
| `Kind` | `string` | Tipo do chunk: `content` (trecho do texto da resposta — padrão quando omitido), `reasoning` (trecho do raciocínio do modelo, ex.: qwen — exiba separado do texto da resposta) ou `boundary` (`Delta` vazio: o parcial acumulado vira mensagem intermediária definitiva e os próximos chunks `content` abrem bolha nova — permite N mensagens do agente por turno). Constantes `ChunkKindContent`/`ChunkKindReasoning`/`ChunkKindBoundary` |
| `Delta` | `string` | Trecho incremental do texto da resposta (ou do raciocínio, quando `Kind=reasoning`) |
| `Reset` | `bool` | `true` no 1º chunk de uma **nova tentativa de modelo** (fallback): zere o parcial acumulado (texto **e** raciocínio) antes de continuar |

### Reconexão e confirmação de envio

- **Reconexão automática** com backoff exponencial (1s dobrando até 30s), mantendo a
  mesma `Session` — o servidor retoma a fila pendente de onde parou.
- **`Send` não enfileira**: durante uma queda/reconexão ele falha na hora com
  `ErrStreamNotConnected`; guarde a mensagem e reenvie (com o mesmo `UUID`) se
  precisar de garantia.
- Toda mensagem enviada gera uma confirmação em `Accepted()` com o mesmo `UUID`
  do frame — `Status` (`queued`/`throttled`) em caso de aceite, `Error` em caso
  de rejeição.
- Envelopes `message`/`event` seguem entrega confirmada (*at-least-once*), como
  no monitor: dedup pelo `UUID` do envelope se isso importar para o consumidor.
- Conexão interna: as opções `BaseURL` (IP) + `Host` do `NewClient` valem também para
  este WebSocket (ver [Conexão interna](#conexão-interna-ip-direto--host-header)).

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
| `synapse.ErrStreamNotConnected`        | —    | `Send` no WebSocket de chat enquanto desconectado/reconectando |

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
