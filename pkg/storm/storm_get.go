package storm

import (
	"github.com/PaulioRandall/randalls-spellbook/pkg/storm/modtab"
)

// List returns all records for the table associated
// with the passed model.
func (st *Storm) List[T any](model T) (result []T, e error) {
	if !st.IsOpen() {
		return nil, st.errNotOpen()
	}

	table, found, e := modtab.Map(st.db, model)
	if e != nil {
		goto Err
	}

	if !found {
		return nil, nil
	}

	result, e = table.Select[T](st.db)
	if e != nil {
		goto Err
	}

	return result, nil

Err:
	return nil, st.errForModel(model, e)
}

// Get returns the record with the given id from
// the table associated with the passed model. If no
// record is found then an error is returned. The model's
// type must match a registered type or an error is
// returned.
func (st *Storm) Get[T, ID any](
	model T,
	id ID,
) (result T, e error) {
	var empty T

	if !st.IsOpen() {
		return empty, st.errNotOpen()
	}

	table, found, e := modtab.Map(st.db, model)
	if e != nil {
		goto Err
	}

	if !found {
		e = ErrObjectNotFound
		goto Err
	}

	result, found, e = table.SelectById[T](st.db, id)
	if e != nil {
		goto Err
	}

	if !found {
		e = ErrObjectNotFound
		goto Err
	}

	return result, nil

Err:
	return empty, st.errForObject(model, e, id)
}
