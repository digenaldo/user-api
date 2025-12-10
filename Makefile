
## Makefile com targets úteis para desenvolvimento
## Comentários (pt-BR) explicam cada alvo e opções para Podman/Docker.
.PHONY: up down logs

## Sobe o projeto (compila e inicia os containers).
## - Usa `podman compose up --build` por ser compatível com Podman.
## - Para rodar em background: adicione `-d` ao comando no terminal.
up: ## Sobe o projeto com Docker Compose
	podman compose up --build

## Para o projeto (remove containers definidos no compose)
## Nota: aqui está usando `docker-compose down` porque o repositório
## pode ter sido escrito originalmente para Docker. No seu sistema com
## Podman você pode executar `podman compose down` em vez disso.
down: ## Para o projeto
	docker-compose down

## Mostra logs dos containers em tempo real
## Recomendação: com Podman use `podman compose logs -f`.
logs: ## Mostra logs dos containers
	docker-compose logs -f

## Observações:
## - Mantenha consistência com o runtime que você usa (Podman vs Docker).
## - Esses alvos são convenientes para desenvolvimento, mas em produção
##   prefira scripts de deploy/CI que controlem lifecycle de forma segura.
