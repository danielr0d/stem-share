package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
)

const usage = `usage:
  stemshare send <file.wav>...
  stemshare receive [-o dir] [-force] <code>
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var err error
	switch os.Args[1] {
	case "send":
		err = send(ctx, os.Args[2:])
	case "receive", "recv":
		fs := flag.NewFlagSet("receive", flag.ExitOnError)
		out := fs.String("o", ".", "output directory")
		force := fs.Bool("force", false, "overwrite existing files")
		fs.Parse(os.Args[2:])
		if fs.NArg() != 1 {
			fmt.Fprint(os.Stderr, usage)
			os.Exit(2)
		}
		err = receive(ctx, fs.Arg(0), *out, *force)
	default:
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", friendly(err))
		os.Exit(1)
	}
}

func friendly(err error) error {
	if strings.Contains(err.Error(), "decrypt message failed") {
		return fmt.Errorf("wrong code (or someone else used it)")
	}
	return err
}
