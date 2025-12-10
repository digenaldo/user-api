# ============================================
# MAKEFILE - AUTOMAÇÃO DE TAREFAS
# ============================================
# Makefile é um arquivo que define "targets" (alvos) e comandos para executá-los
# Facilita executar comandos longos com nomes curtos
#
# COMO USAR:
#   make up      # Executa o target "up"
#   make down    # Executa o target "down"
#   make logs    # Executa o target "logs"
#
# VANTAGENS:
# - Não precisa lembrar comandos longos
# - Padroniza comandos entre desenvolvedores
# - Pode ter dependências entre targets
#
# SOBRE .PHONY:
# - Indica que esses targets não criam arquivos com esses nomes
# - Sem .PHONY, o Make poderia pensar que "up" é um arquivo e não executar
# - É uma boa prática para targets que executam comandos
.PHONY: up down logs

# ============================================
# TARGET: UP
# ============================================
# Sobe o projeto (compila e inicia os containers)
#
# O QUE FAZ:
# 1. Compila a aplicação Go (--build)
# 2. Cria os containers (MongoDB e API)
# 3. Inicia os serviços
# 4. Conecta os containers na mesma rede
#
# SOBRE podman compose:
# - Podman é uma alternativa ao Docker (mesma sintaxe, sem daemon)
# - Se você usa Docker, pode trocar por "docker compose"
# - O comando "compose" é o novo padrão (substitui "docker-compose")
#
# SOBRE --build:
# - Força a recompilação da imagem mesmo se já existir
# - Sem --build, usa a imagem em cache (mais rápido, mas pode estar desatualizada)
#
# PARA RODAR EM BACKGROUND:
# - Adicione -d ao comando: podman compose up --build -d
# - Os containers rodam em background (detached mode)
# - Útil para desenvolvimento contínuo
up: ## Sobe o projeto com Docker Compose
	podman compose up --build

# ============================================
# TARGET: DOWN
# ============================================
# Para o projeto (remove containers definidos no compose)
#
# O QUE FAZ:
# 1. Para todos os containers em execução
# 2. Remove os containers
# 3. Remove a rede criada pelo compose
# 4. NÃO remove volumes (dados do MongoDB são preservados)
#
# SOBRE docker-compose down:
# - Aqui usamos "docker-compose" (com hífen) por compatibilidade
# - Se usar Podman, pode trocar por "podman compose down"
# - O comportamento é o mesmo
#
# PARA REMOVER VOLUMES TAMBÉM:
# - Adicione -v: docker-compose down -v
# - Isso APAGA os dados do MongoDB!
# - Use apenas se quiser começar do zero
down: ## Para o projeto
	docker-compose down

# ============================================
# TARGET: LOGS
# ============================================
# Mostra logs dos containers em tempo real
#
# O QUE FAZ:
# - Exibe os logs de todos os serviços (mongo e api)
# - Atualiza em tempo real (follow mode)
# - Útil para debugar problemas
#
# SOBRE -f (follow):
# - Mantém os logs abertos e mostra novas linhas conforme aparecem
# - Similar ao "tail -f" do Linux
# - Para sair, pressione Ctrl+C
#
# PARA VER LOGS DE UM SERVIÇO ESPECÍFICO:
# - docker-compose logs -f api    # Só logs da API
# - docker-compose logs -f mongo   # Só logs do MongoDB
#
# PARA VER ÚLTIMAS 100 LINHAS:
# - docker-compose logs --tail=100
logs: ## Mostra logs dos containers
	docker-compose logs -f

# ============================================
# OBSERVAÇÕES IMPORTANTES
# ============================================
# - Mantenha consistência: se usar Podman, use em todos os targets
# - Esses targets são para desenvolvimento local
# - Em produção, use scripts de deploy/CI mais robustos
# - O Makefile facilita, mas você pode executar os comandos diretamente também
