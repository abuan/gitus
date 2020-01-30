package commands

import (
	"github.com/abuan/gitus/db"
	"errors"
	"fmt"
	"strconv"
	//"github.com/abuan/gitus/userstory"
	"github.com/abuan/gitus/utils"
	"github.com/spf13/cobra"
)

var(
	projectList []string
)

//Fonction créant une userstory à partir des arguments de la CLI
func runLinkStory(cmd *cobra.Command, args []string) error {
	//Check si la liste passée n'est pas vide retourner un message d'erreur le cas échéant 
	
	//Verification liste non vide
	projectToBeLinked := len(projectList)>0

	if projectToBeLinked == false {
		return errors.New("Provide project's name to be linked")
	}

	//Suppresion des doublons
	projectList = utils.RemoveStringDuplicates(projectList)

	//Sélection de tous les projets en BDD
	projects, err := db.TaskGetAllProjects()
	if err != nil {
		return err
	}
	//Vérifies que tous les projets passés dans la liste existent
	//Garde parmis la liste de tous les projets seulement ceux souhaités
	exists, unknownValues,projects := utils.VerifyProjectNames(projects, projectList)
	if !exists {
		for _, val := range unknownValues {
			fmt.Print(val + " / ")
		}
		return errors.New("Les noms de projet Gitus précédents n'existent pas")
	}
	//Création tableau de IDs des projets
	var projectsIDs [] int
	for _,val := range projects{
		projectsIDs = append(projectsIDs,val.ID)
	}

	//Ajout des liens
	idStory, _ := strconv.Atoi(args[0])
	err = db.TaskLinkProjectsToUS(idStory, projectsIDs)
		if err != nil {
			return err
		}
	return nil
}

// Var Cobra décrivant une commande CLI créant une UserStory
var userStoryLinkCmd = &cobra.Command{
	Use:      "link [ID]",
	Short:    "Link a Story to one or multiple project",
	Args:     cobra.MinimumNArgs(1),
	RunE:     runLinkStory,
	PreRunE:  connexionForData,
	PostRunE: deconnexionForData,
}

func init() {
	userStoryCmd.AddCommand(userStoryLinkCmd)

	userStoryLinkCmd.Flags().SortFlags = false

	//Ajout du flag permettant de link la story à des projets
	userStoryLinkCmd.Flags().StringSliceVarP(&projectList, "projects", "p", nil,
		"Provide a list of Project's name to be linked to the project. Example : \"toto\",\"foo\",\"bar\"",
	)

}
