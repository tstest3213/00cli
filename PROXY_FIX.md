# Solução para Problema de Proxy no Pre-commit

## Problema

O pre-commit está falhando porque as variáveis de ambiente de proxy estão configuradas (`http://127.0.0.1:9001/`), mas o proxy não está rodando.

## Soluções

### Solução 1: Desabilitar Proxy Temporariamente

Para executar pre-commit sem proxy:

```bash
unset http_proxy https_proxy HTTP_PROXY HTTPS_PROXY all_proxy ALL_PROXY
pre-commit run --all-files
```

Ou use o script fornecido:

```bash
./fix-pre-commit.sh
```

### Solução 2: Configurar Pre-commit para Ignorar Proxy

Adicione ao seu `.pre-commit-config.yaml`:

```yaml
default_language_version:
  python: system

repos:
  - repo: https://github.com/pre-commit/mirrors-mypy
    hooks:
      - id: mypy
        additional_dependencies: []
        # Desabilitar proxy para este hook
        env:
          http_proxy: ""
          https_proxy: ""
```

### Solução 3: Iniciar o Proxy (se necessário)

Se você realmente precisa do proxy, inicie-o antes:

```bash
# Exemplo com um proxy local
# (ajuste conforme seu setup)
```

### Solução 4: Configurar Proxy Corretamente

Se o proxy deve estar ativo, verifique:

1. O proxy está rodando?
   ```bash
   netstat -tuln | grep 9001
   ```

2. As variáveis estão corretas?
   ```bash
   echo $http_proxy
   ```

## Para Commits Futuros

Você pode criar um alias no Git:

```bash
git config --global alias.commit-no-proxy '!f(){ unset http_proxy https_proxy HTTP_PROXY HTTPS_PROXY all_proxy ALL_PROXY; git commit "$@"; }; f'
```

Então use:
```bash
git commit-no-proxy -m "sua mensagem"
```
