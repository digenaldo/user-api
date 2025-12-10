
User API — Guia Didático (para alunos)

Este repositório contém uma API REST simples em Go para CRUD de usuários,
com MongoDB como persistência. O objetivo deste README é explicar de forma
didática a arquitetura do projeto, o fluxo de uma requisição e como executar
o sistema localmente (Docker / Podman).

Sumário rápido
- Arquitetura: Clean Architecture (camadas: domain, usecase, repository, handler, infra)
- Pontos principais: estrutura de pastas, como a requisição percorre as camadas, e onde olhar para cada parte.
- Execução: instruções para Docker e Podman, e saúde da aplicação (`/healthz`).

## Arquitetura e fluxo (explicação para alunos)

1) Handlers (HTTP)
- Local: `internal/handler/http` (ex: `user_handler.go`, `health_handler.go`)
- Responsabilidade: traduzir requisições HTTP para chamadas aos usecases;
  tratar e formatar respostas HTTP (status e JSON).

2) Usecases (lógica de negócio)
- Local: `internal/usecase` (ex: `user_usecase.go`)
- Responsabilidade: regras do domínio — validações, orquestração entre
  repositório e handlers. Ex.: validação de e-mail, construção da entidade.

3) Repository (persistência)
- Local: `internal/repository` (ex: `user_mongo_repository.go`)
- Responsabilidade: acesso ao banco (MongoDB). Converte entre entidades
  de domínio e documentos BSON, lida com timeouts via `context`.

4) Infra (infraestrutura)
- Local: `internal/infra/mongo/mongo_client.go`
- Responsabilidade: criação e configuração do cliente MongoDB (connect, ping).

5) Domain (entidades e interfaces)
- Local: `internal/domain` (ex: `user.go`)
- Responsabilidade: definir entidades (User) e interfaces (UserRepository, UserUseCase)
  usadas pelas camadas superiores.

Fluxo de uma requisição (exemplo `POST /api/v1/users`):
- O cliente chama o endpoint `POST /api/v1/users` (handler em `user_handler.go`).
- O handler decodifica o JSON e chama `CreateUser` no usecase.
- O usecase valida valores (ex.: e-mail) e cria a entidade `domain.User`.
- Em seguida chama `repo.Create(user)`; a implementação Mongo (`user_mongo_repository.go`)
  persiste o documento e popula o campo `ID`.
- O usecase retorna a entidade criada ao handler, que a serializa como JSON
  e responde com `201 Created`.

Onde olhar para entender o código (passo a passo):
- `cmd/api/main.go`: ponto de entrada — inicializa Mongo, cria repo/usecase/handler e registra rotas.
- `internal/handler/http/user_handler.go`: endpoints HTTP e tradução para usecases.
- `internal/usecase/user_usecase.go`: regras de negócio e validações.
- `internal/repository/user_mongo_repository.go`: queries Mongo, conversões ObjectID, timeouts.

## Como executar (passo-a-passo)

Pré-requisitos: Docker ou Podman (no macOS, use `podman machine` antes).

1) Clonar o repositório:

```bash
git clone <repository-url>
cd user-api
```

2) Executar com Docker Compose (recomendado):

```bash
docker-compose up --build
```

3) Executar com Podman (alternativa):

```bash
# Se macOS: inicializar a VM do Podman (uma vez)
podman machine init --now
podman machine start

# Subir os serviços (usa o mesmo docker-compose.yml)
podman compose up --build
```

Observação: a aplicação escuta na porta `8080` por padrão. Se a VM do
Podman estiver ativa, `http://localhost:8080` deve apontar para a API.

4) Verificar saúde da aplicação (rápido):

```bash
curl http://localhost:8080/healthz
# Deve retornar: {"status":"ok","time":"2025-..."}
```

## Endpoints (resumo didático)

- `GET  /healthz` — endpoint leve para verificar se a aplicação responde.
- `POST /api/v1/users` — cria um usuário.
- `GET  /api/v1/users` — lista todos os usuários.
- `GET  /api/v1/users/{id}` — busca por ID (ID é ObjectID do Mongo em hex).
- `PUT  /api/v1/users/{id}` — atualiza campos (name, email).
- `DELETE /api/v1/users/{id}` — remove usuário.

Regras importantes:
- Email deve conter `@` — essa validação é feita na camada de usecase.
- IDs esperados são hex strings de ObjectID do MongoDB.

## Exemplos rápidos com cURL

Criar usuário:
```bash
curl -sS -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"João Silva","email":"joao@example.com"}'
```

Listar usuários:
```bash
curl -sS http://localhost:8080/api/v1/users
```

Buscar por ID:
```bash
curl -sS http://localhost:8080/api/v1/users/<id>
```

Atualizar usuário:
```bash
curl -sS -X PUT http://localhost:8080/api/v1/users/<id> \
  -H "Content-Type: application/json" \
  -d '{"name":"João Atualizado","email":"joao.novo@example.com"}'
```

Deletar usuário:
```bash
curl -sS -X DELETE http://localhost:8080/api/v1/users/<id>
```

## Dicas de ensino (para você usar com alunos)

- Comece mostrando `cmd/api/main.go` para explicar inicialização e injeção
  de dependências simples (criamos repo → usecase → handler).
- Peça aos alunos para seguirem o fluxo de uma requisição (handler → usecase → repo).
- Mostre `internal/repository/user_mongo_repository.go` para explicar contextos,
  timeouts e conversões entre `ObjectID` e `string`.
- Explique por que dependemos de interfaces (`domain.UserRepository`) — facilita testes.

## Variáveis de ambiente

- `MONGO_URI` - URI de conexão do MongoDB (padrão: `mongodb://localhost:27017`)
- `PORT` - Porta do servidor HTTP (padrão: `8080`)

No `docker-compose.yml` já configuramos as variáveis para rodar em containers.

## Parar os serviços

```bash
docker-compose down -v   # ou podman compose down -v
```

## Banco de dados (detalhes rápidos)

- Database: `userdb`
- Collection: `users`
- Porta MongoDB (exposta): `27017`
- Credenciais (docker-compose): username `root`, password `root`

---

Se quiser, eu posso também gerar uma versão do README em formato de aula
com slides/questões para os alunos — diga se deseja uma versão "aula".
**Resposta (200 OK):**

