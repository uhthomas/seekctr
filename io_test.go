package seekctr_test

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"errors"
	"fmt"
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
	if _, err := r.Seek(0, io.SeekEnd); err != nil {
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

// TestCipherSeekStart tests that seeking to the start of the seeker works.
func TestCipherSeekStart(t *testing.T) {
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

	data := []byte(`
	Lorem ipsum dolor sit amet, consectetur adipiscing elit. Ut nec libero augue. Mauris tempor lectus in augue semper malesuada. Nullam eget semper neque. Sed facilisis urna non mauris pellentesque, at vehicula augue tristique. Donec feugiat quam at sem faucibus porttitor quis ut dolor. Ut mattis ex quis nisi dictum sodales. Suspendisse potenti. Aenean sed bibendum sem. Sed eget mi nec justo tristique bibendum ut et tellus. Pellentesque vel mi mattis turpis posuere dapibus. Suspendisse fermentum metus ut rhoncus bibendum.
	Etiam rutrum faucibus lobortis. Integer vel varius est. Sed pellentesque libero id ultricies pellentesque. Phasellus volutpat felis feugiat est finibus, eget ultricies velit porttitor. Proin pharetra eros in arcu pulvinar, at porttitor nunc finibus. Curabitur convallis lacus a tortor porttitor, vitae eleifend lorem tristique. Cras ultrices, urna sed egestas condimentum, lorem metus lacinia magna, sit amet hendrerit urna lacus at massa. Maecenas rhoncus malesuada lacus eget convallis. Praesent bibendum ipsum in erat faucibus luctus. Nullam condimentum dolor vitae turpis venenatis, consequat tincidunt turpis hendrerit. In vitae ex in lacus pharetra interdum quis ac nunc. Aliquam et urna malesuada, fermentum tellus eget, hendrerit odio. Fusce rutrum finibus neque, ut sodales lorem ullamcorper non. Duis egestas placerat quam. Nunc pharetra feugiat imperdiet.
	Nulla at magna quam. Suspendisse potenti. Nullam eu urna pretium, tincidunt tellus dignissim, varius neque. Morbi aliquam felis non lectus tempus, nec vehicula quam dictum. In lacinia sit amet sem a porta. Quisque a ante tincidunt, vulputate velit at, suscipit ligula. Mauris pretium aliquam ultrices. In hac habitasse platea dictumst. Sed ut ante augue. Etiam id tempus purus.
	Nulla vitae nisi vehicula, sagittis sem et, mattis nibh. Aenean dignissim cursus enim. Nullam non dapibus metus. Integer vulputate tellus in metus venenatis, id congue orci rutrum. Maecenas ornare sit amet erat iaculis semper. Suspendisse id odio efficitur, porta odio sed, consectetur neque. Proin ligula erat, finibus quis neque quis, venenatis sagittis diam. Mauris at elit sit amet odio tincidunt euismod. Nam nec tellus arcu. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Ut at enim id eros dignissim hendrerit. Phasellus sit amet ante sed nunc vehicula bibendum. Aenean tempor faucibus faucibus.
	Sed vel tortor gravida, dapibus enim non, fringilla libero. Aenean vitae odio eu diam sodales mollis vitae ac ipsum. Morbi eget arcu sem. Curabitur cursus dui orci, lobortis eleifend nibh egestas vel. Nunc euismod tellus eget consectetur viverra. Quisque at ante eu urna eleifend sagittis. Sed id metus dapibus, placerat ante in, facilisis velit. Aliquam erat volutpat. Nunc et mi sit amet enim tincidunt gravida sed vitae metus. Praesent gravida mauris ut sem suscipit, vitae ornare tellus gravida. Nullam fermentum risus eget nisi malesuada finibus. Praesent nisl eros orci aliquam.	
	`)

	fmt.Println("data", len(data))

	bytesReader := bytes.NewReader(data)
	_, err = io.Copy(w, bytesReader)
	if err != nil {
		t.Fatal(err)
	}
	w.Close()

	r, err := seekctr.NewReader(bytes.NewReader(buf.Bytes()), key[:], iv[:])
	if err != nil {
		t.Fatal(err)
	}

	t.Run("readall", func(t *testing.T) {
		readBack, err := io.ReadAll(r)
		if err != nil {
			t.Fatal(err)
		}
		if string(readBack) != string(data) {
			t.Errorf("got %s, expected %s", string(readBack), string(data))
		}
	})

	t.Run("seek start & readall", func(t *testing.T) {
		seekStart, err := r.Seek(0, io.SeekStart)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("seek start", seekStart)

		readBack, err := io.ReadAll(r)
		if err != nil {
			t.Fatal(err)
		}
		if string(readBack) != string(data) {
			t.Errorf("got %s, expected %s", string(readBack), string(data))
		}
	})

	t.Run("seek end -10 & readall", func(t *testing.T) {
		seekEnd, err := r.Seek(-10, io.SeekEnd)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("seek end", seekEnd)

		readBack, err := io.ReadAll(r)
		if err != nil {
			t.Fatal(err)
		}
		if string(readBack) != "liquam.\t\n\t" {
			t.Error("readback length is not 10")
		}
		fmt.Print(readBack)
	})

	t.Run("seek start & readall", func(t *testing.T) {
		_, err = r.Seek(0, io.SeekStart)
		if err != nil {
			t.Fatal(err)
		}

		readBack, err := io.ReadAll(r)
		if err != nil {
			t.Fatal(err)
		}
		if string(readBack) != string(data) {
			t.Errorf("got %s, expected %s", string(readBack), string(data))
		}
	})

	// rdr := seekctr.NewReader()
}

type FakeReadWriteSeeker int64

func (s *FakeReadWriteSeeker) Write(b []byte) (int, error) { return len(b), nil }

func (s *FakeReadWriteSeeker) Read(b []byte) (int, error) { return len(b), nil }

func (s *FakeReadWriteSeeker) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		*s = FakeReadWriteSeeker(offset)
	case io.SeekCurrent:
		*s += FakeReadWriteSeeker(offset)
	default:
		return 0, errors.New("invalid whence")
	}
	return int64(*s), nil
}

func BenchmarkReader(b *testing.B) {
	var key, iv [16]byte
	if _, err := io.ReadFull(rand.Reader, key[:]); err != nil {
		b.Fatal(err)
	}
	if _, err := io.ReadFull(rand.Reader, iv[:]); err != nil {
		b.Fatal(err)
	}

	r, err := seekctr.NewReader(new(FakeReadWriteSeeker), key[:], iv[:])
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

func BenchmarkSeek(b *testing.B) {
	w, err := seekctr.NewWriter(new(FakeReadWriteSeeker), make([]byte, aes.BlockSize), make([]byte, aes.BlockSize))
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < int64(b.N); i++ {
		_, err := w.Seek(i, io.SeekCurrent)
		if err != nil {
			b.Fatal(err)
		}
	}
}
