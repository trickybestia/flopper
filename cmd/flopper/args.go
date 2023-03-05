package main

import (
	"github.com/alexflint/go-arg"
)

type args struct {
	Config string `default:"/etc/flopper.conf" arg:"-c,--config" help:"Set path to config file."`
}

func parseArgs() args {
	var args args

	arg.MustParse(&args)

	return args
}
