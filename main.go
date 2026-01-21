package main

import (
	"fmt"
	"os"

	"github.com/tstest3213/00cli/cmd"
)

var (
	version = "v0.1.0" // Será sobrescrito no build com -ldflags
)

func main() {
	// Verificar atualizações em background (apenas se não for comando update)
	if len(os.Args) > 1 && os.Args[1] != "update" {
		go cmd.CheckForUpdates(version)
	}

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
		os.Exit(1)
	}
}
