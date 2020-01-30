package commands

import (
	"github.com/abuan/gitus/cache"
	"github.com/abuan/gitus/commands/select"
	"github.com/MichaelMure/git-bug/util/interrupt"
	"github.com/spf13/cobra"
)

func runStoryStatusOpen(cmd *cobra.Command, args []string) error {
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

	_, err = s.Open()
	if err != nil {
		return err
	}

	return s.Commit()
}

var openCmdStory = &cobra.Command{
	Use:     "open [<id>]",
	Short:   "Mark a story as open.",
	PreRunE: loadRepoEnsureUser,
	RunE:    runStoryStatusOpen,
}

func init() {
	storyStatusCmd.AddCommand(openCmdStory)
}
