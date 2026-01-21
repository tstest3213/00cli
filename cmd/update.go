package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	githubRepoOwner = "tstest3213"
	githubRepoName  = "00cli"
	githubAPIURL    = "https://api.github.com/repos/%s/%s/releases/latest"
)

// Release representa uma release (GitHub ou servidor customizado)
type Release struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	PublishedAt time.Time `json:"published_at"`
	HTMLURL     string    `json:"html_url"`
	Body        string    `json:"body"`
	Assets      []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
		Size               int    `json:"size"`
	} `json:"assets"`
}

// GitHubRelease mantém compatibilidade
type GitHubRelease = Release

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Atualiza o 00cli para a versão mais recente",
	Long: `Verifica e instala automaticamente a versão mais recente do 00cli.
Pode usar servidor customizado (configurado em .00cli/settings.json ou variável 00CLI_UPDATE_SERVER)
ou GitHub como fallback.`,
	RunE: runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

// CheckForUpdates verifica se há atualizações disponíveis
func CheckForUpdates(currentVersion string) {
	// Aguardar um pouco para não bloquear o início do programa
	time.Sleep(500 * time.Millisecond)

	release, err := getLatestRelease()
	if err != nil {
		return // Falha silenciosa
	}

	// Comparar versões
	if release.TagName != "" && release.TagName != currentVersion {
		fmt.Printf("\n⚠️  Nova versão disponível: %s (atual: %s)\n", release.TagName, currentVersion)
		fmt.Printf("   Execute '00cli update' para atualizar automaticamente\n")
		if release.HTMLURL != "" {
			fmt.Printf("   Ou baixe em: %s\n\n", release.HTMLURL)
		} else {
			fmt.Println()
		}
	}
}

// getUpdateServerURL obtém a URL do servidor de atualizações configurado
func getUpdateServerURL() string {
	// Verificar variável de ambiente primeiro
	if url := os.Getenv("00CLI_UPDATE_SERVER"); url != "" {
		return url
	}

	// Tentar carregar do settings.json do projeto atual (se existir)
	if _, err := os.Getwd(); err == nil {
		root, _ := getProjectRoot()
		if settings, err := loadSettings(root); err == nil && settings.UpdateServer != "" {
			return settings.UpdateServer
		}
	}

	return "" // Usar GitHub como padrão
}

// getLatestRelease obtém a última release (servidor customizado ou GitHub)
func getLatestRelease() (*Release, error) {
	// Verificar se há servidor customizado configurado
	customServer := getUpdateServerURL()
	if customServer != "" {
		release, err := getLatestReleaseFromCustomServer(customServer)
		if err == nil {
			return release, nil
		}
		// Se falhar, tentar GitHub como fallback
	}

	// Usar GitHub como padrão ou fallback
	return getLatestReleaseFromGitHub()
}

// getLatestReleaseFromCustomServer obtém release do servidor customizado
func getLatestReleaseFromCustomServer(serverURL string) (*Release, error) {
	// Garantir que a URL termina com /latest ou /updates/latest
	url := strings.TrimSuffix(serverURL, "/")
	if !strings.HasSuffix(url, "/latest") && !strings.HasSuffix(url, "/updates") {
		url += "/latest"
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "00cli-updater")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var release Release
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, err
	}

	// Se o servidor retornar apenas URL do binário, construir asset
	if len(release.Assets) == 0 && release.TagName != "" {
		binaryName := getBinaryName()
		baseURL := strings.TrimSuffix(url, "/latest")
		release.Assets = []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
			Size               int    `json:"size"`
		}{
			{
				Name:               binaryName,
				BrowserDownloadURL: fmt.Sprintf("%s/download/%s", baseURL, binaryName),
			},
		}
	}

	return &release, nil
}

// getLatestReleaseFromGitHub obtém a última release do GitHub
func getLatestReleaseFromGitHub() (*Release, error) {
	url := fmt.Sprintf(githubAPIURL, githubRepoOwner, githubRepoName)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var release Release
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, err
	}

	return &release, nil
}

// getBinaryName retorna o nome do binário baseado no OS e ARCH
func getBinaryName() string {
	os := runtime.GOOS
	arch := runtime.GOARCH

	var ext string
	if os == "windows" {
		ext = ".exe"
	}

	return fmt.Sprintf("00cli-%s-%s%s", os, arch, ext)
}

