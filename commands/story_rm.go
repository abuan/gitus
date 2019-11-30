package commands

import (
	"github.com/abuan/gitus/db"
	"github.com/spf13/cobra"
	"strconv"
)

//Fonction supprimant une userstory à partir de l'ID passé dans la CLI
func runRmUS(cmd *cobra.Command, args []string) error {
	//Suppression de la US en BDD via ID

	id, _ := strconv.Atoi(args[0])
	err := db.TaskDeleteUserStory(id)
	if err != nil {
		return err
	}
	return nil
}

// Var Cobra décrivant une commande CLI supprimant une UserStory
var userStoryRmCmd = &cobra.Command{
	Use:      "rm [<id>]",
	Short:    "Delete a UserStory from its Id.",
	Args:     cobra.MinimumNArgs(1),
	RunE:     runRmUS,
	PreRunE:  connexionForData,
	PostRunE: deconnexionForData,
}

func init() {
	userStoryCmd.AddCommand(userStoryRmCmd)
}
