package main

import (
	"fmt"
	"os"

	"github.com/tstest3213/00cli/cmd"
)

var version = "v0.1.0"

func main() {
	// Verificar atualizações em background
	go cmd.CheckForUpdates(version)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
		os.Exit(1)
	}
}
