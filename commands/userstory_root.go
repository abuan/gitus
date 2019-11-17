package commands

import(
	"errors"
	"github.com/spf13/cobra"
)

// Var Cobra décrivant une commande CLI de base pour les actions liées aux userStory
var userStroryCmd = &cobra.Command{
	Use:     "userstory",
	Short:   "Root command for user story",
	Long: 	`Root command for user story. 
	
	Start of each command related to a user story. Do nothing alone`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("Provide item to the userstory command")
	},
}

func init() {
	rootCmd.AddCommand(userStroryCmd)
}