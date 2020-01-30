package _select

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/abuan/gitus/cache"
	"github.com/MichaelMure/git-bug/repository"
)

func TestSelect(t *testing.T) {
	repo := repository.CreateTestRepo(false)
	defer repository.CleanupTestRepos(t, repo)

	repoCache, err := cache.NewRepoCache(repo)
	require.NoError(t, err)

	_, _, err = ResolveStory(repoCache, []string{})
	require.Equal(t, ErrNoValidId, err)

	err = Select(repoCache, "invalid")
	require.NoError(t, err)

	// Resolve without a pattern should fail when no story is selected
	_, _, err = ResolveStory(repoCache, []string{})
	require.Error(t, err)

	// generate a bunch of stories

	rene, err := repoCache.NewIdentity("René Descartes", "rene@descartes.fr")
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		_, _, err := repoCache.NewStoryRaw(rene, time.Now().Unix(), "title", "descript",1, nil)
		require.NoError(t, err)
	}

	// and two more for testing
	s1, _, err := repoCache.NewStoryRaw(rene, time.Now().Unix(), "title", "descript",2, nil)
	require.NoError(t, err)
	s2, _, err := repoCache.NewStoryRaw(rene, time.Now().Unix(), "title", "descript",3, nil)
	require.NoError(t, err)

	err = Select(repoCache, s1.Id())
	require.NoError(t, err)

	// normal select without args
	s3, _, err := ResolveStory(repoCache, []string{})
	require.NoError(t, err)
	require.Equal(t, s1.Id(), s3.Id())

	// override selection with same id
	s4, _, err := ResolveStory(repoCache, []string{s1.Id().String()})
	require.NoError(t, err)
	require.Equal(t, s1.Id(), s4.Id())

	// override selection with a prefix
	s5, _, err := ResolveStory(repoCache, []string{s1.Id().Human()})
	require.NoError(t, err)
	require.Equal(t, s1.Id(), s5.Id())

	// args that shouldn't override
	s6, _, err := ResolveStory(repoCache, []string{"arg"})
	require.NoError(t, err)
	require.Equal(t, s1.Id(), s6.Id())

	// override with a different id
	s7, _, err := ResolveStory(repoCache, []string{s2.Id().String()})
	require.NoError(t, err)
	require.Equal(t, s2.Id(), s7.Id())

	err = Clear(repoCache)
	require.NoError(t, err)

	// Resolve without a pattern should error again after clearing the selected story
	_, _, err = ResolveStory(repoCache, []string{})
	require.Error(t, err)
}
