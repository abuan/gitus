package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/MichaelMure/git-bug/entity"
	"github.com/MichaelMure/git-bug/identity"
	"github.com/MichaelMure/git-bug/repository"
	//"github.com/MichaelMure/git-bug/util/git"
	"github.com/MichaelMure/git-bug/util/process"

	"github.com/abuan/gitus/story"
)

const identityCacheFile = "identity-cache"
const storyCacheFile= "story-cache"

// 1: original format
// 2: added cache for identities with a reference in the story cache
const formatVersion = 2

type ErrInvalidCacheFormat struct {
	message string
}

func (e ErrInvalidCacheFormat) Error() string {
	return e.message
}

var _ repository.RepoCommon = &RepoCache{}

// RepoCache is a cache for a Repository. This cache has multiple functions:
//
// 1. After being loaded, a Story is kept in memory in the cache, allowing for fast
// 		access later.
// 2. The cache maintain in memory and on disk a pre-digested excerpt for each story,
// 		allowing for fast querying the whole set of stories without having to load
//		them individually.
// 3. The cache guarantee that a single instance of a Story is loaded at once, avoiding
// 		loss of data that we could have with multiple copies in the same process.
// 4. The same way, the cache maintain in memory a single copy of the loaded identities.
//
// The cache also protect the on-disk data by locking the git repository for its
// own usage, by writing a lock file. Of course, normal git operations are not
// affected, only gitus related one.
type RepoCache struct {
	// the underlying repo
	repo repository.ClockedRepo

	// excerpt of stories data for all stories
	storyExcerpts map[entity.Id]*StoryExcerpt
	// stories loaded in memory
	stories map[entity.Id]*StoryCache

	// excerpt of identities data for all identities
	identitiesExcerpts map[entity.Id]*IdentityExcerpt
	// identities loaded in memory
	identities map[entity.Id]*IdentityCache

	// the user identity's id, if known
	userIdentityId entity.Id
}

func NewRepoCache(r repository.ClockedRepo) (*RepoCache, error) {
	c := &RepoCache{
		repo:       r,
		identities: make(map[entity.Id]*IdentityCache),
		stories:    make(map[entity.Id]*StoryCache),
	}

	err := c.lock()
	if err != nil {
		return &RepoCache{}, err
	}

	err = c.load()
	if err == nil {
		return c, nil
	}
	if _, ok := err.(ErrInvalidCacheFormat); ok {
		return nil, err
	}

	err = c.buildCache()
	if err != nil {
		return nil, err
	}

	return c, c.write()
}

// LocalConfig give access to the repository scoped configuration
func (c *RepoCache) LocalConfig() repository.Config {
	return c.repo.LocalConfig()
}

// GlobalConfig give access to the git global configuration
func (c *RepoCache) GlobalConfig() repository.Config {
	return c.repo.GlobalConfig()
}

// GetPath returns the path to the repo.
func (c *RepoCache) GetPath() string {
	return c.repo.GetPath()
}

// GetCoreEditor returns the name of the editor that the user has used to configure git.
func (c *RepoCache) GetCoreEditor() (string, error) {
	return c.repo.GetCoreEditor()
}

// GetRemotes returns the configured remotes repositories.
func (c *RepoCache) GetRemotes() (map[string]string, error) {
	return c.repo.GetRemotes()
}

// GetUserName returns the name the the user has used to configure git
func (c *RepoCache) GetUserName() (string, error) {
	return c.repo.GetUserName()
}

// GetUserEmail returns the email address that the user has used to configure git.
func (c *RepoCache) GetUserEmail() (string, error) {
	return c.repo.GetUserEmail()
}

func (c *RepoCache) lock() error {
	lockPath := repoLockFilePath(c.repo)

	err := repoIsAvailable(c.repo)
	if err != nil {
		return err
	}

	f, err := os.Create(lockPath)
	if err != nil {
		return err
	}

	pid := fmt.Sprintf("%d", os.Getpid())
	_, err = f.WriteString(pid)
	if err != nil {
		return err
	}

	return f.Close()
}

