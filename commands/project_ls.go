package commands

import (
	"fmt"
	"github.com/abuan/gitus/db"
	"github.com/spf13/cobra"
	"strconv"
)

// Affiche la liste des Project Gitus dans la CLI
func runLsProjects(cmd *cobra.Command, args []string) error {
	//Récupération de la liste des Project en BDD

	pList, err := db.TaskGetAllProjects()
	if err != nil {
		return err
	}

	// Affichage dans la CLI
	fmt.Println("\n*************** Project List ***************")
	fmt.Println("\tID\tName")
	fmt.Println("\t--\t----")
	for _, p := range pList {
		fmt.Println("\t" + strconv.Itoa(p.ID) + "\t" + p.Name)
	}
	return nil
}

// Var Cobra décrivant une commande CLI affichant la liste des projets Gitus
var projectLsCmd = &cobra.Command{
	Use:      "ls",
	Short:    "Display a list of all the Gitus Project.",
	Args:     cobra.NoArgs,
	PreRunE:  connexionForData,
	RunE:     runLsProjects,
	PostRunE: deconnexionForData,
}

func init() {
	projectCmd.AddCommand(projectLsCmd)
}
