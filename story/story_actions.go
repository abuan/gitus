package story

import (
	"fmt"
	"strings"

	"github.com/MichaelMure/git-bug/entity"
	"github.com/MichaelMure/git-bug/repository"
	"github.com/pkg/errors"
)

// Fetch retrieve updates from a remote
// This does not change the local stories state
func Fetch(repo repository.Repo, remote string) (string, error) {
	remoteRefSpec := fmt.Sprintf(storiesRemoteRefPattern, remote)
	fetchRefSpec := fmt.Sprintf("%s*:%s*", storiesRefPattern, remoteRefSpec)

	return repo.FetchRefs(remote, fetchRefSpec)
}

// Push update a remote with the local changes
func Push(repo repository.Repo, remote string) (string, error) {
	return repo.PushRefs(remote, storiesRefPattern+"*")
}

// Pull will do a Fetch + MergeAll
// This function will return an error if a merge fail
func Pull(repo repository.ClockedRepo, remote string) error {
	_, err := Fetch(repo, remote)
	if err != nil {
		return err
	}

	for merge := range MergeAll(repo, remote) {
		if merge.Err != nil {
			return merge.Err
		}
		if merge.Status == entity.MergeStatusInvalid {
			return errors.Errorf("merge failure: %s", merge.Reason)
		}
	}

	return nil
}

// MergeAll will merge all the available remote story:
//
// - If the remote has new commit, the local story is updated to match the same history
//   (fast-forward update)
// - if the local story has new commits but the remote don't, nothing is changed
// - if both local and remote story have new commits (that is, we have a concurrent edition),
//   new local commits are rewritten at the head of the remote history (that is, a rebase)
func MergeAll(repo repository.ClockedRepo, remote string) <-chan entity.MergeResult {
	out := make(chan entity.MergeResult)

	go func() {
		defer close(out)

		remoteRefSpec := fmt.Sprintf(storiesRemoteRefPattern, remote)
		remoteRefs, err := repo.ListRefs(remoteRefSpec)

		if err != nil {
			out <- entity.MergeResult{Err: err}
			return
		}

		for _, remoteRef := range remoteRefs {
			refSplit := strings.Split(remoteRef, "/")
			id := entity.Id(refSplit[len(refSplit)-1])

			if err := id.Validate(); err != nil {
				out <- entity.NewMergeInvalidStatus(id, errors.Wrap(err, "invalid ref").Error())
				continue
			}

			remoteStory, err := readStory(repo, remoteRef)

			if err != nil {
				out <- entity.NewMergeInvalidStatus(id, errors.Wrap(err, "remote story is not readable").Error())
				continue
			}

			// Check for error in remote data
			if err := remoteStory.Validate(); err != nil {
				out <- entity.NewMergeInvalidStatus(id, errors.Wrap(err, "remote story is invalid").Error())
				continue
			}

			localRef := storiesRefPattern + remoteStory.Id().String()
			localExist, err := repo.RefExist(localRef)

			if err != nil {
				out <- entity.NewMergeError(err, id)
				continue
			}

			// the story is not local yet, simply create the reference
			if !localExist {
				err := repo.CopyRef(remoteRef, localRef)

				if err != nil {
					out <- entity.NewMergeError(err, id)
					return
				}

				out <- entity.NewMergeStatus(entity.MergeStatusNew, id, remoteStory)
				continue
			}

			localStory, err := readStory(repo, localRef)

			if err != nil {
				out <- entity.NewMergeError(errors.Wrap(err, "local story is not readable"), id)
				return
			}

			updated, err := localStory.Merge(repo, remoteStory)

			if err != nil {
				out <- entity.NewMergeInvalidStatus(id, errors.Wrap(err, "merge failed").Error())
				return
			}

			if updated {
				out <- entity.NewMergeStatus(entity.MergeStatusUpdated, id, localStory)
			} else {
				out <- entity.NewMergeStatus(entity.MergeStatusNothing, id, localStory)
			}
		}
	}()

	return out
}
