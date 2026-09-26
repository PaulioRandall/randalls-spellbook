package storm

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/PaulioRandall/randalls-spellbook/pkg/storm/schema"
	//"github.com/PaulioRandall/randalls-spellbook/pkg/sprints"
)

func Test_Storm_prepareTable_1(t *testing.T) {
	// GIVEN model that's not in database.
	// WHEN calling Storm.prepareTable.
	// THEN table of the model's type is created.

	db := openCreateInsert(t, nil)
	defer db.Close()

	e := db.prepareTable(testCheeseMaker{})
	require.NoError(t, e)

	tableSchema, e := schema.Query(db.db, "testCheeseMaker")
	require.NoError(t, e)

	//for _, m := range metadata {
	//sprints.Println(m)
	//}

	require.Equal(t, "testCheeseMaker", tableSchema.Name)
}
