package commands

import(
	"github.com/spf13/cobra"
	"github.com/abuan/gitus/db"
	"strconv"
)

// Affiche le contenu d'une US dans la CLI
func runDisplayUS(cmd *cobra.Command, args []string) error{
	//Récupération de la US en BDD via ID
	err := db.InitDB()
	defer db.CloseDB()
	if err != nil{
		return err
	}
	id,_ := strconv.Atoi(args[0])
	u,err:= db.TaskGetUserStory(id);
	if err != nil{
		return err
	}

	//Affichage de la US
	u.Display()
	
	return nil
}

// Var Cobra décrivant une commande CLI modifiant une UserStory
var userStoryDisplayCmd = &cobra.Command{
	Use:     "display [<id>]",
	Short:   "Display a UserStory content from its Id.",
	Args:	 cobra.MinimumNArgs(1),
	RunE:    runDisplayUS,
}

func init() {
	userStoryCmd.AddCommand(userStoryDisplayCmd)
}