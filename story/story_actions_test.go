package story

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/MichaelMure/git-bug/identity"
	"github.com/MichaelMure/git-bug/repository"
)

func TestPushPull(t *testing.T) {
	//Créer différents dépot
	repoA, repoB, remote := repository.SetupReposAndRemote(t)
	defer repository.CleanupTestRepos(t, repoA, repoB, remote)

	//Ajoute identité
	reneA := identity.NewIdentity("René Descartes", "rene@descartes.fr")

	story1, _, err := Create(reneA, time.Now().Unix(), "story1","descript", 3)
	require.NoError(t, err)
	assert.True(t, story1.NeedCommit())
	err = story1.Commit(repoA)
	require.NoError(t, err)
	assert.False(t, story1.NeedCommit())

	// distribute the identity
	_, err = identity.Push(repoA, "origin")
	require.NoError(t, err)
	err = identity.Pull(repoB, "origin")
	require.NoError(t, err)

	// A --> remote --> B
	_, err = Push(repoA, "origin")
	require.NoError(t, err)

	err = Pull(repoB, "origin")
	require.NoError(t, err)

	stories := allStories(t, ReadAllLocalStories(repoB))

	if len(stories) != 1 {
		t.Fatal("Unexpected number of stories")
	}

	// B --> remote --> A
	reneB, err := identity.ReadLocal(repoA, reneA.Id())
	require.NoError(t, err)

	story2, _, err := Create(reneB, time.Now().Unix(), "story2", "descript", 3)
	require.NoError(t, err)
	err = story2.Commit(repoB)
	require.NoError(t, err)

	_, err = Push(repoB, "origin")
	require.NoError(t, err)

	err = Pull(repoA, "origin")
	require.NoError(t, err)

	stories = allStories(t, ReadAllLocalStories(repoA))

	if len(stories) != 2 {
		t.Fatal("Unexpected number of stories")
	}
}

func allStories(t testing.TB, stories <-chan StreamedStory) []*Story {
	var result []*Story
	for streamed := range stories {
		if streamed.Err != nil {
			t.Fatal(streamed.Err)
		}
		result = append(result, streamed.Story)
	}
	return result
}

func TestRebaseTheirs(t *testing.T) {
	_RebaseTheirs(t)
}

func BenchmarkRebaseTheirs(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_RebaseTheirs(b)
	}
}

func _RebaseTheirs(t testing.TB) {
	repoA, repoB, remote := repository.SetupReposAndRemote(t)
	defer repository.CleanupTestRepos(t, repoA, repoB, remote)

	reneA := identity.NewIdentity("René Descartes", "rene@descartes.fr")

	story1, _, err := Create(reneA, time.Now().Unix(), "story1","descript", 3)
	require.NoError(t, err)
	assert.True(t, story1.NeedCommit())
	err = story1.Commit(repoA)
	require.NoError(t, err)
	assert.False(t, story1.NeedCommit())

	// distribute the identity
	_, err = identity.Push(repoA, "origin")
	require.NoError(t, err)
	err = identity.Pull(repoB, "origin")
	require.NoError(t, err)

	// A --> remote

	_, err = Push(repoA, "origin")
	require.NoError(t, err)

	// remote --> B
	err = Pull(repoB, "origin")
	require.NoError(t, err)

	story2, err := ReadLocalStory(repoB, story1.Id())
	require.NoError(t, err)
	assert.False(t, story2.NeedCommit())

	reneB, err := identity.ReadLocal(repoA, reneA.Id())
	require.NoError(t, err)

	
	_, err = SetEffort(story2, reneB, time.Now().Unix(), 2)
	require.NoError(t, err)
	assert.True(t, story2.NeedCommit())
	_, err = SetEffort(story2, reneB, time.Now().Unix(), 5)
	require.NoError(t, err)
	_, err = SetEffort(story2, reneB, time.Now().Unix(), 10)
	require.NoError(t, err)
	
	err = story2.Commit(repoB)
	require.NoError(t, err)
	assert.False(t, story2.NeedCommit())

	// B --> remote
	_, err = Push(repoB, "origin")
	require.NoError(t, err)

	// remote --> A
	err = Pull(repoA, "origin")
	require.NoError(t, err)

	stories := allStories(t, ReadAllLocalStories(repoB))

	if len(stories) != 1 {
		t.Fatal("Unexpected number of stories")
	}

	story3, err := ReadLocalStory(repoA, story1.Id())
	require.NoError(t, err)

	if nbOps(story3) != 4 {
		t.Fatal("Unexpected number of operations")
	}
}

