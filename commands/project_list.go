package commands

import (
	"fmt"
	"strconv"

	"github.com/abuan/gitus/db"
	"github.com/spf13/cobra"
)

// Affiche la liste des Project Gitus dans la CLI
func runListProjects(cmd *cobra.Command, args []string) error {
	//Récupération de la liste des Project en BDD
	err := db.InitDB()
	defer db.CloseDB()
	if err != nil {
		return err
	}
	pList, err := db.TaskGetAllProjects()
	if err != nil {
		return err
	}

	// Affichage dans la CLI
	fmt.Println("\n*************** Project List ***************")
	fmt.Println("\tID\tName\tAuthor")
	fmt.Println("\t--\t----")
	for _, p := range pList {
		fmt.Println("\t" + strconv.Itoa(p.ID) + "\t" + p.Name + "\t" + p.Author)
	}
	return nil
}

// Var Cobra décrivant une commande CLI affichant la liste des projets Gitus
var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "Display a list of all the Gitus Project.",
	Args:  cobra.NoArgs,
	RunE:  runListProjects,
}

func init() {
	projectCmd.AddCommand(projectListCmd)
}
