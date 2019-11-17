package commands

import(
	"github.com/spf13/cobra"
	"github.com/abuan/gitus/db"
	"strconv"
)

//Fonction supprimant un projet à partir de l'ID passé dans la CLI
func runDeleteProject(cmd *cobra.Command, args []string) error {
	//Suppression de la US en BDD via ID
	err := db.InitDB()
	defer db.CloseDB()
	if err != nil{
		return err
	}
	id,_ := strconv.Atoi(args[0])
	err = db.TaskDeleteProject(id);
	if err != nil{
		return err
	}
	return nil
}

// Var Cobra décrivant une commande CLI supprimant un projet
var projectDeleteCmd = &cobra.Command{
	Use:     "delete [<id>]",
	Short:   "Delete a Project from its Id.",
	Args:	 cobra.MinimumNArgs(1),
	RunE:    runDeleteProject,
}

func init() {
	projectCmd.AddCommand(projectDeleteCmd)
}