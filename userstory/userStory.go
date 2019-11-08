package userstory

import (
	"github.com/abuan/gitus/task"
	"time"
)

var counterUserStory int

// UserStory : le type de base de notre projet
type UserStory struct {
	Name		  string
	Description   string
	ID            int
	listTask      []task.Task
	listUserStory []UserStory
	Effort 		  int
	CreationDate  time.Time
}

// SetDescription : Affecte une description à une UserStroy
func (u *UserStory) SetDescription(s string) {
	u.Description = s
}
func (u *UserStory) addTask(s string) {
	u.listTask = append(u.listTask, task.NewTask(s))
}
func (u *UserStory) addUserStory(name , description string, effort int) {
	u.listUserStory = append(u.listUserStory, NewUserStory(name,description,effort))
}
func (u *UserStory) getID() int {
	return u.ID
}

// SetEffort : Affecte une valeur d'effort à une UserStroy
func(u *UserStory) SetEffort(e int){
	u.Effort = e;
}

// SetCreationDate : Affecte une valeur date de création à une USerStory
func(u *UserStory) SetCreationDate(t time.Time){
	t.Format("2006-01-02 15:04:05")
	u.CreationDate = t
}

// SetCreationDateNow : Affecte la date actuelle à une USerStory
func(u *UserStory) SetCreationDateNow(){
	t := time.Now()
	t.Format("2006-01-02 15:04:05")
	u.CreationDate = t
}

//NewUserStory : constructeur structure UserStory, incrémente l'id automatiquement
func NewUserStory(name , description string, effort int) UserStory {
	t := time.Now()
	t.Format("2006-01-02 15:04:05")
	return UserStory{name, description, 0, nil, nil, effort, t}
}
