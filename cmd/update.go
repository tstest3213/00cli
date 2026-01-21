package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	githubRepoOwner = "tstest3213"
	githubRepoName  = "00cli" // Repositório do 00cli (não 00off)
	githubAPIURL    = "https://api.github.com/repos/%s/%s/releases/latest"
)

type GitHubRelease struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	PublishedAt time.Time `json:"published_at"`
	HTMLURL     string    `json:"html_url"`
	Body        string    `json:"body"`
}

// CheckForUpdates verifica se há atualizações disponíveis no GitHub
func CheckForUpdates(currentVersion string) {
	// Aguardar um pouco para não bloquear o início do programa
	time.Sleep(500 * time.Millisecond)

	url := fmt.Sprintf(githubAPIURL, githubRepoOwner, githubRepoName)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return // Falha silenciosa
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return // Falha silenciosa (sem internet, etc)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return
	}

	// Comparar versões (simples comparação de strings)
	if release.TagName != "" && release.TagName != currentVersion {
		fmt.Printf("\n⚠️  Nova versão disponível: %s (atual: %s)\n", release.TagName, currentVersion)
		fmt.Printf("   Baixe em: %s\n\n", release.HTMLURL)
	}
}
