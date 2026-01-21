# Configuração do 00cli

## Estrutura de Arquivos

Após executar `00cli init`, a seguinte estrutura será criada:

```
seu-projeto/
├── .00cli/
│   ├── settings.json    # Configurações do servidor
│   └── deploy.json      # Configurações de deploy
└── ...
```

## settings.json

O arquivo `./.00cli/settings.json` contém as configurações de conexão com o servidor.

### Estrutura Completa

```json
{
  "server": {
    "host": "meuservidor.com",
    "port": 22,
    "user": "deploy",
    "ssh_key": "/home/usuario/.ssh/id_rsa",
    "password": ""
  },
  "current_version": "v0.0.0",
  "project_name": "meu-projeto",
  "update_server": "http://192.168.1.100:8080/updates"
}
```

### Campos

#### `server.host` (obrigatório)
- **Tipo**: `string`
- **Descrição**: Endereço IP ou hostname do servidor
- **Exemplo**: `"192.168.1.100"` ou `"meuservidor.com"`

#### `server.port` (opcional)
- **Tipo**: `integer`
- **Descrição**: Porta SSH do servidor
- **Padrão**: `22`

#### `server.user` (obrigatório)
- **Tipo**: `string`
- **Descrição**: Usuário para conexão SSH
- **Exemplo**: `"deploy"`, `"root"`, `"ubuntu"`

#### `server.ssh_key` (recomendado)
- **Tipo**: `string`
- **Descrição**: Caminho completo para a chave SSH privada
- **Exemplo**: `"/home/usuario/.ssh/id_rsa"`
- **Nota**: Se fornecido, será usado ao invés de senha

#### `server.password` (não recomendado)
- **Tipo**: `string`
- **Descrição**: Senha do usuário (menos seguro que chave SSH)
- **Nota**: ⚠️ **Não commite este arquivo no Git se usar senha!**

#### `current_version` (opcional)
- **Tipo**: `string`
- **Descrição**: Versão atual do projeto
- **Exemplo**: `"v1.0.0"`, `"v0.1.0"`

#### `project_name` (opcional)
- **Tipo**: `string`
- **Descrição**: Nome do projeto
- **Padrão**: Nome do diretório onde `00cli init` foi executado

#### `update_server` (opcional)
- **Tipo**: `string`
- **Descrição**: URL do servidor de atualizações customizado
- **Exemplo**: `"http://192.168.1.100:8080/updates"` ou `"https://updates.seudominio.com"`
- **Nota**: Se não configurado, usa GitHub como padrão. Veja [update-server.md](./update-server.md) para mais detalhes.

### Exemplos

#### Exemplo 1: Usando Chave SSH (Recomendado)

```json
{
  "server": {
    "host": "192.168.1.100",
    "port": 22,
    "user": "deploy",
    "ssh_key": "/home/usuario/.ssh/id_rsa"
  },
  "current_version": "v1.0.0",
  "project_name": "meu-projeto"
}
```

#### Exemplo 2: Usando Senha (Menos Seguro)

```json
{
  "server": {
    "host": "192.168.1.100",
    "port": 22,
    "user": "deploy",
    "password": "sua_senha_aqui"
  },
  "current_version": "v1.0.0",
  "project_name": "meu-projeto"
}
```

⚠️ **Importante**: Se usar senha, adicione `.00cli/settings.json` ao `.gitignore`!

#### Exemplo 3: Com Servidor de Updates Customizado

```json
{
  "server": {
    "host": "192.168.1.100",
    "port": 22,
    "user": "deploy",
    "ssh_key": "/home/usuario/.ssh/id_rsa"
  },
  "current_version": "v1.0.0",
  "project_name": "meu-projeto",
  "update_server": "http://192.168.1.100:8080"
}
```

## deploy.json

O arquivo `./.00cli/deploy.json` contém as configurações de deploy.

### Estrutura

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
    "NODE_ENV": "production"
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

### Campos

#### `type` (obrigatório)
- **Tipo**: `string`
- **Valores**: `"ssh"`, `"docker"`, `"git"`
- **Descrição**: Tipo de deploy a ser executado

#### `commands` (obrigatório para tipo ssh)
- **Tipo**: `array<string>`
- **Descrição**: Lista de comandos a serem executados no servidor

#### `environment` (opcional)
- **Tipo**: `object`
- **Descrição**: Variáveis de ambiente para o deploy

#### `provision` (opcional)
- **Tipo**: `object`
- **Descrição**: Arquivos a serem enviados ao servidor antes do deploy

## Variáveis de Ambiente

O 00cli também suporta configuração via variáveis de ambiente:

| Variável | Descrição |
|----------|-----------|
| `00CLI_UPDATE_SERVER` | URL do servidor de atualizações |
| `00CLI_VERSION` | Versão do CLI (override) |

Exemplo:

```bash
export 00CLI_UPDATE_SERVER="http://192.168.1.100:8080"
00cli update
```
