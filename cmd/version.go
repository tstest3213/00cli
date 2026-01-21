package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Mostra a versão do 00cli",
	Long:  `Mostra a versão atual do 00cli instalada.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("00cli version v0.1.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