func (c *RepoCache) Close() error {
	c.identities = make(map[entity.Id]*IdentityCache)
	c.identitiesExcerpts = nil
	c.stories = make(map[entity.Id]*StoryCache)
	c.storyExcerpts = nil

	lockPath := repoLockFilePath(c.repo)
	return os.Remove(lockPath)
}

// identityUpdated is a callback to trigger when the excerpt of an identity
// changed, that is each time an identity is updated
func (c *RepoCache) identityUpdated(id entity.Id) error {
	i, ok := c.identities[id]
	if !ok {
		panic("missing identity in the cache")
	}

	c.identitiesExcerpts[id] = NewIdentityExcerpt(i.Identity)

	// we only need to write the identity cache
	return c.writeIdentityCache()
}

// load will try to read from the disk all the cache files
func (c *RepoCache) load() error {

	err := c.loadStoryCache()
	if err != nil {
		return err
	}

	return c.loadIdentityCache()
}

// load will try to read from the disk the identity cache file
func (c *RepoCache) loadIdentityCache() error {
	f, err := os.Open(identityCacheFilePath(c.repo))
	if err != nil {
		return err
	}

	decoder := gob.NewDecoder(f)

	aux := struct {
		Version  uint
		Excerpts map[entity.Id]*IdentityExcerpt
	}{}

	err = decoder.Decode(&aux)
	if err != nil {
		return err
	}

	if aux.Version != 2 {
		return ErrInvalidCacheFormat{
			message: fmt.Sprintf("unknown cache format version %v", aux.Version),
		}
	}

	c.identitiesExcerpts = aux.Excerpts
	return nil
}

// write will serialize on disk all the cache files
func (c *RepoCache) write() error {

	err := c.writeStoryCache()
	if err != nil {
		return err
	}

	return c.writeIdentityCache()
}

// write will serialize on disk the identity cache file
func (c *RepoCache) writeIdentityCache() error {
	var data bytes.Buffer

	aux := struct {
		Version  uint
		Excerpts map[entity.Id]*IdentityExcerpt
	}{
		Version:  formatVersion,
		Excerpts: c.identitiesExcerpts,
	}

	encoder := gob.NewEncoder(&data)

	err := encoder.Encode(aux)
	if err != nil {
		return err
	}

	f, err := os.Create(identityCacheFilePath(c.repo))
	if err != nil {
		return err
	}

	_, err = f.Write(data.Bytes())
	if err != nil {
		return err
	}

	return f.Close()
}

func identityCacheFilePath(repo repository.Repo) string {
	//Le dossier de stockage des caches se se nomme "git-bug" car la création de ce dossier est géré par le package github.com/MichaelMure/git-bug/repository
	return path.Join(repo.GetPath(), "git-bug", identityCacheFile)
}

func (c *RepoCache) buildCache() error {
	_, _ = fmt.Fprintf(os.Stderr, "Building identity cache... ")

	c.identitiesExcerpts = make(map[entity.Id]*IdentityExcerpt)

	allIdentities := identity.ReadAllLocalIdentities(c.repo)

	for i := range allIdentities {
		if i.Err != nil {
			return i.Err
		}

		c.identitiesExcerpts[i.Identity.Id()] = NewIdentityExcerpt(i.Identity)
	}

	_, _ = fmt.Fprintln(os.Stderr, "Done.")

	_, _ = fmt.Fprintf(os.Stderr, "Building story cache... ")

	c.storyExcerpts = make(map[entity.Id]*StoryExcerpt)

	allStories := story.ReadAllLocalStories(c.repo)

	for s := range allStories {
		if s.Err != nil {
			return s.Err
		}

		snap := s.Story.Compile()
		c.storyExcerpts[s.Story.Id()] = NewStoryExcerpt(s.Story, &snap)
	}

	_, _ = fmt.Fprintln(os.Stderr, "Done.")

	return nil
}


