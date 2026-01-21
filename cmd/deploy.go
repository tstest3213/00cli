package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tstest3213/00cli/internal/deploy"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Faz deploy da aplicação no servidor configurado",
	Long: `Faz deploy da aplicação usando as configurações em ./.00cli/settings.json
e ./.00cli/deploy.json. O deploy utiliza o diretório /provision/ se disponível.`,
	RunE: runDeploy,
}

func init() {
	rootCmd.AddCommand(deployCmd)
}

func runDeploy(cmd *cobra.Command, args []string) error {
	root, err := getProjectRoot()
	if err != nil {
		return err
	}

	if err := checkProjectStructure(root); err != nil {
		return fmt.Errorf("estrutura do projeto inválida: %w", err)
	}

	// Carregar settings.json
	settings, err := loadSettings(root)
	if err != nil {
		return fmt.Errorf("erro ao carregar settings.json: %w", err)
	}

	// Carregar deploy.json
	deployConfig, err := loadDeployConfig(root)
	if err != nil {
		return fmt.Errorf("erro ao carregar deploy.json: %w", err)
	}

	if verbose {
		fmt.Printf("📦 Projeto: %s\n", root)
		fmt.Printf("🖥️  Servidor: %s@%s:%d\n", settings.Server.User, settings.Server.Host, settings.Server.Port)
		fmt.Printf("📋 Versão atual no servidor: %s\n", settings.CurrentVersion)
	}

	// Verificar se existe diretório provision
	provisionPath := filepath.Join(root, "provision")
	if _, err := os.Stat(provisionPath); os.IsNotExist(err) {
		if verbose {
			fmt.Printf("⚠️  Diretório /provision/ não encontrado\n")
		}
	} else {
		if verbose {
			fmt.Printf("✅ Diretório /provision/ encontrado\n")
		}
	}

	fmt.Println("\n🚀 Iniciando deploy...")
	fmt.Printf("   Tipo: %s\n", deployConfig.Type)

	// Criar configuração para o deployer
	config := deploy.ConfigMap{
		"project_path": root,
	}

	// Configurar baseado no tipo de deploy
	switch deployConfig.Type {
	case "ssh":
		config["host"] = settings.Server.Host
		config["port"] = settings.Server.Port
		config["user"] = settings.Server.User
		if settings.Server.SSHKey != "" {
			config["ssh_key"] = settings.Server.SSHKey
		}
		if settings.Server.Password != "" {
			config["password"] = settings.Server.Password
		}

	case "git":
		// Git deploy pode ser local ou remoto
		config["repository"] = "" // Será detectado automaticamente se for repo local
		config["branch"] = "main"
		config["commands"] = deployConfig.Commands

	default:
		return fmt.Errorf("tipo de deploy não suportado: %s. Tipos suportados: ssh, git", deployConfig.Type)
	}

	// Criar deployer
	deployer, err := deploy.NewDeployer(deployConfig.Type, config)
	if err != nil {
		return fmt.Errorf("erro ao criar deployer: %w", err)
	}

	// Executar deploy
	if err := deployer.Execute(deployConfig.Commands); err != nil {
		return fmt.Errorf("erro durante deploy: %w", err)
	}

	// Atualizar versão no settings.json (opcional - pode ser feito manualmente)
	fmt.Println("\n✅ Deploy concluído com sucesso!")
	return nil
}
