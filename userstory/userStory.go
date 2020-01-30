package userstory

import (
	"fmt"
	"strconv"
	"github.com/abuan/gitus/namesgenerator"
	"time"
	"github.com/abuan/gitus/task"
)

var counterUserStory int

// UserStory : le type de base de notre projet
type UserStory struct {
	Title		  string
	Description   string
	Author        string
	ID            int
	Effort        int
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
func (u *UserStory) addUserStory(description, author string, effort int) {
	u.ListUserStory = append(u.ListUserStory, NewUserStory(description, author, effort))
}
func (u *UserStory) getID() int {
	return u.ID
}

// SetEffort : Affecte une valeur d'effort à une UserStroy
func (u *UserStory) SetEffort(e int) {
	u.Effort = e
}

// SetCreationDate : Affecte une valeur date de création à une USerStory
func (u *UserStory) SetCreationDate(t time.Time) {
	t.Format("2006-01-02 15:04:05")
	u.CreationDate = t
}

// SetCreationDateNow : Affecte la date actuelle à une USerStory
func (u *UserStory) SetCreationDateNow() {
	t := time.Now()
	t.Format("2006-01-02 15:04:05")
	u.CreationDate = t
}

//NewUserStory : constructeur structure UserStory, incrémente l'id automatiquement
func NewUserStory(description, author string, effort int) UserStory {
	t := time.Now()
	t.Format("2006-01-02 15:04:05")
	title := namesgenerator.GetRandomName(0);
	return UserStory{title, description, author, 0, effort, t, nil, nil}
}

// Display : Affiche le contenu de la US sur le terminal
func (u *UserStory) Display() {
	fmt.Println("\n**************************************** User Story ****************************************")
	fmt.Println("\tId:\t\t" + strconv.Itoa(u.ID))
	fmt.Println("\tTitle :\t\t"+u.Title)
	fmt.Println("\tEffort:\t\t" + strconv.Itoa(u.Effort))
	fmt.Println("\tAuthor:\t\t" + u.Author)
	fmt.Println("\tCreation Date :\t" + u.CreationDate.Format("2006-01-02 15:04:05"))
	fmt.Println("\tDescription :\t" + u.Description)
}
