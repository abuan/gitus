package db

import(
	"log"
	"strconv"
)

/*Ce fichier contient l'ensemble des taches liées à la base de données.
* Pour chaque interraction spécifique avec la base comme par exemple récupérer l'ensemble des
* projets on va créer une fonction dédiée, dans notre cas "TaskGetProjects", contenant l'ensemble 
* des instructions SQL permettant des récupérer ces informations.
* Lire le tutoriel suivant pour les différentes indications sur comment réaliser une Query et
* quel type de Query utiliser : http://go-database-sql.org/index.html
*/

// TaskTestDB : Fonction de test pour la BDD
//Sert d'exemple dans un premier temps à supprimer par la suite
func TaskTestDB(){
	if db != nil{
		// Query d'insertion/mise à jour de données, on utilise les fonctions Prepare puis Exec
		statement, _ := db.Prepare("INSERT INTO projet(name) VALUES(?)")
		statement.Exec("Test")
		// Query de sélection d'éléments, on récupère un "row" contenant tous les résultat, on utilise la fonction Query
		rows, _ := db.Query("SELECT id, name FROM projet")
		var id int
		var name string
		// On parcourt l'ensemble des résultat de du row puis on les traite
		for rows.Next() {
			rows.Scan(&id, &name)
			log.Println(strconv.Itoa(id) + " : " + name)
		}
	}else{
		log.Println("nil db")
	}
	
}
