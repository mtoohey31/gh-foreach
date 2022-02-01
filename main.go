package main

import (
	"reflect"
	"regexp"

	"github.com/alecthomas/kong"
)

type Repo struct {
	Name      string
	URL       string
	Clone_URL string
}

var cmd struct {
	Clean Clean `cmd:"" help:"Remove cache."`
	Run   Run   `cmd:"" help:"Execute a command."`
}

type regexpMapper bool

func (m regexpMapper) Decode(ctx *kong.DecodeContext, target reflect.Value) error {
	var regexString string
	if err := ctx.Scan.PopValueInto("regex", &regexString); err != nil {
		return err
	}
	r, err := regexp.CompilePOSIX(regexString)
	if err != nil {
		return err
	}
	target.Set(reflect.ValueOf(r))
	return nil
}

func main() {
	ctx := kong.Parse(&cmd, kong.TypeMapper(reflect.TypeOf(&regexp.Regexp{}), regexpMapper(false)))
	// err := ctx.Run()
	err := ctx.Run(kong.TypeMapper(reflect.TypeOf(&regexp.Regexp{}), regexpMapper(false)))
	ctx.FatalIfErrorf(err)
}
