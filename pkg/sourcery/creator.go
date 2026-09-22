package sourcery

import (
	"maps"
	"net/http"
	"reflect"

	"github.com/crgimenes/glaze"

	"github.com/PaulioRandall/randalls-spellbook/pkg/curse"
)

type Portal interface {
	Init(w *World)
	Free()
}
type PortalMap = map[string]Portal
type HandlerMap = map[string]http.Handler
type FuncMap = map[string]any

// ********************************************************

type httpOnlyPortal struct {
	handler http.Handler
}

func (httpOnlyPortal) Init(w *World) {}
func (httpOnlyPortal) Free()         {}
func (hop httpOnlyPortal) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	hop.handler.ServeHTTP(w, r)
}

// ********************************************************

type Creator struct {
	options    AppOptions
	portalMap  PortalMap
	handlerMap HandlerMap
	funcMap    FuncMap
}

func NewCreator() *Creator {
	return &Creator{
		options: AppOptions{
			Hint: glaze.HintNone,
		},
		portalMap:  PortalMap{},
		handlerMap: HandlerMap{},
		funcMap:    FuncMap{},
	}
}

func (wb *Creator) Debug() *Creator {
	wb.options.Debug = true
	return wb
}

func (wb *Creator) Name(
	name string,
) *Creator {
	wb.options.Title = name
	return wb
}

func (wb *Creator) Size(
	width, height int,
) *Creator {
	wb.options.Width = width
	wb.options.Height = height
	return wb
}

func (wb *Creator) AddPortal(port Portal) *Creator {
	funcMap := identifyFuncs(port)
	validateFuncMap(funcMap)
	maps.Copy(wb.funcMap, funcMap)
	wb.portalMap[nameOfPortal(port)] = port
	return wb
}

func validateFuncMap(funcMap FuncMap) {
	for name, f := range funcMap {
		e := ValidateValErrFunc(f)
		if e != nil {
			e = curse.Fmt(
				"'%s' has invalid function signature",
				name,
			).Wraps(e)
			panic(e)
		}
	}
}

func nameOfPortal(port Portal) string {
	typ := reflect.TypeOf(port)

	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	return typ.Name()
}

func (wb *Creator) AddServer(
	path string,
	handler http.Handler,
) *Creator {
	if _, ok := handler.(Portal); ok {
		wb.handlerMap[path] = handler
	} else {
		wb.handlerMap[path] = httpOnlyPortal{
			handler: handler,
		}
	}
	return wb
}

func (wb *Creator) BuildWorld() *World {
	options := wb.options
	options.Handler = marryHandlers(wb.handlerMap)

	for name, _ := range wb.portalMap {
		println(name)
	}

	return buildWorld(options, wb.portalMap, wb.funcMap)
}

func marryHandlers(handlerMap HandlerMap) *http.ServeMux {
	mux := http.NewServeMux()

	for path, handler := range handlerMap {
		mux.Handle(path, handler)
	}

	return mux
}

func gatherFuncs(portalMap PortalMap) FuncMap {
	result := FuncMap{}

	for _, port := range portalMap {
		funcMap := identifyFuncs(port)
		maps.Copy(result, funcMap)
	}

	return result
}

func identifyFuncs(port Portal) FuncMap {
	portVal := reflect.ValueOf(port)
	portTyp := portVal.Type()
	funcMap := FuncMap{}

	for i := 0; i < portVal.NumMethod(); i++ {
		funcVal := portVal.Method(i)
		funcTyp := portTyp.Method(i)

		if !funcTyp.IsExported() {
			continue
		}

		// Ignore Portal & Handler functions.
		n := funcTyp.Name
		if n == "Init" || n == "Free" || n == "ServeHttp" {
			continue
		}

		name := portTyp.Elem().Name() + "." + funcTyp.Name
		funcMap[name] = funcVal.Interface()
	}

	return funcMap
}
