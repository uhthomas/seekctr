package seekctr

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"
)

func TestCipher(t *testing.T) {
	var raw [4 << 10]byte
	if _, err := io.ReadFull(rand.Reader, raw[:]); err != nil {
		t.Fatal(err)
	}
	var key, iv [16]byte
	if _, err := io.ReadFull(rand.Reader, key[:]); err != nil {
		t.Fatal(err)
	}
	if _, err := io.ReadFull(rand.Reader, iv[:]); err != nil {
		t.Fatal(err)
	}
	rr, ww := io.Pipe()
	w, err := NewWriter(ww, key[:], iv[:])
	if err != nil {
		t.Fatal(err)
	}
	r, err := NewReader(rr, key[:], iv[:])
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		defer w.Close()
		if _, err := io.Copy(w, bytes.NewReader(raw[:])); err != nil {
			t.Fatal(err)
		}
	}()
	var out bytes.Buffer
	if _, err := io.Copy(&out, r); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(raw[:], out.Bytes()) {
		t.Fatal("Result is not equal")
	}
}

func TestCipherSeek(t *testing.T) {
	var raw [8 << 10]byte
	if _, err := io.ReadFull(rand.Reader, raw[:]); err != nil {
		t.Fatal(err)
	}
	var key, iv [16]byte
	if _, err := io.ReadFull(rand.Reader, key[:]); err != nil {
		t.Fatal(err)
	}
	if _, err := io.ReadFull(rand.Reader, iv[:]); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	w, err := NewWriter(&buf, key[:], iv[:])
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(w, bytes.NewReader(raw[:])); err != nil {
		t.Fatal(err)
	}
	var b [8 << 10]byte
	if _, err := io.ReadFull(&buf, b[:]); err != nil {
		t.Fatal(err)
	}
	r, err := NewReader(bytes.NewReader(b[:]), key[:], iv[:])
	if err != nil {
		t.Fatal(err)
	}
	if _, err := r.Seek(8<<10, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	if _, err := r.Seek(4<<10, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if _, err := io.Copy(&out, r); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(raw[4<<10:], out.Bytes()) {
		t.Fatal("Result is not equal")
	}
}