// Fetch retrieve updates from a remote
// This does not change the local stories or identities state
func (c *RepoCache) Fetch(remote string) (string, error) {
	stdout1, err := identity.Fetch(c.repo, remote)
	if err != nil {
		return stdout1, err
	}

	stdout2, err := story.Fetch(c.repo, remote)
	if err != nil {
		return stdout2, err
	}

	return stdout1 + stdout2, nil
}

// MergeAll will merge all the available remote stories and identities
func (c *RepoCache) MergeAll(remote string) <-chan entity.MergeResult {
	out := make(chan entity.MergeResult)

	// Intercept merge results to update the cache properly
	go func() {
		defer close(out)

		results := identity.MergeAll(c.repo, remote)
		for result := range results {
			out <- result

			if result.Err != nil {
				continue
			}

			switch result.Status {
			case entity.MergeStatusNew, entity.MergeStatusUpdated:
				i := result.Entity.(*identity.Identity)
				c.identitiesExcerpts[result.Id] = NewIdentityExcerpt(i)
			}
		}

		results = story.MergeAll(c.repo, remote)
		for result := range results {
			out <- result

			if result.Err != nil {
				continue
			}

			switch result.Status {
			case entity.MergeStatusNew, entity.MergeStatusUpdated:
				s := result.Entity.(*story.Story)
				snap := s.Compile()
				c.storyExcerpts[result.Id] = NewStoryExcerpt(s, &snap)
			}
		}

		err := c.write()

		// No easy way out here ..
		if err != nil {
			panic(err)
		}
	}()

	return out
}

// Push update a remote with the local changes
func (c *RepoCache) Push(remote string) (string, error) {
	stdout1, err := identity.Push(c.repo, remote)
	if err != nil {
		return stdout1, err
	}

	stdout2, err := story.Push(c.repo, remote)
	if err != nil {
		return stdout2, err
	}

	return stdout1 + stdout2, nil
}

// Pull will do a Fetch + MergeAll
// This function will return an error if a merge fail
func (c *RepoCache) Pull(remote string) error {
	_, err := c.Fetch(remote)
	if err != nil {
		return err
	}

	for merge := range c.MergeAll(remote) {
		if merge.Err != nil {
			return merge.Err
		}
		if merge.Status == entity.MergeStatusInvalid {
			return errors.Errorf("merge failure: %s", merge.Reason)
		}
	}

	return nil
}

func repoLockFilePath(repo repository.Repo) string {
	//Le dossier de stockage des caches se se nomme "git-bug" car la création de ce dossier est géré par le package github.com/MichaelMure/git-bug/repository
	return path.Join(repo.GetPath(), "git-bug", lockfile)
}

// repoIsAvailable check is the given repository is locked by a Cache.
// Note: this is a smart function that will cleanup the lock file if the
// corresponding process is not there anymore.
// If no error is returned, the repo is free to edit.
func repoIsAvailable(repo repository.Repo) error {
	lockPath := repoLockFilePath(repo)

	// Todo: this leave way for a racey access to the repo between the test
	// if the file exist and the actual write. It's probably not a problem in
	// practice because using a repository will be done from user interaction
	// or in a context where a single instance of gitus is already guaranteed
	// (say, a server with the web UI running). But still, that might be nice to
	// have a mutex or something to guard that.

	// Todo: this will fail if somehow the filesystem is shared with another
	// computer. Should add a configuration that prevent the cleaning of the
	// lock file

	f, err := os.Open(lockPath)

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if err == nil {
		// lock file already exist
		buf, err := ioutil.ReadAll(io.LimitReader(f, 10))
		if err != nil {
			return err
		}
		if len(buf) == 10 {
			return fmt.Errorf("the lock file should be < 10 bytes")
		}

		pid, err := strconv.Atoi(string(buf))
		if err != nil {
			return err
		}

		if process.IsRunning(pid) {
			return fmt.Errorf("the repository you want to access is already locked by the process pid %d", pid)
		}

		// The lock file is just laying there after a crash, clean it

		fmt.Println("A lock file is present but the corresponding process is not, removing it.")
		err = f.Close()
		if err != nil {
			return err
		}

		err = os.Remove(lockPath)
		if err != nil {
			return err
		}
	}

	return nil
}

