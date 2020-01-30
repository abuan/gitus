package commands

import (
	"github.com/abuan/gitus/cache"
	"github.com/abuan/gitus/commands/select"
	"github.com/MichaelMure/git-bug/util/interrupt"
	"github.com/spf13/cobra"
)

func runStoryStatusClose(cmd *cobra.Command, args []string) error {
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

	_, err = s.Close()
	if err != nil {
		return err
	}

	return s.Commit()
}

var storyCloseCmd = &cobra.Command{
	Use:     "close [<id>]",
	Short:   "Mark a story as closed.",
	PreRunE: loadRepoEnsureUser,
	RunE:    runStoryStatusClose,
}

func init() {
	storyStatusCmd.AddCommand(storyCloseCmd)
}
