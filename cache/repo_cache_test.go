package cache

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MichaelMure/git-bug/repository"
)

//Ce test ne passe pas sur ma machine à cause d'un problème d'accès à un fichier
func TestCache(t *testing.T) {
	repo := repository.CreateTestRepo(false)
	defer repository.CleanupTestRepos(t, repo)

	cache, err := NewRepoCache(repo)
	require.NoError(t, err)

	// Create, set and get user identity
	iden1, err := cache.NewIdentity("René Descartes", "rene@descartes.fr")
	require.NoError(t, err)
	err = cache.SetUserIdentity(iden1)
	require.NoError(t, err)
	userIden, err := cache.GetUserIdentity()
	require.NoError(t, err)
	require.Equal(t, iden1.Id(), userIden.Id())

	// it's possible to create two identical identities
	iden2, err := cache.NewIdentity("René Descartes", "rene@descartes.fr")
	require.NoError(t, err)

	// Two identical identities yield a different id
	require.NotEqual(t, iden1.Id(), iden2.Id())

	// There is now two identities in the cache
	require.Len(t, cache.AllIdentityIds(), 2)
	require.Len(t, cache.identitiesExcerpts, 2)
	require.Len(t, cache.identities, 2)

	// Create a story
	story1, _, err := cache.NewStory("title", "descript",1)
	require.NoError(t, err)

	// It's possible to create two identical stories
	story2, _, err := cache.NewStory("title", "descript",1)
	require.NoError(t, err)

	// two identical stories yield a different id
	require.NotEqual(t, story1.Id(), story2.Id())

	// There is now two stories in the cache
	require.Len(t, cache.AllStoriesIds(), 2)
	require.Len(t, cache.storyExcerpts, 2)
	require.Len(t, cache.stories, 2)

	// Resolving
	_, err = cache.ResolveIdentity(iden1.Id())
	require.NoError(t, err)
	_, err = cache.ResolveIdentityExcerpt(iden1.Id())
	require.NoError(t, err)
	_, err = cache.ResolveIdentityPrefix(iden1.Id().String()[:10])
	require.NoError(t, err)

	_, err = cache.ResolveStory(story1.Id())
	require.NoError(t, err)
	_, err = cache.ResolveStoryExcerpt(story1.Id())
	require.NoError(t, err)
	_, err = cache.ResolveStoryPrefix(story1.Id().String()[:10])
	require.NoError(t, err)

	// Querying
	query, err := ParseQuery("status:open author:descartes sort:edit-asc")
	require.NoError(t, err)
	require.Len(t, cache.QueryStories(query), 2)

	// Close
	require.NoError(t, cache.Close())
	require.Empty(t, cache.stories)
	require.Empty(t, cache.storyExcerpts)
	require.Empty(t, cache.identities)
	require.Empty(t, cache.identitiesExcerpts)

	// Reload, only excerpt are loaded
	require.NoError(t, cache.load())
	require.Empty(t, cache.stories)
	require.Empty(t, cache.identities)
	require.Len(t, cache.storyExcerpts, 2)
	require.Len(t, cache.identitiesExcerpts, 2)

	// Resolving load from the disk
	_, err = cache.ResolveIdentity(iden1.Id())
	require.NoError(t, err)
	_, err = cache.ResolveIdentityExcerpt(iden1.Id())
	require.NoError(t, err)
	_, err = cache.ResolveIdentityPrefix(iden1.Id().String()[:10])
	require.NoError(t, err)

	_, err = cache.ResolveStory(story1.Id())
	require.NoError(t, err)
	_, err = cache.ResolveStoryExcerpt(story1.Id())
	require.NoError(t, err)
	_, err = cache.ResolveStoryPrefix(story1.Id().String()[:10])
	require.NoError(t, err)
}

func TestPushPull(t *testing.T) {
	repoA, repoB, remote := repository.SetupReposAndRemote(t)
	defer repository.CleanupTestRepos(t, repoA, repoB, remote)

	cacheA, err := NewRepoCache(repoA)
	require.NoError(t, err)

	cacheB, err := NewRepoCache(repoB)
	require.NoError(t, err)

	// Create, set and get user identity
	reneA, err := cacheA.NewIdentity("René Descartes", "rene@descartes.fr")
	require.NoError(t, err)
	err = cacheA.SetUserIdentity(reneA)
	require.NoError(t, err)

	// distribute the identity
	_, err = cacheA.Push("origin")
	require.NoError(t, err)
	err = cacheB.Pull("origin")
	require.NoError(t, err)

	// Create a story in A
	_, _, err = cacheA.NewStory("story1", "descript",1)
	require.NoError(t, err)

	// A --> remote --> B
	_, err = cacheA.Push("origin")
	require.NoError(t, err)

	err = cacheB.Pull("origin")
	require.NoError(t, err)

	require.Len(t, cacheB.AllStoriesIds(), 1)

	// retrieve and set identity
	reneB, err := cacheB.ResolveIdentity(reneA.Id())
	require.NoError(t, err)

	err = cacheB.SetUserIdentity(reneB)
	require.NoError(t, err)

	// B --> remote --> A
	_, _, err = cacheB.NewStory("story2", "descript",1)
	require.NoError(t, err)

	_, err = cacheB.Push("origin")
	require.NoError(t, err)

	err = cacheA.Pull("origin")
	require.NoError(t, err)

	require.Len(t, cacheA.AllStoriesIds(), 2)
}
