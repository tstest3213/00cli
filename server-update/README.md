# Server Update - 00cli

Servidor de atualizações automáticas para o 00cli.

## Instalação

```bash
cd server-update
pip install -r requirements.txt
```

## Uso

### Iniciar o servidor

```bash
# Porta padrão 8080
python api.py

# Porta customizada
python api.py --port 9000

# Modo debug
python api.py --debug
```

### Compilar binários (sem iniciar servidor)

```bash
python api.py --build
```

## Endpoints

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/` | Info do servidor |
| GET | `/latest` | Retorna info da última versão |
| GET | `/download/<binary>` | Download do binário |
| GET | `/health` | Health check |
| POST | `/build` | Compila todos os binários |
| POST | `/build/<platform>` | Compila plataforma específica |

## Plataformas Suportadas

- `linux-amd64`
- `linux-arm64`
- `darwin-amd64`
- `darwin-arm64`
- `windows-amd64`

## Integração com Make

O Makefile já está configurado para chamar o servidor automaticamente após o build.

```bash
# Compilar e atualizar servidor
make build-server

# Ou manualmente
curl -X POST http://localhost:8080/build
```

## Configuração no 00cli

Para usar este servidor de atualizações, configure no `settings.json`:

```json
{
  "update_server": "http://seu-servidor:8080"
}
```

Ou via variável de ambiente:

```bash
export 00CLI_UPDATE_SERVER="http://seu-servidor:8080"
```

## Estrutura de Arquivos

```
server-update/
├── api.py              # Servidor Flask
├── requirements.txt    # Dependências Python
├── metadata.json       # Info da release (gerado automaticamente)
└── binaries/           # Binários compilados
    ├── 00cli-linux-amd64
    ├── 00cli-linux-arm64
    ├── 00cli-darwin-amd64
    ├── 00cli-darwin-arm64
    └── 00cli-windows-amd64.exe
```

## Executar como Serviço (Systemd)

Crie `/etc/systemd/system/00cli-update.service`:

```ini
[Unit]
Description=00cli Update Server
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/00cli-go/server-update
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
