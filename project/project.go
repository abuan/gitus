package project

import (
	"fmt"
	"time"
	"github.com/abuan/gitus/userstory"
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
	return Project{0, name, description, nil, t}
}

//modifier le nom du projet
func (p *Project) setName(s string) {
	p.Name = s
}

//modifier la description
func (p *Project) setDescription(s string) {
	p.Description = s
}

//ajouter une us à un projet
func (p *Project) addUserStory(s string) {
	p.ListUserStory = append(p.ListUserStory, userstory.NewUserStory(s))
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

//DisplayProject : afficher un projet entier
func (p *Project) DisplayProject() {

	fmt.Printf("The project %v is %v\n %v\n Created the %v\n", p.id, p.name, p.description, p.creationDate.Format(time.ANSIC))
}

func (p *Project) getUserStories() {}
