package seekctr_test

import (
	"bytes"
	"crypto/rand"
	"io"
	"io/ioutil"
	"testing"

	"github.com/uhthomas/seekctr"
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
	w, err := seekctr.NewWriter(ww, key[:], iv[:])
	if err != nil {
		t.Fatal(err)
	}
	r, err := seekctr.NewReader(rr, key[:], iv[:])
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
	w, err := seekctr.NewWriter(&buf, key[:], iv[:])
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
	r, err := seekctr.NewReader(bytes.NewReader(b[:]), key[:], iv[:])
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

type infiniteReader struct{}

func (infiniteReader) Read(b []byte) (int, error) { return len(b), nil }

func BenchmarkReader(b *testing.B) {
	var key, iv [16]byte
	if _, err := io.ReadFull(rand.Reader, key[:]); err != nil {
		b.Fatal(err)
	}
	if _, err := io.ReadFull(rand.Reader, iv[:]); err != nil {
		b.Fatal(err)
	}

	r, err := seekctr.NewReader(new(infiniteReader), key[:], iv[:])
	if err != nil {
		b.Fatal(err)
	}

	var out [32 << 10]byte
	b.SetBytes(int64(len(out)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := r.Read(out[:])
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWriter(b *testing.B) {
	var key, iv [16]byte
	if _, err := io.ReadFull(rand.Reader, key[:]); err != nil {
		b.Fatal(err)
	}
	if _, err := io.ReadFull(rand.Reader, iv[:]); err != nil {
		b.Fatal(err)
	}

	w, err := seekctr.NewWriter(ioutil.Discard, key[:], iv[:])
	if err != nil {
		b.Fatal(err)
	}

	var out [32 << 10]byte
	b.SetBytes(int64(len(out)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := w.Write(out[:])
		if err != nil {
			b.Fatal(err)
		}
	}
}
