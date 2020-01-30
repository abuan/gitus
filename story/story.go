package story

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/MichaelMure/git-bug/entity"
	"github.com/MichaelMure/git-bug/identity"
	"github.com/MichaelMure/git-bug/repository"
	"github.com/MichaelMure/git-bug/util/git"
	"github.com/MichaelMure/git-bug/util/lamport"

)

const storiesRefPattern = "refs/stories/"
const storiesRemoteRefPattern = "refs/remotes/%s/stories/"

const opsEntryName = "ops"
const rootEntryName = "root"

const createClockEntryPrefix = "create-clock-"
const createClockEntryPattern = "create-clock-%d"
const editClockEntryPrefix = "edit-clock-"
const editClockEntryPattern = "edit-clock-%d"

var ErrStoryNotExist = errors.New("story doesn't exist")

func NewErrMultipleMatchStory(matching []entity.Id) *entity.ErrMultipleMatch {
	return entity.NewErrMultipleMatch("story", matching)
}

func NewErrMultipleMatchOp(matching []entity.Id) *entity.ErrMultipleMatch {
	return entity.NewErrMultipleMatch("operation", matching)
}

var _ Interface = &Story{}
var _ entity.Interface = &Story{}

// Story : le type de base de notre projet
type Story struct {

	// A Lamport clock is a logical clock that allow to order event
	// inside a distributed system.
	// It must be the first field in this struct due to https://github.com/golang/go/issues/599
	createTime lamport.Time
	editTime   lamport.Time

	// Id used as unique identifier
	id entity.Id

	lastCommit git.Hash
	rootPack   git.Hash

	// all the committed operations
	packs []OperationPack

	// a temporary pack of operations used for convenience to pile up new operations
	// before a commit
	staging OperationPack
}

// NewStory create a new Story
func NewStory() *Story {
	return &Story{}
}

// Append an operation into the staging area, to be committed later
func (story *Story) Append(op Operation) {
	story.staging.Append(op)
}


func (story *Story) NeedCommit() bool {
	return !story.staging.IsEmpty()
}

// Id return the Story identifier
func (story *Story) Id() entity.Id {
	if story.id == "" {
		// simply panic as it would be a coding error
		// (using an id of a story not stored yet)
		panic("no id yet")
	}
	return story.id
}

// Lookup for the very first operation of the story.
// For a valid Story, this operation should be a CreateOp
func (story *Story) FirstOp() Operation {
	for _, pack := range story.packs {
		for _, op := range pack.Operations {
			return op
		}
	}

	if !story.staging.IsEmpty() {
		return story.staging.Operations[0]
	}

	return nil
}

// Compile a story in a easily usable snapshot
func (story *Story) Compile() Snapshot {
	snap := Snapshot{
		id:     story.id,
	}

	it := NewOperationIterator(story)

	for it.Next() {
		op := it.Value()
		op.Apply(&snap)
		snap.Operations = append(snap.Operations, op)
	}

	return snap
}

// CreateLamportTime return the Lamport time of creation
func (story *Story) CreateLamportTime() lamport.Time {
	return story.createTime
}

// EditLamportTime return the Lamport time of the last edit
func (story *Story) EditLamportTime() lamport.Time {
	return story.editTime
}

func (story *Story) CommitAsNeeded(repo repository.ClockedRepo) error {
	if !story.NeedCommit() {
		return nil
	}
	return story.Commit(repo)
}

// Lookup for the very last operation of the story.
// For a valid Story, should never be nil
func (story *Story) LastOp() Operation {
	if !story.staging.IsEmpty() {
		return story.staging.Operations[len(story.staging.Operations)-1]
	}

	if len(story.packs) == 0 {
		return nil
	}

	lastPack := story.packs[len(story.packs)-1]

	if len(lastPack.Operations) == 0 {
		return nil
	}

	return lastPack.Operations[len(lastPack.Operations)-1]
}

// FindLocalStory find an existing Story matching a prefix
func FindLocalStory(repo repository.ClockedRepo, prefix string) (*Story, error) {
	ids, err := ListLocalIds(repo)

	if err != nil {
		return nil, err
	}

	// preallocate but empty
	matching := make([]entity.Id, 0, 5)

	for _, id := range ids {
		if id.HasPrefix(prefix) {
			matching = append(matching, id)
		}
	}

	if len(matching) == 0 {
		return nil, errors.New("no matching story found.")
	}

	if len(matching) > 1 {
		return nil, NewErrMultipleMatchStory(matching)
	}

	return ReadLocalStory(repo, matching[0])
}

