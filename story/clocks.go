package story

import (
	"github.com/MichaelMure/git-bug/repository"
)

// Witnesser will read all the available Story to recreate the different logical
// clocks
func Witnesser(repo repository.ClockedRepo) error {
	for s := range ReadAllLocalStories(repo) {
		if s.Err != nil {
			return s.Err
		}

		err := repo.WitnessCreate(s.Story.createTime)
		if err != nil {
			return err
		}

		err = repo.WitnessEdit(s.Story.editTime)
		if err != nil {
			return err
		}
	}

	return nil
}
