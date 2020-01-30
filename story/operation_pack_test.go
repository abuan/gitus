package story

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/MichaelMure/git-bug/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOperationPackSerialize(t *testing.T) {
	opp := &OperationPack{}

	rene := identity.NewBare("René Descartes", "rene@descartes.fr")
	createOp := NewCreateOp(rene, time.Now().Unix(), "title", "descript", 3, OpenStatus)
	setTitleOp := NewSetTitleOp(rene, time.Now().Unix(), "title2", "title1")
	setStatusOp := NewSetStatusOp(rene, time.Now().Unix(), ClosedStatus)

	opp.Append(createOp)
	opp.Append(setTitleOp)
	opp.Append(setStatusOp)

	opMeta := NewSetTitleOp(rene, time.Now().Unix(), "title3", "title2")
	opMeta.SetMetadata("key", "value")
	opp.Append(opMeta)

	assert.Equal(t, 1, len(opMeta.Metadata))

	data, err := json.Marshal(opp)
	assert.NoError(t, err)

	var opp2 *OperationPack
	err = json.Unmarshal(data, &opp2)
	assert.NoError(t, err)

	ensureIDs(t, opp)

	assert.Equal(t, opp, opp2)
}

func ensureIDs(t *testing.T, opp *OperationPack) {
	for _, op := range opp.Operations {
		id := op.Id()
		require.NoError(t, id.Validate())
		id = op.GetAuthor().Id()
		require.NoError(t, id.Validate())
	}
}
