# Aula: Construindo uma API REST em Go com Clean Architecture

## 📋 Objetivo da Aula

Construir uma API REST completa para gerenciamento de usuários (CRUD) utilizando:
- **Go** como linguagem
- **MongoDB** como banco de dados
- **Clean Architecture** como padrão arquitetural
- **Docker** para containerização

---

## 🎯 Slide 1: Apresentação do Projeto

### O que vamos construir?

Uma API REST que permite:
- ✅ Criar usuários
- ✅ Listar todos os usuários
- ✅ Buscar usuário por ID
- ✅ Atualizar usuário
- ✅ Deletar usuário

### Endpoints que teremos:
```
POST   /api/v1/users      - Criar usuário
GET    /api/v1/users      - Listar usuários
GET    /api/v1/users/{id} - Buscar usuário
PUT    /api/v1/users/{id} - Atualizar usuário
DELETE /api/v1/users/{id} - Deletar usuário
```

---

## 🏗️ Slide 2: Clean Architecture - Conceitos Básicos

### Por que Clean Architecture?

- **Separação de responsabilidades**
- **Testabilidade**
- **Independência de frameworks**
- **Facilita manutenção**

### Camadas que vamos criar:

```
┌─────────────────────────────────┐
│   Handler (HTTP)                │  ← Interface com o mundo externo
├─────────────────────────────────┤
│   UseCase (Lógica de Negócio)   │  ← Regras de negócio
├─────────────────────────────────┤
│   Repository (Dados)            │  ← Acesso a dados
├─────────────────────────────────┤
│   Domain (Entidades)            │  ← Modelos e interfaces
└─────────────────────────────────┘
```

**Regra de ouro:** As camadas internas NÃO conhecem as externas!

---

## 📁 Slide 3: Estrutura de Pastas

### Vamos criar a seguinte estrutura:

```
user-api/
├── cmd/
│   └── api/
│       └── main.go              # Ponto de entrada
├── internal/
│   ├── domain/
│   │   └── user.go              # Entidades e interfaces
│   ├── usecase/
│   │   └── user_usecase.go      # Lógica de negócio
│   ├── repository/
│   │   └── user_mongo_repository.go  # Implementação MongoDB
│   ├── handler/
│   │   └── http/
│   │       └── user_handler.go  # Handlers HTTP
│   └── infra/
│       └── mongo/
│           └── mongo_client.go  # Cliente MongoDB
├── go.mod
├── Dockerfile
└── docker-compose.yml
```

**Pergunta para os alunos:** Por que `internal/`? (Resposta: Go não permite importar pacotes de `internal/` de fora do módulo)

---

## 🚀 Slide 4: Setup Inicial

### Passo 1: Inicializar o projeto

```bash
mkdir user-api
cd user-api
go mod init user-api
```

### Passo 2: Instalar dependências

```bash
go get github.com/go-chi/chi/v5
go get go.mongodb.org/mongo-driver/mongo
```

### Passo 3: Criar estrutura de pastas

```bash
mkdir -p cmd/api
mkdir -p internal/{domain,usecase,repository,handler/http,infra/mongo}
```

**Explicar:** Vamos criar os arquivos vazios primeiro, depois preencher.

---

## 🎨 Slide 5: Domain Layer - O Coração da Aplicação

### O que é o Domain?

- **Entidades**: Modelos de dados (User)
- **Interfaces**: Contratos que outras camadas devem implementar

### Vamos criar `internal/domain/user.go`

**Conceitos a explicar:**
1. **Struct User**: Representa nossa entidade
2. **Interface UserRepository**: Define o que um repositório DEVE fazer
3. **Interface UserUseCase**: Define o que a lógica de negócio DEVE fazer

**Pergunta:** Por que usar interfaces? (Resposta: Permite trocar implementações sem mudar o código que usa)

### Estrutura básica:

```go
// Entidade
type User struct {
    ID    string
    Name  string
    Email string
}

// Interface do Repository
type UserRepository interface {
    Create(user *User) error
    GetByID(id string) (*User, error)
    List() ([]*User, error)
    Update(user *User) error
    Delete(id string) error
}

// Interface do UseCase
type UserUseCase interface {
    CreateUser(name, email string) (*User, error)
    GetUser(id string) (*User, error)
    ListUsers() ([]*User, error)
    UpdateUser(id, name, email string) (*User, error)
    DeleteUser(id string) error
}
```

