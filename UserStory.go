package main

var counterUserStory int

// UserStory : le type de base de notre projet
type UserStory struct {
	description   string
	id            int
	listTache     []Tache
	listUserStory []UserStory
}

func (u *UserStory) setDescription(s string) {
	u.description = s
}
func (u *UserStory) addTache(s string) {
	u.listTache = append(u.listTache, newTache(s))
}
func (u *UserStory) addUserStory(s string) {
	u.listUserStory = append(u.listUserStory, newUserStory(s))
}
func (u *UserStory) getID() int {
	return u.id
}
func newUserStory(description string) UserStory {
	counterUserStory++
	return UserStory{description, counterUserStory, nil, nil}
}
