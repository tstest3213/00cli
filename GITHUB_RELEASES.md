# Como Criar Releases no GitHub para 00cli

Este documento explica como criar releases no GitHub para habilitar a verificação automática de atualizações do `00cli`.

## 📋 Pré-requisitos

1. Repositório criado no GitHub: `https://github.com/tstest3213/00cli`
2. Acesso de escrita ao repositório
3. Git configurado localmente

## 🚀 Processo de Release

### 1. Preparar o Código

Certifique-se de que todas as mudanças estão commitadas:

```bash
git add .
git commit -m "Preparar release v0.2.0"
git push origin main
```

### 2. Criar uma Tag

As tags Git são usadas para identificar releases:

```bash
# Criar tag anotada
git tag -a v0.2.0 -m "Release v0.2.0: Adiciona suporte a deploy SSH, Docker e Git"

# Enviar tag para o GitHub
git push origin v0.2.0
```

### 3. Criar Release no GitHub

#### Opção A: Via Interface Web do GitHub

1. Acesse: `https://github.com/tstest3213/00cli/releases/new`
2. Selecione a tag que você acabou de criar (ex: `v0.2.0`)
3. Preencha:
   - **Title**: `v0.2.0` (ou um título descritivo)
   - **Description**: Descreva as mudanças desta versão
4. Opcionalmente, anexe binários compilados
5. Clique em **"Publish release"**

#### Opção B: Via GitHub CLI (gh)

Se você tem o GitHub CLI instalado:

```bash
# Criar release
gh release create v0.2.0 \
  --title "v0.2.0" \
  --notes "Release v0.2.0: Adiciona suporte a deploy SSH, Docker e Git" \
  --target main
```

### 4. Anexar Binários (Opcional mas Recomendado)

Para facilitar a instalação, você pode anexar binários compilados:

```bash
# Compilar para diferentes plataformas
GOOS=linux GOARCH=amd64 go build -o 00cli-linux-amd64 .
GOOS=linux GOARCH=arm64 go build -o 00cli-linux-arm64 .
GOOS=darwin GOARCH=amd64 go build -o 00cli-darwin-amd64 .
GOOS=darwin GOARCH=arm64 go build -o 00cli-darwin-arm64 .
GOOS=windows GOARCH=amd64 go build -o 00cli-windows-amd64.exe .

# Anexar ao release via GitHub CLI
gh release upload v0.2.0 00cli-linux-amd64 00cli-linux-arm64 \
  00cli-darwin-amd64 00cli-darwin-arm64 00cli-windows-amd64.exe
```

## 🔄 Automatização com GitHub Actions

Você pode automatizar a criação de releases usando GitHub Actions. Crie `.github/workflows/release.yml`:

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Build
        run: |
          GOOS=linux GOARCH=amd64 go build -o 00cli-linux-amd64 .
          GOOS=linux GOARCH=arm64 go build -o 00cli-linux-arm64 .
          GOOS=darwin GOARCH=amd64 go build -o 00cli-darwin-amd64 .
          GOOS=darwin GOARCH=arm64 go build -o 00cli-darwin-arm64 .
          GOOS=windows GOARCH=amd64 go build -o 00cli-windows-amd64.exe .
      
      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          files: |
            00cli-linux-amd64
            00cli-linux-arm64
            00cli-darwin-amd64
            00cli-darwin-arm64
            00cli-windows-amd64.exe
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

## 📝 Formato de Versionamento

O `00cli` usa [Semantic Versioning](https://semver.org/):

- **MAJOR** (v1.0.0): Mudanças incompatíveis
- **MINOR** (v0.1.0): Novas funcionalidades compatíveis
- **PATCH** (v0.0.1): Correções de bugs

## ✅ Verificação

Após criar o release, o `00cli` verificará automaticamente na próxima execução:

```bash
00cli --help
```

Se houver uma nova versão, você verá:

```
⚠️  Nova versão disponível: v0.2.0 (atual: v0.1.0)
   Baixe em: https://github.com/tstest3213/00cli/releases/tag/v0.2.0
```

## 🔗 Links Úteis

- [GitHub Releases API](https://docs.github.com/en/rest/releases/releases)
- [Semantic Versioning](https://semver.org/)
- [GitHub Actions](https://docs.github.com/en/actions)
