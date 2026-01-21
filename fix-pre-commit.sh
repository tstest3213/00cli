#!/bin/bash
# Script para executar pre-commit sem proxy
# Uso: ./fix-pre-commit.sh [comando]

# Desabilitar proxy temporariamente
unset http_proxy https_proxy HTTP_PROXY HTTPS_PROXY all_proxy ALL_PROXY

# Executar comando ou pre-commit por padrão
if [ $# -eq 0 ]; then
    pre-commit run --all-files
else
    "$@"
fi
