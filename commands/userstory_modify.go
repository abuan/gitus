package commands

import(
	"github.com/spf13/cobra"
	"github.com/abuan/gitus/db"
	"strconv"
)

// Variable passé au flag de la commande Cobra stockant les valeurs attribuées
var (
	newEffort     int
	newName		  string
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
	if len(args)>1 {
		var s string
		for i := 1; i < len(args); i++ {
			s += args[i] + " "
		}
		u.Description = s
	}

	// Mise à jour des champs selon flags
	if newEffort != -1{
		u.Effort = newEffort
	}
	if newName != "" {
		u.Name = newName
	}

	//Update en BDD de la US
	err = db.TaskUpdateUserStory(u)
	if err != nil{
		return err
	}

	return nil
}

// Var Cobra décrivant une commande CLI modifiant une UserStory
var userStroryModifyCmd = &cobra.Command{
	Use:     "modify [<id>] <description>[...]",
	Short:   "Modify a UserStory from its Id.",
	Args:	 cobra.MinimumNArgs(1),
	RunE:    runModifyUS,
}

func init() {
	userStroryCmd.AddCommand(userStroryModifyCmd)

	userStroryModifyCmd.Flags().SortFlags = false

	userStroryModifyCmd.Flags().IntVarP(&newEffort, "effort", "e", -1,
		"Provide a new effort value to the User Story",
	)
	userStroryModifyCmd.Flags().StringVarP(&newName, "name", "n", "",
		"Provide a new name to the User Story",
	)
}