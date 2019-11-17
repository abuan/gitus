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