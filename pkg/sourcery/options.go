package sourcery

import (
	"net/http"
	"reflect"
)

// DefaultOptions is a slice of options accessible from any
// importing package. Use the [Append] function to add to
// it. Then construct your app with the [DefaultApp]
// function.
var DefaultOptions []Option

// Append adds the options to the default options object.
func Append(options ...Option) {
	DefaultOptions = append(DefaultOptions, options...)
}

// Option is passed to [NewApp] to apply a configuration.
type Option func(*config)

type config struct {
	debug     bool
	title     string
	width     int
	height    int
	html      string
	handlers  map[string]http.Handler
	entities  map[string]entityEntry
	listeners []Listener
}

// OptionDebug enables or disables debugging.
func OptionDebug(enable bool) Option {
	return func(c *config) {
		c.debug = enable
	}
}

// OptionTitle sets application title.
func OptionTitle(title string) Option {
	return func(c *config) {
		c.title = title
	}
}

// OptionSize sets the initial size of the application
// window.
func OptionSize(width, height int) Option {
	return func(c *config) {
		c.width = width
		c.height = height
	}
}

// OptionHtml sets the initial HTML of the webview.
func OptionHtml(html string) Option {
	return func(c *config) {
		c.html = html
	}
}

// OptionServer adds a http.handler for the specified path.
func OptionServer(path string, handler http.Handler) Option {
	return func(c *config) {
		c.handlers[path] = handler
	}
}

// OptionEntity adds an entity who's functions will be
// bound to the webview. It uses the name of the
// concrete type as the entity name.
func OptionEntity(entity any) Option {
	t := reflect.ValueOf(entity).Type()

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return OptionNamedEntity(t.Name(), entity)
}

// OptionNamedEntity adds a named entity who's functions
// are accessible from the webview.
func OptionNamedEntity(name string, entity any) Option {
	return func(c *config) {
		funcs := parseEntity(entity)
		c.entities[name] = entityEntry{
			entity: entity,
			funcs:  funcs,
		}
	}
}

// OptionListener adds a listener of applications events.
// Upon event, listeners are called in the order they are
// registered.
func OptionListener(l Listener) Option {
	return func(c *config) {
		c.listeners = append(c.listeners, l)
	}
}

func parseEntity(entity any) entityFuncs {
	entityVal := reflect.ValueOf(entity)
	functions := entityFuncs{}

	for i := 0; i < entityVal.NumMethod(); i++ {
		funcTyp := entityVal.Type().Method(i)

		// We always ignore unexported functions.
		if !funcTyp.IsExported() {
			continue
		}

		// Check the function signature.
		f := entityVal.Method(i).Interface()
		e := ValidateFunction(f)
		if e != nil {
			panic("Function '" + funcTyp.Name + "' has invalid signature")
		}

		functions[funcTyp.Name] = f
	}

	return functions
}
