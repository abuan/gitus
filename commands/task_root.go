package commands

import (
	"errors"
	"github.com/spf13/cobra"
)

// Var Cobra décrivant une commande CLI de base pour les actions liées aux tâches
var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Root command for task",
	Long: `Root command for task. 
	
	Start of each command related to a task. Do nothing alone`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("Provide item to the task command")
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)
}
