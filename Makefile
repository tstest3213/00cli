.PHONY: build install uninstall clean test build-server build-all update-server

# Nome do binário
BINARY_NAME=00cli
INSTALL_PATH=/usr/local/bin

# Versão
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "v0.1.0")

# Flags de build
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# Servidor de atualizações
UPDATE_SERVER_URL ?= http://localhost:8080
UPDATE_SERVER_DIR = server-update

build:
	@echo "🔨 Compilando $(BINARY_NAME)..."
	@go build $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "✅ Build concluído: ./$(BINARY_NAME)"

install: build
	@echo "📦 Instalando $(BINARY_NAME) em $(INSTALL_PATH)..."
	@sudo cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@sudo chmod +x $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✅ Instalado com sucesso! Execute '$(BINARY_NAME)' de qualquer diretório."

uninstall:
	@echo "🗑️  Removendo $(BINARY_NAME)..."
	@sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✅ Removido com sucesso."

clean:
	@echo "🧹 Limpando..."
	@rm -f $(BINARY_NAME)
	@go clean
	@echo "✅ Limpeza concluída."

test:
	@echo "🧪 Executando testes..."
	@go test -v ./...

run:
	@go run .

# Instalação local para desenvolvimento
dev-install: build
	@cp $(BINARY_NAME) $(HOME)/.local/bin/$(BINARY_NAME) 2>/dev/null || mkdir -p $(HOME)/.local/bin && cp $(BINARY_NAME) $(HOME)/.local/bin/$(BINARY_NAME)
	@chmod +x $(HOME)/.local/bin/$(BINARY_NAME)
	@echo "✅ Instalado em $(HOME)/.local/bin/$(BINARY_NAME)"
	@echo "   Certifique-se de que $(HOME)/.local/bin está no seu PATH"

# ============ Servidor de Atualizações ============

# Compilar todos os binários para todas as plataformas
build-all:
	@echo "🚀 Compilando para todas as plataformas..."
	@cd $(UPDATE_SERVER_DIR) && python3 api.py --build
	@echo "✅ Todos os binários compilados em $(UPDATE_SERVER_DIR)/binaries/"

# Compilar e notificar servidor de atualizações (se estiver rodando)
build-server: build-all
	@echo "📡 Notificando servidor de atualizações..."
	@curl -s -X POST $(UPDATE_SERVER_URL)/build > /dev/null 2>&1 || echo "⚠️  Servidor não está rodando (opcional)"
	@echo "✅ Build completo!"

# Atualizar servidor de atualizações (POST /build)
update-server:
	@echo "📡 Atualizando servidor de atualizações..."
	@curl -s -X POST $(UPDATE_SERVER_URL)/build | python3 -m json.tool 2>/dev/null || echo "❌ Servidor não está rodando"

# Iniciar servidor de atualizações
server:
	@echo "🚀 Iniciando servidor de atualizações..."
	@cd $(UPDATE_SERVER_DIR) && python3 api.py

# Instalar dependências do servidor
server-deps:
	@echo "📦 Instalando dependências do servidor..."
	@pip install -r $(UPDATE_SERVER_DIR)/requirements.txt