// ReadLocalStory will read a local story from its hash
func ReadLocalStory(repo repository.ClockedRepo, id entity.Id) (*Story, error) {
	ref := storiesRefPattern + id.String()
	return readStory(repo, ref)
}

// ReadRemoteStory will read a remote story from its hash
func ReadRemoteStory(repo repository.ClockedRepo, remote string, id string) (*Story, error) {
	ref := fmt.Sprintf(storiesRemoteRefPattern, remote) + id
	return readStory(repo, ref)
}

// readStory will read and parse a Story from git
func readStory(repo repository.ClockedRepo, ref string) (*Story, error) {
	refSplit := strings.Split(ref, "/")
	id := entity.Id(refSplit[len(refSplit)-1])

	if err := id.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid ref ")
	}

	hashes, err := repo.ListCommits(ref)

	// TODO: this is not perfect, it might be a command invoke error
	if err != nil {
		return nil, ErrStoryNotExist
	}

	story := Story{
		id:       id,
		editTime: 0,
	}

	// Load each OperationPack
	for _, hash := range hashes {
		entries, err := repo.ListEntries(hash)
		if err != nil {
			return nil, errors.Wrap(err, "can't list git tree entries")
		}

		story.lastCommit = hash

		var opsEntry repository.TreeEntry
		opsFound := false
		var rootEntry repository.TreeEntry
		rootFound := false
		var createTime uint64
		var editTime uint64

		for _, entry := range entries {
			if entry.Name == opsEntryName {
				opsEntry = entry
				opsFound = true
				continue
			}
			if entry.Name == rootEntryName {
				rootEntry = entry
				rootFound = true
			}
			if strings.HasPrefix(entry.Name, createClockEntryPrefix) {
				n, err := fmt.Sscanf(entry.Name, createClockEntryPattern, &createTime)
				if err != nil {
					return nil, errors.Wrap(err, "can't read create lamport time")
				}
				if n != 1 {
					return nil, fmt.Errorf("could not parse create time lamport value")
				}
			}
			if strings.HasPrefix(entry.Name, editClockEntryPrefix) {
				n, err := fmt.Sscanf(entry.Name, editClockEntryPattern, &editTime)
				if err != nil {
					return nil, errors.Wrap(err, "can't read edit lamport time")
				}
				if n != 1 {
					return nil, fmt.Errorf("could not parse edit time lamport value")
				}
			}
		}

		if !opsFound {
			return nil, errors.New("invalid tree, missing the ops entry")
		}
		if !rootFound {
			return nil, errors.New("invalid tree, missing the root entry")
		}

		if story.rootPack == "" {
			story.rootPack = rootEntry.Hash
			story.createTime = lamport.Time(createTime)
		}

		// Due to rebase, edit Lamport time are not necessarily ordered
		if editTime > uint64(story.editTime) {
			story.editTime = lamport.Time(editTime)
		}

		// Update the clocks
		if err := repo.WitnessCreate(story.createTime); err != nil {
			return nil, errors.Wrap(err, "failed to update create lamport clock")
		}
		if err := repo.WitnessEdit(story.editTime); err != nil {
			return nil, errors.Wrap(err, "failed to update edit lamport clock")
		}

		data, err := repo.ReadData(opsEntry.Hash)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read git blob data")
		}

		opp := &OperationPack{}
		err = json.Unmarshal(data, &opp)

		if err != nil {
			return nil, errors.Wrap(err, "failed to decode OperationPack json")
		}

		// tag the pack with the commit hash
		opp.commitHash = hash

		story.packs = append(story.packs, *opp)
	}

	// Make sure that the identities are properly loaded
	resolver := identity.NewSimpleResolver(repo)
	err = story.EnsureIdentities(resolver)
	if err != nil {
		return nil, err
	}

	return &story, nil
}

type StreamedStory struct {
	Story *Story
	Err error
}

