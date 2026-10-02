package storm

// Create creates tables, represented by the passed models,
// within the database. If a table already exists then
// the model is ignored. It's safe to call Create at
// anytime, but doing it all upfront is recommended.
//
//   - Model (struct) name becomes the table name.
//   - The exported model fields become columns.
//   - Field name is the column name.
//   - Field type is mapped to a SQLite column type.
//   - Only primitive types may be used as field types.
//   - All columns have NOT NULL constraint.
//   - All columns have DEFAULT set to the zero value of
//     the field type.
//   - By default, the first field in the model is
//     designated the PRIMARY KEY.
func (st *Storm) Create(models ...any) error {
	if !st.IsOpen() {
		return st.errNotOpen()
	}

	st.mutex.Lock()
	defer st.mutex.Unlock()

	st.prepareMapper()

	for _, m := range models {
		table, exists, e := st.mapModel(m)
		if e != nil {
			return st.errForModel(m, e)
		}

		if exists {
			continue
		}

		e = table.Create(st.db)
		if e != nil {
			return st.errForModel(m, e)
		}
	}

	return nil
}

// Put inserts or updates (upserts) every object in the
// passed list. Inserting via an object that only
// represents part of a table will cause the unfilled
// columns to default to their zero value. When updating,
// the primary key (ID) field is used to identify the row
// but the primary key column is not updated.
func (st *Storm) Put[T any](objects ...T) error {
	if !st.IsOpen() {
		return st.errNotOpen()
	}

	st.mutex.Lock()
	defer st.mutex.Unlock()

	st.prepareMapper()

	for _, o := range objects {
		table, exists, e := st.mapModel(o)
		if e != nil {
			return st.errForModel(o, e)
		}

		if !exists {
			e = table.Create(st.db)
			if e != nil {
				return st.errForModel(o, e)
			}
		}

		if e == nil {
			e = table.Upsert(st.db, o)
		}

		if e != nil {
			return st.errForModel(o, e)
		}
	}

	return nil
}

// List returns all records for the table associated
// with the passed model.
func (st *Storm) List[T any](model T) (result []T, e error) {
	if !st.IsOpen() {
		return nil, st.errNotOpen()
	}

	st.mutex.Lock()
	defer st.mutex.Unlock()

	st.prepareMapper()

	table, found, e := st.mapModel(model)
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
func (st *Storm) Get[T, ID any](model T, id ID) (result T, e error) {
	var empty T

	if !st.IsOpen() {
		return empty, st.errNotOpen()
	}

	st.mutex.Lock()
	defer st.mutex.Unlock()

	st.prepareMapper()

	table, found, e := st.mapModel(model)
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

// Delete removes the records with the given IDs from
// the table associated with the passed model. If no
// record is found then nothing happens. The model's
// type must match a registered type or an error is
// returned.
func (st *Storm) Delete[T, ID any](model T, ids ...ID) (e error) {
	if !st.IsOpen() {
		return st.errNotOpen()
	}

	st.mutex.Lock()
	defer st.mutex.Unlock()

	st.prepareMapper()
	table, found, e := st.mapModel(model)
	if e != nil {
		return st.errForModel(model, e)
	}

	if !found {
		return nil
	}

	for _, id := range ids {
		e = table.DeleteById(st.db, id)
		if e != nil {
			return st.errForObject(model, e, id)
		}
	}

	return nil
}

// Drop removes a table from the database. If the target
// table doesn't exist then nothing happens and no error is
// returned. All model cache entries associated with the
// table are also removed. All data is deleted in the
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

		st.cachedMapper.ClearTable(table.SqlName)

		e = table.Drop(st.db)
		if e != nil {
			return st.errForModel(m, e)
		}
	}

	return nil
}
