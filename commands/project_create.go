package commands

import (
	"errors"


	"fmt"
	"strconv"

	"github.com/abuan/gitus/db"
	"github.com/abuan/gitus/project"
	"github.com/abuan/gitus/utils"
	"github.com/spf13/cobra"
)

// Variable passée au flag
var (
	// COntient la liste des US à lier avec le projet
	usList []int
)

//Fonction créant une userstory à partir des arguments de la CLI
func runCreateProject(cmd *cobra.Command, args []string) error {
	//Création du projet
	p := project.NewProject(args[0], "", addAuthor)

	// Affectation description
	if len(args) > 1 {
		var s string
		for i := 1; i < len(args); i++ {
			s += args[i] + " "
		}
		p.Description = s
	}

	//Vérification des US passées en argument
	usToLink := len(usList) > 0
	if usToLink {
		//Suppression des doublons
		usList = utils.RemoveDuplicates(usList)

		//Sélection de tous les IDS des US
		ids, err := db.TaskGetAllUserStoryID()
		if err != nil {
			return err
		}
		// Vérifie si les ids passés en argument éxistent en BDD
		exists, unknownValues := utils.Contains(ids, usList)
		//Si des valeurs d'IDs n'éxistent pas alors on informe l'utilisateur
		if !exists {
			for _, val := range unknownValues {
				fmt.Print(strconv.Itoa(val) + " / ")
			}
			return errors.New("Les User Story liées aux IDs précédents n'existent pas")
		}
	}

	//Sauvegarde en BDD du projet
	projectID, err := db.TaskAddProject(&p)
	if err != nil {
		return err
	}
	if usToLink {
		//Ajoute le liens entre les US et le projet
		// Ce lien est fait via une table de jointure. Voir la méthodologie "Many to Many" liée à la gestion de BDD
		err = db.TaskLinkUsToProject(usList, projectID)
		if err != nil {
			return err
		}
	}

	return nil
}

// Var Cobra décrivant une commande CLI créant une UserStory
var projectCreateCmd = &cobra.Command{
	Use:      "create [<name>] <description>[...]",
	Short:    "Create a new project.",
	Args:     cobra.MinimumNArgs(1),
	PreRunE:  connexionForData,
	RunE:     runCreateProject,
	PostRunE: deconnexionForData,
}

func init() {
	projectCmd.AddCommand(projectCreateCmd)

	projectCreateCmd.Flags().SortFlags = false

	//Ajout du flag permettant le liage de US
	projectCreateCmd.Flags().IntSliceVarP(&usList, "userStories", "u", nil,
		"Provide a list of User Story's ID to be linked to the project. Example : 1,3,8",
	)
	projectCreateCmd.Flags().StringVarP(&addAuthor, "author", "a", "unknow",
		"Provide an author to the Project",
	)
}
