package nidorino

type Formatter struct {
	template string
}

func Given(template string) *Formatter {
	return &Formatter{
		template: template,
	}
}

func (form *Formatter) Fmt(key, value any) *Formatter {
	// TODO
	return form
}

func (form *Formatter) Join(
	key, delim string,
	values ...any,
) *Formatter {
	// TODO
	return form
}

func (form *Formatter) String() string {
	// TODO
	return ""
}
