package commands

import(
	"github.com/spf13/cobra"
	"github.com/abuan/gitus/db"
	"fmt"
	"strconv"
)

// Affiche la liste des US Gitus dans la CLI
func runListUS(cmd *cobra.Command, args []string) error{
	//Récupération de la liste des Project en BDD
	err := db.InitDB()
	defer db.CloseDB()
	if err != nil{
		return err
	}
	usList,err:= db.TaskGetAllUserStory();
	if err != nil{
		return err
	}

	// Affichage dans la CLI
	fmt.Println("\n*************** User Story List ***************")
	fmt.Println("\tID\tEffort\tTitle")
	fmt.Println("\t--\t------\t----")
	for _,u:= range usList{
		fmt.Println("\t"+strconv.Itoa(u.ID)+"\t"+strconv.Itoa(u.Effort)+"\t"+u.Title)
	}
	return nil
}

// Var Cobra décrivant une commande CLI affichant la liste des US Gitus
var userStoryListCmd = &cobra.Command{
	Use:     "list",
	Short:   "Display a list of all the Gitus US.",
	Args:	 cobra.NoArgs,
	RunE:    runListUS,
}

func init() {
	userStoryCmd.AddCommand(userStoryListCmd)
}