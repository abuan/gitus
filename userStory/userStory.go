package userstory

import (
	"github.com/abuan/proto_gitus/task"
)

var counterUserStory int

// UserStory : le type de base de notre projet
type UserStory struct {
	description   string
	id            int
	listTask      []task.Task
	listUserStory []UserStory
}

func (u *UserStory) setDescription(s string) {
	u.description = s
}
func (u *UserStory) addTask(s string) {
	u.listTask = append(u.listTask, task.NewTask(s))
}
func (u *UserStory) addUserStory(s string) {
	u.listUserStory = append(u.listUserStory, NewUserStory(s))
}
func (u *UserStory) getID() int {
	return u.id
}

//NewUserStory : constructeur structure UserStory, incrémente l'id automatiquement
func NewUserStory(description string) UserStory {
	counterUserStory++
	return UserStory{description, counterUserStory, nil, nil}
}
