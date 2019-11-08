package project

import (
	"fmt"
	"time"

	"github.com/abuan/gitus/userstory"
)

var counterProject int

// Project : un projet contient contient un id,
// un nom, une description, une liste de user stories
type Project struct {
	id            int
	name          string
	description   string
	listUserStory []userstory.UserStory
	creationDate  time.Time
}

//NewProject : constructeur de Projet, implémente automatiquement la fonction projet
func NewProject(name, description string) Project {
	counterProject++
	return Project{counterProject, name, description, nil, time.Now()}
}

//modifier le nom du projet
func (p *Project) setName(s string) {
	p.name = s
}

//modifier la description
func (p *Project) setDescription(s string) {
	p.description = s
}

//ajouter une us à un projet
func (p *Project) addUserStory(s string) {
	p.listUserStory = append(p.listUserStory, userstory.NewUserStory(s))
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
