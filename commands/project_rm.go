package commands

import (
	"github.com/abuan/gitus/db"
	"github.com/spf13/cobra"
	"strconv"
)

//Fonction supprimant un projet à partir de l'ID passé dans la CLI
func runRmProject(cmd *cobra.Command, args []string) error {
	//Suppression de la US en BDD via ID

	id, _ := strconv.Atoi(args[0])
	err := db.TaskDeleteProject(id)
	if err != nil {
		return err
	}
	return nil
}

// Var Cobra décrivant une commande CLI supprimant un projet
var projectRmCmd = &cobra.Command{
	Use:      "rm [<id>]",
	Short:    "Delete a Project from its Id.",
	Args:     cobra.MinimumNArgs(1),
	PreRunE:  connexionForData,
	RunE:     runRmProject,
	PostRunE: deconnexionForData,
}

func init() {
	projectCmd.AddCommand(projectRmCmd)
}
