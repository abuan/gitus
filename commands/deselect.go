package commands

import (
	"github.com/abuan/gitus/cache"
	"github.com/abuan/gitus/commands/select"
	"github.com/MichaelMure/git-bug/util/interrupt"
	"github.com/spf13/cobra"
)

func runDeselect(cmd *cobra.Command, args []string) error {
	backend, err := cache.NewRepoCache(repo)
	if err != nil {
		return err
	}
	defer backend.Close()
	interrupt.RegisterCleaner(backend.Close)

	err = _select.Clear(backend)
	if err != nil {
		return err
	}

	return nil
}

var deselectCmd = &cobra.Command{
	Use:   "deselect",
	Short: "Clear the implicitly selected story.",
	Example: `gitus  select 2f15
gitus status
gitus show
gitus deselect
`,
	PreRunE: loadRepo,
	RunE:    runDeselect,
}

func init() {
	RootCmd.AddCommand(deselectCmd)
	deselectCmd.Flags().SortFlags = false
}
