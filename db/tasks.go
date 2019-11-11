package db

import(
	"log"
	"strconv"
	"github.com/abuan/gitus/userstory"
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
	err := db.Ping()
	if err != nil{
		// Query d'insertion/mise à jour de données, on utilise les fonctions Prepare puis Exec
		statement, _ := db.Prepare("INSERT INTO projet(name) VALUES(?)")
		statement.Exec("Test")
		// Query de sélection de plusieurs éléments, on récupère un "Rows" contenant tous les résultat, on utilise la fonction Query
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

// TaskAddUserStory : Ajoute une nouvelle UserStory à la BDD
func TaskAddUserStory(u *userstory.UserStory)error{
	// Vérifie que la connection est bien établie
	err := db.Ping()
	if err != nil{
		return err
		}
	stmt, err := db.Prepare("INSERT INTO UserStory(name,descript,effort,creation_date) VALUES(?,?,?,?)")
	if err != nil{
		return err
		}
	defer stmt.Close()
	_, err = stmt.Exec(u.Name,u.Description,u.Effort,u.CreationDate)
	if err != nil{
		return err
	}
	return nil
}

// TaskGetUserStory : Récupère une UserStory dans la BDD à partir de son ID
func TaskGetUserStory(id int)(*userstory.UserStory,error){
	err := db.Ping()
	if err != nil{
		return nil,err
		}
	// Query de sélection d'un élément, on récupère un "Row" contenant tous les résultat, on utilise la fonction QueryRow
	stmt, err := db.Prepare("SELECT name,descript,creation_date,effort FROM UserStory WHERE id = ?")
	if err != nil{
		return nil,err
		}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	//Utilisation des Keys de la structure pour indiquer le champs que l'on affecte
	u:= userstory.UserStory{ID:id}
	//Scan des résultat du Row
	err = row.Scan(&u.Name, &u.Description,&u.CreationDate,&u.Effort)
	if err != nil{
		return nil,err
	}
	return &u,err
}

// TaskUpdateUserStory : Met à jour une UserStory en BDD
func TaskUpdateUserStory(u*userstory.UserStory)error{
	err := db.Ping()
	if err != nil{
		return err
		}
	stmt, err := db.Prepare("UPDATE UserStory SET name = ?,descript = ?,effort = ? WHERE id = ?")
	if err != nil{
		return err
		}
	defer stmt.Close()
	_, err = stmt.Exec(u.Name,u.Description,u.Effort,u.ID)
	if err != nil{
		return err
	}
	return err
}