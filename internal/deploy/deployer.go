package deploy

import (
	"fmt"
)

// Deployer interface para diferentes tipos de deploy
type Deployer interface {
	Execute(commands []string) error
}

// NewDeployer cria um deployer baseado no tipo
func NewDeployer(deployType string, config interface{}) (Deployer, error) {
	switch deployType {
	case "ssh":
		return createSSHDeployer(config)
	case "docker":
		return createDockerDeployer(config)
	case "git":
		return createGitDeployer(config)
	default:
		return nil, fmt.Errorf("tipo de deploy não suportado: %s", deployType)
	}
}

// ConfigMap é um mapa genérico para configurações
type ConfigMap map[string]interface{}

func createSSHDeployer(config interface{}) (Deployer, error) {
	cfg, ok := config.(ConfigMap)
	if !ok {
		return nil, fmt.Errorf("configuração SSH inválida")
	}

	deployer := &SSHDeployer{}

	if host, ok := cfg["host"].(string); ok {
		deployer.Host = host
	}
	if port, ok := cfg["port"].(int); ok {
		deployer.Port = port
	}
	if user, ok := cfg["user"].(string); ok {
		deployer.User = user
	}
	if sshKey, ok := cfg["ssh_key"].(string); ok {
		deployer.SSHKey = sshKey
	}
	if password, ok := cfg["password"].(string); ok {
		deployer.Password = password
	}

	return deployer, nil
}

func createDockerDeployer(config interface{}) (Deployer, error) {
	cfg, ok := config.(ConfigMap)
	if !ok {
		return nil, fmt.Errorf("configuração Docker inválida")
	}

	deployer := &DockerDeployer{}

	if projectPath, ok := cfg["project_path"].(string); ok {
		deployer.ProjectPath = projectPath
	}
	if composeFile, ok := cfg["compose_file"].(string); ok {
		deployer.ComposeFile = composeFile
	}
	if env, ok := cfg["environment"].(map[string]string); ok {
		deployer.Environment = env
	} else {
		deployer.Environment = make(map[string]string)
	}

	return deployer, nil
}

func createGitDeployer(config interface{}) (Deployer, error) {
	cfg, ok := config.(ConfigMap)
	if !ok {
		return nil, fmt.Errorf("configuração Git inválida")
	}

	deployer := &GitDeployer{}

	if repo, ok := cfg["repository"].(string); ok {
		deployer.Repository = repo
	}
	if branch, ok := cfg["branch"].(string); ok {
		deployer.Branch = branch
	}
	if projectPath, ok := cfg["project_path"].(string); ok {
		deployer.ProjectPath = projectPath
	}
	if commands, ok := cfg["commands"].([]string); ok {
		deployer.Commands = commands
	}

	return deployer, nil
}