// ResolveIdentity retrieve an identity matching the exact given id
func (c *RepoCache) ResolveIdentity(id entity.Id) (*IdentityCache, error) {
	cached, ok := c.identities[id]
	if ok {
		return cached, nil
	}

	i, err := identity.ReadLocal(c.repo, id)
	if err != nil {
		return nil, err
	}

	cached = NewIdentityCache(c, i)
	c.identities[id] = cached

	return cached, nil
}

// ResolveIdentityExcerpt retrieve a IdentityExcerpt matching the exact given id
func (c *RepoCache) ResolveIdentityExcerpt(id entity.Id) (*IdentityExcerpt, error) {
	e, ok := c.identitiesExcerpts[id]
	if !ok {
		return nil, identity.ErrIdentityNotExist
	}

	return e, nil
}

// ResolveIdentityPrefix retrieve an Identity matching an id prefix.
// It fails if multiple identities match.
func (c *RepoCache) ResolveIdentityPrefix(prefix string) (*IdentityCache, error) {
	// preallocate but empty
	matching := make([]entity.Id, 0, 5)

	for id := range c.identitiesExcerpts {
		if id.HasPrefix(prefix) {
			matching = append(matching, id)
		}
	}

	if len(matching) > 1 {
		return nil, identity.NewErrMultipleMatch(matching)
	}

	if len(matching) == 0 {
		return nil, identity.ErrIdentityNotExist
	}

	return c.ResolveIdentity(matching[0])
}

// ResolveIdentityImmutableMetadata retrieve an Identity that has the exact given metadata on
// one of it's version. If multiple version have the same key, the first defined take precedence.
func (c *RepoCache) ResolveIdentityImmutableMetadata(key string, value string) (*IdentityCache, error) {
	// preallocate but empty
	matching := make([]entity.Id, 0, 5)

	for id, i := range c.identitiesExcerpts {
		if i.ImmutableMetadata[key] == value {
			matching = append(matching, id)
		}
	}

	if len(matching) > 1 {
		return nil, identity.NewErrMultipleMatch(matching)
	}

	if len(matching) == 0 {
		return nil, identity.ErrIdentityNotExist
	}

	return c.ResolveIdentity(matching[0])
}

// AllIdentityIds return all known identity ids
func (c *RepoCache) AllIdentityIds() []entity.Id {
	result := make([]entity.Id, len(c.identitiesExcerpts))

	i := 0
	for _, excerpt := range c.identitiesExcerpts {
		result[i] = excerpt.Id
		i++
	}

	return result
}

func (c *RepoCache) SetUserIdentity(i *IdentityCache) error {
	err := identity.SetUserIdentity(c.repo, i.Identity)
	if err != nil {
		return err
	}

	// Make sure that everything is fine
	if _, ok := c.identities[i.Id()]; !ok {
		panic("SetUserIdentity while the identity is not from the cache, something is wrong")
	}

	c.userIdentityId = i.Id()

	return nil
}

func (c *RepoCache) GetUserIdentity() (*IdentityCache, error) {
	if c.userIdentityId != "" {
		i, ok := c.identities[c.userIdentityId]
		if ok {
			return i, nil
		}
	}

	i, err := identity.GetUserIdentity(c.repo)
	if err != nil {
		return nil, err
	}

	cached := NewIdentityCache(c, i)
	c.identities[i.Id()] = cached
	c.userIdentityId = i.Id()

	return cached, nil
}

func (c *RepoCache) IsUserIdentitySet() (bool, error) {
	return identity.IsUserIdentitySet(c.repo)
}

