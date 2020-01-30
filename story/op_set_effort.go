package story

import (
	"encoding/json"
	"fmt"

	"github.com/MichaelMure/git-bug/entity"
	"github.com/MichaelMure/git-bug/identity"
	"github.com/MichaelMure/git-bug/util/timestamp"

)

var _ Operation = &SetEffortOperation{}

// SetEffortOperation will change the title of a story
type SetEffortOperation struct {
	OpBase
	Effort int `json:"effort"`
	Was   int `json:"was"`
}

func (op *SetEffortOperation) base() *OpBase {
	return &op.OpBase
}

func (op *SetEffortOperation) Id() entity.Id {
	return idOperation(op)
}

func (op *SetEffortOperation) Apply(snapshot *Snapshot) {
	snapshot.Effort = op.Effort
	snapshot.addActor(op.Author)

	item := &SetEffortTimelineItem{
		id:       op.Id(),
		Author:   op.Author,
		UnixTime: timestamp.Timestamp(op.UnixTime),
		Effort:   op.Effort,
		Was:      op.Was,
	}

	snapshot.Timeline = append(snapshot.Timeline, item)
}

func (op *SetEffortOperation) Validate() error {
	if err := opBaseValidate(op, SetEffortOp); err != nil {
		return err
	}

	//Vérification des efforts 
	if op.Effort < 0 {
		return fmt.Errorf("effort value can't be negative")
	}

	if op.Was < 0 {
		return fmt.Errorf("previous effort value can't be negative")
	}

	return nil
}

// UnmarshalJSON is a two step JSON unmarshaling
// This workaround is necessary to avoid the inner OpBase.MarshalJSON
// overriding the outer op's MarshalJSON
func (op *SetEffortOperation) UnmarshalJSON(data []byte) error {
	// Unmarshal OpBase and the op separately

	base := OpBase{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return err
	}

	aux := struct {
		Effort int `json:"effort"`
		Was   int `json:"was"`
	}{}

	err = json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	op.OpBase = base
	op.Effort = aux.Effort
	op.Was = aux.Was

	return nil
}

// Sign post method for gqlgen
func (op *SetEffortOperation) IsAuthored() {}

func NewSetEffortOp(author identity.Interface, unixTime int64, effort int, was int) *SetEffortOperation {
	return &SetEffortOperation{
		OpBase: newOpBase(SetEffortOp, author, unixTime),
		Effort:  effort,
		Was:    was,
	}
}

type SetEffortTimelineItem struct {
	id       entity.Id
	Author   identity.Interface
	UnixTime timestamp.Timestamp
	Effort   int
	Was      int
}

func (s SetEffortTimelineItem) Id() entity.Id {
	return s.id
}

// Sign post method for gqlgen
func (s *SetEffortTimelineItem) IsAuthored() {}

// Convenience function to apply the operation
func SetEffort(s Interface, author identity.Interface, unixTime int64, effort int) (*SetEffortOperation, error) {
	it := NewOperationIterator(s)

	var lastEffortOp Operation
	for it.Next() {
		op := it.Value()
		if op.base().OperationType == SetEffortOp {
			lastEffortOp = op
		}
	}

	var was int
	if lastEffortOp != nil {
		was = lastEffortOp.(*SetEffortOperation).Effort
	} else {
		was = s.FirstOp().(*CreateOperation).Effort
	}

	setEffortOp := NewSetEffortOp(author, unixTime, effort, was)

	if err := setEffortOp.Validate(); err != nil {
		return nil, err
	}

	s.Append(setEffortOp)
	return setEffortOp, nil
}
