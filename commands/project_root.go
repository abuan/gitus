package commands

import (
	"errors"
	"github.com/spf13/cobra"
)

// Var Cobra décrivant une commande CLI de base pour les actions liées aux project
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Root command for project",
	Long: `Root command for project. 
	
	Start of each command related to a project. Do nothing alone`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("Provide item to the project command")
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)
}