// NewIdentity create a new identity
// The new identity is written in the repository (commit)
func (c *RepoCache) NewIdentity(name string, email string) (*IdentityCache, error) {
	return c.NewIdentityRaw(name, email, "", "", nil)
}

// NewIdentityFull create a new identity
// The new identity is written in the repository (commit)
func (c *RepoCache) NewIdentityFull(name string, email string, login string, avatarUrl string) (*IdentityCache, error) {
	return c.NewIdentityRaw(name, email, login, avatarUrl, nil)
}

func (c *RepoCache) NewIdentityRaw(name string, email string, login string, avatarUrl string, metadata map[string]string) (*IdentityCache, error) {
	i := identity.NewIdentityFull(name, email, login, avatarUrl)

	for key, value := range metadata {
		i.SetMetadata(key, value)
	}

	err := i.Commit(c.repo)
	if err != nil {
		return nil, err
	}

	if _, has := c.identities[i.Id()]; has {
		return nil, fmt.Errorf("identity %s already exist in the cache", i.Id())
	}

	cached := NewIdentityCache(c, i)
	c.identities[i.Id()] = cached

	// force the write of the excerpt
	err = c.identityUpdated(i.Id())
	if err != nil {
		return nil, err
	}

	return cached, nil
}

//**************************  Fonction pour Story  **************************
// NewStory create a new story
// The new story is written in the repository (commit)
func (c *RepoCache) NewStory(title string, description string, effort int) (*StoryCache, *story.CreateOperation, error) {
	author, err := c.GetUserIdentity()
	if err != nil {
		return nil, nil, err
	}

	return c.NewStoryRaw(author, time.Now().Unix(), title, description,effort, nil)
}


// NewStoryRaw create a new story with metadata for the Create operation.
// The new story is written in the repository (commit)
func (c *RepoCache) NewStoryRaw(author *IdentityCache, unixTime int64, title string, description string, effort int, metadata map[string]string) (*StoryCache, *story.CreateOperation, error) {
	s, op, err := story.Create(author.Identity, unixTime, title, description, effort)
	if err != nil {
		return nil, nil, err
	}

	for key, value := range metadata {
		op.SetMetadata(key, value)
	}

	err = s.Commit(c.repo)
	if err != nil {
		return nil, nil, err
	}

	if _, has := c.stories[s.Id()]; has {
		return nil, nil, fmt.Errorf("story %s already exist in the cache", s.Id())
	}

	cached := NewStoryCache(c, s)
	c.stories[s.Id()] = cached

	// force the write of the excerpt
	err = c.storyUpdated(s.Id())
	if err != nil {
		return nil, nil, err
	}

	return cached, op, nil
}

// storyUpdated is a callback to trigger when the excerpt of a story changed,
// that is each time a story is updated
func (c *RepoCache) storyUpdated(id entity.Id) error {
	s, ok := c.stories[id]
	if !ok {
		panic("missing story in the cache")
	}

	c.storyExcerpts[id] = NewStoryExcerpt(s.story, s.Snapshot())

	// we only need to write the story cache
	return c.writeStoryCache()
}

// load will try to read from the disk the Story cache file
func (c *RepoCache) loadStoryCache() error {
	f, err := os.Open(storyCacheFilePath(c.repo))
	if err != nil {
		return err
	}

	decoder := gob.NewDecoder(f)

	aux := struct {
		Version  uint
		Excerpts map[entity.Id]*StoryExcerpt
	}{}

	err = decoder.Decode(&aux)
	if err != nil {
		return err
	}

	if aux.Version != 2 {
		return ErrInvalidCacheFormat{
			message: fmt.Sprintf("unknown cache format version %v", aux.Version),
		}
	}

	c.storyExcerpts = aux.Excerpts
	return nil
}

func storyCacheFilePath(repo repository.Repo) string {
	//Le dossier de stockage des caches se se nomme "git-bug" car la création de ce dossier est géré par le package github.com/MichaelMure/git-bug/repository
	return path.Join(repo.GetPath(), "git-bug", storyCacheFile)
}

