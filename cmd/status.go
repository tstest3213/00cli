package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Mostra o status atual do projeto e configurações",
	Long:  `Mostra informações sobre o projeto atual, configurações do servidor e versão deployada.`,
	RunE:  runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	root, err := getProjectRoot()
	if err != nil {
		return err
	}

	fmt.Printf("📁 Diretório do projeto: %s\n\n", root)

	// Verificar estrutura
	if err := checkProjectStructure(root); err != nil {
		return fmt.Errorf("estrutura do projeto inválida: %w", err)
	}

	fmt.Println("✅ Estrutura do projeto válida")

	// Carregar e mostrar settings
	settings, err := loadSettings(root)
	if err != nil {
		return fmt.Errorf("erro ao carregar settings: %w", err)
	}

	fmt.Println("\n📋 Configurações do Servidor:")
	fmt.Printf("   Host: %s\n", settings.Server.Host)
	fmt.Printf("   Porta: %d\n", settings.Server.Port)
	fmt.Printf("   Usuário: %s\n", settings.Server.User)
	fmt.Printf("   Versão atual: %s\n", settings.CurrentVersion)

	if settings.ProjectName != "" {
		fmt.Printf("   Nome do projeto: %s\n", settings.ProjectName)
	}

	// Verificar provision
	provisionPath := filepath.Join(root, "provision")
	if info, err := os.Stat(provisionPath); err == nil {
		fmt.Printf("\n📦 Diretório /provision/ existe")
		if info.IsDir() {
			files, _ := os.ReadDir(provisionPath)
			fmt.Printf(" (%d arquivos)\n", len(files))
		}
	} else {
		fmt.Println("\n⚠️  Diretório /provision/ não encontrado")
	}

	// Carregar deploy config
	deployConfig, err := loadDeployConfig(root)
	if err != nil {
		return fmt.Errorf("erro ao carregar deploy.json: %w", err)
	}

	fmt.Printf("\n🚀 Configuração de Deploy:")
	fmt.Printf("   Tipo: %s\n", deployConfig.Type)
	if deployConfig.Provision.Path != "" {
		fmt.Printf("   Path provision: %s\n", deployConfig.Provision.Path)
	}

	return nil
}
