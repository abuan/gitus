package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/abuan/gitus/cache"
	"github.com/MichaelMure/git-bug/util/interrupt"
)

func runLsID(cmd *cobra.Command, args []string) error {

	backend, err := cache.NewRepoCache(repo)
	if err != nil {
		return err
	}
	defer backend.Close()
	interrupt.RegisterCleaner(backend.Close)

	var prefix = ""
	if len(args) != 0 {
		prefix = args[0]
	}

	for _, id := range backend.AllStoriesIds() {
		if prefix == "" || id.HasPrefix(prefix) {
			fmt.Println(id.Human())
		}
	}

	return nil
}

var listStoryIDCmd = &cobra.Command{
	Use:     "ls-id [<prefix>]",
	Short:   "List story identifiers.",
	PreRunE: loadRepo,
	RunE:    runLsID,
}

func init() {
	RootCmd.AddCommand(listStoryIDCmd)
}
