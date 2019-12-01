package commands

import (
	"github.com/abuan/gitus/db"
	"github.com/spf13/cobra"
	"strconv"
)

// Affiche le contenu d'une US dans la CLI
func runInspectUS(cmd *cobra.Command, args []string) error {
	//Récupération de la US en BDD via ID
	id, _ := strconv.Atoi(args[0])
	u, err := db.TaskGetUserStory(id)
	if err != nil {
		return err
	}

	//Affichage de la US
	u.Display()

	return nil
}

// Var Cobra décrivant une commande CLI modifiant une UserStory
var userStoryInspectCmd = &cobra.Command{
	Use:      "inspect [<id>]",
	Short:    "Display a UserStory content from its Id.",
	Args:     cobra.MinimumNArgs(1),
	RunE:     runInspectUS,
	PreRunE:  connexionForData,
	PostRunE: deconnexionForData,
}

func init() {
	userStoryCmd.AddCommand(userStoryInspectCmd)
}
