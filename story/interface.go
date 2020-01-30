package story

import (
	"github.com/MichaelMure/git-bug/entity"
	"github.com/MichaelMure/git-bug/repository"
	"github.com/MichaelMure/git-bug/util/lamport"
)

type Interface interface {
	// Id return the Story identifier
	Id() entity.Id

	// Validate check if the Story data is valid
	Validate() error

	// Append an operation into the staging area, to be committed later
	Append(op Operation)

	// Indicate that the in-memory state changed and need to be commit in the repository
	NeedCommit() bool

	// Commit write the staging area in Git and move the operations to the packs
	Commit(repo repository.ClockedRepo) error

	// Merge a different version of the same story by rebasing operations of this story
	// that are not present in the other on top of the chain of operations of the
	// other version.
	//Merge(repo repository.Repo, other Interface) (bool, error)

	// Lookup for the very first operation of the story.
	// For a valid Story, this operation should be a CreateOp
	FirstOp() Operation

	// Lookup for the very last operation of the story.
	// For a valid Story, should never be nil
	LastOp() Operation

	// Compile a story in a easily usable snapshot
	Compile() Snapshot

	// CreateLamportTime return the Lamport time of creation
	CreateLamportTime() lamport.Time

	// EditLamportTime return the Lamport time of the last edit
	EditLamportTime() lamport.Time
}

func storyFromInterface(story Interface) *Story {
	switch story.(type) {
	case *Story:
		return story.(*Story)
	case *WithSnapshot:
		return story.(*WithSnapshot).Story
	default:
		panic("missing type case")
	}
}
