# seekctr [![GoDoc](https://godoc.org/github.com/uhthomas/seekctr?status.svg)](https://godoc.org/github.com/uhthomas/seekctr)

[cipher.Stream](https://pkg.go.dev/crypto/cipher#Stream) does not implement [io.Seeker](https://pkg.go.dev/io#Seeker) despite XOR stream ciphers being seekable.

## Usage

```go
package main

import (
	"fmt"
	"io"
	"log"
	"os"

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

Recreating the original stream cipher with a new initialization vector (where `iv += offset / block size`) and discarding the remaining bytes (`offset % block size`) may be preferrable.

```go
var key, iv [16]byte

b, err := aes.NewCipher(key[:])
if err != nil { ... }

offset := uint64(4 << 10)

// offset in chunks
chunks := uint64(int(offset) / b.BlockSize())

// iv += offset
var c uint16
for i := len(iv[:]) - 1; i >= 0; i-- {
	c = uint16(iv[i]) + uint16(chunks & 0xFF) + c
	iv[i], c, chunks = byte(c), c >> 8, chunks >> 8
}

// Reinitialize cipher
s := cipher.NewCTR(b, iv[:])

// Discard n bytes
d := make([]byte, int(offset) % b.BlockSize())
s.XORKeyStream(d, d)
```
