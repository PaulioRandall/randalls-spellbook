package sourcery

// Initialiser is implemented by entities that need to be
// initialised when the application starts up.
type Initialiser interface {
	// Init is called once when the application starts up.
	Init(w *World) error
}

// Unitialiser is implemented by entities that need to be
// cleaned up, unitialised, when the application exits.
type Unitialiser interface {
	// Unit is called once when the application exits.
	// Implement Init if the World object is needed during
	// unitialisation.
	Unit() error
}

// entityList is a list of entities where each may
// optionally implement [Initialiser] and [Unitialiser].
type entityList []any

// Append appends an item to the end of the list the same
// as Go's standard append function, thus, always assign
// the output.
func (el entityList) Append(item any) entityList {
	return append(el, item)
}

// Init initialises all list items that implement
// [Initialiser]. If an error occurs all items that were
// initialised, except the item that produced the error,
// will be unitialised, as though [entityList.Unit] was
// called exclusively for them. Any errors from
// unitialisation are ignored. A panic will bypass
// unitialisation.
func (el entityList) Init(w *World) error {
	var list []any = []any(el)
	var e error

	for i := 0; i < len(list); i++ {
		v, ok := list[i].(Initialiser)
		if !ok {
			continue
		}

		e = v.Init(w)
		if e != nil {
			// Unit all those initialised, from back to front.
			unit(list[:i])
			return e
		}
	}

	return nil
}

// Unit unitialises all list items that implement
// [Unitialiser] in reverse order. If an error occurs then
// unitialisation will continue for the remaining items.
// Only the first error encountered will be returned. A
// panic will bypass remaining unitialisation.
func (el entityList) Unit() error {
	return unit([]any(el))
}

func unit(list []any) (err error) {
	// Unit last item first.
	for i := len(list) - 1; i >= 0; i-- {
		v, ok := list[i].(Unitialiser)
		if !ok {
			continue
		}

		e := v.Unit()
		if err == nil && e != nil {
			err = e
		}
	}

	return err
}

var (
	_ Initialiser = &entityList{}
	_ Unitialiser = &entityList{}
)
