package db

import(
	"log"
	"strconv"
	"strings"
	"github.com/abuan/gitus/userstory"
	"github.com/abuan/gitus/project"
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

// TaskGetUserStoryList : Récupère une list de UserStory dans la BDD à partir de leur IDs
func TaskGetUserStoryList(ids []int)([]*userstory.UserStory,error){
	//Passage de la slice en "empty interface"
	args := make([]interface{}, len(ids))
	for i, id := range ids {
    	args[i] = id
	}

	//Création de la slice des US
	usList := make([]*userstory.UserStory,len(ids))
	// Preparation de la query pour récupérer toutes les US matchant les ids en param
	stmt, err := db.Prepare("SELECT name,descript,creation_date,effort FROM UserStory WHERE id in(?"+ strings.Repeat(",?",len(args)-1) + ")")
	if err != nil{
		return nil,err
		}
	defer stmt.Close()
	
	//Exécution de la query
	rows,err := stmt.Query(args...)
	if err != nil{
		return nil,err
	}

	//Variable d'incrémention pour accès ids
	i := 0
	for rows.Next(){
		//Création d'une US
		u:= userstory.UserStory{ID:ids[i]}
		//Scan des résultat du Row
		err = rows.Scan(&u.Name, &u.Description,&u.CreationDate,&u.Effort)
		if err != nil{
			return nil,err
		}
		usList = append(usList, &u)
		i++
	}

	return usList,err
}

// TaskGetAllUserStoryID : Récupère l'ensemble des IDs des UserStory existantes
func TaskGetAllUserStoryID()([]int,error){
	// Query de sélection d'un élément, on récupère un "Row" contenant tous les résultat, on utilise la fonction QueryRow
	stmt, err := db.Prepare("SELECT id FROM UserStory")
	if err != nil{
		return nil,err
		}
	defer stmt.Close()
	rows,err := stmt.Query()
	if err != nil{
		return nil,err
	}
	var (
		ids [] int
		id int
	)
	for rows.Next(){
		err = rows.Scan(&id)
		if err != nil{
			return nil,err
		}
		ids = append(ids, id)
	}
	return ids,err
}

// TaskAddProject : Ajoute un nouveau projet à la BDD
func TaskAddProject(p *project.Project)(int,error){
	// Insert le projet dans la BDD
	stmt, err := db.Prepare("INSERT INTO Project(name,descript,creation_date) VALUES(?,?,?)")
	if err != nil{
		return 0,err
	}
	defer stmt.Close()
	_, err = stmt.Exec(p.Name,p.Description,p.CreationDate)
	if err != nil{
		return 0,err
	}

	// Récupération de l'ID généré par Mysql
	stmt, err = db.Prepare("SELECT id FROM Project WHERE Project.name = ? AND Project.descript = ?")
	if err != nil{
		return 0,err
	}
	row := stmt.QueryRow(p.Name,p.Description)
	var id int
	err =row.Scan(&id)
	if err != nil{
		return 0,err
	}
	return id,nil
}

// TaskLinkUsToProject : Lie des US story à un projet
func TaskLinkUsToProject(usList []int, projectID int)error{
	// Query string to be completed
	sqlStr := "INSERT INTO Project_structure(project_id,userstory_id) VALUES "
	// Création et remplissage du tableau des couples (project_id,userstory_id)
	vals := []interface{}{}
	for _, usID := range usList {
		sqlStr += "(?,?),"
		vals = append(vals, projectID, usID)
	}
	//Supprime la dernière ,
	sqlStr = strings.TrimSuffix(sqlStr, ",")
	// Insert dans la table l'ensemble des couples de valeurs
	stmt, err := db.Prepare(sqlStr)
	if err != nil{
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(vals...)
	if err != nil{
		return err
	}
	return nil
}