package _select

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/abuan/gitus/story"
	"github.com/abuan/gitus/cache"
	"github.com/MichaelMure/git-bug/entity"
	"github.com/MichaelMure/git-bug/repository"
)

const selectFile = "select"

var ErrNoValidId = errors.New("you must provide a story id or use the \"select\" command first")

func selectFilePath(repo repository.RepoCommon) string {
	return path.Join(repo.GetPath(), "git-bug", selectFile)
}


func ResolveStory(repo *cache.RepoCache, args []string) (*cache.StoryCache, []string, error) {
	// At first, try to use the first argument as a story prefix
	if len(args) > 0 {
		s, err := repo.ResolveStoryPrefix(args[0])

		if err == nil {
			return s, args[1:], nil
		}

		if err != story.ErrStoryNotExist {
			return nil, nil, err
		}
	}

	// first arg is not a valid story prefix, we can safely use the preselected story if any
	s, err := selected(repo)

	// selected story is invalid
	if err == story.ErrStoryNotExist {
		// we clear the selected story
		err = Clear(repo)
		if err != nil {
			return nil, nil, err
		}
		return nil, nil, ErrNoValidId
	}

	// another error when reading the story
	if err != nil {
		return nil, nil, err
	}

	// story is successfully retrieved
	if s != nil {
		return s, args, nil
	}

	// no selected story and no valid first argument
	return nil, nil, ErrNoValidId
}

// Select will select a story for future use
func Select(repo *cache.RepoCache, id entity.Id) error {
	selectPath := selectFilePath(repo)

	f, err := os.OpenFile(selectPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	_, err = f.WriteString(id.String())
	if err != nil {
		return err
	}

	return f.Close()
}

// Clear will clear the selected story, if any
func Clear(repo *cache.RepoCache) error {
	selectPath := selectFilePath(repo)

	return os.Remove(selectPath)
}

func selected(repo *cache.RepoCache) (*cache.StoryCache, error) {
	selectPath := selectFilePath(repo)

	f, err := os.Open(selectPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	buf, err := ioutil.ReadAll(io.LimitReader(f, 100))
	if err != nil {
		return nil, err
	}
	if len(buf) == 100 {
		return nil, fmt.Errorf("the select file should be < 100 bytes")
	}

	id := entity.Id(buf)
	if err := id.Validate(); err != nil {
		err = os.Remove(selectPath)
		if err != nil {
			return nil, errors.Wrap(err, "error while removing invalid select file")
		}

		return nil, fmt.Errorf("select file in invalid, removing it")
	}

	s, err := repo.ResolveStory(id)
	if err != nil {
		return nil, err
	}

	err = f.Close()
	if err != nil {
		return nil, err
	}

	return s, nil
}
