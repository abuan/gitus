package commands

import (
	"fmt"

	"github.com/abuan/gitus/cache"
	"github.com/abuan/gitus/commands/select"
	"github.com/MichaelMure/git-bug/util/interrupt"
	"github.com/spf13/cobra"
)

func runTitle(cmd *cobra.Command, args []string) error {
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

	fmt.Println(snap.Title)

	return nil
}

var titleCmd = &cobra.Command{
	Use:     "title [<id>]",
	Short:   "Display or change a title of a story.",
	PreRunE: loadRepo,
	RunE:    runTitle,
}

func init() {
	RootCmd.AddCommand(titleCmd)

	titleCmd.Flags().SortFlags = false
}
