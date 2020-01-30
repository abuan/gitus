package story

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/MichaelMure/git-bug/identity"
	"github.com/MichaelMure/git-bug/repository"
)

func ExampleOperationIterator() {
	s := NewStory()

	// add operations

	it := NewOperationIterator(s)

	for it.Next() {
		// do something with each operations
		_ = it.Value()
	}
}

func TestOpIterator(t *testing.T) {
	mockRepo := repository.NewMockRepoForTest()

	rene := identity.NewIdentity("René Descartes", "rene@descartes.fr")
	unix := time.Now().Unix()

	createOp := NewCreateOp(rene, unix, "title", "descript", 3, OpenStatus)
	setStatusOp := NewSetStatusOp(rene, unix, ClosedStatus)

	story1 := NewStory()

	// first pack
	story1.Append(createOp)
	story1.Append(setStatusOp)
	err := story1.Commit(mockRepo)
	assert.NoError(t, err)

	it := NewOperationIterator(story1)

	counter := 0
	for it.Next() {
		_ = it.Value()
		counter++
	}

	assert.Equal(t, 2, counter)
}
