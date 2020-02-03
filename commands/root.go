// Package commands contains the CLI commands
package commands

import (
	"fmt"
	"os"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"

	"github.com/abuan/gitus/story"
	"github.com/MichaelMure/git-bug/identity"
	"github.com/MichaelMure/git-bug/repository"
)

const rootCommandName = "gitus"

// package scoped var to hold the repo after the PreRun execution
var repo repository.ClockedRepo

var ErrNoIdentitySet = errors.New("No identity is set.\n" +
	"To interact with stories, an identity first needs to be created using " +
	"\"gitus user create\"")

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   rootCommandName,
	Short: "A User Story manager embedded in Git.",
	Long: `gitus is a User Story manager embedded in git.

gitus use git objects to store the User Story separated from the files
history. As Story are regular git objects, they can be pushed and pulled from/to
the same git remote your are already using to collaborate with other peoples.

`,

	// For the root command, force the execution of the PreRun
	// even if we just display the help. This is to make sure that we check
	// the repository and give the user early feedback.
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			os.Exit(1)
		}
	},

	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// loadRepo is a pre-run function that load the repository for use in a command
func loadRepo(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to get the current working directory: %q", err)
	}

	repo, err = repository.NewGitRepo(cwd, story.Witnesser)
	if err == repository.ErrNotARepo {
		return fmt.Errorf("%s must be run from within a git repo", rootCommandName)
	}

	if err != nil {
		return err
	}

	return nil
}

// loadRepoEnsureUser is the same as loadRepo, but also ensure that the user has configured
// an identity. Use this pre-run function when an error after using the configured user won't
// do.
func loadRepoEnsureUser(cmd *cobra.Command, args []string) error {
	err := loadRepo(cmd, args)
	if err != nil {
		return err
	}
	
	// GetUserIdentity(repo) retourne un message d'erreur indiquant une commande git bug à utiliser pour 
	// résoudre le problème en cas d'une identité non créée. 
	//Les lignes suivantes permettent de capter l'erreur afin d'afficher la commande propre à gitus
	configs, err := repo.LocalConfig().ReadAll("git-bug.identity")
	if err != nil {
		return err
	}

	if len(configs) == 0 {
		return ErrNoIdentitySet
	}

	_, err = identity.GetUserIdentity(repo)
	if err != nil {
		return err
	}

	return nil
}
