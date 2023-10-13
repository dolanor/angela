# File Server

## Implementation log

I started to implement the merkle tree construction as it is the part I never did before and might take more time at first. I had a result quite fast, but no way to test it properly without the merkle proofs.

Since there were different ways to implement odd numbers of files, I wanted to check my tree had the correct structure and implemented a fmt.Stringer interface to be able to print it, for debugging purposes. As hashes were hard to read, I created a small library to print hashes in more visual way. Then I had a base for understanding my algorithm in a different way. It confirmed my understanding, so I moved on to the proofs.

---

For now, my Go tests were only a way for me to visually ensure it's okay. But now, I want to create to test that the proof is correct and the tree as well. Then, I'll be to move to the client/server part.

---

I implemented the proof generation and proof verification and tested it was OK.

---

I'm now reorganizing the code to support the client/server service

---

I implemented an HTTP server with 2 endpoints:
- 1 for upload of many files in a JSON format (not the best, but it was the fastest)
- 1 for getting the i-th file + the proof in a JSON format as well

I split the file storage into a specific FileServer, the handlers handle the HTTP, decoding, encoding part only

---

I used the recent improvement in Docker to use build mounts to allow build without copying and caching modules download + build cache. Thus reducing the build from 10s -> 3s

---

add the client implementation in CLI, using the web API

---

## Improvements

- create a generic version of the merkle tree that could handle any type of data, not just []byte
- prune the tree to avoid long simple branch with 1 hash at the end in case of odd amount of leaves
- use context.Context along the service to allow for request cancelling and useless processing
- save the files content on disk
- use multipart/form-data for file uploads
- create one endpoint to get file data and one for file merkle proof (instead of the current bundle of json + base64 encoded data)
- shard data among different servers (we can shard by bucket/upload)
  - use a distributed hash table to store where the buckets are stored (which server)
- use TLS for the HTTP server interacting with the client
- use mTLS to encrypt connection between servers and ensure trust
- harmonize the client UX
  - order of params
  - format of files
  - destination of files
- improve the API
