package main

import (
	"github.com/abuan/proto_gitus/db"
	_ "github.com/abuan/proto_gitus/task"
	_ "github.com/abuan/proto_gitus/userstory"
)

//Variable globale contenant le logine pour la connection avec la BDD MySQL
// Solution temporaire, il faudra penser à mettre au point un système de login avec un fichier de config
//Changer la valeur pour correspondre à votre login MySQL
//Si vous avez un mdp il faut utiliser le format suivant : "username:password"
var login = "abuan"

func main() {

	//Connexion avec la BDD
	db.InitDB(login)
	//Ferme la connexion une fois la fonction "main" terminée
	defer db.CloseDB()
	// Début du code de l'application
	db.TaskTestDB()
}