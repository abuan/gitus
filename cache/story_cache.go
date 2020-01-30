package cache

import (
	"fmt"
	"time"

	"github.com/abuan/gitus/story"
	"github.com/MichaelMure/git-bug/entity"
)

var ErrNoMatchingOp = fmt.Errorf("no matching operation found")

// StoryCache is a wrapper around a Story. It provide multiple functions:
//
// 1. Provide a higher level API to use than the raw API from Story.
// 2. Maintain an up to date Snapshot available.
type StoryCache struct {
	repoCache *RepoCache
	story       *story.WithSnapshot
}

func NewStoryCache(repoCache *RepoCache, s *story.Story) *StoryCache {
	return &StoryCache{
		repoCache: repoCache,
		story:       &story.WithSnapshot{Story: s},
	}
}

func (c *StoryCache) Snapshot() *story.Snapshot {
	return c.story.Snapshot()
}

func (c *StoryCache) Id() entity.Id {
	return c.story.Id()
}

func (c *StoryCache) notifyUpdated() error {
	return c.repoCache.storyUpdated(c.story.Id())
}

// ResolveOperationWithMetadata will find an operation that has the matching metadata
func (c *StoryCache) ResolveOperationWithMetadata(key string, value string) (entity.Id, error) {
	// preallocate but empty
	matching := make([]entity.Id, 0, 5)

	it := story.NewOperationIterator(c.story)
	for it.Next() {
		op := it.Value()
		opValue, ok := op.GetMetadata(key)
		if ok && value == opValue {
			matching = append(matching, op.Id())
		}
	}

	if len(matching) == 0 {
		return "", ErrNoMatchingOp
	}

	if len(matching) > 1 {
		return "", story.NewErrMultipleMatchOp(matching)
	}

	return matching[0], nil
}


// fonction lier aux opérations non implémentées pour l'instant
/*
func (c *StoryCache) SetMetadata(target entity.Id, newMetadata map[string]string) (*story.SetMetadataOperation, error) {
	author, err := c.repoCache.GetUserIdentity()
	if err != nil {
		return nil, err
	}

	return c.SetMetadataRaw(author, time.Now().Unix(), target, newMetadata)
}

func (c *StoryCache) SetMetadataRaw(author *IdentityCache, unixTime int64, target entity.Id, newMetadata map[string]string) (*story.SetMetadataOperation, error) {
	op, err := story.SetMetadata(c.story, author.Identity, unixTime, target, newMetadata)
	if err != nil {
		return nil, err
	}

	return op, c.notifyUpdated()
}
 */

 func (c *StoryCache) Open() (*story.SetStatusOperation, error) {
	author, err := c.repoCache.GetUserIdentity()
	if err != nil {
		return nil, err
	}

	return c.OpenRaw(author, time.Now().Unix(), nil)
}

func (c *StoryCache) OpenRaw(author *IdentityCache, unixTime int64, metadata map[string]string) (*story.SetStatusOperation, error) {
	op, err := story.Open(c.story, author.Identity, unixTime)
	if err != nil {
		return nil, err
	}

	for key, value := range metadata {
		op.SetMetadata(key, value)
	}

	return op, c.notifyUpdated()
}

func (c *StoryCache) Close() (*story.SetStatusOperation, error) {
	author, err := c.repoCache.GetUserIdentity()
	if err != nil {
		return nil, err
	}

	return c.CloseRaw(author, time.Now().Unix(), nil)
}

func (c *StoryCache) CloseRaw(author *IdentityCache, unixTime int64, metadata map[string]string) (*story.SetStatusOperation, error) {
	op, err := story.Close(c.story, author.Identity, unixTime)
	if err != nil {
		return nil, err
	}

	for key, value := range metadata {
		op.SetMetadata(key, value)
	}

	return op, c.notifyUpdated()
}

func (c *StoryCache) EditEffort(effort int) (*story.SetEffortOperation, error) {
	author, err := c.repoCache.GetUserIdentity()
	if err != nil {
		return nil, err
	}

	return c.EditEffortRaw(author, time.Now().Unix(), effort, nil)
}

func (c *StoryCache) EditEffortRaw(author *IdentityCache, unixTime int64, effort int, metadata map[string]string) (*story.SetEffortOperation, error) {
	op, err := story.SetEffort(c.story, author.Identity, unixTime, effort)
	if err != nil {
		return nil, err
	}

	for key, value := range metadata {
		op.SetMetadata(key, value)
	}

	return op, c.notifyUpdated()
}

func (c *StoryCache) SetTitle(title string) (*story.SetTitleOperation, error) {
	author, err := c.repoCache.GetUserIdentity()
	if err != nil {
		return nil, err
	}

	return c.SetTitleRaw(author, time.Now().Unix(), title, nil)
}

func (c *StoryCache) SetTitleRaw(author *IdentityCache, unixTime int64, title string, metadata map[string]string) (*story.SetTitleOperation, error) {
	op, err := story.SetTitle(c.story, author.Identity, unixTime, title)
	if err != nil {
		return nil, err
	}

	for key, value := range metadata {
		op.SetMetadata(key, value)
	}

	return op, c.notifyUpdated()
}


func (c *StoryCache) Commit() error {
	err := c.story.Commit(c.repoCache.repo)
	if err != nil {
		return err
	}
	return c.notifyUpdated()
}

func (c *StoryCache) CommitAsNeeded() error {
	err := c.story.CommitAsNeeded(c.repoCache.repo)
	if err != nil {
		return err
	}
	return c.notifyUpdated()
}

func (c *StoryCache) NeedCommit() bool {
	return c.story.NeedCommit()
}
