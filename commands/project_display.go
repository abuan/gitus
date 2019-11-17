package commands

import(
	"github.com/spf13/cobra"
	"github.com/abuan/gitus/db"
	"strconv"
)

// Affiche le contenu d'un Project dans la CLI
func runDisplayProject(cmd *cobra.Command, args []string) error{
	//Récupération du projet et de la liste des US liées en BDD via son ID
	err := db.InitDB()
	defer db.CloseDB()
	if err != nil{
		return err
	}
	id,_ := strconv.Atoi(args[0])
	p,ids,err:= db.TaskGetProject(id);
	if err != nil{
		return err
	}

	//Affichage du projet
	p.Display(ids)
	
	return nil
}

// Var Cobra décrivant une commande CLI affichant la composition d'un projet
var projectDisplayCmd = &cobra.Command{
	Use:     "display [<id>]",
	Short:   "Display a Project content from its Id.",
	Args:	 cobra.MinimumNArgs(1),
	RunE:    runDisplayProject,
}

func init() {
	projectCmd.AddCommand(projectDisplayCmd)
}