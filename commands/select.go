package commands

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/abuan/gitus/cache"
	"github.com/abuan/gitus/commands/select"
	"github.com/MichaelMure/git-bug/util/interrupt"
)

func runSelect(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("You must provide a story id")
	}

	backend, err := cache.NewRepoCache(repo)
	if err != nil {
		return err
	}
	defer backend.Close()
	interrupt.RegisterCleaner(backend.Close)

	prefix := args[0]

	s, err := backend.ResolveStoryPrefix(prefix)
	if err != nil {
		return err
	}

	err = _select.Select(backend, s.Id())
	if err != nil {
		return err
	}

	fmt.Printf("selected story %s: %s\n", s.Id().Human(), s.Snapshot().Title)

	return nil
}

var selectCmd = &cobra.Command{
	Use:   "select <id>",
	Short: "Select a story for implicit use in future commands.",
	Example: `gitus select 2f15
gitus status
`,
	Long: `Select a story for implicit use in future commands.

This command allows you to omit any story <id> argument, for example:
  gitus show
instead of
  gitus show 2f153ca

The complementary command is "gitus deselect" performing the opposite operation.
`,
	PreRunE: loadRepo,
	RunE:    runSelect,
}

func init() {
	RootCmd.AddCommand(selectCmd)
	selectCmd.Flags().SortFlags = false
}
