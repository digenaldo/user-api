# Primeira etapa: compilar a aplicação
# Usamos uma imagem com Go para ter acesso ao compilador
FROM golang:1.22 AS builder

WORKDIR /app

# Copia apenas os arquivos de dependências primeiro
# Isso permite que o Docker reutilize o cache se as dependências não mudarem
COPY go.mod go.sum ./
RUN go mod download

# Agora copia todo o código e compila
COPY . .
# CGO_ENABLED=0 cria um binário estático (não precisa de libs C do sistema)
# GOOS=linux garante que compila para Linux mesmo se você estiver no Mac/Windows
RUN CGO_ENABLED=0 GOOS=linux go build -o user-api ./cmd/api

# Segunda etapa: criar a imagem final (bem menor)
# Alpine é uma distro Linux minimalista, perfeita para containers
FROM alpine:3.19

WORKDIR /app

# Copia apenas o binário compilado da etapa anterior
# Não precisamos do código fonte nem do Go na imagem final
COPY --from=builder /app/user-api /app/user-api

# Informa que a aplicação usa a porta 8080
EXPOSE 8080

# Comando que será executado quando o container iniciar
CMD ["./user-api"]
