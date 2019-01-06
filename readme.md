# seekctr [![GoDoc](https://godoc.org/github.com/uhthomas/seekctr?status.svg)](https://godoc.org/github.com/uhthomas/seekctr)

## Why?
The [native Go implementation](https://golang.org/pkg/crypto/cipher) does not provide a `Seek` method and since these stream ciphers _are_ seekable since they are xor stream ciphers, I implemented a seekable CTR cipher.

## Usage
```go
package main

import (
	"fmt"
	"log"

	"github.com/uhthomas/seekctr"
)

func main() {
	// your key and initialization vector
	var key, iv []byte
	// open the encrypted file or stream
	f, err := os.Open("encrypted file")
	if err != nil {
		log.Fatal(err)
	}
	r, err := seekctr.NewReader(f, key, iv)
	if err != nil {
		log.Fatal(err)
	}
	// Seek past the first 1Kb
	if _, err := r.Seek(1 << 10); err != nil {
		log.Fatal(err)
	}
	// copy the remaining contents to stdout
	if _, err := io.Copy(os.Stdout, r); err != nil {
		log.Fatal(err)
	}
}
```