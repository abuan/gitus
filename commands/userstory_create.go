package commands

import(
	"github.com/abuan/gitus/userstory"
	"github.com/spf13/cobra"
	"github.com/abuan/gitus/db"
)

var (
	addEffort     int
)

func runCreateUS(cmd *cobra.Command, args []string) error {
	us := userstory.NewUserStory(args[0],"",addEffort)

	// Effectation description
	var s string
	for i := 1; i < len(args); i++ {
		s += args[i] + " "
	}
	us.SetDescription(s)

	//Sauvegarde en BDD
	err := db.InitDB("abuan")
	defer db.CloseDB()
	if err != nil{
		return err
	}
	err = db.TaskAddUserStory(&us)
	if err != nil{
		return err
	}
	return nil
}

var userStroryCreateCmd = &cobra.Command{
	Use:     "create [<name>] <description>[...]",
	Short:   "Create a new UserStory.",
	Args:	 cobra.MinimumNArgs(1),
	RunE:    runCreateUS,
}

func init() {
	userStroryCmd.AddCommand(userStroryCreateCmd)

	userStroryCreateCmd.Flags().SortFlags = false

	userStroryCreateCmd.Flags().IntVarP(&addEffort, "effort", "e", 0,
		"Provide an effort to the User Story",
	)
}