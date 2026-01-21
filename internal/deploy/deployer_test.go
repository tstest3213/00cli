package deploy

import (
	"testing"
)

func TestNewDeployer(t *testing.T) {
	tests := []struct {
		name        string
		deployType  string
		config      ConfigMap
		expectError bool
	}{
		{
			name:       "SSH deployer válido",
			deployType: "ssh",
			config: ConfigMap{
				"host": "example.com",
				"port": 22,
				"user": "deploy",
			},
			expectError: false,
		},
		{
			name:       "Docker deployer válido",
			deployType: "docker",
			config: ConfigMap{
				"project_path": "/tmp/test",
				"environment":  map[string]string{"NODE_ENV": "production"},
			},
			expectError: false,
		},
		{
			name:       "Git deployer válido",
			deployType: "git",
			config: ConfigMap{
				"project_path": "/tmp/test",
				"branch":       "main",
			},
			expectError: false,
		},
		{
			name:        "Tipo inválido",
			deployType:  "invalid",
			config:      ConfigMap{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deployer, err := NewDeployer(tt.deployType, tt.config)
			if tt.expectError {
				if err == nil {
					t.Errorf("esperado erro, mas não houve erro")
				}
				if deployer != nil {
					t.Errorf("esperado deployer nil quando há erro")
				}
			} else {
				if err != nil {
					t.Errorf("não esperado erro: %v", err)
				}
				if deployer == nil {
					t.Errorf("esperado deployer não-nil")
				}
			}
		})
	}
}

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name     string
		cmd      string
		expected []string
	}{
		{
			name:     "Comando simples",
			cmd:      "git pull",
			expected: []string{"git", "pull"},
		},
		{
			name:     "Comando com aspas",
			cmd:      `docker run -d --name "my container"`,
			expected: []string{"docker", "run", "-d", "--name", "my container"},
		},
		{
			name:     "Comando vazio",
			cmd:      "",
			expected: []string{},
		},
		{
			name:     "Comando com múltiplos espaços",
			cmd:      "git   pull   origin",
			expected: []string{"git", "pull", "origin"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseCommand(tt.cmd)
			if len(result) != len(tt.expected) {
				t.Errorf("esperado %d partes, obtido %d", len(tt.expected), len(result))
				return
			}
			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("parte %d: esperado '%s', obtido '%s'", i, expected, result[i])
				}
			}
		})
	}
}
