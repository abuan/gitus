package story

import (
	"testing"
	"time"

	"github.com/MichaelMure/git-bug/identity"
	"github.com/MichaelMure/git-bug/repository"
	"github.com/stretchr/testify/assert"
)

func TestStoryId(t *testing.T) {
	mockRepo := repository.NewMockRepoForTest()

	story1 := NewStory()

	rene := identity.NewIdentity("René Descartes", "rene@descartes.fr")
	createOp := NewCreateOp(rene, time.Now().Unix(), "title", "descript", 3, OpenStatus)

	story1.Append(createOp)

	err := story1.Commit(mockRepo)

	if err != nil {
		t.Fatal(err)
	}

	story1.Id()
}

func TestStoryValidity(t *testing.T) {
	mockRepo := repository.NewMockRepoForTest()

	story1 := NewStory()

	rene := identity.NewIdentity("René Descartes", "rene@descartes.fr")
	createOp := NewCreateOp(rene, time.Now().Unix(), "title", "descript", 3, OpenStatus)

	if story1.Validate() == nil {
		t.Fatal("Empty story should be invalid")
	}

	story1.Append(createOp)

	if story1.Validate() != nil {
		t.Fatal("Story with just a CreateOp should be valid")
	}

	err := story1.Commit(mockRepo)
	if err != nil {
		t.Fatal(err)
	}

	story1.Append(createOp)

	if story1.Validate() == nil {
		t.Fatal("Story with multiple CreateOp should be invalid")
	}

	err = story1.Commit(mockRepo)
	if err == nil {
		t.Fatal("Invalid story should not commit")
	}
}

func TestStoryCommitLoad(t *testing.T) {
	story1 := NewStory()

	rene := identity.NewIdentity("René Descartes", "rene@descartes.fr")
	createOp := NewCreateOp(rene, time.Now().Unix(), "title", "descript", 3, OpenStatus)
	setTitleOp := NewSetTitleOp(rene, time.Now().Unix(), "title2", "title1")

	story1.Append(createOp)
	story1.Append(setTitleOp)

	repo := repository.NewMockRepoForTest()

	assert.True(t, story1.NeedCommit())

	err := story1.Commit(repo)
	assert.Nil(t, err)
	assert.False(t, story1.NeedCommit())

	story2, err := ReadLocalStory(repo, story1.Id())
	assert.NoError(t, err)
	equivalentStory(t, story1, story2)

}

func equivalentStory(t *testing.T, expected, actual *Story) {
	assert.Equal(t, len(expected.packs), len(actual.packs))

	for i := range expected.packs {
		for j := range expected.packs[i].Operations {
			actual.packs[i].Operations[j].base().id = expected.packs[i].Operations[j].base().id
		}
	}

	assert.Equal(t, expected, actual)
}
