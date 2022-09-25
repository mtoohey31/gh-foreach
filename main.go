// cspell:ignore kongplete

package main

import (
	"os"
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
	parser := kong.Must(&cmd,
		kong.Description("Automatically clone and execute commands across multiple GitHub repositories."),
		kong.TypeMapper(reflect.TypeOf(&regexp.Regexp{}), regexpMapper(false)))
	ctx, err := parser.Parse(os.Args[1:])
	parser.FatalIfErrorf(err)
	err = ctx.Run()
	ctx.FatalIfErrorf(err)
}
