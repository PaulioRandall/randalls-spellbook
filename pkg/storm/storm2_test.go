package storm

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/PaulioRandall/randalls-spellbook/pkg/sin"
	//"github.com/PaulioRandall/randalls-spellbook/pkg/storm/mapper"
	"github.com/PaulioRandall/randalls-spellbook/pkg/storm/schema"
)

func openStormDatabase(t *testing.T) *Storm {
	st := New(":memory:")

	e := st.Open()
	require.NoError(t, e)

	return st
}

func Test_Storm_Open_Close_1(t *testing.T) {
	st := openStormDatabase(t)
	defer st.Close()

	require.Equal(t, true, st.IsOpen())

	e := st.Close()
	require.NoError(t, e)
	require.Equal(t, false, st.IsOpen())
}

func Test_Storm_Create_1(t *testing.T) {
	type Player struct {
		Id      int
		Name    string
		Rating  float64
		ignored *int
	}

	st := openStormDatabase(t)
	defer st.Close()

	e := st.Create(Player{})
	require.NoError(t, sin.Stack(e, false, true))

	tableSchema, e := schema.QuerySqliteSchema(st.db, "Player")
	require.NoError(t, sin.Stack(e, false, true))
	require.Equal(t, "Player", tableSchema.Name)
	require.Equal(t, "Player", tableSchema.TableName)
}

/*
func Test_Storm_Table_1(t *testing.T) {
	type Player struct {
		Id      int
		Name    string
		Rating  float64
		ignored *int
	}

	st := openStormDatabase(t)
	defer st.Close()
	e := st.Create(Player{})
	require.NoError(t, sin.Stack(e, false, true))

	act, e := st.Table(Player{})
	require.NoError(t, sin.Stack(e, false, true))

	exp := mapper.SqlTable{
		Name: "Player",
		Columns: []mapper.SqlColumn{
			mapper.SqlColumn{
				Name:       "Id",
				Type:       "INTEGER",
				Default:    int(0),
				PrimaryKey: true,
			},
			mapper.SqlColumn{
				Name:       "Name",
				Type:       "TEXT",
				Default:    "''",
				PrimaryKey: false,
			},
			mapper.SqlColumn{
				Name:       "Rating",
				Type:       "REAL",
				Default:    float64(0),
				PrimaryKey: false,
			},
		},
	}

	require.Equal(t, exp, act)
}
*/
