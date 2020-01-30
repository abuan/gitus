package commands

import (
	"fmt"
	"strings"

	"github.com/abuan/gitus/cache"
	_select "github.com/abuan/gitus/commands/select"
	"github.com/MichaelMure/git-bug/util/interrupt"
	"github.com/spf13/cobra"
)

var (
	showFieldsQuery string
)

func runShowStory(cmd *cobra.Command, args []string) error {
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

	snapshot := s.Snapshot()


	if showFieldsQuery != "" {
		switch showFieldsQuery {
		case "author":
			fmt.Printf("%s\n", snapshot.Author.DisplayName())
		case "authorEmail":
			fmt.Printf("%s\n", snapshot.Author.Email())
		case "createTime":
			fmt.Printf("%s\n", snapshot.CreatedAt)
		case "humanId":
			fmt.Printf("%s\n", snapshot.Id().Human())
		case "id":
			fmt.Printf("%s\n", snapshot.Id())
		case "actors":
			for _, a := range snapshot.Actors {
				fmt.Printf("%s\n", a.DisplayName())
			}
		case "participants":
			for _, p := range snapshot.Participants {
				fmt.Printf("%s\n", p.DisplayName())
			}
		case "shortId":
			fmt.Printf("%s\n", snapshot.Id().Human())
		case "status":
			fmt.Printf("%s\n", snapshot.Status)
		case "title":
			fmt.Printf("%s\n", snapshot.Title)
		default:
			return fmt.Errorf("\nUnsupported field: %s\n", showFieldsQuery)
		}

		return nil
	}

	// Header
	fmt.Printf("[%s] %s %s\nDescription : %s \nEffort : %d\n\n",
		snapshot.Status,
		snapshot.Id().Human(),
		snapshot.Title,
		snapshot.Description,
		snapshot.Effort,
	)

	fmt.Printf("%s opened this story %s\n\n",
		snapshot.Author.DisplayName(),
		snapshot.CreatedAt,
	)

	// Actors
	var actors = make([]string, len(snapshot.Actors))
	for i := range snapshot.Actors {
		actors[i] = snapshot.Actors[i].DisplayName()
	}

	fmt.Printf("actors: %s\n",
		strings.Join(actors, ", "),
	)

	// Participants
	var participants = make([]string, len(snapshot.Participants))
	for i := range snapshot.Participants {
		participants[i] = snapshot.Participants[i].DisplayName()
	}

	fmt.Printf("participants: %s\n\n",
		strings.Join(participants, ", "),
	)

	return nil
}

var showCmd = &cobra.Command{
	Use:     "show [<id>]",
	Short:   "Display the details of a Story.",
	PreRunE: loadRepo,
	RunE:    runShowStory,
}

func init() {
	RootCmd.AddCommand(showCmd)
	showCmd.Flags().StringVarP(&showFieldsQuery, "field", "f", "",
		"Select field to display. Valid values are [author,authorEmail,createTime,humanId,id,shortId,status,title,actors,participants]")
}
