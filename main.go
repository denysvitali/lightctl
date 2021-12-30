package main

import (
	"github.com/alexflint/go-arg"
	"github.com/sirupsen/logrus"
	"os"
)

type Args struct {
	Device   string       `arg:"-d"`
	Debug    *bool        `arg:"-D,--debug"`
	Increase *IncreaseCmd `arg:"subcommand:increase"`
	Decrease *DecreaseCmd `arg:"subcommand:decrease"`
	Set      *SetCmd      `arg:"subcommand:set"`
	List     *ListCmd     `arg:"subcommand:list"`
}

var args Args

func main() {
	parser := arg.MustParse(&args)

	if args.Debug != nil && *args.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	var err error
	if args.Increase != nil {
		err = doIncrease(&args)
		if err != nil {
			logrus.Fatalf("unable to increase brightness: %v", err)
		}
		return
	}

	if args.Decrease != nil {
		err = doDecrease(&args)
		if err != nil {
			logrus.Fatalf("unable to decrease brightness: %v", err)
		}
		return
	}

	if args.List != nil {
		err = doList(&args)
		if err != nil {
			logrus.Fatalf("unable to list devices: %v", err)
		}
		return
	}

	if args.Set != nil {
		err = doSet(&args)
		if err != nil {
			logrus.Fatalf("unable to set brightness: %v", err)
		}
		return
	}

	parser.WriteHelp(os.Stderr)
}
