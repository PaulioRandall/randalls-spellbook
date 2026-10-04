package stormy

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Stormy_Custom_1(t *testing.T) {
	// When performing a custom operation
	// then the operation is called
	// and the context is not dead
	// and error is returned

	st := openStormyDatabase(t)
	defer st.Close()

	errMeh := errors.New("Meh")

	_, e := st.Custom(func(ctx OperationContext) (any, error) {
		// Calling this will cause panic if dead.
		ctx.Database()
		return nil, errMeh
	})

	require.ErrorIs(t, e, ErrForDatabase)
	require.ErrorIs(t, e, errMeh)
}

func Test_Stormy_Custom_2(t *testing.T) {
	// When performing a custom operation
	// if the operation does not return an error
	// then the result value will be returned

	st := openStormyDatabase(t)
	defer st.Close()

	exp := TestTable{
		Id:     1,
		Name:   "Alice",
		Rating: 1.11,
	}

	act, e := st.Custom(func(ctx OperationContext) (TestTable, error) {
		return exp, nil
	})

	require.NoError(t, e)
	require.Equal(t, exp, act)
}

func Test_Stormy_Custom_3(t *testing.T) {
	// When performing a custom operation
	// if an operation context is called after the operation
	// ends
	// then panic ensues

	st := openStormyDatabase(t)
	defer st.Close()

	var escapedCtx OperationContext

	st.Custom(func(ctx OperationContext) (any, error) {
		escapedCtx = ctx
		return nil, nil
	})

	require.Panics(t, func() {
		escapedCtx.Database()
	})
}