func TestRebaseOurs(t *testing.T) {
	_RebaseOurs(t)
}

func BenchmarkRebaseOurs(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_RebaseOurs(b)
	}
}

func _RebaseOurs(t testing.TB) {
	repoA, repoB, remote := repository.SetupReposAndRemote(t)
	defer repository.CleanupTestRepos(t, repoA, repoB, remote)

	reneA := identity.NewIdentity("René Descartes", "rene@descartes.fr")

	story1, _, err := Create(reneA, time.Now().Unix(), "story1", "descript", 3)
	require.NoError(t, err)
	err = story1.Commit(repoA)
	require.NoError(t, err)

	// distribute the identity
	_, err = identity.Push(repoA, "origin")
	require.NoError(t, err)
	err = identity.Pull(repoB, "origin")
	require.NoError(t, err)

	// A --> remote
	_, err = Push(repoA, "origin")
	require.NoError(t, err)

	// remote --> B
	err = Pull(repoB, "origin")
	require.NoError(t, err)


	_, err = SetEffort(story1, reneA, time.Now().Unix(), 2)
	require.NoError(t, err)
	_, err = SetEffort(story1, reneA, time.Now().Unix(), 1)
	require.NoError(t, err)
	_, err = SetEffort(story1, reneA, time.Now().Unix(), 4)
	require.NoError(t, err)
	
	
	err = story1.Commit(repoA)
	require.NoError(t, err)

	
	_, err = SetEffort(story1, reneA, time.Now().Unix(), 6)
	require.NoError(t, err)
	_, err = SetEffort(story1, reneA, time.Now().Unix(), 5)
	require.NoError(t, err)
	_, err = SetEffort(story1, reneA, time.Now().Unix(), 7)
	require.NoError(t, err)
	
	err = story1.Commit(repoA)
	require.NoError(t, err)

	
	_, err = SetEffort(story1, reneA, time.Now().Unix(), 8)
	require.NoError(t, err)
	_, err = SetEffort(story1, reneA, time.Now().Unix(), 9)
	require.NoError(t, err)
	_, err = SetEffort(story1, reneA, time.Now().Unix(), 10)
	require.NoError(t, err)
	
	err = story1.Commit(repoA)
	require.NoError(t, err)

	// remote --> A
	err = Pull(repoA, "origin")
	require.NoError(t, err)

	stories := allStories(t, ReadAllLocalStories(repoA))

	if len(stories) != 1 {
		t.Fatal("Unexpected number of stories")
	}

	story2, err := ReadLocalStory(repoA, story1.Id())
	require.NoError(t, err)

	if nbOps(story2) != 10 {
		t.Fatal("Unexpected number of operations")
	}
}

func nbOps(b *Story) int {
	it := NewOperationIterator(b)
	counter := 0
	for it.Next() {
		counter++
	}
	return counter
}

func TestRebaseConflict(t *testing.T) {
	_RebaseConflict(t)
}

func BenchmarkRebaseConflict(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_RebaseConflict(b)
	}
}

