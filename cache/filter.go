package cache

import (
	"strings"

	"github.com/abuan/gitus/story"
)

// Filter is a predicate that match a subset of stories
type Filter func(repoCache *RepoCache, excerpt *StoryExcerpt) bool

// StatusFilter return a Filter that match a story status
func StatusFilter(query string) (Filter, error) {
	status, err := story.StatusFromString(query)
	if err != nil {
		return nil, err
	}

	return func(repoCache *RepoCache, excerpt *StoryExcerpt) bool {
		return excerpt.Status == status
	}, nil
}

// AuthorFilter return a Filter that match a story author
func AuthorFilter(query string) Filter {
	return func(repoCache *RepoCache, excerpt *StoryExcerpt) bool {
		query = strings.ToLower(query)

		// Normal identity
		if excerpt.AuthorId != "" {
			author, ok := repoCache.identitiesExcerpts[excerpt.AuthorId]
			if !ok {
				panic("missing identity in the cache")
			}

			return author.Match(query)
		}

		// Legacy identity support
		return strings.Contains(strings.ToLower(excerpt.LegacyAuthor.Name), query) ||
			strings.Contains(strings.ToLower(excerpt.LegacyAuthor.Login), query)
	}
}

// ActorFilter return a Filter that match a story actor
func ActorFilter(query string) Filter {
	return func(repoCache *RepoCache, excerpt *StoryExcerpt) bool {
		query = strings.ToLower(query)

		for _, id := range excerpt.Actors {
			identityExcerpt, ok := repoCache.identitiesExcerpts[id]
			if !ok {
				panic("missing identity in the cache")
			}

			if identityExcerpt.Match(query) {
				return true
			}
		}
		return false
	}
}

// ParticipantFilter return a Filter that match a story participant
func ParticipantFilter(query string) Filter {
	return func(repoCache *RepoCache, excerpt *StoryExcerpt) bool {
		query = strings.ToLower(query)

		for _, id := range excerpt.Participants {
			identityExcerpt, ok := repoCache.identitiesExcerpts[id]
			if !ok {
				panic("missing identity in the cache")
			}

			if identityExcerpt.Match(query) {
				return true
			}
		}
		return false
	}
}

// TitleFilter return a Filter that match if the title contains the given query
func TitleFilter(query string) Filter {
	return func(repo *RepoCache, excerpt *StoryExcerpt) bool {
		return strings.Contains(
			strings.ToLower(excerpt.Title),
			strings.ToLower(query),
		)
	}
}

// Filters is a collection of Filter that implement a complex filter
type Filters struct {
	Status      []Filter
	Author      []Filter
	Actor       []Filter
	Participant []Filter
	Title       []Filter
}

// Match check if a story match the set of filters
func (f *Filters) Match(repoCache *RepoCache, excerpt *StoryExcerpt) bool {
	if match := f.orMatch(f.Status, repoCache, excerpt); !match {
		return false
	}

	if match := f.orMatch(f.Author, repoCache, excerpt); !match {
		return false
	}

	if match := f.orMatch(f.Participant, repoCache, excerpt); !match {
		return false
	}

	if match := f.orMatch(f.Actor, repoCache, excerpt); !match {
		return false
	}

	if match := f.andMatch(f.Title, repoCache, excerpt); !match {
		return false
	}

	return true
}

// Check if any of the filters provided match the story
func (*Filters) orMatch(filters []Filter, repoCache *RepoCache, excerpt *StoryExcerpt) bool {
	if len(filters) == 0 {
		return true
	}

	match := false
	for _, f := range filters {
		match = match || f(repoCache, excerpt)
	}

	return match
}

// Check if all of the filters provided match the story
func (*Filters) andMatch(filters []Filter, repoCache *RepoCache, excerpt *StoryExcerpt) bool {
	if len(filters) == 0 {
		return true
	}

	match := true
	for _, f := range filters {
		match = match && f(repoCache, excerpt)
	}

	return match
}
