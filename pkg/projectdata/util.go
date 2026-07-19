package projectdata

type entityWithId interface {
	EntityId() EntityId
}

func deleteFromSlice[T any](slice []T, i int) []T {
	return append(slice[0:i], slice[i+1:]...)
}

func findEntityIndexById[T entityWithId](slice []T, id EntityId) int {
	for i, v := range slice {
		if v.EntityId() == id {
			return i
		}
	}

	return -1
}

func findEntityById[T entityWithId](slice []T, id EntityId, defaultItem T) T {
	i := findEntityIndexById(slice, id)

	if i < 0 {
		return defaultItem
	}

	return slice[i]
}

func updateEntity[T entityWithId](slice []T, entity T) bool {
	i := findEntityIndexById(slice, entity.EntityId())

	if i < 0 {
		return false
	}

	slice[i] = entity
	return true
}
