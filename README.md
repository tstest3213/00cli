# 00cli

CLI de deploy e automação para projetos. Simplifica o deploy via SSH, Docker e Git.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## ✨ Funcionalidades

- 🚀 **Deploy Rápido** - Deploy via SSH com comandos customizáveis
- 🐳 **Suporte Docker** - Deploy de containers Docker
- 📦 **Git Integration** - Deploy via Git pull
- 🔄 **Auto-Update** - Verifica e atualiza automaticamente
- ⚙️ **Configuração Simples** - Arquivo JSON para configuração
- 📁 **Provisionamento** - Envie arquivos de configuração para o servidor

## 📦 Instalação

### Instalação Rápida (Linux/macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/tstest3213/00cli/main/install.sh | bash
```

### Compilar do Fonte

```bash
git clone https://github.com/tstest3213/00cli.git
cd 00cli
make build
sudo make install
```

### Verificar Instalação

```bash
00cli version
```

Para mais opções de instalação, veja [docs/install.md](docs/install.md).

## 🚀 Início Rápido

### 1. Inicializar Projeto

```bash
cd /seu/projeto
00cli init
```

Isso criará a estrutura `.00cli/` com arquivos de configuração.

### 2. Configurar Servidor

Edite `.00cli/settings.json`:

```json
{
  "server": {
    "host": "meuservidor.com",
    "port": 22,
    "user": "deploy",
    "ssh_key": "~/.ssh/id_rsa"
  },
  "project_name": "meu-projeto"
}
```

### 3. Configurar Deploy

Edite `.00cli/deploy.json`:

```json
{
  "type": "ssh",
  "commands": [
    "cd /var/www/meu-projeto",
    "git pull origin main",
    "npm install",
    "npm run build",
    "pm2 restart app"
  ]
}
```

### 4. Fazer Deploy

```bash
00cli deploy
```

## 📖 Comandos

| Comando | Descrição |
|---------|-----------|
| `00cli init` | Inicializa estrutura de configuração |
| `00cli deploy` | Executa deploy no servidor |
| `00cli status` | Mostra status do servidor |
| `00cli version` | Mostra versão do CLI |
| `00cli update` | Atualiza para versão mais recente |

### Flags Globais

```bash
00cli deploy --verbose  # Modo verboso
00cli deploy --dry-run  # Simular sem executar
```

## 📚 Documentação

| Documento | Descrição |
|-----------|-----------|
| [Instalação](docs/install.md) | Guia completo de instalação |
| [Configuração](docs/settings.md) | Referência do settings.json |
| [Exemplos](docs/examples.md) | Exemplos de uso |
| [Servidor de Updates](docs/update-server.md) | Servidor customizado de atualizações |
| [GitHub Releases](docs/github-releases.md) | Como criar releases no GitHub |

## 🔄 Atualizações Automáticas

O 00cli verifica automaticamente por novas versões. Quando disponível:

```
⚠️  Nova versão disponível: v0.2.0 (atual: v0.1.0)
   Execute '00cli update' para atualizar automaticamente
```

Para atualizar:

```bash
00cli update
```

### Servidor Customizado

Você pode usar seu próprio servidor de atualizações:

```json
{
  "update_server": "http://seu-servidor:8080"
}
```

Veja [docs/update-server.md](docs/update-server.md) para detalhes.

## 🛠️ Desenvolvimento

### Pré-requisitos

- Go 1.21+
- Make

### Comandos do Makefile

```bash
make build          # Compilar binário
make install        # Instalar no sistema
make test           # Rodar testes
make clean          # Limpar builds
make build-all      # Compilar para todas plataformas
make server         # Iniciar servidor de updates
```

### Estrutura do Projeto

```
00cli/
├── cmd/              # Comandos CLI
│   ├── root.go       # Comando raiz
│   ├── init.go       # 00cli init
│   ├── deploy.go     # 00cli deploy
│   ├── status.go     # 00cli status
│   ├── update.go     # 00cli update
│   └── version.go    # 00cli version
├── internal/         # Código interno
│   └── deploy/       # Lógica de deploy
├── server-update/    # Servidor de atualizações
├── docs/             # Documentação
└── Makefile
```

## 📝 Licença

MIT License - veja [LICENSE](LICENSE) para detalhes.

## 🤝 Contribuindo

1. Fork o repositório
2. Crie uma branch (`git checkout -b feature/minha-feature`)
3. Commit suas mudanças (`git commit -am 'Adiciona minha feature'`)
4. Push para a branch (`git push origin feature/minha-feature`)
5. Abra um Pull Request

## 📞 Suporte

- 📖 [Documentação](docs/)
- 🐛 [Issues](https://github.com/tstest3213/00cli/issues)
