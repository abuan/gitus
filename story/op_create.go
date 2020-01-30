package story

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/MichaelMure/git-bug/entity"
	"github.com/MichaelMure/git-bug/identity"
	"github.com/MichaelMure/git-bug/util/text"
)

var _ Operation = &CreateOperation{}

// CreateOperation define the initial creation of a story 
type CreateOperation struct {
	OpBase
	Title   string     `json:"title"`
	Description string `json:"description"`
	Effort		int	   `json:"effort"`
	Status		Status `json:"status"`
}

func (op *CreateOperation) base() *OpBase {
	return &op.OpBase
}

func (op *CreateOperation) Id() entity.Id {
	return idOperation(op)
}

func (op *CreateOperation) Apply(snapshot *Snapshot) {
	snapshot.addActor(op.Author)
	snapshot.addParticipant(op.Author)

	snapshot.Title = op.Title
	snapshot.Description = op.Description
	snapshot.Effort = op.Effort
	snapshot.Author = op.Author
	snapshot.Status = op.Status
	snapshot.CreatedAt = op.Time()
}

//Validate : Vérifie si les informations passées lors de la création sont correctes
func (op *CreateOperation) Validate() error {
	if err := opBaseValidate(op, CreateOp); err != nil {
		return err
	}

	if text.Empty(op.Title) {
		return fmt.Errorf("title is empty")
	}

	if strings.Contains(op.Title, "\n") {
		return fmt.Errorf("title should be a single line")
	}

	if !text.Safe(op.Title) {
		return fmt.Errorf("title is not fully printable")
	}

	if !text.Safe(op.Description) {
		return fmt.Errorf("message is not fully printable")
	}
	//Vérification sur valeur effort (to improve to match fibonacci value)
	if op.Effort < 0 {
		return fmt.Errorf("effort value can't be negative")
	}
	//Vérification du statut
	if isValid := op.Status.Validate(); isValid != nil {
		return fmt.Errorf("status")
	}

	return nil
}

// UnmarshalJSON is a two step JSON unmarshaling
// This workaround is necessary to avoid the inner OpBase.MarshalJSON
// overriding the outer op's MarshalJSON
func (op *CreateOperation) UnmarshalJSON(data []byte) error {
	// Unmarshal OpBase and the op separately

	base := OpBase{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return err
	}

	aux := struct {
		Title   string     `json:"title"`
		Description string `json:"description"`
		Effort int 		   `json:"effort"`
		Status	Status 	   `json:"status"`
	}{}

	err = json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	op.OpBase = base
	op.Title = aux.Title
	op.Description = aux.Description
	op.Effort = aux.Effort
	op.Status = aux.Status

	return nil
}

// Sign post method for gqlgen
func (op *CreateOperation) IsAuthored() {}

func NewCreateOp(author identity.Interface, unixTime int64, title, message string, effort int, status Status) *CreateOperation {
	return &CreateOperation{
		OpBase:  newOpBase(CreateOp, author, unixTime),
		Title:   title,
		Description: message,
		Effort: effort,
		Status: status,
	}
}

// Convenience function to apply the operation
func Create(author identity.Interface, unixTime int64, title, message string, effort int) (*Story, *CreateOperation, error) {
	newStory := NewStory()
	createOp := NewCreateOp(author, unixTime, title, message, effort,OpenStatus)
	if err := createOp.Validate(); err != nil {
		return nil, createOp, err
	}

	newStory.Append(createOp)

	return newStory, createOp, nil
}