// findBinaryAsset encontra o asset correto para a plataforma atual
func findBinaryAsset(release *Release) string {
	binaryName := getBinaryName()

	for _, asset := range release.Assets {
		if asset.Name == binaryName {
			return asset.BrowserDownloadURL
		}
	}

	// Se não encontrar, tentar construir URL baseada no padrão do servidor customizado
	if customServer := getUpdateServerURL(); customServer != "" {
		baseURL := strings.TrimSuffix(customServer, "/latest")
		baseURL = strings.TrimSuffix(baseURL, "/updates")
		return fmt.Sprintf("%s/download/%s", baseURL, binaryName)
	}

	return ""
}

// runUpdate executa a atualização
func runUpdate(cmd *cobra.Command, args []string) error {
	fmt.Println("🔍 Verificando atualizações...")

	release, err := getLatestRelease()
	if err != nil {
		return fmt.Errorf("erro ao verificar atualizações: %w", err)
	}

	// Obter versão atual
	currentVersion := getCurrentVersion()
	if release.TagName == currentVersion {
		fmt.Printf("✅ Você já está na versão mais recente: %s\n", currentVersion)
		return nil
	}

	fmt.Printf("📦 Nova versão encontrada: %s (atual: %s)\n", release.TagName, currentVersion)
	fmt.Println("🚀 Iniciando atualização...")

	// Encontrar o binário correto
	downloadURL := findBinaryAsset(release)
	if downloadURL == "" {
		return fmt.Errorf("binário não encontrado para %s/%s. Baixe manualmente em: %s", runtime.GOOS, runtime.GOARCH, release.HTMLURL)
	}

	// Obter caminho do binário atual
	currentBinary, err := os.Executable()
	if err != nil {
		return fmt.Errorf("erro ao obter caminho do binário: %w", err)
	}

	// Criar arquivo temporário
	tmpFile := filepath.Join(os.TempDir(), "00cli-update")
	if runtime.GOOS == "windows" {
		tmpFile += ".exe"
	}

	fmt.Printf("⬇️  Baixando %s...\n", release.TagName)

	// Baixar novo binário
	if err := downloadFile(downloadURL, tmpFile); err != nil {
		return fmt.Errorf("erro ao baixar: %w", err)
	}

	// Tornar executável (Unix)
	if runtime.GOOS != "windows" {
		if err := os.Chmod(tmpFile, 0755); err != nil {
			return fmt.Errorf("erro ao tornar executável: %w", err)
		}
	}

	fmt.Println("📦 Instalando nova versão...")

	// Substituir binário antigo
	if runtime.GOOS == "windows" {
		// Windows: precisa fechar o processo primeiro
		oldFile := currentBinary + ".old"
		if err := os.Rename(currentBinary, oldFile); err != nil {
			return fmt.Errorf("erro ao renomear binário antigo: %w", err)
		}
		if err := os.Rename(tmpFile, currentBinary); err != nil {
			os.Rename(oldFile, currentBinary) // Reverter em caso de erro
			return fmt.Errorf("erro ao instalar novo binário: %w", err)
		}
		os.Remove(oldFile)
	} else {
		// Unix: pode substituir diretamente
		if err := os.Rename(tmpFile, currentBinary); err != nil {
			return fmt.Errorf("erro ao instalar novo binário: %w", err)
		}
	}

	fmt.Printf("✅ Atualização concluída! Nova versão: %s\n", release.TagName)
	fmt.Println("   Execute '00cli version' para verificar.")

	return nil
}

// downloadFile baixa um arquivo de uma URL
func downloadFile(url, dest string) error {
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

// getCurrentVersion obtém a versão atual do binário
func getCurrentVersion() string {
	// Tentar ler da variável de ambiente (setada no build)
	if v := os.Getenv("00CLI_VERSION"); v != "" {
		return v
	}

	// Tentar executar o comando version
	exe, err := os.Executable()
	if err != nil {
		return "v0.0.0"
	}

	cmd := exec.Command(exe, "version")
	output, err := cmd.Output()
	if err != nil {
		return "v0.0.0"
	}

	// Extrair versão da saída (formato: "00cli version v0.1.0")
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "version") {
			parts := strings.Fields(line)
			for _, part := range parts {
				if strings.HasPrefix(part, "v") {
					return part
				}
			}
		}
	}

	return "v0.0.0"
}
