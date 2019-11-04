package task

var counterTask int

// Task : une tache contient la description d'une tache et d'un id
type Task struct {
	description string
	id          int
}

func NewTask(description string) Task {
	counterTask++
	return Task{description, counterTask}
}