**Dica:** Começar simples, depois adicionar validações.

---

## 💾 Slide 6: Infrastructure Layer - Cliente MongoDB

### Por que separar a infraestrutura?

- Facilita testes (mock)
- Permite trocar banco de dados
- Isola detalhes técnicos

### Vamos criar `internal/infra/mongo/mongo_client.go`

**Conceitos:**
- Conexão com MongoDB
- Context com timeout
- Ping para verificar conexão

**Estrutura básica:**
```go
func NewClient(uri string) *mongo.Client {
    // Criar contexto com timeout
    // Configurar opções de conexão
    // Conectar
    // Fazer ping
    // Retornar cliente
}
```

**Explicar:** Por que usar context? (Controle de timeout, cancelamento)

---

## 📦 Slide 7: Repository Layer - Acesso aos Dados

### O que faz o Repository?

- **Abstrai** o acesso ao banco de dados
- **Implementa** a interface do Domain
- **Converte** entre entidades do domínio e documentos do MongoDB

### Vamos criar `internal/repository/user_mongo_repository.go`

**Conceitos importantes:**
1. **userDoc**: Estrutura para MongoDB (usa `primitive.ObjectID`)
2. **User**: Estrutura do domínio (usa `string` para ID)
3. **Conversão** entre os dois formatos

**Operações CRUD:**
- `Create`: Insere e retorna ID gerado
- `GetByID`: Busca por ID (converte string → ObjectID)
- `List`: Retorna todos
- `Update`: Atualiza por ID
- `Delete`: Remove por ID

**Pergunta:** Por que não usar User diretamente no MongoDB? (Resposta: MongoDB usa ObjectID, nosso domínio usa string)

---

## 🧠 Slide 8: UseCase Layer - Lógica de Negócio

### O que faz o UseCase?

- **Orquestra** as operações
- **Aplica regras de negócio**
- **Valida dados**
- **Trata erros**

### Vamos criar `internal/usecase/user_usecase.go`

**Regras de negócio que vamos implementar:**
1. Email deve conter "@"
2. Verificar se usuário existe antes de atualizar/deletar

**Estrutura:**
```go
type userUseCase struct {
    repo domain.UserRepository  // Dependência
}

// Implementa domain.UserUseCase
```

**Conceitos:**
- **Injeção de dependência**: Repository vem de fora
- **Erros customizados**: `ErrInvalidEmail`, `ErrNotFound`
- **Validações**: Antes de salvar/atualizar

**Explicar:** Por que validações aqui e não no handler? (Lógica de negócio pertence ao UseCase)

---

## 🌐 Slide 9: Handler Layer - Interface HTTP

### O que faz o Handler?

- **Recebe** requisições HTTP
- **Valida** formato dos dados
- **Chama** o UseCase
- **Retorna** respostas HTTP

### Vamos criar `internal/handler/http/user_handler.go`

**Conceitos:**
1. **Chi Router**: Framework para rotas
2. **JSON encoding/decoding**
3. **Status codes HTTP**
4. **Tratamento de erros**

**Estrutura de um handler:**
```go
type UserHandler struct {
    uc domain.UserUseCase  // Dependência
}

// Métodos:
// - createUser
// - listUsers
// - getUser
// - updateUser
// - deleteUser
// - RegisterRoutes (configura rotas)
```

**Funções auxiliares:**
- `writeJSON`: Escreve resposta JSON
- `writeError`: Escreve erro JSON

