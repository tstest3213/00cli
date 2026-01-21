# Exemplos de Uso do 00cli

## Inicialização de um Projeto

```bash
# Navegue até o diretório do seu projeto
cd /caminho/do/seu/projeto

# Inicialize a estrutura 00cli
00cli init
```

Isso criará:
- `./.00cli/settings.json`
- `./.00cli/deploy.json`

## Configuração do settings.json

Edite `./.00cli/settings.json` com as informações do seu servidor:

### Exemplo Básico (com chave SSH):

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

### Exemplo com Senha (menos seguro):

```json
{
  "server": {
    "host": "meuservidor.com",
    "port": 22,
    "user": "deploy",
    "password": "sua_senha_aqui"
  },
  "current_version": "v1.0.0",
  "project_name": "meu-projeto"
}
```

**⚠️ Nota de Segurança:** Prefira usar `ssh_key` ao invés de `password`. Se usar senha, certifique-se de que o arquivo `settings.json` não seja commitado no Git (adicione `.00cli/settings.json` ao `.gitignore`).

## Configuração do deploy.json

Edite `./.00cli/deploy.json` com os comandos de deploy:

```json
{
  "type": "ssh",
  "commands": [
    "cd /var/www/meu-projeto",
    "git pull origin main",
    "npm install",
    "npm run build",
    "pm2 restart app"
  ],
  "environment": {
    "NODE_ENV": "production",
    "DATABASE_URL": "postgresql://..."
  },
  "provision": {
    "path": "./provision",
    "files": [
      "nginx.conf",
      "app.conf"
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
seu-projeto/
├── .00cli/
│   ├── settings.json
│   └── deploy.json
├── provision/
│   ├── nginx.conf       # Config do nginx
│   ├── app.conf         # Config do PM2/supervisor
│   └── .env.production  # Variáveis de ambiente
└── ...
```

### Exemplo de nginx.conf

```nginx
server {
    listen 80;
    server_name meusite.com;

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

## Cenários Comuns

### Deploy de Aplicação Node.js

```json
{
  "type": "ssh",
  "commands": [
    "cd /var/www/minha-app",
    "git pull origin main",
    "npm ci --production",
    "npm run build",
    "pm2 restart ecosystem.config.js --env production"
  ]
}
```

### Deploy de Aplicação Python/Django

```json
{
  "type": "ssh",
  "commands": [
    "cd /var/www/django-app",
    "git pull origin main",
    "source venv/bin/activate",
    "pip install -r requirements.txt",
    "python manage.py migrate --noinput",
    "python manage.py collectstatic --noinput",
    "sudo systemctl restart gunicorn"
  ]
}
```

### Deploy com Docker

```json
{
  "type": "docker",
  "commands": [
    "cd /var/www/minha-app",
    "git pull origin main",
    "docker-compose pull",
    "docker-compose up -d --build"
  ]
}
```

### Deploy Simples (só Git Pull)

```json
{
  "type": "git",
  "commands": [
    "cd /var/www/site-estatico",
    "git pull origin main"
  ]
}
```

## Workflows Avançados

### Deploy com Backup Prévio

```json
{
  "type": "ssh",
  "commands": [
    "cd /var/www/minha-app",
    "tar -czf backup-$(date +%Y%m%d-%H%M%S).tar.gz .",
    "git pull origin main",
    "npm ci --production",
    "npm run build",
    "pm2 restart app"
  ]
}
```

### Deploy com Testes

```json
{
  "type": "ssh",
  "commands": [
    "cd /var/www/minha-app",
    "git pull origin main",
    "npm ci",
    "npm test",
    "npm run build",
    "pm2 restart app"
  ]
}
```

## Dicas

1. **Sempre use chave SSH** em vez de senha para maior segurança
2. **Adicione `.00cli/settings.json` ao `.gitignore`** se contiver informações sensíveis
3. **Use `--verbose`** para debugar problemas de deploy
4. **Teste comandos manualmente** antes de adicioná-los ao deploy.json
