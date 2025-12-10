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
- Docker ou Podman instalado (para rodar com containers)
- Go 1.22+ instalado (para rodar localmente)
- MongoDB instalado (se rodar localmente sem Docker)

### Opção 1: Usando Docker/Podman (Recomendado)

**Com Docker:**
```bash
# Na raiz do projeto
docker-compose up --build
```

**Com Podman:**
```bash
# Na raiz do projeto
podman compose up --build
```

**Verificar se está funcionando:**
```bash
curl http://localhost:8080/healthz
```

### Opção 2: Rodar Localmente (sem Docker)

1. **Inicie o MongoDB localmente** (se não tiver, use Docker apenas para o MongoDB):
   ```bash
   # Com Docker (apenas MongoDB)
   docker run -d -p 27017:27017 --name mongo mongo:7.0
   ```

2. **Configure as variáveis de ambiente:**
   ```bash
   # Linux/Mac
   export MONGO_URI="mongodb://localhost:27017"
   export PORT="8080"
   
   # Windows (PowerShell)
   $env:MONGO_URI="mongodb://localhost:27017"
   $env:PORT="8080"
   
   # Windows (CMD)
   set MONGO_URI=mongodb://localhost:27017
   set PORT=8080
   ```

3. **Execute a aplicação:**
   ```bash
   # Na raiz do projeto
   go run cmd/api/main.go
   ```

**Verificar se está funcionando:**
```bash
curl http://localhost:8080/healthz
```

A API estará disponível em `http://localhost:8080`

## Documentação Swagger (UI Interativa)

A API possui documentação interativa usando Swagger UI, que permite testar todos os endpoints diretamente no navegador.

### Instalação do Swag

Instale a ferramenta `swag` que gera a documentação:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

**Se o comando não for encontrado, adicione ao PATH:**
```bash
# Linux/Mac
export PATH=$PATH:$(go env GOPATH)/bin

# Windows (PowerShell)
$env:Path += ";$(go env GOPATH)\bin"
```

### Gerar a Documentação

**Na raiz do projeto, execute:**
```bash
swag init -g cmd/api/main.go -o docs
```

**O que este comando faz:**
- `-g cmd/api/main.go`: Especifica o arquivo main.go com as anotações Swagger
- `-o docs`: Gera os arquivos na pasta `docs/` na raiz do projeto
- Cria: `docs/docs.go`, `docs/swagger.json`, `docs/swagger.yaml`

### Acessar a UI do Swagger

1. **Inicie a aplicação** (Docker ou localmente)

2. **Acesse no navegador:**
   ```
   http://localhost:8080/swagger/index.html
   ```

3. **Na interface você pode:**
   - Ver todos os endpoints disponíveis
   - Ver exemplos de requisições e respostas
   - Testar os endpoints diretamente no navegador
   - Ver os modelos de dados (User)

### Atualizar a Documentação

Após modificar endpoints ou adicionar novos, regenere:
```bash
# Na raiz do projeto
swag init -g cmd/api/main.go -o docs
```

E reinicie a aplicação.

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

## Resolução de Problemas

### Erro: "package user-api is not in GOROOT" (Windows)

**Problema:** Ao tentar compilar ou executar, aparece erro dizendo que o package `user-api` não foi encontrado.

**Causa:** No Windows, o Go precisa que você esteja dentro do diretório do módulo ou que o `go.mod` esteja configurado corretamente.

**Solução:**
1. Certifique-se de estar na raiz do projeto (onde está o arquivo `go.mod`)
2. Execute `go mod tidy` para sincronizar as dependências
3. Se ainda não funcionar, verifique se o nome do módulo no `go.mod` está correto:
   ```bash
   cat go.mod
   ```
   Deve mostrar `module user-api` na primeira linha
4. Tente compilar novamente:
   ```bash
   go build ./cmd/api
   ```

### Erro: "docker-credential-desktop: executable file not found in PATH" (Windows)

**Problema:** Ao executar `docker-compose up --build`, aparece o erro:
```
error getting credentials - err: exec: "docker-credential-desktop": executable file not found in %PATH%
```

**Causa:** O Docker Desktop não está configurado corretamente ou não está rodando.

**Soluções:**

**Opção 1: Iniciar o Docker Desktop**
1. Abra o Docker Desktop no Windows
2. Aguarde até que ele esteja totalmente iniciado (ícone na bandeja do sistema)
3. Tente novamente o comando `docker-compose up --build`

**Opção 2: Remover credenciais do Docker**
Se o Docker Desktop estiver rodando e ainda assim der erro, tente remover a configuração de credenciais:

1. Abra o arquivo `~/.docker/config.json` (ou `%USERPROFILE%\.docker\config.json` no Windows)
2. Remova ou comente a linha que contém `"credsStore": "desktop"` ou `"credHelpers"`
3. Salve o arquivo e tente novamente

**Opção 3: Usar Docker sem credenciais**
Configure o Docker para não usar credenciais:
```bash
# No PowerShell ou CMD
docker config --help
```

### Aviso: "the attribute `version` is obsolete" no docker-compose.yml

**Problema:** Ao executar `docker-compose up`, aparece um aviso:
```
the attribute `version` is obsolete, it will be ignored
```

**Causa:** A partir do Docker Compose v2, o campo `version` não é mais necessário.

**Solução:** Remova a primeira linha `version: '3.8'` do arquivo `docker-compose.yml`. O arquivo funcionará normalmente sem ela.

### Problemas comuns ao rodar localmente (sem Docker)

Se você quiser rodar sem Docker, precisa ter o MongoDB instalado localmente:

1. Instale o MongoDB no seu sistema
2. Inicie o serviço do MongoDB
3. Execute a aplicação:
   ```bash
   go run cmd/api/main.go
   ```
4. Ou defina as variáveis de ambiente:
   ```bash
   # Windows (PowerShell)
   $env:MONGO_URI="mongodb://localhost:27017"
   $env:PORT="8080"
   go run cmd/api/main.go
   
   # Windows (CMD)
   set MONGO_URI=mongodb://localhost:27017
   set PORT=8080
   go run cmd/api/main.go
   
   # Linux/Mac
   export MONGO_URI="mongodb://localhost:27017"
   export PORT="8080"
   go run cmd/api/main.go
   ```

