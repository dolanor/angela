package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/dolanor/angela/merkle"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
)

func main() {
	fs := ff.NewFlagSet("filec")
	filecCmd := &ff.Command{
		Name:      "filec",
		Usage:     "filec SUBCOMMAND ...",
		ShortHelp: "filec let you save files to server in a secure way",
	}

	putCmd := &ff.Command{
		Name:      "put",
		Usage:     "filec put FILE [FILE...]",
		ShortHelp: "send files to the server",
	}

	prepareCmd := &ff.Command{
		Name:      "prepare",
		Usage:     "filec prepare FILE [FILE...]",
		ShortHelp: "generate the merkle root from the files",
		Exec: func(ctx context.Context, args []string) error {
			if len(args) == 0 {
				return errors.New("no files listed")
			}

			var content []merkle.Content
			for _, path := range args {
				fi, err := os.Stat(path)
				if err != nil {
					return err
				}

				if fi.IsDir() {
					// TODO include files in the directory as well, possibly with -r?
					continue
				}

				b, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				content = append(content, b)

			}

			tree := merkle.FromContentSlice(content)
			err := os.WriteFile("root.hash", tree.Root.Hash, 0o600)
			if err != nil {
				return err
			}

			return nil
		},
	}

	filecCmd.Subcommands = append(filecCmd.Subcommands, putCmd, prepareCmd)

	err := filecCmd.ParseAndRun(context.Background(), os.Args[1:])
	if errors.Is(err, ff.ErrHelp) {
		fmt.Fprintf(os.Stderr, "%s\n%s\n", ffhelp.Command(filecCmd), ffhelp.Flags(fs))
		os.Exit(0)
	}
	if err != nil {
		panic(err)
	}
}
