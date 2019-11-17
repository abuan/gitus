package utils

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