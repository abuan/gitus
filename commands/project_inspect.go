package commands

import (
	"github.com/abuan/gitus/db"
	"github.com/spf13/cobra"
	"strconv"
)

// Affiche le contenu d'un Project dans la CLI
func runInspectProject(cmd *cobra.Command, args []string) error {
	//Récupération du projet et de la liste des US liées en BDD via son ID
	id, _ := strconv.Atoi(args[0])
	p, ids, err := db.TaskGetProject(id)
	if err != nil {
		return err
	}

	//Affichage du projet
	p.Display(ids)

	return nil
}

// Var Cobra décrivant une commande CLI affichant la composition d'un projet
var projectInspectCmd = &cobra.Command{
	Use:      "inspect [<id>]",
	Short:    "Display a Project content from its Id.",
	Args:     cobra.MinimumNArgs(1),
	PreRunE:  connexionForData,
	RunE:     runInspectProject,
	PostRunE: deconnexionForData,
}

func init() {
	projectCmd.AddCommand(projectInspectCmd)
}