**Explicar:** Por que separar em funções auxiliares? (DRY - Don't Repeat Yourself)

---

## 🔌 Slide 10: Main - Composição e Inicialização

### O que faz o main.go?

- **Conecta todas as camadas**
- **Configura o servidor HTTP**
- **Inicia a aplicação**

### Vamos criar `cmd/api/main.go`

**Fluxo de inicialização:**
1. Ler variáveis de ambiente
2. Criar cliente MongoDB
3. Criar repository
4. Criar usecase
5. Criar handler
6. Configurar rotas
7. Iniciar servidor

**Conceitos:**
- **Dependency Injection**: Passar dependências manualmente
- **Defer**: Garantir desconexão do MongoDB
- **Variáveis de ambiente**: Configuração flexível

**Estrutura:**
```go
func main() {
    // 1. Config
    // 2. MongoDB Client
    // 3. Repository
    // 4. UseCase
    // 5. Handler
    // 6. Router
    // 7. Server
}
```

---

## 🐳 Slide 11: Docker - Containerização

### Por que Docker?

- **Ambiente consistente**
- **Fácil de rodar**
- **Isola dependências**

### Vamos criar `Dockerfile`

**Conceitos:**
- Multi-stage build (otimização)
- Build da aplicação
- Imagem final minimalista

### Vamos criar `docker-compose.yml`

**Serviços:**
1. **MongoDB**: Banco de dados
2. **API**: Nossa aplicação

**Conceitos:**
- Networks (comunicação entre containers)
- Volumes (persistência de dados)
- Environment variables
- Ports mapping

**Explicar:** Por que usar docker-compose? (Orquestra múltiplos serviços)

---

## ✅ Slide 12: Testando a API

### Como testar?

**Opção 1: cURL**
```bash
# Criar usuário
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"João","email":"joao@example.com"}'

# Listar usuários
curl http://localhost:8080/api/v1/users
```

**Opção 2: Postman/Insomnia**
- Interface gráfica
- Mais fácil para testes

**Opção 3: Testes automatizados**
- Unit tests
- Integration tests
- (Tópico para outra aula)

---

## 📝 Slide 13: Resumo e Próximos Passos

### O que aprendemos?

✅ Clean Architecture em Go  
✅ Separação de responsabilidades  
✅ CRUD completo  
✅ Integração com MongoDB  
✅ Docker e Docker Compose  

### Possíveis melhorias (para próximas aulas):

- 🔐 Autenticação e autorização
- ✅ Validações mais robustas
- 📊 Logging estruturado
- 🧪 Testes unitários e de integração
- 📈 Métricas e observabilidade
- 🔄 Middleware (CORS, rate limiting)
- 📄 Documentação com Swagger

---

## 🎓 Conceitos-Chave para Reforçar

### 1. Clean Architecture
- **Dependências apontam para dentro**
- **Domain não depende de nada**
- **Interfaces no domain, implementações fora**

### 2. Injeção de Dependência
- **Não criar dependências dentro das funções**
- **Receber por parâmetro**
- **Facilita testes**

### 3. Tratamento de Erros
- **Erros customizados**
- **Propagação correta**
- **Status HTTP apropriados**

### 4. Context em Go
- **Timeout**
- **Cancelamento**
- **Propagação de valores**

---

## 💡 Dicas para o Professor

### Durante o desenvolvimento:

1. **Começar pelo Domain** - É o coração, não depende de nada
2. **Implementar de dentro para fora** - Domain → Repository → UseCase → Handler
3. **Testar cada camada** - Antes de passar para a próxima
4. **Explicar o "porquê"** - Não só o "como"
5. **Mostrar erros comuns** - E como resolver

### Perguntas para engajar:

- "Por que não colocar tudo no main.go?"
- "O que acontece se mudarmos de MongoDB para PostgreSQL?"
- "Como testaríamos isso sem banco de dados?"
- "Onde colocaríamos validação de email mais robusta?"

---

## 📚 Recursos Adicionais

### Documentação:
- [Go Documentation](https://go.dev/doc/)
- [MongoDB Go Driver](https://www.mongodb.com/docs/drivers/go/current/)
- [Chi Router](https://github.com/go-chi/chi)

### Conceitos para estudar:
- Clean Architecture (Robert C. Martin)
- SOLID Principles
- Dependency Injection
- Repository Pattern

---

## 🎯 Checklist de Implementação

Use este checklist durante o desenvolvimento:

- [ ] Setup inicial (go.mod, estrutura de pastas)
- [ ] Domain layer (User, interfaces)
- [ ] Infrastructure (MongoDB client)
- [ ] Repository (implementação CRUD)
- [ ] UseCase (lógica de negócio)
- [ ] Handler (HTTP endpoints)
- [ ] Main (composição)
- [ ] Dockerfile
- [ ] docker-compose.yml
- [ ] Testes manuais
- [ ] Documentação

---

**Boa aula! 🚀**

