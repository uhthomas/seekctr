package seekctr

import (
	"crypto/aes"
	"io"
)

// Reader wraps and io.Reader and will decrypt anything read from it using the
// ctr.
type Reader struct {
	*ctr
	r io.Reader
}

// NewReader returns an io.ReadSeekCloser. The key must be 16, 24 or 32 bytes
// in size and the iv must be 16 bytes.
func NewReader(r io.Reader, key, iv []byte) (*Reader, error) {
	b, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return &Reader{newCTR(b, iv), r}, nil
}

// Read will read len(p) bytes from r and decrypt them using ctr.
func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	r.XORKeyStream(p[:n], p[:n])
	return
}

// Seek will, if r is an io.Seeker, Seek to the offset given the whence and
// will then seek the ctr cipher.
func (r *Reader) Seek(offset int64, whence int) (ret int64, err error) {
	if s, ok := r.r.(io.Seeker); ok {
		ret, err = s.Seek(offset, whence)
		r.seek(ret)
	}
	return
}

// Close closes the underlying reader if it is an io.Closer.
func (r *Reader) Close() (err error) {
	if c, ok := r.r.(io.Closer); ok {
		err = c.Close()
	}
	return
}

// Writer wraps an io.Writer and will encrypt anything written to it using the
// ctr.
type Writer struct {
	*ctr
	w io.Writer
}

// NewWriter returns an io.WriteSeekCloser. The key must be 16, 24 or 32 bytes
// in size and the iv must be 16 bytes.
func NewWriter(w io.Writer, key, iv []byte) (*Writer, error) {
	b, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return &Writer{newCTR(b, iv), w}, nil
}

// Write encrypts b using ctr and then writes to w.
func (w *Writer) Write(b []byte) (n int, err error) {
	c := make([]byte, len(b))
	w.XORKeyStream(c, b)
	return w.w.Write(c)
}

// Seek will, if w is an io.Seeker, Seek to the offset given the whence and
// will then seek the ctr cipher.
func (w *Writer) Seek(offset int64, whence int) (ret int64, err error) {
	if s, ok := w.w.(io.Seeker); ok {
		ret, err = s.Seek(offset, whence)
		w.seek(ret)
	}
	return
}

// Close closes the underlying writer if it is an io.Closer.
func (w *Writer) Close() (err error) {
	if c, ok := w.w.(io.Closer); ok {
		err = c.Close()
	}
	return
}
