package utils
import(
	"github.com/abuan/gitus/project"
)

//RemoveDuplicates : supprime les doublons d'une Slice d'entiers
func RemoveDuplicates(intSlice []int) []int {
	//MAp contenant une valeur et un booléen idiquant si il est déjà apparut
	keys := make(map[int]bool)
	// Liste contenant les valeurs sans les doublons
    list := []int{} 
    for _, entry := range intSlice {
		// On récupère dans la map le booléen associé à la valeur courante de la slice
		// Si il est à false ( par défaut si il n'est pas dans la liste) alors on le set à true et rentre dans liste de sortie
		// Si il est à true on ne fait rien
        if _, value := keys[entry]; !value {
            keys[entry] = true
            list = append(list, entry)
        }
    }    
    return list
}

// Contains : Vérifie si `slice` contient tous les éléments de `values`. Retourne false et les valeurs non trouvées le cas échéant
func Contains(slice ,values []int)(bool,[]int){
	keys := make(map[int]bool)
	list := []int{} 
	for _, entry := range slice{
		keys[entry] = true
	}
	contains := true
	for _,entry := range values{
		if _, value := keys[entry]; !value {
            contains = false
            list = append(list, entry)
        }
	}

	return contains, list
}

//RemoveStringDuplicates : supprime les doublons d'une Slice d'entiers
func RemoveStringDuplicates(intSlice []string) []string {
	//MAp contenant une valeur et un booléen idiquant si il est déjà apparut
	keys := make(map[string]bool)
	// Liste contenant les valeurs sans les doublons
    list := []string{} 
    for _, entry := range intSlice {
		// On récupère dans la map le booléen associé à la valeur courante de la slice
		// Si il est à false ( par défaut si il n'est pas dans la liste) alors on le set à true et rentre dans liste de sortie
		// Si il est à true on ne fait rien
        if _, value := keys[entry]; !value {
            keys[entry] = true
            list = append(list, entry)
        }
    }    
    return list
}

// VerifyProjectNames : Vérifie si `slice` contient tous les éléments de `values`. Retourne false et les valeurs non trouvées le cas échéant
func VerifyProjectNames(projects []*project.Project, values []string)(bool , []string,[]*project.Project) {
	//Map contenant tous les noms de projets pour vérifier leurs existance
	keys := make(map[string]bool)
	//Map contenant les noms des projets existant
	link := make(map[string]bool)
	//List des noms de projets qui n'existent pas
	list := []string{} 
	//On remplit la map avec le nom de tous les projets
	for _, entry := range projects{
		keys[entry.Name] = true
	}
	//Verifie que tous les noms passés en argument existent
	contains := true
	for _,entry := range values{
		if _, value := keys[entry]; !value {
            contains = false
			list = append(list, entry)
        }else{
			//On ajoute dans une map le nom du projet
			link[entry]=true
		}
	}
	//Si liste contains == false certains nom passés en param n'éxistent pas
	if !contains {
		return contains, list, nil
	}
	//Parcours slice projets, si nom présent dans values ont garde sinon on supprime
	for i, entry := range projects{
		if _, value := keys[entry.Name]; !value {
            if i < len(projects)-1 {
				copy(projects[i:], projects[i+1:])
			  }
			  projects[len(projects)-1] = nil // or the zero value of T
			  projects = projects[:len(projects)-1]
        }
	}
	return contains, list,projects
}