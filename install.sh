#!/bin/bash

# Script de instalação do 00cli
# Uso: ./install.sh

set -e

BINARY_NAME="00cli"
INSTALL_PATH="/usr/local/bin"
REPO_URL="https://github.com/tstest3213/00cli.git"
TEMP_DIR=$(mktemp -d)

echo "🚀 Instalando 00cli..."

# Verificar se Go está instalado
if ! command -v go &> /dev/null; then
    echo "❌ Go não está instalado. Por favor, instale o Go primeiro."
    echo "   Visite: https://golang.org/dl/"
    exit 1
fi

# Clonar repositório (se não estiver no diretório atual)
if [ ! -f "go.mod" ] || [ ! -f "main.go" ]; then
    echo "📦 Clonando repositório..."
    cd "$TEMP_DIR"
    git clone "$REPO_URL" .
else
    echo "📦 Usando repositório local..."
    TEMP_DIR=$(pwd)
fi

# Compilar
echo "🔨 Compilando $BINARY_NAME..."
cd "$TEMP_DIR"
go build -o "$BINARY_NAME" .

# Instalar
echo "📦 Instalando em $INSTALL_PATH..."
sudo cp "$BINARY_NAME" "$INSTALL_PATH/$BINARY_NAME"
sudo chmod +x "$INSTALL_PATH/$BINARY_NAME"

# Limpar
if [ "$TEMP_DIR" != "$(pwd)" ]; then
    rm -rf "$TEMP_DIR"
fi

echo "✅ 00cli instalado com sucesso!"
echo ""
echo "Execute '00cli --help' para ver os comandos disponíveis."
echo "Execute '00cli init' em um projeto para inicializar a estrutura."