// ReadAllLocalStories read and parse all local stories
func ReadAllLocalStories(repo repository.ClockedRepo) <-chan StreamedStory {
	return readAllStories(repo, storiesRefPattern)
}

// ReadAllRemoteStories read and parse all remote stories for a given remote
func ReadAllRemoteStories(repo repository.ClockedRepo, remote string) <-chan StreamedStory {
	refPrefix := fmt.Sprintf(storiesRemoteRefPattern, remote)
	return readAllStories(repo, refPrefix)
}

// Read and parse all available story with a given ref prefix
func readAllStories(repo repository.ClockedRepo, refPrefix string) <-chan StreamedStory {
	out := make(chan StreamedStory)

	go func() {
		defer close(out)

		refs, err := repo.ListRefs(refPrefix)
		if err != nil {
			out <- StreamedStory{Err: err}
			return
		}

		for _, ref := range refs {
			s, err := readStory(repo, ref)

			if err != nil {
				out <- StreamedStory{Err: err}
				return
			}

			out <- StreamedStory{Story: s}
		}
	}()

	return out
}

// ListLocalIds list all the available local story ids
func ListLocalIds(repo repository.Repo) ([]entity.Id, error) {
	refs, err := repo.ListRefs(storiesRefPattern)
	if err != nil {
		return nil, err
	}

	return refsToIds(refs), nil
}

func refsToIds(refs []string) []entity.Id {
	ids := make([]entity.Id, len(refs))

	for i, ref := range refs {
		split := strings.Split(ref, "/")
		ids[i] = entity.Id(split[len(split)-1])
	}

	return ids
}


// Commit write the staging area in Git and move the operations to the packs
func (story *Story) Commit(repo repository.ClockedRepo) error {

	if !story.NeedCommit() {
		return fmt.Errorf("can't commit a story with no pending operation")
	}

	if err := story.Validate(); err != nil {
		return errors.Wrap(err, "can't commit a story with invalid data")
	}

	// Write the Ops as a Git blob containing the serialized array
	hash, err := story.staging.Write(repo)
	if err != nil {
		return err
	}

	if story.rootPack == "" {
		story.rootPack = hash
	}

	// Make a Git tree referencing this blob
	tree := []repository.TreeEntry{
		// the last pack of ops
		{ObjectType: repository.Blob, Hash: hash, Name: opsEntryName},
		// always the first pack of ops (might be the same)
		{ObjectType: repository.Blob, Hash: story.rootPack, Name: rootEntryName},
	}

	// Store the logical clocks as well
	// --> edit clock for each OperationPack/commits
	// --> create clock only for the first OperationPack/commits
	//
	// To avoid having one blob for each clock value, clocks are serialized
	// directly into the entry name
	emptyBlobHash, err := repo.StoreData([]byte{})
	if err != nil {
		return err
	}

	story.editTime, err = repo.EditTimeIncrement()
	if err != nil {
		return err
	}

	tree = append(tree, repository.TreeEntry{
		ObjectType: repository.Blob,
		Hash:       emptyBlobHash,
		Name:       fmt.Sprintf(editClockEntryPattern, story.editTime),
	})
	if story.lastCommit == "" {
		story.createTime, err = repo.CreateTimeIncrement()
		if err != nil {
			return err
		}

		tree = append(tree, repository.TreeEntry{
			ObjectType: repository.Blob,
			Hash:       emptyBlobHash,
			Name:       fmt.Sprintf(createClockEntryPattern, story.createTime),
		})
	}

	// Store the tree
	hash, err = repo.StoreTree(tree)
	if err != nil {
		return err
	}

	// Write a Git commit referencing the tree, with the previous commit as parent
	if story.lastCommit != "" {
		hash, err = repo.StoreCommitWithParent(hash, story.lastCommit)
	} else {
		hash, err = repo.StoreCommit(hash)
	}

	if err != nil {
		return err
	}

	story.lastCommit = hash

	// if it was the first commit, use the commit hash as story id
	if story.id == "" {
		story.id = entity.Id(hash)
	}

	// Create or update the Git reference for this story
	// When pushing later, the remote will ensure that this ref update
	// is fast-forward, that is no data has been overwritten
	ref := fmt.Sprintf("%s%s", storiesRefPattern, story.id)
	err = repo.UpdateRef(ref, hash)

	if err != nil {
		return err
	}

	story.staging.commitHash = hash
	story.packs = append(story.packs, story.staging)
	story.staging = OperationPack{}

	return nil
}

