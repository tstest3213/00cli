package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSettings(t *testing.T) {
	// Criar diretório temporário
	tmpDir := t.TempDir()
	cliDir := filepath.Join(tmpDir, ".00cli")
	if err := os.MkdirAll(cliDir, 0755); err != nil {
		t.Fatalf("erro ao criar diretório: %v", err)
	}

	// Criar settings.json de teste
	settings := Settings{
		CurrentVersion: "v1.0.0",
		ProjectName:    "test-project",
	}
	settings.Server.Host = "test.com"
	settings.Server.Port = 22
	settings.Server.User = "testuser"

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		t.Fatalf("erro ao serializar settings: %v", err)
	}

	settingsPath := filepath.Join(cliDir, "settings.json")
	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		t.Fatalf("erro ao escrever settings.json: %v", err)
	}

	// Testar carregamento
	loaded, err := loadSettings(tmpDir)
	if err != nil {
		t.Fatalf("erro ao carregar settings: %v", err)
	}

	if loaded.Server.Host != "test.com" {
		t.Errorf("esperado host 'test.com', obtido '%s'", loaded.Server.Host)
	}

	if loaded.Server.Port != 22 {
		t.Errorf("esperado port 22, obtido %d", loaded.Server.Port)
	}

	if loaded.CurrentVersion != "v1.0.0" {
		t.Errorf("esperado version 'v1.0.0', obtido '%s'", loaded.CurrentVersion)
	}
}

func TestLoadDeployConfig(t *testing.T) {
	// Criar diretório temporário
	tmpDir := t.TempDir()
	cliDir := filepath.Join(tmpDir, ".00cli")
	if err := os.MkdirAll(cliDir, 0755); err != nil {
		t.Fatalf("erro ao criar diretório: %v", err)
	}

	// Criar deploy.json de teste
	deployConfig := DeployConfig{
		Type: "ssh",
		Commands: []string{
			"git pull",
			"npm install",
			"npm run build",
		},
	}
	deployConfig.Provision.Path = "./provision"
	deployConfig.Environment = map[string]string{
		"NODE_ENV": "production",
	}

	data, err := json.MarshalIndent(deployConfig, "", "  ")
	if err != nil {
		t.Fatalf("erro ao serializar deploy config: %v", err)
	}

	deployPath := filepath.Join(cliDir, "deploy.json")
	if err := os.WriteFile(deployPath, data, 0644); err != nil {
		t.Fatalf("erro ao escrever deploy.json: %v", err)
	}

	// Testar carregamento
	loaded, err := loadDeployConfig(tmpDir)
	if err != nil {
		t.Fatalf("erro ao carregar deploy config: %v", err)
	}

	if loaded.Type != "ssh" {
		t.Errorf("esperado type 'ssh', obtido '%s'", loaded.Type)
	}

	if len(loaded.Commands) != 3 {
		t.Errorf("esperado 3 comandos, obtido %d", len(loaded.Commands))
	}

	if loaded.Provision.Path != "./provision" {
		t.Errorf("esperado provision path './provision', obtido '%s'", loaded.Provision.Path)
	}
}

func TestCheckProjectStructure(t *testing.T) {
	// Criar diretório temporário
	tmpDir := t.TempDir()
	cliDir := filepath.Join(tmpDir, ".00cli")
	if err := os.MkdirAll(cliDir, 0755); err != nil {
		t.Fatalf("erro ao criar diretório: %v", err)
	}

	// Testar sem arquivos (deve falhar)
	err := checkProjectStructure(tmpDir)
	if err == nil {
		t.Error("esperado erro quando arquivos não existem")
	}

	// Criar settings.json
	settingsPath := filepath.Join(cliDir, "settings.json")
	if err := os.WriteFile(settingsPath, []byte("{}"), 0644); err != nil {
		t.Fatalf("erro ao criar settings.json: %v", err)
	}

	// Testar sem deploy.json (deve falhar)
	err = checkProjectStructure(tmpDir)
	if err == nil {
		t.Error("esperado erro quando deploy.json não existe")
	}

	// Criar deploy.json
	deployPath := filepath.Join(cliDir, "deploy.json")
	if err := os.WriteFile(deployPath, []byte("{}"), 0644); err != nil {
		t.Fatalf("erro ao criar deploy.json: %v", err)
	}

	// Testar com ambos os arquivos (deve passar)
	err = checkProjectStructure(tmpDir)
	if err != nil {
		t.Errorf("não esperado erro quando ambos arquivos existem: %v", err)
	}
}
