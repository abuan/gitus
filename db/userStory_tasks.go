package db

import(
	"strings"
	"github.com/abuan/gitus/userstory"
)

// TaskAddUserStory : Ajoute une nouvelle UserStory à la BDD
func TaskAddUserStory(u *userstory.UserStory)error{
	stmt, err := db.Prepare("INSERT INTO UserStory(title,descript,effort,creation_date) VALUES(?,?,?,?)")
	if err != nil{
		return err
		}
	defer stmt.Close()
	_, err = stmt.Exec(u.Title,u.Description,u.Effort,u.CreationDate)
	if err != nil{
		return err
	}
	return nil
}

// TaskGetUserStory : Récupère une UserStory dans la BDD à partir de son ID
func TaskGetUserStory(id int)(*userstory.UserStory,error){
	// Query de sélection d'un élément, on récupère un "Row" contenant tous les résultat, on utilise la fonction QueryRow
	stmt, err := db.Prepare("SELECT title,descript,creation_date,effort FROM UserStory WHERE id = ?")
	if err != nil{
		return nil,err
		}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	//Utilisation des Keys de la structure pour indiquer le champs que l'on affecte
	u:= userstory.UserStory{ID:id}
	//Scan des résultat du Row
	err = row.Scan(&u.Title, &u.Description,&u.CreationDate,&u.Effort)
	if err != nil{
		return nil,err
	}
	return &u,err
}

// TaskUpdateUserStory : Met à jour une UserStory en BDD
func TaskUpdateUserStory(u*userstory.UserStory)error{
	stmt, err := db.Prepare("UPDATE UserStory SET descript = ?,effort = ? WHERE id = ?")
	if err != nil{
		return err
		}
	defer stmt.Close()
	_, err = stmt.Exec(u.Description,u.Effort,u.ID)
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
	stmt, err := db.Prepare("SELECT title,descript,creation_date,effort FROM UserStory WHERE id in(?"+ strings.Repeat(",?",len(args)-1) + ")")
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
		err = rows.Scan(&u.Title, &u.Description,&u.CreationDate,&u.Effort)
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

// TaskDeleteUserStory : Supprime une UserStory dans la BDD à partir de son ID
func TaskDeleteUserStory(id int) error{
	//Suppression des liens dela  US vers les projets
	stmt, err := db.Prepare("DELETE FROM Project_structure WHERE userstory_id = ?")
	if err != nil{
		return err
	}
	defer stmt.Close()
	
	_,err = stmt.Exec(id)
	if err != nil{
		return err
	}	
	
	//Suppression de la US
	stmt, err = db.Prepare("DELETE FROM UserStory WHERE id = ?")
	if err != nil{
		return err
	}
	_,err = stmt.Exec(id)
	if err != nil{
		return err
	}
	
	return nil
}

// TaskGetAllUserStory : Récupère l'ensemble des projets Gitus
func TaskGetAllUserStory()([]*userstory.UserStory,error){
	// Récupère les infos des US dans la BDD
	stmt, err := db.Prepare("SELECT id,title,descript,effort,creation_date FROM UserStory")
	if err != nil{
		return nil,err
		}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil{
		return nil,err
	}

	var usList [] *userstory.UserStory

	for rows.Next(){
		var u userstory.UserStory
		err = rows.Scan(&u.ID,&u.Title,&u.Description,&u.Effort,&u.CreationDate)
		if err != nil{
			return nil,err
		}
		usList = append(usList,&u)
	}
	return usList,nil
}