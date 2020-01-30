package commands

import (
	"github.com/abuan/gitus/cache"
	"github.com/abuan/gitus/commands/select"
	"github.com/MichaelMure/git-bug/util/interrupt"
	"github.com/spf13/cobra"
)
var(
	effortNewEffort int
)

func runStoryEffortEdit(cmd *cobra.Command, args []string) error {
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

	_, err = s.EditEffort(effortNewEffort)
	if err != nil {
		return err
	}

	return s.Commit()
}

var editEffortCmd = &cobra.Command{
	Use:     "edit [<id>]",
	Short:   "Change story effort.",
	PreRunE: loadRepoEnsureUser,
	RunE:    runStoryEffortEdit,
}

func init() {
	storyEffortCmd.AddCommand(editEffortCmd)

	editEffortCmd.Flags().SortFlags = false

	editEffortCmd.Flags().IntVarP(&effortNewEffort, "effort", "e", 0,
		"Provide an effort to evaluate the story",
	)
}
