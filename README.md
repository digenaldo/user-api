# User API — Guia Didático

API REST em Go para CRUD de usuários usando MongoDB. Este projeto segue Clean Architecture para manter o código organizado e testável.

## Arquitetura

O projeto está dividido em camadas bem definidas:

**1. Handlers (HTTP)** - `internal/handler/http`
- Recebe requisições HTTP
- Traduz para chamadas aos usecases
- Formata respostas JSON

**2. Usecases (Lógica de Negócio)** - `internal/usecase`
- Contém as regras do domínio
- Faz validações (ex: email deve ter '@')
- Orquestra chamadas ao repositório

**3. Repository (Persistência)** - `internal/repository`
- Acessa o banco de dados MongoDB
- Converte entre entidades do domínio e documentos do MongoDB
- Usa context para controlar timeouts

**4. Infra (Infraestrutura)** - `internal/infra/mongo`
- Cria e configura o cliente MongoDB
- Faz conexão e ping

**5. Domain (Entidades)** - `internal/domain`
- Define a entidade User
- Define interfaces (UserRepository, UserUseCase)
- Não depende de nada externo

## Fluxo de uma Requisição

Exemplo: `POST /api/v1/users` para criar um usuário

1. Handler recebe a requisição HTTP e decodifica o JSON
2. Handler chama `uc.CreateUser()` do usecase
3. Usecase valida o email e cria a entidade `domain.User`
4. Usecase chama `repo.Create(user)` do repositório
5. Repository converte para formato MongoDB e salva no banco
6. Repository popula o ID na entidade
7. Usecase retorna o usuário criado
8. Handler serializa para JSON e retorna `201 Created`

## Como Executar

### Pré-requisitos
- Docker ou Podman instalado

### Passo a passo

1. Clone o repositório:
```bash
git clone <repository-url>
cd user-api
```

2. Execute com Docker Compose:
```bash
docker-compose up --build
```

3. Ou com Podman (macOS):
```bash
podman machine init --now
podman machine start
podman compose up --build
```

4. Teste se está funcionando:
```bash
curl http://localhost:8080/healthz
```

A API estará disponível em `http://localhost:8080`

## Endpoints

- `GET  /healthz` - Verifica se a aplicação está respondendo
- `POST /api/v1/users` - Cria um novo usuário
- `GET  /api/v1/users` - Lista todos os usuários
- `GET  /api/v1/users/{id}` - Busca usuário por ID
- `PUT  /api/v1/users/{id}` - Atualiza um usuário
- `DELETE /api/v1/users/{id}` - Remove um usuário

**Regras:**
- Email deve conter `@` (validação no usecase)
- IDs são strings hexadecimais do ObjectID do MongoDB

## Exemplos com cURL

**Criar usuário:**
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"João Silva","email":"joao@example.com"}'
```

**Listar usuários:**
```bash
curl http://localhost:8080/api/v1/users
```

**Buscar por ID:**
```bash
curl http://localhost:8080/api/v1/users/507f1f77bcf86cd799439011
```

**Atualizar usuário:**
```bash
curl -X PUT http://localhost:8080/api/v1/users/507f1f77bcf86cd799439011 \
  -H "Content-Type: application/json" \
  -d '{"name":"João Atualizado","email":"joao.novo@example.com"}'
```

**Deletar usuário:**
```bash
curl -X DELETE http://localhost:8080/api/v1/users/507f1f77bcf86cd799439011
```

## Onde Começar a Entender o Código

1. **`cmd/api/main.go`** - Ponto de entrada. Mostra como tudo é montado (Mongo → Repository → Usecase → Handler)

2. **`internal/handler/http/user_handler.go`** - Endpoints HTTP e como traduzem para usecases

3. **`internal/usecase/user_usecase.go`** - Regras de negócio e validações

4. **`internal/repository/user_mongo_repository.go`** - Como acessamos o MongoDB, conversões de ObjectID, uso de context

## Variáveis de Ambiente

- `MONGO_URI` - URI do MongoDB (padrão: `mongodb://localhost:27017`)
- `PORT` - Porta do servidor (padrão: `8080`)

No `docker-compose.yml` essas variáveis já estão configuradas.

## Parar os Serviços

```bash
docker-compose down -v
# ou
podman compose down -v
```

## Banco de Dados

- **Database:** `userdb`
- **Collection:** `users`
- **Porta:** `27017`
- **Credenciais (docker-compose):** `root` / `root`

## Dicas para Estudar

- Siga o fluxo de uma requisição do handler até o banco
- Veja como as interfaces permitem trocar implementações
- Entenda por que usamos ponteiros em Go
- Observe como o context controla timeouts

