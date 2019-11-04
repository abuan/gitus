package task

var counterTask int

// Tache : une tache contient la description d'une tache et d'un id
type Task struct {
	description string
	id          int
}

func newTask(description string) Task {
	counterTask++
	return Task{description, counterTask}
}
