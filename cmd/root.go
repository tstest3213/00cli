package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	projectPath string
	verbose     bool
)

// Settings representa as configurações do servidor
type Settings struct {
	Server struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		SSHKey   string `json:"ssh_key,omitempty"`
		Password string `json:"password,omitempty"`
	} `json:"server"`
	CurrentVersion string `json:"current_version"`
	ProjectName    string `json:"project_name,omitempty"`
}

// DeployConfig representa a configuração de deploy
type DeployConfig struct {
	Type        string            `json:"type"` // "docker", "ssh", "git", etc
	Commands    []string          `json:"commands,omitempty"`
	Scripts     []string          `json:"scripts,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
	Provision   struct {
		Path  string   `json:"path"`
		Files []string `json:"files,omitempty"`
	} `json:"provision"`
}

var rootCmd = &cobra.Command{
	Use:   "00cli",
	Short: "00cli - Ferramenta CLI para deploy e gerenciamento de projetos",
	Long: `00cli é uma ferramenta CLI inspirada no agent-cursor para gerenciar
deploys e configurações de projetos. O programa verifica automaticamente
por atualizações e requer arquivos de configuração em ./00cli/`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&projectPath, "project", "p", "", "Caminho do projeto (padrão: diretório atual)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Modo verboso")
}

// getProjectRoot retorna o diretório raiz do projeto
func getProjectRoot() (string, error) {
	if projectPath != "" {
		return projectPath, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("erro ao obter diretório atual: %w", err)
	}

	return wd, nil
}

// checkProjectStructure verifica se o projeto tem a estrutura correta
func checkProjectStructure(root string) error {
	settingsPath := filepath.Join(root, "00cli", "settings.json")
	deployPath := filepath.Join(root, "00cli", "deploy.json")

	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		return fmt.Errorf("arquivo não encontrado: %s", settingsPath)
	}

	if _, err := os.Stat(deployPath); os.IsNotExist(err) {
		return fmt.Errorf("arquivo não encontrado: %s", deployPath)
	}

	return nil
}

// loadSettings carrega o arquivo settings.json
func loadSettings(root string) (*Settings, error) {
	path := filepath.Join(root, "00cli", "settings.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, err
	}

	return &settings, nil
}

// loadDeployConfig carrega o arquivo deploy.json
func loadDeployConfig(root string) (*DeployConfig, error) {
	path := filepath.Join(root, "00cli", "deploy.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config DeployConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
