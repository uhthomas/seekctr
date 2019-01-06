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

## Note
Alternatively, although seeking would be a lot less efficient, to use the original implementation, the CTR can be re-initialized with the modified iv and then n bytes discarded. For example:
```go
offset := uint64(4 << 10)
var key, iv, boffset [16]byte
b, err := aes.NewCipher(key[:])
if err != nil { ... }
// Convert offset to [16]byte
binary.BigEndian.PutUint64(boffset[8:], offset)
// Add offset to iv
var c uint16
for i := len(b) - 1; i >= 0; i-- {
	c = uint16(iv[i]) + uint16(boffset[i]) + c
	iv[i] = byte(c)
	c >>= 8
}

// Reinitialize cipher
s := cipher.NewCTR(b, iv)
```