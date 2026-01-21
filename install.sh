#!/bin/bash

# Script de instalação do 00cli
# Uso: ./install.sh [--from-source]

set -e

BINARY_NAME="00cli"
INSTALL_PATH="/usr/local/bin"
REPO_OWNER="tstest3213"
REPO_NAME="00cli"
GITHUB_API="https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest"

# Detectar OS e ARCH
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Normalizar ARCH
case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    i386|i686) ARCH="386" ;;
esac

# Normalizar OS
case "$OS" in
    linux) OS="linux" ;;
    darwin) OS="darwin" ;;
    *) echo "❌ OS não suportado: $OS"; exit 1 ;;
esac

BINARY_FILE="${BINARY_NAME}-${OS}-${ARCH}"
if [ "$OS" = "windows" ]; then
    BINARY_FILE="${BINARY_FILE}.exe"
fi

echo "🚀 Instalando 00cli..."
echo "   OS: $OS"
echo "   ARCH: $ARCH"
echo ""

# Verificar se --from-source foi passado
FROM_SOURCE=false
if [ "$1" = "--from-source" ]; then
    FROM_SOURCE=true
fi

if [ "$FROM_SOURCE" = true ]; then
    echo "📦 Instalação a partir do código fonte..."

    # Verificar se Go está instalado
    if ! command -v go &> /dev/null; then
        echo "❌ Go não está instalado. Por favor, instale o Go primeiro."
        echo "   Visite: https://golang.org/dl/"
        exit 1
    fi

    # Verificar se estamos no diretório do projeto
    if [ ! -f "go.mod" ] || [ ! -f "main.go" ]; then
        echo "❌ Execute este script no diretório do projeto 00cli"
        exit 1
    fi

    # Compilar
    echo "🔨 Compilando $BINARY_NAME..."
    go build -o "$BINARY_NAME" .

    # Instalar
    echo "📦 Instalando em $INSTALL_PATH..."
    sudo cp "$BINARY_NAME" "$INSTALL_PATH/$BINARY_NAME"
    sudo chmod +x "$INSTALL_PATH/$BINARY_NAME"

    echo "✅ 00cli instalado com sucesso!"
else
    echo "📦 Baixando última versão do GitHub..."

    # Verificar se curl ou wget está disponível
    if command -v curl &> /dev/null; then
        DOWNLOAD_CMD="curl -L -s"
    elif command -v wget &> /dev/null; then
        DOWNLOAD_CMD="wget -q -O -"
    else
        echo "❌ curl ou wget não encontrado. Instale um deles primeiro."
        exit 1
    fi

    # Obter URL do release
    RELEASE_INFO=$($DOWNLOAD_CMD "$GITHUB_API" 2>/dev/null)

    if [ -z "$RELEASE_INFO" ]; then
        echo "⚠️  Não foi possível obter informações do release."
        echo "   Tentando instalação a partir do código fonte..."
        echo ""
        exec "$0" --from-source
        exit $?
    fi

    # Extrair tag e download URL
    TAG=$(echo "$RELEASE_INFO" | grep -o '"tag_name": "[^"]*' | cut -d'"' -f4)
    DOWNLOAD_URL=$(echo "$RELEASE_INFO" | grep -o "\"browser_download_url\": \"[^\"]*${BINARY_FILE}[^\"]*" | cut -d'"' -f4)

    if [ -z "$DOWNLOAD_URL" ]; then
        echo "⚠️  Binário pré-compilado não encontrado para $OS/$ARCH"
        echo "   Tentando instalação a partir do código fonte..."
        echo ""
        exec "$0" --from-source
        exit $?
    fi

    echo "   Versão: $TAG"
    echo "   Baixando: $DOWNLOAD_URL"

    # Criar diretório temporário
    TEMP_DIR=$(mktemp -d)
    TEMP_FILE="$TEMP_DIR/$BINARY_NAME"

    # Baixar binário
    if command -v curl &> /dev/null; then
        curl -L -o "$TEMP_FILE" "$DOWNLOAD_URL" || {
            echo "❌ Erro ao baixar. Tentando instalação a partir do código fonte..."
            exec "$0" --from-source
            exit $?
        }
    else
        wget -O "$TEMP_FILE" "$DOWNLOAD_URL" || {
            echo "❌ Erro ao baixar. Tentando instalação a partir do código fonte..."
            exec "$0" --from-source
            exit $?
        }
    fi

    # Tornar executável
    chmod +x "$TEMP_FILE"

    # Instalar
    echo "📦 Instalando em $INSTALL_PATH..."
    sudo cp "$TEMP_FILE" "$INSTALL_PATH/$BINARY_NAME"
    sudo chmod +x "$INSTALL_PATH/$BINARY_NAME"

    # Limpar
    rm -rf "$TEMP_DIR"

    echo "✅ 00cli $TAG instalado com sucesso!"
fi

echo ""
echo "Execute '00cli --help' para ver os comandos disponíveis."
echo "Execute '00cli init' em um projeto para inicializar a estrutura."
echo "Execute '00cli update' para atualizar para a versão mais recente."
