# Servidor de Atualizações Customizado

O `00cli` suporta usar um servidor web próprio para verificar e baixar atualizações, ao invés do GitHub.

## Configuração

### Opção 1: Via settings.json (Recomendado)

Edite `./.00cli/settings.json` e adicione o campo `update_server`:

```json
{
  "server": {
    "host": "meuservidor.com",
    "port": 22,
    "user": "deploy"
  },
  "current_version": "v0.1.0",
  "project_name": "meu-projeto",
  "update_server": "http://192.168.1.100:8080/updates"
}
```

### Opção 2: Via Variável de Ambiente

```bash
export 00CLI_UPDATE_SERVER="http://192.168.1.100:8080/updates"
00cli update
```

## Formato da API do Servidor

O servidor precisa implementar um endpoint que retorne JSON no seguinte formato:

### Endpoint: `GET /latest` ou `GET /updates/latest`

**Resposta esperada:**

```json
{
  "tag_name": "v0.2.0",
  "name": "Release v0.2.0",
  "published_at": "2024-01-15T10:00:00Z",
  "html_url": "http://192.168.1.100:8080/releases/v0.2.0",
  "body": "Descrição da release",
  "assets": [
    {
      "name": "00cli-linux-amd64",
      "browser_download_url": "http://192.168.1.100:8080/download/00cli-linux-amd64",
      "size": 12345678
    },
    {
      "name": "00cli-linux-arm64",
      "browser_download_url": "http://192.168.1.100:8080/download/00cli-linux-arm64",
      "size": 12345678
    }
  ]
}
```

### Campos Obrigatórios

- `tag_name`: Versão da release (ex: "v0.2.0")
- `assets`: Array com os binários disponíveis
  - `name`: Nome do binário (ex: "00cli-linux-amd64")
  - `browser_download_url`: URL para download direto

### Campos Opcionais

- `name`: Nome da release
- `published_at`: Data de publicação
- `html_url`: URL da página da release
- `body`: Descrição da release

## Usando o Servidor Incluído

O 00cli inclui um servidor de atualizações pronto em `server-update/api.py`.

### Instalação

```bash
cd server-update
pip install -r requirements.txt
```

### Iniciar Servidor

```bash
# Via Makefile (recomendado)
make server

# Ou diretamente
python server-update/api.py --port 8080
```

### Compilar e Publicar

```bash
# Compilar para todas as plataformas e atualizar servidor
make build-server

# Ou apenas compilar
make build-all
```

### Endpoints do Servidor

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/` | Info do servidor |
| GET | `/latest` | Info da última versão |
| GET | `/download/<binary>` | Download de binário |
| GET | `/health` | Health check |
| POST | `/build` | Compilar todos os binários |
| POST | `/build/<platform>` | Compilar plataforma específica |

## Estrutura de Diretórios no Servidor

```
/var/www/updates/
├── latest.json          # Endpoint /latest retorna este arquivo
└── download/
    ├── 00cli-linux-amd64
    ├── 00cli-linux-arm64
    ├── 00cli-darwin-amd64
    ├── 00cli-darwin-arm64
    └── 00cli-windows-amd64.exe
```

## Exemplo de Implementação (Nginx)

### Configuração Nginx

```nginx
server {
    listen 8080;
    server_name 192.168.1.100;

    root /var/www/updates;
    index latest.json;

    location /latest {
        alias /var/www/updates/latest.json;
        add_header Content-Type application/json;
    }

    location /download/ {
        alias /var/www/updates/download/;
        add_header Content-Type application/octet-stream;
    }
}
```

### Exemplo de latest.json

```json
{
  "tag_name": "v0.2.0",
  "name": "Release v0.2.0",
  "published_at": "2024-01-15T10:00:00Z",
  "html_url": "http://192.168.1.100:8080/releases/v0.2.0",
  "body": "Nova versão com melhorias",
  "assets": [
    {
      "name": "00cli-linux-amd64",
      "browser_download_url": "http://192.168.1.100:8080/download/00cli-linux-amd64",
      "size": 12345678
    },
    {
      "name": "00cli-linux-arm64",
      "browser_download_url": "http://192.168.1.100:8080/download/00cli-linux-arm64",
      "size": 12345678
    }
  ]
}
```

## Fallback para GitHub

Se o servidor customizado não estiver disponível ou retornar erro, o `00cli` automaticamente tentará usar o GitHub como fallback.

## Testando

```bash
# Configurar servidor customizado
export 00CLI_UPDATE_SERVER="http://192.168.1.100:8080"

# Verificar atualizações
00cli update

# Ou verificar automaticamente (em background)
00cli version
```

## Segurança

- Use HTTPS se possível (recomendado para produção)
- Valide o tamanho dos arquivos baixados
- Considere adicionar autenticação se necessário
- Mantenha os binários atualizados e verificados

## Nomes de Binários Suportados

O sistema detecta automaticamente a plataforma e procura por:

- `00cli-linux-amd64`
- `00cli-linux-arm64`
- `00cli-darwin-amd64`
- `00cli-darwin-arm64`
- `00cli-windows-amd64.exe`

## Executar como Serviço (Systemd)

Crie `/etc/systemd/system/00cli-update.service`:

```ini
[Unit]
Description=00cli Update Server
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/00cli/server-update
ExecStart=/usr/bin/python3 api.py --port 8080
Restart=always

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable 00cli-update
sudo systemctl start 00cli-update
```
