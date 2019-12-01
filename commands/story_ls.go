package commands

import (
	"fmt"
	"strconv"

	"github.com/abuan/gitus/db"
	"github.com/spf13/cobra"
)

// Affiche la liste des US Gitus dans la CLI
func runLsUS(cmd *cobra.Command, args []string) error {
	//Récupération de la liste des Project en BDD

	usList, err := db.TaskGetAllUserStory()
	if err != nil {
		return err
	}

	// Affichage dans la CLI
	fmt.Println("\n*************** User Story List ***************")
	fmt.Println("\tID\tEffort\tName")
	fmt.Println("\t--\t------\t----")
	for _, u := range usList {
		fmt.Println("\t" + strconv.Itoa(u.ID) + "\t" + strconv.Itoa(u.Effort) + "\t" + u.Title)
	}
	return nil
}

// Var Cobra décrivant une commande CLI affichant la liste des US Gitus
var userStoryLsCmd = &cobra.Command{
	Use:      "ls",
	Short:    "Display a list of all the Gitus US.",
	Args:     cobra.NoArgs,
	RunE:     runLsUS,
	PreRunE:  connexionForData,
	PostRunE: deconnexionForData,
}

func init() {
	userStoryCmd.AddCommand(userStoryLsCmd)
}