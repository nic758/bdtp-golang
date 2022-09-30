package cli

import (
	"os"
	"strings"
)

type Flags []Flag
type Args []string

type Flag interface {
	Get(name string) bool
	Set(val string)
	SetEnv()
}

func (flags Flags) SetEnv() {
	for _, f := range flags {
		f.SetEnv()
	}
}

func (flags Flags) Last() Flag {
	if len(flags) > 0 {
		return flags[len(flags)-1]
	}

	return nil
}

func (flags Flags) Get(name string) *Flag {
	for _, f := range flags {
		if f.Get(name) {
			return &f
		}
	}

	return nil
}

// StringFlag is a flag with type string
type StringFlag struct {
	Name   string
	Usage  string
	EnvVar string
	Hidden bool
	Value  string
}

func (f *StringFlag) Get(name string) bool {
	return strings.Contains(name, f.Name)
}
func (f *StringFlag) Set(val string) {
	f.Value = val
}

func (f *StringFlag) SetEnv() {
	v, _ := os.LookupEnv(f.EnvVar)
	if f.EnvVar != "" && v == "" {
		err := os.Setenv(f.EnvVar, f.Value)
		if err != nil {
			panic(err)
		}
	}
}

func (args Args) Last() string {
	if len(args) > 0 {
		return args[len(args)-1]
	}
	return ""
}
