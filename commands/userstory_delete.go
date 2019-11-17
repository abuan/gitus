package commands

import(
	"github.com/spf13/cobra"
	"github.com/abuan/gitus/db"
	"strconv"
)

//Fonction supprimant une userstory à partir de l'ID passé dans la CLI
func runDeleteUS(cmd *cobra.Command, args []string) error {
	//Suppression de la US en BDD via ID
	err := db.InitDB()
	defer db.CloseDB()
	if err != nil{
		return err
	}
	id,_ := strconv.Atoi(args[0])
	err = db.TaskDeleteUserStory(id);
	if err != nil{
		return err
	}
	return nil
}

// Var Cobra décrivant une commande CLI supprimant une UserStory
var userStroryDeleteCmd = &cobra.Command{
	Use:     "delete [<id>]",
	Short:   "Delete a UserStory from its Id.",
	Args:	 cobra.MinimumNArgs(1),
	RunE:    runDeleteUS,
}

func init() {
	userStroryCmd.AddCommand(userStroryDeleteCmd)
}