func _RebaseConflict(t testing.TB) {
	repoA, repoB, remote := repository.SetupReposAndRemote(t)
	defer repository.CleanupTestRepos(t, repoA, repoB, remote)

	reneA := identity.NewIdentity("René Descartes", "rene@descartes.fr")

	story1, _, err := Create(reneA, time.Now().Unix(), "story1", "descript", 3)
	require.NoError(t, err)
	err = story1.Commit(repoA)
	require.NoError(t, err)

	// distribute the identity
	_, err = identity.Push(repoA, "origin")
	require.NoError(t, err)
	err = identity.Pull(repoB, "origin")
	require.NoError(t, err)

	// A --> remote
	_, err = Push(repoA, "origin")
	require.NoError(t, err)

	// remote --> B
	err = Pull(repoB, "origin")
	require.NoError(t, err)


	_, err = SetEffort(story1, reneA, time.Now().Unix(), 2)
	require.NoError(t, err)
	_, err = SetEffort(story1, reneA, time.Now().Unix(), 4)
	require.NoError(t, err)
	_, err = SetEffort(story1, reneA, time.Now().Unix(), 5)
	require.NoError(t, err)
	
	err = story1.Commit(repoA)
	require.NoError(t, err)

	_, err = SetEffort(story1, reneA, time.Now().Unix(), 6)
	require.NoError(t, err)
	_, err = SetEffort(story1, reneA, time.Now().Unix(), 7)
	require.NoError(t, err)
	_, err = SetEffort(story1, reneA, time.Now().Unix(), 8)
	require.NoError(t, err)
	
	err = story1.Commit(repoA)
	require.NoError(t, err)

	
	_, err = SetEffort(story1, reneA, time.Now().Unix(), 9)
	require.NoError(t, err)
	_, err = SetEffort(story1, reneA, time.Now().Unix(), 10)
	require.NoError(t, err)
	_, err = SetEffort(story1, reneA, time.Now().Unix(), 11)
	require.NoError(t, err)
	
	err = story1.Commit(repoA)
	require.NoError(t, err)

	story2, err := ReadLocalStory(repoB, story1.Id())
	require.NoError(t, err)

	reneB, err := identity.ReadLocal(repoA, reneA.Id())
	require.NoError(t, err)

	
	_, err = SetEffort(story2, reneB, time.Now().Unix(), 12)
	require.NoError(t, err)
	_, err = SetEffort(story2, reneB, time.Now().Unix(), 13)
	require.NoError(t, err)
	_, err = SetEffort(story2, reneB, time.Now().Unix(), 14)
	require.NoError(t, err)
	
	err = story2.Commit(repoB)
	require.NoError(t, err)

	_, err = SetEffort(story2, reneB, time.Now().Unix(), 15)
	require.NoError(t, err)
	_, err = SetEffort(story2, reneB, time.Now().Unix(), 16)
	require.NoError(t, err)
	_, err = SetEffort(story2, reneB, time.Now().Unix(), 17)
	require.NoError(t, err)

	err = story2.Commit(repoB)
	require.NoError(t, err)

	_, err = SetEffort(story2, reneB, time.Now().Unix(), 18)
	require.NoError(t, err)
	_, err = SetEffort(story2, reneB, time.Now().Unix(), 19)
	require.NoError(t, err)
	_, err = SetEffort(story2, reneB, time.Now().Unix(), 20)
	require.NoError(t, err)

	err = story2.Commit(repoB)
	require.NoError(t, err)

	// A --> remote
	_, err = Push(repoA, "origin")
	require.NoError(t, err)

	// remote --> B
	err = Pull(repoB, "origin")
	require.NoError(t, err)

	stories := allStories(t, ReadAllLocalStories(repoB))

	if len(stories) != 1 {
		t.Fatal("Unexpected number of stories")
	}

	story3, err := ReadLocalStory(repoB, story1.Id())
	require.NoError(t, err)

	if nbOps(story3) != 19 {
		t.Fatal("Unexpected number of operations")
	}

	// B --> remote
	_, err = Push(repoB, "origin")
	require.NoError(t, err)

	// remote --> A
	err = Pull(repoA, "origin")
	require.NoError(t, err)

	stories = allStories(t, ReadAllLocalStories(repoA))

	if len(stories) != 1 {
		t.Fatal("Unexpected number of stories")
	}

	story4, err := ReadLocalStory(repoA, story1.Id())
	require.NoError(t, err)

	if nbOps(story4) != 19 {
		t.Fatal("Unexpected number of operations")
	}
}
