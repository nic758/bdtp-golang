package cli

type Context struct {
	App     *App
	Command Command
	Args    Args
}

func NewContext(app *App) *Context {
	c := &Context{App: app}

	return c
}
