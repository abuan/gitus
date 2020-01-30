package commands

import (
	"fmt"
	//"strings"
	"github.com/spf13/cobra"

	"github.com/abuan/gitus/cache"
	"github.com/MichaelMure/git-bug/util/interrupt"
)

var (
	lsStatusQuery      []string
	lsAuthorQuery      []string
	lsParticipantQuery []string
	lsTitleQuery       []string
	lsActorQuery       []string
	lsSortBy           string
)

func runLsStory(cmd *cobra.Command, args []string) error {
	backend, err := cache.NewRepoCache(repo)
	if err != nil {
		return err
	}
	defer backend.Close()
	interrupt.RegisterCleaner(backend.Close)

	
	var query *cache.Query
	query, err = lsQueryFromFlags()
	if err != nil {
		return err
	}

	allIds := backend.QueryStories(query)

	if len(allIds)>0{
		fmt.Println("ID\t\tTitle\t\tEffort\tAuthor\tStatus")
		fmt.Println("--\t\t-----\t\t------\t------\t------")
	}
	for _, id := range allIds {
		s, err := backend.ResolveStoryExcerpt(id)
		if err != nil {
			return err
		}

		var name string
		if s.AuthorId != "" {
			author, err := backend.ResolveIdentityExcerpt(s.AuthorId)
			if err != nil {
				name = "<missing author data>"
			} else {
				name = author.DisplayName()
			}
		} else {
			name = s.LegacyAuthor.DisplayName()
		}
	
		fmt.Printf("%s\t\t%s\t\t%d\t%s\t[%s]\n",
			s.Id.Human(),
			s.Title,
			s.Effort,
			name,
			s.Status,
		)
	}

	return nil
}

// Transform the command flags into a query
func lsQueryFromFlags() (*cache.Query, error) {
	query := cache.NewQuery()

	for _, status := range lsStatusQuery {
		f, err := cache.StatusFilter(status)
		if err != nil {
			return nil, err
		}
		query.Status = append(query.Status, f)
	}

	for _, title := range lsTitleQuery {
		f := cache.TitleFilter(title)
		query.Title = append(query.Title, f)
	}

	for _, author := range lsAuthorQuery {
		f := cache.AuthorFilter(author)
		query.Author = append(query.Author, f)
	}

	for _, actor := range lsActorQuery {
		f := cache.ActorFilter(actor)
		query.Actor = append(query.Actor, f)
	}

	for _, participant := range lsParticipantQuery {
		f := cache.ParticipantFilter(participant)
		query.Participant = append(query.Participant, f)
	}

	switch lsSortBy {
	case "id":
		query.OrderBy = cache.OrderById
	case "creation":
		query.OrderBy = cache.OrderByCreation
	case "edit":
		query.OrderBy = cache.OrderByEdit
	default:
		return nil, fmt.Errorf("unknown sort flag %s", lsSortBy)
	}

	return query, nil
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List stories.",
	Long: `Display a summary of each stroies.

You can pass an additional query to filter and order the list. This query can be expressed with flags.`,
	Example: `List closed stories sorted by creation with flags:
gitus ls --status closed --by creation
`,
	PreRunE: loadRepo,
	RunE:    runLsStory,
}

func init() {
	RootCmd.AddCommand(lsCmd)

	lsCmd.Flags().SortFlags = false

	lsCmd.Flags().StringSliceVarP(&lsStatusQuery, "status", "s", nil,
		"Filter by status. Valid values are [open,closed]")
	lsCmd.Flags().StringSliceVarP(&lsAuthorQuery, "author", "a", nil,
		"Filter by author")
	lsCmd.Flags().StringSliceVarP(&lsParticipantQuery, "participant", "p", nil,
		"Filter by participant")
	lsCmd.Flags().StringSliceVarP(&lsActorQuery, "actor", "A", nil,
		"Filter by actor")
	lsCmd.Flags().StringSliceVarP(&lsTitleQuery, "title", "t", nil,
		"Filter by title")
	lsCmd.Flags().StringVarP(&lsSortBy, "by", "b", "creation",
		"Sort the results by a characteristic. Valid values are [id,creation,edit]")
}
