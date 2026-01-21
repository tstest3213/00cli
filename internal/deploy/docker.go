package deploy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DockerDeployer implementa deploy via Docker/Docker Compose
type DockerDeployer struct {
	ComposeFile string
	ProjectPath string
	Environment map[string]string
}

// Execute executa deploy Docker
func (d *DockerDeployer) Execute(commands []string) error {
	// Se não houver comandos específicos, usar comandos padrão do Docker Compose
	if len(commands) == 0 {
		commands = []string{
			"docker-compose down",
			"docker-compose pull",
			"docker-compose up -d --build",
		}
	}

	// Procurar docker-compose.yml
	composeFile := d.ComposeFile
	if composeFile == "" {
		composeFile = filepath.Join(d.ProjectPath, "docker-compose.yml")
		if _, err := os.Stat(composeFile); os.IsNotExist(err) {
			// Tentar em provision/
			composeFile = filepath.Join(d.ProjectPath, "provision", "docker-compose.yml")
		}
	}

	// Verificar se docker-compose existe
	if _, err := os.Stat(composeFile); os.IsNotExist(err) {
		return fmt.Errorf("arquivo docker-compose.yml não encontrado em %s", d.ProjectPath)
	}

	fmt.Printf("📦 Usando docker-compose: %s\n", composeFile)

	// Executar comandos
	for i, cmd := range commands {
		fmt.Printf("  [%d/%d] Executando: %s\n", i+1, len(commands), cmd)

		// Parsear comando
		parts := parseCommand(cmd)
		if len(parts) == 0 {
			continue
		}

		// Criar comando
		command := exec.Command(parts[0], parts[1:]...)
		command.Dir = d.ProjectPath
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		// Adicionar variáveis de ambiente
		command.Env = os.Environ()
		for k, v := range d.Environment {
			command.Env = append(command.Env, fmt.Sprintf("%s=%s", k, v))
		}

		// Executar
		if err := command.Run(); err != nil {
			return fmt.Errorf("erro ao executar '%s': %w", cmd, err)
		}
	}

	return nil
}

// parseCommand divide uma string de comando em partes
func parseCommand(cmd string) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false

	for _, char := range cmd {
		if char == '"' || char == '\'' {
			inQuotes = !inQuotes
		} else if char == ' ' && !inQuotes {
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
		} else {
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}