// Validate check if the Story data is valid
func (story *Story) Validate() error {
	// non-empty
	if len(story.packs) == 0 && story.staging.IsEmpty() {
		return fmt.Errorf("story has no operations")
	}

	// check if each pack and operations are valid
	for _, pack := range story.packs {
		if err := pack.Validate(); err != nil {
			return err
		}
	}

	// check if staging is valid if needed
	if !story.staging.IsEmpty() {
		if err := story.staging.Validate(); err != nil {
			return errors.Wrap(err, "staging")
		}
	}

	// The very first Op should be a CreateOp
	firstOp := story.FirstOp()
	if firstOp == nil || firstOp.base().OperationType != CreateOp {
		return fmt.Errorf("first operation should be a Create op")
	}

	// The story Id should be the hash of the first commit
	if len(story.packs) > 0 && string(story.packs[0].commitHash) != story.id.String() {
		return fmt.Errorf("story id should be the first commit hash")
	}

	// Check that there is no more CreateOp op
	// Check that there is no colliding operation's ID
	it := NewOperationIterator(story)
	createCount := 0
	ids := make(map[entity.Id]struct{})
	for it.Next() {
		if it.Value().base().OperationType == CreateOp {
			createCount++
		}
		if _, ok := ids[it.Value().Id()]; ok {
			return fmt.Errorf("id collision: %s", it.Value().Id())
		}
		ids[it.Value().Id()] = struct{}{}
	}

	if createCount != 1 {
		return fmt.Errorf("only one Create op allowed")
	}

	return nil
}

// Merge a different version of the same story by rebasing operations of this story
// that are not present in the other on top of the chain of operations of the
// other version.
func (story *Story) Merge(repo repository.Repo, other Interface) (bool, error) {
	var otherStory = storyFromInterface(other)

	// Note: a faster merge should be possible without actually reading and parsing
	// all operations pack of our side.
	// Reading the other side is still necessary to validate remote data, at least
	// for new operations

	if story.id != otherStory.id {
		return false, errors.New("merging unrelated stories is not supported")
	}

	if len(otherStory.staging.Operations) > 0 {
		return false, errors.New("merging a story with a non-empty staging is not supported")
	}

	if story.lastCommit == "" || otherStory.lastCommit == "" {
		return false, errors.New("can't merge a story that has never been stored")
	}

	ancestor, err := repo.FindCommonAncestor(story.lastCommit, otherStory.lastCommit)
	if err != nil {
		return false, errors.Wrap(err, "can't find common ancestor")
	}

	ancestorIndex := 0
	newPacks := make([]OperationPack, 0, len(story.packs))

	// Find the root of the rebase
	for i, pack := range story.packs {
		newPacks = append(newPacks, pack)

		if pack.commitHash == ancestor {
			ancestorIndex = i
			break
		}
	}

	if len(otherStory.packs) == ancestorIndex+1 {
		// Nothing to rebase, return early
		return false, nil
	}

	// get other story's extra packs
	for i := ancestorIndex + 1; i < len(otherStory.packs); i++ {
		// clone is probably not necessary
		newPack := otherStory.packs[i].Clone()

		newPacks = append(newPacks, newPack)
		story.lastCommit = newPack.commitHash
	}

	// rebase our extra packs
	for i := ancestorIndex + 1; i < len(story.packs); i++ {
		pack := story.packs[i]

		// get the referenced git tree
		treeHash, err := repo.GetTreeHash(pack.commitHash)

		if err != nil {
			return false, err
		}

		// create a new commit with the correct ancestor
		hash, err := repo.StoreCommitWithParent(treeHash, story.lastCommit)

		if err != nil {
			return false, err
		}

		// replace the pack
		newPack := pack.Clone()
		newPack.commitHash = hash
		newPacks = append(newPacks, newPack)

		// update the story
		story.lastCommit = hash
	}

	story.packs = newPacks

	// Update the git ref
	err = repo.UpdateRef(storiesRefPattern+story.id.String(), story.lastCommit)
	if err != nil {
		return false, err
	}

	return true, nil
}
