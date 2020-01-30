package story

import (

	"github.com/MichaelMure/git-bug/entity"
)

type TimelineItem interface {
	// ID return the identifier of the item
	Id() entity.Id
}