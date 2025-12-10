
# Imagem de build (multi-stage): usa a imagem oficial do Go
# - Compila o binário dentro de uma imagem completa (com toolchain)
# - Mantém o artefato final pequeno copiando só o binário para a imagem final
FROM golang:1.22 AS builder

WORKDIR /app

## Copiamos apenas os arquivos de mod primeiro para aproveitar cache de layer
## e executar `go mod download` só quando as dependências mudarem.
COPY go.mod go.sum ./
RUN go mod download

## Copia todo o código e compila o binário para Linux.
## CGO_ENABLED=0 -> gera binário estático (sem dependências C), facilitando
## execução em imagens base mínimas como Alpine.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o user-api ./cmd/api

## Imagem final, pequena, contendo apenas o binário e o runtime necessário.
## Alpine é uma escolha comum para imagens pequenas; alternativa: distroless.
FROM alpine:3.19

WORKDIR /app
## Copia o binário construído na etapa anterior
COPY --from=builder /app/user-api /app/user-api

## Porta que a aplicação usa internamente (documentativa). O mapping é feito
## no docker-compose / podman-compose ou ao executar `podman run -p`.
EXPOSE 8080

## Comando padrão para iniciar a aplicação
CMD ["./user-api"]

## - Usamos multi-stage para manter a imagem final enxuta e segura.
## - `CGO_ENABLED=0` evita dependências C, útil para compatibilidade em contêineres.
## - Se precisar de depuração local, execute `go build` localmente e rode sem Docker.

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o user-api ./cmd/api

FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/user-api /app/user-api

EXPOSE 8080
CMD ["./user-api"]

