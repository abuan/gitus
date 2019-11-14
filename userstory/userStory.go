package userstory

import (
	"github.com/abuan/gitus/task"
	"time"
	"fmt"
	"strconv"
)

var counterUserStory int

// UserStory : le type de base de notre projet
type UserStory struct {
	Name		  string
	Description   string
	ID            int
	Effort 		  int
	CreationDate  time.Time
	ListTask      []task.Task
	ListUserStory []UserStory
}

// SetDescription : Affecte une description à une UserStroy
func (u *UserStory) SetDescription(s string) {
	u.Description = s
}
func (u *UserStory) addTask(s string) {
	u.ListTask = append(u.ListTask, task.NewTask(s))
}
func (u *UserStory) addUserStory(name , description string, effort int) {
	u.ListUserStory = append(u.ListUserStory, NewUserStory(name,description,effort))
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
	return UserStory{name, description, 0, effort, t, nil, nil}
}

// Display : Affiche le contenu de la US sur le terminal
func(u *UserStory)Display(){
	fmt.Println("**************************************** User Story ****************************************")
	fmt.Println("\tId:\t\t"+strconv.Itoa(u.ID))
	fmt.Println("\tName :\t\t"+u.Name)
	fmt.Println("\tEffort:\t\t"+strconv.Itoa(u.Effort))
	fmt.Println("\tDescription :\t"+u.Description)
}