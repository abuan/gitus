package commands

import (
	"fmt"

	"github.com/abuan/gitus/cache"
	"github.com/abuan/gitus/commands/select"
	"github.com/MichaelMure/git-bug/util/interrupt"
	"github.com/spf13/cobra"
)

var (
	titleEditTitle string
)

func runTitleEdit(cmd *cobra.Command, args []string) error {
	backend, err := cache.NewRepoCache(repo)
	if err != nil {
		return err
	}
	defer backend.Close()
	interrupt.RegisterCleaner(backend.Close)

	s, args, err := _select.ResolveStory(backend, args)
	if err != nil {
		return err
	}

	snap := s.Snapshot()

	if titleEditTitle == "" {
			fmt.Println("Empty title, you must provide a new title to the story, aborting.")
			return nil
	}

	if titleEditTitle == snap.Title {
		fmt.Println("No change, aborting.")
	}

	_, err = s.SetTitle(titleEditTitle)
	if err != nil {
		return err
	}

	return s.Commit()
}

var titleEditCmd = &cobra.Command{
	Use:     "edit [<id>]",
	Short:   "Edit a title of a story.",
	PreRunE: loadRepoEnsureUser,
	RunE:    runTitleEdit,
}

func init() {
	titleCmd.AddCommand(titleEditCmd)

	titleEditCmd.Flags().SortFlags = false

	titleEditCmd.Flags().StringVarP(&titleEditTitle, "title", "t", "",
		"Provide a title to describe the issue",
	)
}
