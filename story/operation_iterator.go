package story

type OperationIterator struct {
	story       *Story
	packIndex int
	opIndex   int
}

func NewOperationIterator(story Interface) *OperationIterator {
	return &OperationIterator{
		story:       storyFromInterface(story),
		packIndex: 0,
		opIndex:   -1,
	}
}

func (it *OperationIterator) Next() bool {
	// Special case of the staging area
	if it.packIndex == len(it.story.packs) {
		pack := it.story.staging
		it.opIndex++
		return it.opIndex < len(pack.Operations)
	}

	if it.packIndex >= len(it.story.packs) {
		return false
	}

	pack := it.story.packs[it.packIndex]

	it.opIndex++

	if it.opIndex < len(pack.Operations) {
		return true
	}

	// Note: this iterator doesn't handle the empty pack case
	it.opIndex = 0
	it.packIndex++

	// Special case of the non-empty staging area
	if it.packIndex == len(it.story.packs) && len(it.story.staging.Operations) > 0 {
		return true
	}

	return it.packIndex < len(it.story.packs)
}

func (it *OperationIterator) Value() Operation {
	// Special case of the staging area
	if it.packIndex == len(it.story.packs) {
		pack := it.story.staging

		if it.opIndex >= len(pack.Operations) {
			panic("Iterator is not valid anymore")
		}

		return pack.Operations[it.opIndex]
	}

	if it.packIndex >= len(it.story.packs) {
		panic("Iterator is not valid anymore")
	}

	pack := it.story.packs[it.packIndex]

	if it.opIndex >= len(pack.Operations) {
		panic("Iterator is not valid anymore")
	}

	return pack.Operations[it.opIndex]
}
