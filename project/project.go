package project

import (
	"fmt"
	"time"
	"github.com/abuan/gitus/userstory"
	"strconv"
)

// Project : un projet contient contient un id,
// un nom, une description, une liste de user stories
type Project struct {
	ID            int
	Name          string
	Description   string
	CreationDate  time.Time
	ListUserStory []userstory.UserStory
}

//NewProject : constructeur de Projet, implémente automatiquement la fonction projet
func NewProject(name, description string) Project {
	t := time.Now()
	t.Format("2006-01-02 15:04:05")
	return Project{0, name, description, t,nil}
}

//modifier le nom du projet
func (p *Project) setName(s string) {
	p.Name = s
}

//modifier la description
func (p *Project) setDescription(s string) {
	p.Description = s
}

//AddUserStory : ajoute une us à un projet
func (p *Project) AddUserStory(u userstory.UserStory) {
	p.ListUserStory = append(p.ListUserStory, u)
}

/*
// trouver une user story dans un projet
func (p *Project) findUserStory(u userstory.UserStory) int {
	for i, n := range p.listUserStory {
		if u == n {
			return i
		}
	}
	return len(p.listUserStory)
}

//supprimer une us à un projet
func (p *Project) deleteUserStory(u userstory.UserStory) {
	i = p.findUserStory(u)
	p.listUserStory[i] = p.listUserStory[len(p.listUserStory)-1]
	p.listUserStory[len(p.listUserStory)-1] = "" //supprimer dernier element
	p.listUserStory = p.listUserStory[:len(p.listUserStory)-1]
} */

//Display : afficher un projet entier
func (p *Project) Display(usList []int) {
	fmt.Println("\n**************************************** Project ****************************************")
	fmt.Println("\tId:\t\t"+strconv.Itoa(p.ID))
	fmt.Println("\tName :\t\t"+p.Name)
	fmt.Println("\tDescription :\t"+p.Description)
	fmt.Println("\tCreation Date :\t"+p.CreationDate.Format("2006-01-02 15:04:05")+"\n")
	fmt.Print("\tLinked User Story (id) :\t")
	if len(usList)>0{
		for _,value := range usList{
			fmt.Print(strconv.Itoa(value)+" / ")
		}
	}else{
		fmt.Println("No User Story link to this project")
	}
}

func (p *Project) getUserStories() {}
