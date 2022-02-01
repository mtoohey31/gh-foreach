package main

import (
	"reflect"
	"regexp"

	"github.com/alecthomas/kong"
)

var cmd struct {
	Clean clean `cmd:"" help:"Remove cache."`
	Run   run   `cmd:"" help:"Execute a command."`
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
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
