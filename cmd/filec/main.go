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
	"path"

	"github.com/dolanor/angela/api"
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

	getCmd := &ff.Command{
		Name:      "get",
		Usage:     "filec get ENDPOINT BUCKET_NAME FILE_NUMBER",
		ShortHelp: "get the i-th file from the server with the proof",
		LongHelp:  "ENDPOINT is in the format http://host:port/files",
		Exec:      get,
	}

	verifyCmd := &ff.Command{
		Name:      "verify",
		Usage:     "filec verify FILE MERKLE_PROOF_FILE MERKLE_ROOT_FILE",
		ShortHelp: "verify the file content with the merkle root file",
		Exec:      verify,
	}

	filecCmd.Subcommands = append(filecCmd.Subcommands, putCmd, prepareCmd, getCmd, verifyCmd)

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

	req := api.CreateFilesRequest{
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

func get(ctx context.Context, args []string) error {
	if len(args) < 3 {
		return errors.Join(ff.ErrHelp, errors.New("not enough arguments"))
	}
	endpoint := args[0]
	bucketName := args[1]
	fileNumberStr := args[2]

	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return err
	}

	endpointURL.Path = path.Join(endpointURL.Path, bucketName, fileNumberStr)

	resp, err := http.Get(endpointURL.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("getting file failed")
	}

	var getFileResp api.GetFileResponse
	err = json.NewDecoder(resp.Body).Decode(&getFileResp)
	if err != nil {
		return err
	}

	err = os.WriteFile(fmt.Sprintf("file-%s.data", fileNumberStr), getFileResp.Content, 0o600)
	if err != nil {
		return err
	}

	var proofJSON bytes.Buffer
	err = json.NewEncoder(&proofJSON).Encode(getFileResp.Proof)
	if err != nil {
		return err
	}

	err = os.WriteFile(fmt.Sprintf("file-%s.proof", fileNumberStr), proofJSON.Bytes(), 0o600)
	if err != nil {
		return err
	}

	return nil
}

func verify(ctx context.Context, args []string) error {
	if len(args) < 3 {
		return errors.Join(ff.ErrHelp, errors.New("not enough arguments"))
	}

	merkleRootFilePath := args[0]
	merkleProofFilePath := args[1]
	filePath := args[2]

	merkleRootHash, err := os.ReadFile(merkleRootFilePath)
	if err != nil {
		return err
	}

	f, err := os.Open(merkleProofFilePath)
	if err != nil {
		return err
	}

	var merkleProof []merkle.ProofStep
	err = json.NewDecoder(f).Decode(&merkleProof)
	if err != nil {
		return err
	}

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = merkle.Verify(merkleRootHash, merkleProof, fileData)
	if err != nil {
		return err
	}

	fmt.Printf("file %q is correct\n", filePath)

	return nil
}
