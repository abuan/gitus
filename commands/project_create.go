package commands

import(
	"github.com/abuan/gitus/project"
	"github.com/spf13/cobra"
	"github.com/abuan/gitus/db"
	"github.com/abuan/gitus/utils"
	"fmt"
	"errors"
	"strconv"
)

// Variable passée au flag 
var (
	// COntient la liste des US à lier avec le projet
	usList			[]int
)

//Fonction créant une userstory à partir des arguments de la CLI
func runCreateProject(cmd *cobra.Command, args []string) error {
	//Création du projet
	p := project.NewProject(args[0],"")

	// Affectation description
	if len(args)>1{
		var s string
		for i := 1; i < len(args); i++ {
			s += args[i] + " "
		}
		p.Description = s
	}
	//Connection avec la DB
	err := db.InitDB()
	defer db.CloseDB()
	if err != nil{
		return err
	}
	//Vérification des US passées en argument
	//Suppression des doublons
	usList = utils.RemoveDuplicates(usList)

	//Sélection de tous les IDS des US
	ids,err := db.TaskGetAllUserStoryID()
	if err != nil{
		return err
	}

	// Vérifie si les ids passés en argument éxistent en BDD
	exists,unknownValues :=utils.Contains(ids,usList)
	//Si des valeurs d'IDs n'éxistent pas alors on informe l'utilisateur
	if !exists{
		for _,val := range unknownValues{
			fmt.Print(strconv.Itoa(val) + " ; ")
		}
		return errors.New("Les User Story liées aux IDs Prédent n'éxistent pas")
	}
	//Sauvegarde en BDD du projet
	projectID, err := db.TaskAddProject(&p)
	if err != nil{
		return err
	}
	//Ajoute le liens entre les US et le projet
	// Ce lien est fait via une table de jointure. Voir la méthodologie "Many to Many" liée à la gestion de BDD
	err = db.TaskLinkUsToProject(usList,projectID)
	if err != nil{
		return err
	}

	return nil
}

// Var Cobra décrivant une commande CLI créant une UserStory
var projectCreateCmd = &cobra.Command{
	Use:     "create [<name>] <description>[...]",
	Short:   "Create a new project.",
	Args:	 cobra.MinimumNArgs(1),
	RunE:    runCreateProject,
}

func init() {
	projectCmd.AddCommand(projectCreateCmd)

	projectCreateCmd.Flags().SortFlags = false

	//Ajout du flag permettant le liage de US
	projectCreateCmd.Flags().IntSliceVarP(&usList, "userStories", "u", nil,
		"Provide a list of User Story's ID to be linked to the project. Example : 1,3,8",
	)
}