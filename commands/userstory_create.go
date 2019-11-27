package commands

import (
	"github.com/abuan/gitus/db"
	"github.com/abuan/gitus/userstory"
	"github.com/spf13/cobra"
)

// Variable passé au flag --effort stockant la valeur de l'effort attribué
var (
	addEffort int
	addDescription string
	//flag autheur en global, qui est aussi le même utiliser pour les project
	addAuthor string
)

//Fonction créant une userstory à partir des arguments de la CLI
func runCreateUS(cmd *cobra.Command, args []string) error {
	us := userstory.NewUserStory(addDescription,addAuthor,addEffort)

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
	Use:     "create",
	Short: "Create a new UserStory.",
	Args:	 cobra.NoArgs,
	RunE:  runCreateUS,
}

func init() {
	userStoryCmd.AddCommand(userStoryCreateCmd)

	userStoryCreateCmd.Flags().SortFlags = false

	userStoryCreateCmd.Flags().IntVarP(&addEffort, "effort", "e", 0,
		"Provide an effort to the User Story",
	)
	userStoryCreateCmd.Flags().StringVarP(&addDescription, "description", "d", "No description so far",
		"Provide a description to the User Story\nDescription have to be between quotation marks.\nExample : \"My Beautiful description\"",
	)
	userStoryCreateCmd.Flags().StringVarP(&addAuthor, "author", "a", "unknow",
		"Provide an author to the User Story",
	)

}
