package db

import(
	"github.com/abuan/gitus/project"
)

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


// TaskDeleteProject : Supprime un projet dans la BDD à partir de son ID
func TaskDeleteProject(id int) error{
	//Suppression des liens du projet vers les US
	stmt, err := db.Prepare("DELETE FROM Project_structure WHERE project_id = ?")
	if err != nil{
		return err
	}
	defer stmt.Close()

	_,err = stmt.Exec(id)
	if err != nil{
		return err
	}
	
	//Suppression du projet
	stmt, err = db.Prepare("DELETE FROM Project WHERE id = ?")
	if err != nil{
		return err
	}
	_,err = stmt.Exec(id)
	if err != nil{
		return err
	}
	
	return nil
}

// TaskGetProject : Récupère un projet dans la BDD à partir de son ID et la liste des iD des US liées à ce projet
func TaskGetProject(id int)(*project.Project,[]int,error){
	// Récupère les infos du projet dans la BDD
	stmt, err := db.Prepare("SELECT name,descript,creation_date FROM Project WHERE id = ?")
	if err != nil{
		return nil,nil,err
		}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	//Utilisation des Keys de la structure pour indiquer le champs que l'on affecte
	p:= project.Project{ID:id}
	//Scan des résultats du Row
	err = row.Scan(&p.Name, &p.Description,&p.CreationDate)
	if err != nil{
		return nil,nil,err
	}

	//Récupération de la liste des ID des US
	stmt, err = db.Prepare("SELECT userstory_id FROM Project_structure WHERE project_id = ?")
	if err != nil{
		return nil,nil,err
	}
	rows, err := stmt.Query(id)
	if err != nil{
		return nil,nil,err
	}
	var (
		ids [] int
		usID int
	)
	for rows.Next(){
		err = rows.Scan(&usID)
		if err != nil{
			return nil,nil,err
		}
		ids = append(ids, usID)
	}

	return &p,ids,nil
}