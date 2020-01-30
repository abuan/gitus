package commands

import (
	"fmt"
	"os"

	"github.com/MichaelMure/git-bug/bridge/core/auth"
	"github.com/abuan/gitus/cache"
	"github.com/MichaelMure/git-bug/identity"
	"github.com/MichaelMure/git-bug/util/interrupt"
	"github.com/spf13/cobra"
)

func runUserAdopt(cmd *cobra.Command, args []string) error {
	backend, err := cache.NewRepoCache(repo)
	if err != nil {
		return err
	}
	defer backend.Close()
	interrupt.RegisterCleaner(backend.Close)

	prefix := args[0]

	i, err := backend.ResolveIdentityPrefix(prefix)
	if err != nil {
		return err
	}

	_, err = backend.GetUserIdentity()
	if err == identity.ErrNoIdentitySet {
		err = auth.ReplaceDefaultUser(repo, i.Id())
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	err = backend.SetUserIdentity(i)
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintf(os.Stderr, "Your identity is now: %s\n", i.DisplayName())

	return nil
}

var userAdoptCmd = &cobra.Command{
	Use:     "adopt <user-id>",
	Short:   "Adopt an existing identity as your own.",
	PreRunE: loadRepo,
	RunE:    runUserAdopt,
	Args:    cobra.ExactArgs(1),
}

func init() {
	userCmd.AddCommand(userAdoptCmd)
	userAdoptCmd.Flags().SortFlags = false
}
