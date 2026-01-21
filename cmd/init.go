package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Inicializa a estrutura 00cli no projeto atual",
	Long:  `Cria os arquivos settings.json e deploy.json no diretório ./00cli/`,
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	root, err := getProjectRoot()
	if err != nil {
		return err
	}

	cliDir := filepath.Join(root, "00cli")
	
	// Criar diretório 00cli se não existir
	if err := os.MkdirAll(cliDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório 00cli: %w", err)
	}

	// Criar settings.json padrão
	settingsPath := filepath.Join(cliDir, "settings.json")
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		settings := Settings{
			CurrentVersion: "v0.0.0",
			ProjectName:    filepath.Base(root),
		}
		settings.Server.Host = "example.com"
		settings.Server.Port = 22
		settings.Server.User = "deploy"

		data, err := json.MarshalIndent(settings, "", "  ")
		if err != nil {
			return err
		}

		if err := os.WriteFile(settingsPath, data, 0644); err != nil {
			return fmt.Errorf("erro ao criar settings.json: %w", err)
		}

		fmt.Printf("✅ Criado: %s\n", settingsPath)
	} else {
		fmt.Printf("⚠️  Arquivo já existe: %s\n", settingsPath)
	}

	// Criar deploy.json padrão
	deployPath := filepath.Join(cliDir, "deploy.json")
	if _, err := os.Stat(deployPath); os.IsNotExist(err) {
		deployConfig := DeployConfig{
			Type: "ssh",
		}
		deployConfig.Commands = []string{
			"git pull",
			"docker-compose up -d --build",
		}
		deployConfig.Provision.Path = "./provision"
		deployConfig.Environment = map[string]string{
			"NODE_ENV": "production",
		}

		data, err := json.MarshalIndent(deployConfig, "", "  ")
		if err != nil {
			return err
		}

		if err := os.WriteFile(deployPath, data, 0644); err != nil {
			return fmt.Errorf("erro ao criar deploy.json: %w", err)
		}

		fmt.Printf("✅ Criado: %s\n", deployPath)
	} else {
		fmt.Printf("⚠️  Arquivo já existe: %s\n", deployPath)
	}

	fmt.Println("\n✅ Estrutura 00cli inicializada com sucesso!")
	fmt.Println("   Edite os arquivos em ./00cli/ para configurar seu projeto.")

	return nil
}
