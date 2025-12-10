.PHONY: up down logs

up: ## Sobe o projeto com Docker Compose
	podman compose up --build

down: ## Para o projeto
	docker-compose down

logs: ## Mostra logs dos containers
	docker-compose logs -f
