package cli

import "strings"

type Command struct {
	Name        string
	ShortName   string
	Aliases     []string
	Usage       string
	Description string
	Action      ActionFunc
	Flags       Flags
}

func (c Command) HasName(name string) bool {
	for _, n := range c.Names() {
		if n == name {
			return true
		}
	}

	return false
}

func (c Command) Names() []string {
	names := []string{c.Name}
	return append(names, c.Aliases...)
}

func (c Command) init(args Args) {
	for _, a := range args {
		kv := strings.Split(a, "=")
		f := c.Flags.Get(kv[0])

		if f != nil && len(kv) > 1 {
			(*f).Set(kv[1])
		}
	}

	c.Flags.SetEnv()
}

func (c Command) Run(ctx *Context, args Args) (err error) {
	c.init(args)
	ctx.Command = c
	flagRemover := func(args Args) Args {
		temp := make(Args, 0)
		for _, a := range args {
			if !strings.Contains(a, "--") {
				temp = append(temp, a)
			}
		}

		return temp
	}

	ctx.Args = flagRemover(args)

	return HandleAction(c.Action, ctx)
}
