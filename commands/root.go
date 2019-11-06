// Package commands contains the CLI commands
package commands

import (
  "fmt"
  "os"
  "github.com/spf13/cobra"
)

// Name of root command
const rootCommandName = "gitus"


// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
  Use:   rootCommandName,
  Short: "A Git version of Jira",
  Long: `Gitus is an Agile project manager with Git fonctionnalities.
  
  @TODO print a longer description of Gitus's aim`,
  // Uncomment the following line if your bare application
  // has an action associated with it:
  Run: func(cmd *cobra.Command, args []string) { 
    fmt.Println("Hello gitus application")
  },
}


// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
