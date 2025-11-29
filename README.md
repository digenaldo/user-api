# User API

API REST de CRUD de usuários desenvolvida em Go, utilizando MongoDB como banco de dados. O projeto segue os princípios de Clean Architecture com camadas bem separadas.

## Tecnologias

- **Go 1.22+** - Linguagem de programação
- **MongoDB 7.0** - Banco de dados NoSQL
- **Chi Router** - Framework HTTP para rotas
- **Docker & Docker Compose** - Containerização e orquestração

## Pré-requisitos

- Docker
- Docker Compose

## Como Executar

### 1. Clone o repositório (se aplicável)

```bash
git clone <repository-url>
cd user-api
```

### 2. Execute com Docker Compose

```bash
docker-compose up --build
```

Este comando irá:
- Construir a imagem da API
- Subir o container do MongoDB
- Subir o container da API
- Conectar automaticamente os serviços

### 3. Verifique se está funcionando

A API estará disponível em: `http://localhost:8080`

### Alternativa: Usando Makefile

O projeto inclui um Makefile simples para executar o projeto:

```bash
# Subir o projeto
make up

# Parar o projeto
make down

# Ver logs
make logs
```

## Endpoints da API

Base URL: `http://localhost:8080/api/v1/users`

### Criar Usuário

```bash
POST /api/v1/users
Content-Type: application/json

{
  "name": "João Silva",
  "email": "joao@example.com"
}
```

**Resposta (201 Created):**
```json
{
  "id": "507f1f77bcf86cd799439011",
  "name": "João Silva",
  "email": "joao@example.com"
}
```

### Listar Todos os Usuários

```bash
GET /api/v1/users
```

**Resposta (200 OK):**
```json
[
  {
    "id": "507f1f77bcf86cd799439011",
    "name": "João Silva",
    "email": "joao@example.com"
  },
  {
    "id": "507f1f77bcf86cd799439012",
    "name": "Maria Santos",
    "email": "maria@example.com"
  }
]
```

### Buscar Usuário por ID

```bash
GET /api/v1/users/{id}
```

**Resposta (200 OK):**
```json
{
  "id": "507f1f77bcf86cd799439011",
  "name": "João Silva",
  "email": "joao@example.com"
}
```

**Resposta (404 Not Found):**
```json
{
  "error": "User not found"
}
```

### Atualizar Usuário

```bash
PUT /api/v1/users/{id}
Content-Type: application/json

{
  "name": "João Silva Atualizado",
  "email": "joao.novo@example.com"
}
```

**Resposta (200 OK):**
```json
{
  "id": "507f1f77bcf86cd799439011",
  "name": "João Silva Atualizado",
  "email": "joao.novo@example.com"
}
```

### Deletar Usuário

```bash
DELETE /api/v1/users/{id}
```

**Resposta (204 No Content):**

## Exemplos de Uso com cURL

### Criar usuário
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"João Silva","email":"joao@example.com"}'
```

### Listar usuários
```bash
curl http://localhost:8080/api/v1/users
```

### Buscar usuário por ID
```bash
curl http://localhost:8080/api/v1/users/507f1f77bcf86cd799439011
```

### Atualizar usuário
```bash
curl -X PUT http://localhost:8080/api/v1/users/507f1f77bcf86cd799439011 \
  -H "Content-Type: application/json" \
  -d '{"name":"João Atualizado","email":"joao.novo@example.com"}'
```

### Deletar usuário
```bash
curl -X DELETE http://localhost:8080/api/v1/users/507f1f77bcf86cd799439011
```

## Estrutura do Projeto

```
user-api/
├── cmd/
│   └── api/
│       └── main.go              # Entrypoint da aplicação
├── internal/
│   ├── domain/
│   │   └── user.go              # Entidades e interfaces do domínio
│   ├── usecase/
│   │   └── user_usecase.go      # Lógica de negócio
│   ├── repository/
│   │   └── user_mongo_repository.go  # Implementação do repositório MongoDB
│   ├── handler/
│   │   └── http/
│   │       └── user_handler.go  # Handlers HTTP
│   └── infra/
│       └── mongo/
│           └── mongo_client.go  # Cliente MongoDB
├── Dockerfile                   # Build da aplicação
├── docker-compose.yml           # Orquestração dos serviços
├── go.mod                       # Dependências do projeto
└── go.sum                       # Checksums das dependências
```

## Variáveis de Ambiente

A aplicação utiliza as seguintes variáveis de ambiente (com valores padrão):

- `MONGO_URI` - URI de conexão do MongoDB (padrão: `mongodb://localhost:27017`)
- `PORT` - Porta do servidor HTTP (padrão: `8080`)

No `docker-compose.yml`, essas variáveis são configuradas automaticamente.

## Parar os Serviços

Para parar os containers:

```bash
docker-compose down
```

Para parar e remover os volumes (dados do MongoDB):

```bash
docker-compose down -v
```

## Banco de Dados

- **Database:** `userdb`
- **Collection:** `users`
- **Porta MongoDB:** `27017`
- **Credenciais (Docker):**
  - Username: `root`
  - Password: `root`

## Validações

- O email deve conter o caractere `@`
- Campos obrigatórios: `name` e `email` (na criação)
- IDs devem ser ObjectIDs válidos do MongoDB

## Tratamento de Erros

A API retorna os seguintes códigos HTTP:

- `200 OK` - Operação bem-sucedida
- `201 Created` - Recurso criado com sucesso
- `204 No Content` - Recurso deletado com sucesso
- `400 Bad Request` - Dados inválidos (ex: email sem @)
- `404 Not Found` - Recurso não encontrado
- `500 Internal Server Error` - Erro interno do servidor

## Licença

Este projeto é um exemplo de implementação e pode ser usado livremente.

