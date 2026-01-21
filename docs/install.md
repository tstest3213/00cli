# Instalação do 00cli

## 🚀 Instalação Rápida

### Linux/macOS

```bash
# Baixar e executar o instalador
curl -fsSL https://raw.githubusercontent.com/tstest3213/00cli/main/install.sh | bash
```

Ou usando wget:

```bash
wget -qO- https://raw.githubusercontent.com/tstest3213/00cli/main/install.sh | bash
```

### Instalação Manual

1. **Baixar binário pré-compilado** (recomendado):

```bash
# Detectar OS e ARCH
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Normalizar ARCH
case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
esac

# Baixar última versão
LATEST=$(curl -s https://api.github.com/repos/tstest3213/00cli/releases/latest | grep "browser_download_url.*00cli-${OS}-${ARCH}" | cut -d'"' -f4)
curl -L -o 00cli "$LATEST"
chmod +x 00cli
sudo mv 00cli /usr/local/bin/
```

2. **Compilar a partir do código fonte**:

```bash
git clone https://github.com/tstest3213/00cli.git
cd 00cli
make build
sudo make install
```

## 🔄 Atualização Automática

O `00cli` verifica automaticamente por atualizações toda vez que é executado. Se houver uma nova versão disponível, você verá:

```
⚠️  Nova versão disponível: v0.2.0 (atual: v0.1.0)
   Execute '00cli update' para atualizar automaticamente
```

### Atualizar Manualmente

```bash
# Atualizar para a versão mais recente
00cli update
```

O comando `update` irá:
1. Verificar a última versão no GitHub
2. Baixar o binário correto para sua plataforma
3. Substituir o binário atual automaticamente

## 📦 Verificar Instalação

```bash
# Verificar versão
00cli version

# Verificar ajuda
00cli --help
```

## 🗑️ Desinstalar

```bash
sudo rm /usr/local/bin/00cli
```

Ou usando o Makefile (se tiver o código fonte):

```bash
make uninstall
```

## 🔧 Requisitos

- **Linux/macOS**: Nenhum requisito adicional (binário estático)
- **Compilação a partir do código**: Go 1.21 ou superior

## 📝 Notas

- O instalador tenta baixar o binário pré-compilado primeiro
- Se não encontrar binário para sua plataforma, tenta compilar do código fonte
- Binários disponíveis para: Linux (amd64/arm64), macOS (amd64/arm64), Windows (amd64)
