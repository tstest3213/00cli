.PHONY: build install uninstall clean test

# Nome do binário
BINARY_NAME=00cli
INSTALL_PATH=/usr/local/bin

# Versão
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "v0.1.0")

# Flags de build
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

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
