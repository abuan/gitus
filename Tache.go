package main

var counterTache int

// Tache : une tache contient la description d'une tache et d'un id
type Tache struct {
	description string
	id          int
}

func newTache(description string) Tache {
	counterTache++
	return Tache{description, counterTache}
}
