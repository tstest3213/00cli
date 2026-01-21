package deploy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// GitDeployer implementa deploy via Git
type GitDeployer struct {
	Repository  string
	Branch      string
	ProjectPath string
	Commands    []string
}

// Execute executa deploy via Git
func (d *GitDeployer) Execute(commands []string) error {
	// Usar comandos fornecidos se não houver comandos configurados
	if len(commands) > 0 {
		d.Commands = commands
	}
	// Se não houver repositório configurado, assumir que já está em um repo Git
	if d.Repository == "" {
		// Verificar se é um repositório Git
		gitDir := filepath.Join(d.ProjectPath, ".git")
		if _, err := os.Stat(gitDir); os.IsNotExist(err) {
			return fmt.Errorf("diretório não é um repositório Git e nenhum repositório foi especificado")
		}

		fmt.Println("📦 Usando repositório Git local")
	} else {
		// Clonar ou atualizar repositório
		if err := d.cloneOrUpdate(); err != nil {
			return err
		}
	}

	// Executar comandos pós-deploy
	if len(d.Commands) > 0 {
		for i, cmd := range d.Commands {
			fmt.Printf("  [%d/%d] Executando: %s\n", i+1, len(d.Commands), cmd)

			parts := parseCommand(cmd)
			if len(parts) == 0 {
				continue
			}

			command := exec.Command(parts[0], parts[1:]...)
			command.Dir = d.ProjectPath
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr

			if err := command.Run(); err != nil {
				return fmt.Errorf("erro ao executar '%s': %w", cmd, err)
			}
		}
	}

	return nil
}

func (d *GitDeployer) cloneOrUpdate() error {
	gitDir := filepath.Join(d.ProjectPath, ".git")

	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		// Clonar repositório
		fmt.Printf("📥 Clonando repositório: %s\n", d.Repository)

		branch := d.Branch
		if branch == "" {
			branch = "main"
		}

		cmd := exec.Command("git", "clone", "-b", branch, d.Repository, d.ProjectPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("erro ao clonar repositório: %w", err)
		}
	} else {
		// Atualizar repositório existente
		fmt.Println("🔄 Atualizando repositório Git...")

		// Pull
		cmd := exec.Command("git", "pull")
		cmd.Dir = d.ProjectPath
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("erro ao fazer pull: %w", err)
		}

		// Se branch especificada, fazer checkout
		if d.Branch != "" {
			cmd = exec.Command("git", "checkout", d.Branch)
			cmd.Dir = d.ProjectPath
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				return fmt.Errorf("erro ao fazer checkout da branch %s: %w", d.Branch, err)
			}
		}
	}

	return nil
}
