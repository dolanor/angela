# Angela

Send files to a server and verify their integrity thanks to Merkle Trees

## Test

```shell
go test ./...
```

## Use

```shell
# build the file server daemon: filed
docker compose build

# run it in background
docker compose up -d

# build the filed client: filec
go build ./cmd/filec

# create the file you want to send
echo 000 > file0
echo 001 > file1
echo 002 > file2

# prepare files for sending (generate the merkle root hash in the "root.hash" file)
./filec prepare file0 file1 file2

# send them to filed
./filec put http://localhost:7777/files myBucketName file0 file1 file2

# delete the files
rm file0 file1 file2

# get file1 (will save in "file-1.data" and the proof in "file-1.proof"
./filec get http://localhost:7777/files myBucketName 1

# verify the integrity of the file
./filec verify root.hash file-1.proof file-1.data
```
