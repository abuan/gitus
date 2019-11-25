package commands

import (
	"github.com/abuan/gitus/db"
	"github.com/abuan/gitus/userstory"
	"github.com/spf13/cobra"
)

// Variable passé au flag --effort stockant la valeur de l'effort attribué
var (
	addEffort int
	addAuthor string
)

//Fonction créant une userstory à partir des arguments de la CLI
func runCreateUS(cmd *cobra.Command, args []string) error {
	us := userstory.NewUserStory(args[0], "", addAuthor, addEffort)

	// Affectation description
	var s string
	for i := 1; i < len(args); i++ {
		s += args[i] + " "
	}
	us.SetDescription(s)

	//Sauvegarde en BDD
	err := db.InitDB()
	defer db.CloseDB()
	if err != nil {
		return err
	}
	err = db.TaskAddUserStory(&us)
	if err != nil {
		return err
	}
	return nil
}

// Var Cobra décrivant une commande CLI créant une UserStory
var userStoryCreateCmd = &cobra.Command{
	Use:   "create [<name>] <description>[...]",
	Short: "Create a new UserStory.",
	Args:  cobra.MinimumNArgs(2),
	RunE:  runCreateUS,
}

func init() {
	userStoryCmd.AddCommand(userStoryCreateCmd)

	userStoryCreateCmd.Flags().SortFlags = false

	userStoryCreateCmd.Flags().IntVarP(&addEffort, "effort", "e", 0,
		"Provide an effort to the User Story",
	)
	userStoryCreateCmd.Flags().StringVarP(&addAuthor, "author", "a", "unknow",
		"Provide an author to the User Story",
	)
}
