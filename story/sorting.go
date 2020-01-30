package story

type StoriesByCreationTime []*Story

func (b StoriesByCreationTime) Len() int {
	return len(b)
}

func (b StoriesByCreationTime) Less(i, j int) bool {
	if b[i].createTime < b[j].createTime {
		return true
	}

	if b[i].createTime > b[j].createTime {
		return false
	}

	// When the logical clocks are identical, that means we had a concurrent
	// edition. In this case we rely on the timestamp. While the timestamp might
	// be incorrect due to a badly set clock, the drift in sorting is bounded
	// by the first sorting using the logical clock. That means that if users
	// synchronize their stories regularly, the timestamp will rarely be used, and
	// should still provide a kinda accurate sorting when needed.
	return b[i].FirstOp().Time().Before(b[j].FirstOp().Time())
}

func (b StoriesByCreationTime) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

type StoriesByEditTime []*Story

func (b StoriesByEditTime) Len() int {
	return len(b)
}

func (b StoriesByEditTime) Less(i, j int) bool {
	if b[i].editTime < b[j].editTime {
		return true
	}

	if b[i].editTime > b[j].editTime {
		return false
	}

	// When the logical clocks are identical, that means we had a concurrent
	// edition. In this case we rely on the timestamp. While the timestamp might
	// be incorrect due to a badly set clock, the drift in sorting is bounded
	// by the first sorting using the logical clock. That means that if users
	// synchronize their stories regularly, the timestamp will rarely be used, and
	// should still provide a kinda accurate sorting when needed.
	return b[i].LastOp().Time().Before(b[j].LastOp().Time())
}

func (b StoriesByEditTime) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
