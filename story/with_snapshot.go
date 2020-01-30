package story

import "github.com/MichaelMure/git-bug/repository"

var _ Interface = &WithSnapshot{}

// WithSnapshot encapsulate a Story and maintain the corresponding Snapshot efficiently
type WithSnapshot struct {
	*Story
	snap *Snapshot
}

// Snapshot return the current snapshot
func (s *WithSnapshot) Snapshot() *Snapshot {
	if s.snap == nil {
		snap := s.Story.Compile()
		s.snap = &snap
	}
	return s.snap
}

// Append intercept Story.Append() to update the snapshot efficiently
func (s *WithSnapshot) Append(op Operation) {
	s.Story.Append(op)

	if s.snap == nil {
		return
	}

	op.Apply(s.snap)
	s.snap.Operations = append(s.snap.Operations, op)
}

// Commit intercept Story.Commit() to update the snapshot efficiently
func (s *WithSnapshot) Commit(repo repository.ClockedRepo) error {
	err := s.Story.Commit(repo)

	if err != nil {
		s.snap = nil
		return err
	}

	// Commit() shouldn't change anything of the story state apart from the
	// initial ID set

	if s.snap == nil {
		return nil
	}

	s.snap.id = s.Story.id
	return nil
}

// Merge intercept Story.Merge() and clear the snapshot
//func (s *WithSnapshot) Merge(repo repository.Repo, other Interface) (bool, error) {
//	s.snap = nil
//	return s.Story.Merge(repo, other)
//}
