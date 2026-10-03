package storm

// Create creates tables from models of the passed objects.
// If a table already exists then the model is skipped. A
// few rules:
//
//   - Object type name becomes the table name.
//   - The exported model fields become columns.
//   - Field name is the column name.
//   - Field type is mapped to a SQLite column type.
//   - Only primitive Go types may be used as field types.
//   - All columns have NOT NULL constraint.
//   - All columns have DEFAULT set to the zero value of
//     the field's Go type.
//   - By default, the first field in the model is
//     designated the key (ID and PRIMARY KEY).
func (st *Storm) Create(objects ...any) error {
	if !st.IsOpen() {
		return st.errNotOpen()
	}

	st.mutex.Lock()
	defer st.mutex.Unlock()

	st.prepareMapper()

	for _, obj := range objects {
		model, exists, e := st.mapModel(obj)
		if e != nil {
			return st.errForModel(obj, e)
		}

		if exists {
			continue
		}

		e = model.Create(st.db)
		if e != nil {
			return st.errForModel(obj, e)
		}
	}

	return nil
}

// Put inserts or updates (upserts) every object in the
// passed list. Inserting via an object that only
// represents part of a table will cause the unfilled
// columns to default to their zero value. When updating,
// the model's key property is used to identify the row
// but the primary key column is not updated. If a table
// doesn't exist for an object, it will be created.
// However, if an object only maps to part of a table
// you'll need to call [Storm.Create] with an object
// representing the full model first to ensure all columns
// are created.
func (st *Storm) Put[T any](objects ...T) error {
	if !st.IsOpen() {
		return st.errNotOpen()
	}

	st.mutex.Lock()
	defer st.mutex.Unlock()

	st.prepareMapper()

	for _, o := range objects {
		model, exists, e := st.mapModel(o)
		if e != nil {
			return st.errForModel(o, e)
		}

		if !exists {
			e = model.Create(st.db)
			if e != nil {
				return st.errForModel(o, e)
			}
		}

		if e == nil {
			e = model.Upsert(st.db, o)
		}

		if e != nil {
			return st.errForModel(o, e)
		}
	}

	return nil
}

// List returns all rows (objects) in the table associated
// with the passed object. If the table doesn't exist a
// nil or empty result set is returned, not an error.
func (st *Storm) List[T any](object T) (result []T, e error) {
	if !st.IsOpen() {
		return nil, st.errNotOpen()
	}

	st.mutex.Lock()
	defer st.mutex.Unlock()

	st.prepareMapper()

	model, found, e := st.mapModel(object)
	if e != nil {
		goto Err
	}

	if !found {
		return nil, nil
	}

	result, e = model.SelectAll[T](st.db)
	if e != nil {
		goto Err
	}

	return result, nil

Err:
	return nil, st.errForModel(object, e)
}

// Get returns the object (row) with the given ID from
// the table the passed object maps to. If no record is
// found then an error is returned.
func (st *Storm) Get[T, ID any](object T, id ID) (result T, e error) {
	var empty T

	if !st.IsOpen() {
		return empty, st.errNotOpen()
	}

	st.mutex.Lock()
	defer st.mutex.Unlock()

	st.prepareMapper()

	model, found, e := st.mapModel(object)
	if e != nil {
		goto Err
	}

	if !found {
		e = ErrObjectNotFound
		goto Err
	}

	result, found, e = model.SelectById[T](st.db, id)
	if e != nil {
		goto Err
	}

	if !found {
		e = ErrObjectNotFound
		goto Err
	}

	return result, nil

Err:
	return empty, st.errForObject(object, e, id)
}

// Delete removes the object (row) with the given IDs from
// the table the passed object maps to. If no record is
// found then nothing happens.
func (st *Storm) Delete[T, ID any](object T, ids ...ID) error {
	if !st.IsOpen() {
		return st.errNotOpen()
	}

	st.mutex.Lock()
	defer st.mutex.Unlock()

	st.prepareMapper()
	model, found, e := st.mapModel(object)
	if e != nil {
		return st.errForModel(object, e)
	}

	if !found {
		return nil
	}

	for _, id := range ids {
		e = model.DeleteById(st.db, id)
		if e != nil {
			return st.errForObject(object, e, id)
		}
	}

	return nil
}

// Drop removes a table from the database. If the target
// table doesn't exist then nothing happens and no error is
// returned. All model cache entries associated with the
// table are also removed. All table data is deleted in the
// process and there's no way to restore it. To protect
// data, create backups of the database file.
func (st *Storm) Drop(models ...any) error {
	if !st.IsOpen() {
		return st.errNotOpen()
	}

	st.mutex.Lock()
	defer st.mutex.Unlock()

	st.prepareMapper()

	for _, m := range models {
		table, found, e := st.mapModel(m)
		if e != nil {
			return st.errForModel(m, e)
		}

		if !found {
			continue
		}

		st.mapper.ClearTable(table.SqlName)

		e = table.Drop(st.db)
		if e != nil {
			return st.errForModel(m, e)
		}
	}

	return nil
}
