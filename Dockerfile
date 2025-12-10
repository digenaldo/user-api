# ============================================
# DOCKERFILE - MULTI-STAGE BUILD
# ============================================
# Dockerfile define como construir uma imagem Docker
# Multi-stage build: usa duas etapas para criar uma imagem final menor
#
# POR QUE MULTI-STAGE?
# - Etapa 1 (builder): tem Go e compila o código (imagem grande ~800MB)
# - Etapa 2 (runtime): só tem o binário compilado (imagem pequena ~20MB)
# - Resultado: imagem final 40x menor!
#
# VANTAGENS:
# - Imagem menor = download mais rápido, menos espaço, mais seguro
# - Imagem final não tem ferramentas de desenvolvimento (Go, compilador)
# - Reduz superfície de ataque (menos código = menos vulnerabilidades)

# ============================================
# ETAPA 1: BUILD (COMPILAÇÃO)
# ============================================
# Esta etapa compila a aplicação Go
# Usamos uma imagem com Go instalado

# FROM define a imagem base
# golang:1.22 é a imagem oficial do Go versão 1.22
# AS builder dá um nome a esta etapa (podemos referenciar depois)
FROM golang:1.22 AS builder

# WORKDIR define o diretório de trabalho dentro do container
# Todos os comandos seguintes executam neste diretório
# Se não existir, o Docker cria automaticamente
WORKDIR /app

# ============================================
# OTIMIZAÇÃO: CACHE DE DEPENDÊNCIAS
# ============================================
# Copia APENAS os arquivos de dependências primeiro
# Isso permite que o Docker reutilize o cache se as dependências não mudarem
#
# POR QUE FAZER ISSO?
# - go.mod e go.sum raramente mudam
# - Se não mudaram, Docker usa cache e pula go mod download (economiza tempo)
# - Se mudaram, baixa novas dependências
#
# ESTRATÉGIA DE CACHE:
# 1. Copiar arquivos que mudam pouco primeiro (go.mod, go.sum)
# 2. Baixar dependências (pode ser cacheado)
# 3. Copiar código fonte (muda sempre)
# 4. Compilar (precisa refazer se código mudou)
COPY go.mod go.sum ./

# Baixa todas as dependências listadas em go.mod
# go mod download baixa os pacotes para o cache local
# Se go.mod/go.sum não mudaram, Docker usa cache desta etapa
RUN go mod download

# ============================================
# COMPILAÇÃO DA APLICAÇÃO
# ============================================
# Agora copia todo o código fonte
# Como o código muda sempre, esta etapa raramente usa cache
COPY . .

# Compila a aplicação
# 
# SOBRE AS VARIÁVEIS DE AMBIENTE:
# - CGO_ENABLED=0: desabilita CGO (C bindings)
#   * CGO permite chamar código C do Go
#   * Desabilitado = binário estático (não precisa de libs C do sistema)
#   * Mais portável, funciona em qualquer Linux
#
# - GOOS=linux: define o sistema operacional alvo
#   * Mesmo compilando no Mac/Windows, gera binário para Linux
#   * Containers Docker rodam Linux (mesmo no Mac/Windows via VM)
#
# - go build -o user-api ./cmd/api:
#   * -o user-api: nome do binário gerado
#   * ./cmd/api: diretório com o main.go
#
# RESULTADO: binário "user-api" pronto para Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o user-api ./cmd/api

# ============================================
# ETAPA 2: RUNTIME (IMAGEM FINAL)
# ============================================
# Esta etapa cria a imagem final (só com o binário)
# Não precisa do Go, compilador, ou código fonte

# Alpine Linux é uma distro minimalista
# - Muito pequena (~5MB base)
# - Segura (poucos pacotes = menos vulnerabilidades)
# - Perfeita para containers (só o essencial)
FROM alpine:3.19

# Define o diretório de trabalho
WORKDIR /app

# ============================================
# COPIA O BINÁRIO DA ETAPA ANTERIOR
# ============================================
# COPY --from=builder copia arquivos da etapa "builder"
# /app/user-api é o binário compilado na etapa 1
# /app/user-api é o destino na etapa 2
#
# POR QUE SÓ O BINÁRIO?
# - Não precisamos do código fonte (já está compilado)
# - Não precisamos do Go (só precisava para compilar)
# - Não precisamos das dependências (já linkadas no binário)
# - Resultado: imagem final muito menor!
COPY --from=builder /app/user-api /app/user-api

# ============================================
# CONFIGURAÇÃO DO CONTAINER
# ============================================
# EXPOSE informa que a aplicação usa a porta 8080
# Isso é apenas documentação - não abre a porta automaticamente
# A porta é aberta no docker-compose.yml ou docker run -p
EXPOSE 8080

# CMD define o comando padrão quando o container iniciar
# Quando você faz "docker run", este comando é executado
# ["./user-api"] executa o binário user-api no diretório atual
#
# SOBRE A SINTAXE:
# - Forma de lista ["cmd", "arg1", "arg2"] é preferida
# - Forma de string "cmd arg1 arg2" também funciona, mas menos segura
# - A forma de lista não passa pelo shell (mais seguro, mais rápido)
CMD ["./user-api"]
