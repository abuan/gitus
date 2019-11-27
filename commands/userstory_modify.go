package commands

import(
	"github.com/spf13/cobra"
	"github.com/abuan/gitus/db"
	"strconv"
)

// Variable passé au flag de la commande Cobra stockant les valeurs attribuées
var (
	newEffort     int
	newStoryDescription string
)

//Fonction modifiant une userstory à partir des arguments de la CLI
func runModifyUS(cmd *cobra.Command, args []string) error {
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

	// Affectation description
	if newStoryDescription != "" {
		u.Description = newStoryDescription
	}

	// Mise à jour des champs selon flags
	if newEffort != -1{
		u.Effort = newEffort
	}

	//Update en BDD de la US
	err = db.TaskUpdateUserStory(u)
	if err != nil{
		return err
	}

	return nil
}

// Var Cobra décrivant une commande CLI modifiant une UserStory
var userStoryModifyCmd = &cobra.Command{
	Use:     "modify [<id>]",
	Short:   "Modify a UserStory from its Id.",
	Args:	 cobra.MinimumNArgs(1),
	RunE:    runModifyUS,
}

func init() {
	userStoryCmd.AddCommand(userStoryModifyCmd)

	userStoryModifyCmd.Flags().SortFlags = false

	userStoryModifyCmd.Flags().IntVarP(&newEffort, "effort", "e", -1,
		"Provide a new effort value to the User Story",
	)
	userStoryModifyCmd.Flags().StringVarP(&newStoryDescription, "description", "d", "",
	"Provide a new description to the User Story\nDescription have to be between quotation marks.\nExample : \"My Beautiful description\"",
	)
}