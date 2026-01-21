# Exemplos de Uso do 00cli

## Inicialização de um Projeto

```bash
# Navegue até o diretório do seu projeto
cd /caminho/do/seu/projeto

# Inicialize a estrutura 00cli
00cli init
```

Isso criará:
- `./00cli/settings.json`
- `./00cli/deploy.json`

## Configuração do settings.json

Edite `./00cli/settings.json` com as informações do seu servidor:

```json
{
  "server": {
    "host": "meuservidor.com",
    "port": 22,
    "user": "deploy",
    "ssh_key": "/home/usuario/.ssh/id_rsa"
  },
  "current_version": "v1.0.0",
  "project_name": "meu-projeto"
}
```

## Configuração do deploy.json

Edite `./00cli/deploy.json` com os comandos de deploy:

```json
{
  "type": "ssh",
  "commands": [
    "cd /var/www/meu-projeto",
    "git pull origin main",
    "docker-compose down",
    "docker-compose up -d --build",
    "docker-compose exec app npm run migrate"
  ],
  "environment": {
    "NODE_ENV": "production",
    "DATABASE_URL": "postgresql://..."
  },
  "provision": {
    "path": "./provision",
    "files": [
      "docker-compose.yml",
      "nginx.conf"
    ]
  }
}
```

## Verificar Status

```bash
00cli status
```

## Fazer Deploy

```bash
# Deploy normal
00cli deploy

# Deploy com modo verboso
00cli deploy --verbose
```

## Trabalhando com o Diretório /provision/

O diretório `/provision/` é usado para armazenar arquivos de configuração e scripts de provisionamento:

```
projeto/
├── provision/
│   ├── docker-compose.yml
│   ├── nginx.conf
│   ├── deploy.sh
│   └── ...
├── 00cli/
│   ├── settings.json
│   └── deploy.json
└── ...
```

O `00cli` verifica automaticamente se o diretório `/provision/` existe e pode referenciar arquivos dele no `deploy.json`.

## Verificação de Atualizações

O `00cli` verifica automaticamente por atualizações toda vez que é executado. Se houver uma nova versão disponível no GitHub, você verá uma mensagem como:

```
⚠️  Nova versão disponível: v0.2.0 (atual: v0.1.0)
   Baixe em: https://github.com/tstest3213/00cli/releases/tag/v0.2.0
```
