package commands

import (
	"fmt"

	"github.com/abuan/gitus/cache"
	"github.com/MichaelMure/git-bug/util/interrupt"
	"github.com/spf13/cobra"
)

var (
	addDescription     string
	addEffort		   int 
)

func runCreateStory(cmd *cobra.Command, args []string) error {
	var err error

	backend, err := cache.NewRepoCache(repo)
	if err != nil {
		return err
	}
	defer backend.Close()
	interrupt.RegisterCleaner(backend.Close)

	s, _, err := backend.NewStory(args[0], addDescription, addEffort)
	if err != nil {
		return err
	}

	fmt.Printf("%s created\n", s.Id().Human())

	return nil
}

var createStoryCmd = &cobra.Command{
	Use:     "create <title>",
	Short:   "Create a new story.",
	Args:     cobra.MinimumNArgs(1),
	PreRunE: loadRepoEnsureUser,
	RunE:    runCreateStory,
}

func init() {
	RootCmd.AddCommand(createStoryCmd)

	createStoryCmd.Flags().SortFlags = false

	createStoryCmd.Flags().StringVarP(&addDescription, "description", "d", "",
		"Provide a message to describe the story",
	)
	createStoryCmd.Flags().IntVarP(&addEffort, "effort", "e", 0,
		"Provide an effort value to evaluate the story",
	)
}
