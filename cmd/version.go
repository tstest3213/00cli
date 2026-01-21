package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Mostra a versão do 00cli",
	Long:  `Mostra a versão atual do 00cli instalada.`,
	Run: func(cmd *cobra.Command, args []string) {
		version := getVersion()
		fmt.Printf("00cli version %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

// getVersion obtém a versão do binário
func getVersion() string {
	// Tentar ler da variável de ambiente (setada no build)
	if v := os.Getenv("00CLI_VERSION"); v != "" {
		return v
	}

	// Tentar executar o binário com --version
	exe, err := os.Executable()
	if err != nil {
		return "v0.0.0"
	}

	cmd := exec.Command(exe, "version")
	output, err := cmd.Output()
	if err != nil {
		return "v0.0.0"
	}

	// Extrair versão da saída
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "version") {
			parts := strings.Fields(line)
			for _, part := range parts {
				if strings.HasPrefix(part, "v") {
					return part
				}
			}
		}
	}

	return "v0.0.0"
}
