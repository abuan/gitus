package main

import "fmt"
import "task"

//Variable globale contenant le logine pour la connection avec la BDD MySQL
// Solution temporaire, il faudra penser à mettre au point un système de login avec un fichier de config
//Changer la valeur pour correspondre à votre login MySQL
//Si vous avez un mdp il faut utiliser le format suivant : "username:password"
var login = "abuan"

func main() {
/* 	u0 := newUserStory("description1")
	fmt.Println(u0)
	u0.addTache("description de la tache 1")
	fmt.Println(u0)
	u0.addUserStory("description user story ajouté")
	fmt.Println(u0) */
	
	//Connexion avec la BDD
	bdd.InitDB(login)
	//Ferme la connexion une fois la fonction "main" terminée
	defer bdd.CloseDB()
	// Début du code de l'application
	bdd.TaskTestDB()
}
