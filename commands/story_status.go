package commands

import (
	"fmt"

	"github.com/abuan/gitus/cache"
	"github.com/abuan/gitus/commands/select"
	"github.com/MichaelMure/git-bug/util/interrupt"
	"github.com/spf13/cobra"
)

func runStoryStatus(cmd *cobra.Command, args []string) error {
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

	fmt.Println(snap.Status)

	return nil
}

var storyStatusCmd = &cobra.Command{
	Use:     "status [<id>]",
	Short:   "Display or change a story status.",
	PreRunE: loadRepo,
	RunE:    runStoryStatus,
}

func init() {
	RootCmd.AddCommand(storyStatusCmd)
}