// write will serialize on disk the story cache file
func (c *RepoCache) writeStoryCache() error {
	var data bytes.Buffer

	aux := struct {
		Version  uint
		Excerpts map[entity.Id]*StoryExcerpt
	}{
		Version:  formatVersion,
		Excerpts: c.storyExcerpts,
	}

	encoder := gob.NewEncoder(&data)

	err := encoder.Encode(aux)
	if err != nil {
		return err
	}

	f, err := os.Create(storyCacheFilePath(c.repo))
	if err != nil {
		return err
	}

	_, err = f.Write(data.Bytes())
	if err != nil {
		return err
	}

	return f.Close()
}


// ResolveStory retrieve a story matching the exact given id
func (c *RepoCache) ResolveStory(id entity.Id) (*StoryCache, error) {
	cached, ok := c.stories[id]
	if ok {
		return cached, nil
	}

	b, err := story.ReadLocalStory(c.repo, id)
	if err != nil {
		return nil, err
	}

	cached = NewStoryCache(c, b)
	c.stories[id] = cached

	return cached, nil
}

// ResolveStoryExcerpt retrieve a StoryExcerpt matching the exact given id
func (c *RepoCache) ResolveStoryExcerpt(id entity.Id) (*StoryExcerpt, error) {
	e, ok := c.storyExcerpts[id]
	if !ok {
		return nil, story.ErrStoryNotExist
	}

	return e, nil
}

// ResolveStoryPrefix retrieve a story matching an id prefix. It fails if multiple
// stories match.
func (c *RepoCache) ResolveStoryPrefix(prefix string) (*StoryCache, error) {
	// preallocate but empty
	matching := make([]entity.Id, 0, 5)

	for id := range c.storyExcerpts {
		if id.HasPrefix(prefix) {
			matching = append(matching, id)
		}
	}

	if len(matching) > 1 {
		return nil, story.NewErrMultipleMatchStory(matching)
	}

	if len(matching) == 0 {
		return nil, story.ErrStoryNotExist
	}

	return c.ResolveStory(matching[0])
}

// ResolveStoryCreateMetadata retrieve a story that has the exact given metadata on
// its Create operation, that is, the first operation. It fails if multiple stories
// match.
func (c *RepoCache) ResolveStoryCreateMetadata(key string, value string) (*StoryCache, error) {
	// preallocate but empty
	matching := make([]entity.Id, 0, 5)

	for id, excerpt := range c.storyExcerpts {
		if excerpt.CreateMetadata[key] == value {
			matching = append(matching, id)
		}
	}

	if len(matching) > 1 {
		return nil, story.NewErrMultipleMatchStory(matching)
	}

	if len(matching) == 0 {
		return nil, story.ErrStoryNotExist
	}

	return c.ResolveStory(matching[0])
}


// QueryStories return the id of all Story matching the given Query
func (c *RepoCache) QueryStories(query *Query) []entity.Id {
	if query == nil {
		return c.AllStoriesIds()
	}

	var filtered []*StoryExcerpt

	for _, excerpt := range c.storyExcerpts {
		if query.Match(c, excerpt) {
			filtered = append(filtered, excerpt)
		}
	}

	var sorter sort.Interface

	switch query.OrderBy {
	case OrderById:
		sorter = StoriesById(filtered)
	case OrderByCreation:
		sorter = StoriesByCreationTime(filtered)
	case OrderByEdit:
		sorter = StoriesByEditTime(filtered)
	default:
		panic("missing sort type")
	}

	sort.Sort(sorter)

	result := make([]entity.Id, len(filtered))

	for i, val := range filtered {
		result[i] = val.Id
	}

	return result
}

// AllStoriesIds return all known story ids
func (c *RepoCache) AllStoriesIds() []entity.Id {
	result := make([]entity.Id, len(c.storyExcerpts))

	i := 0
	for _, excerpt := range c.storyExcerpts {
		result[i] = excerpt.Id
		i++
	}

	return result
}