package application

type Application struct {
	tokens []string
}

func New() *Application {
	return &Application{}
}

func (app *Application) ConsumeToken(token string) bool {
	// TODO: Replace with actual checker.
	return true
}

func (app *Application) GetMediaPath(id string) string {
	// TODO
	return ""
}
