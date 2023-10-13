package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/dolanor/angela/merkle"
	"github.com/dolanor/angela/web"
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
		Usage:     "filec put ENDPOINT BUCKET_NAME FILE [FILE...]",
		ShortHelp: "send files to the server",
		LongHelp:  "ENDPOINT is in the format http://host:port/files",
		Exec:      put,
	}

	prepareCmd := &ff.Command{
		Name:      "prepare",
		Usage:     "filec prepare FILE [FILE...]",
		ShortHelp: "generate the merkle root from the files",
		Exec:      prepare,
	}

	verifyCmd := &ff.Command{
		Name:      "verify",
		Usage:     "filec verify FILE MERKLE_PROOF_FILE MERKLE_ROOT_FILE",
		ShortHelp: "verify the file content with the merkle root file",
		Exec:      verify,
	}

	filecCmd.Subcommands = append(filecCmd.Subcommands, putCmd, prepareCmd, verifyCmd)

	err := filecCmd.ParseAndRun(context.Background(), os.Args[1:])
	if errors.Is(err, ff.ErrHelp) {
		fmt.Fprintf(os.Stderr, "%s\n%s\n", ffhelp.Command(filecCmd), ffhelp.Flags(fs))
		os.Exit(0)
	}
	if err != nil {
		panic(err)
	}
}

func prepare(ctx context.Context, args []string) error {
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
}

func put(ctx context.Context, args []string) error {
	if len(args) < 3 {
		return errors.Join(ff.ErrHelp, errors.New("not enough arguments"))
	}
	endpoint := args[0]
	bucketName := args[1]
	filePaths := args[2:]

	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return err
	}

	var files []merkle.Content
	for _, path := range filePaths {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		files = append(files, data)
	}

	req := web.CreateFilesRequest{
		BucketName: bucketName,
		Files:      files,
	}

	var b bytes.Buffer
	err = json.NewEncoder(&b).Encode(req)
	if err != nil {
		return err
	}

	resp, err := http.Post(endpointURL.String(), "application/json", &b)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("sending files failed")
	}

	return nil
}

func verify(ctx context.Context, args []string) error {
	return errors.New("not implemented")
}
