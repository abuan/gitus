package cache

import (
	"encoding/gob"
	"fmt"

	"github.com/abuan/gitus/story"
	"github.com/MichaelMure/git-bug/entity"
	"github.com/MichaelMure/git-bug/identity"
	"github.com/MichaelMure/git-bug/util/lamport"
)

// Package initialisation used to register the type for (de)serialization
func init() {
	gob.Register(StoryExcerpt{})
}

// StoryExcerpt hold a subset of the story values to be able to sort and filter stories
// efficiently without having to read and compile each raw stories.
type StoryExcerpt struct {
	Id entity.Id

	CreateLamportTime lamport.Time
	EditLamportTime   lamport.Time
	CreateUnixTime    int64
	EditUnixTime      int64

	Title        string
	Effort 		 int
	Status 		 story.Status
	Actors       []entity.Id
	Participants []entity.Id

	// If author is identity.Bare, LegacyAuthor is set
	// If author is identity.Identity, AuthorId is set and data is deported
	// in a IdentityExcerpt
	LegacyAuthor LegacyAuthorExcerpt
	AuthorId     entity.Id

	CreateMetadata map[string]string
}

// identity.Bare data are directly embedded in the story excerpt
 type LegacyAuthorExcerpt struct {
	Name  string
	Login string
}

func (l LegacyAuthorExcerpt) DisplayName() string {
	switch {
	case l.Name == "" && l.Login != "":
		return l.Login
	case l.Name != "" && l.Login == "":
		return l.Name
	case l.Name != "" && l.Login != "":
		return fmt.Sprintf("%s (%s)", l.Name, l.Login)
	}

	panic("invalid person data")
} 

func NewStoryExcerpt(s story.Interface, snap *story.Snapshot) *StoryExcerpt {
	participantsIds := make([]entity.Id, len(snap.Participants))
	for i, participant := range snap.Participants {
		participantsIds[i] = participant.Id()
	}

	actorsIds := make([]entity.Id, len(snap.Actors))
	for i, actor := range snap.Actors {
		actorsIds[i] = actor.Id()
	}

	e := &StoryExcerpt{
		Id:                s.Id(),
		CreateLamportTime: s.CreateLamportTime(),
		EditLamportTime:   s.EditLamportTime(),
		CreateUnixTime:    s.FirstOp().GetUnixTime(),
		EditUnixTime:      snap.LastEditUnix(),
		Actors:            actorsIds,
		Participants:      participantsIds,
		Title:             snap.Title,
		Effort:			   snap.Effort,
		Status:			   snap.Status,
		CreateMetadata:    s.FirstOp().AllMetadata(),
	}

	switch snap.Author.(type) {
	case *identity.Identity:
		e.AuthorId = snap.Author.Id()
	case *identity.Bare:
		e.LegacyAuthor = LegacyAuthorExcerpt{
			Login: snap.Author.Login(),
			Name:  snap.Author.Name(),
		}
	default:
		panic("unhandled identity type")
	}

	return e
}

/*
 * Sorting
 */

type StoriesById []*StoryExcerpt

func (s StoriesById) Len() int {
	return len(s)
}

func (s StoriesById) Less(i, j int) bool {
	return s[i].Id < s[j].Id
}

func (s StoriesById) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type StoriesByCreationTime []*StoryExcerpt

func (s StoriesByCreationTime) Len() int {
	return len(s)
}

func (s StoriesByCreationTime) Less(i, j int) bool {
	if s[i].CreateLamportTime < s[j].CreateLamportTime {
		return true
	}

	if s[i].CreateLamportTime > s[j].CreateLamportTime {
		return false
	}

	// When the logical clocks are identical, that means we had a concurrent
	// edition. In this case we rely on the timestamp. While the timestamp might
	// be incorrect due to a badly set clock, the drift in sorting is bounded
	// by the first sorting using the logical clock. That means that if users
	// synchronize their stories regularly, the timestamp will rarely be used, and
	// should still provide a kinda accurate sorting when needed.
	return s[i].CreateUnixTime < s[j].CreateUnixTime
}

func (s StoriesByCreationTime) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type StoriesByEditTime []*StoryExcerpt

func (s StoriesByEditTime) Len() int {
	return len(s)
}

func (s StoriesByEditTime) Less(i, j int) bool {
	if s[i].EditLamportTime < s[j].EditLamportTime {
		return true
	}

	if s[i].EditLamportTime > s[j].EditLamportTime {
		return false
	}

	// When the logical clocks are identical, that means we had a concurrent
	// edition. In this case we rely on the timestamp. While the timestamp might
	// be incorrect due to a badly set clock, the drift in sorting is bounded
	// by the first sorting using the logical clock. That means that if users
	// synchronize their stories regularly, the timestamp will rarely be used, and
	// should still provide a kinda accurate sorting when needed.
	return s[i].EditUnixTime < s[j].EditUnixTime
}

func (s StoriesByEditTime) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
