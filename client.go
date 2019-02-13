package kasalink

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"time"
)

var errTooLarge = errors.New("bytes.Buffer: too large")

// TalkToPlug attempts to connect to the Kasa Device and retrieve it's system info
func TalkToPlug(KasaCommand string) (response string, err error) {
	var (
		tplinkClient  net.Conn
		hs300Location = "10.0.0.25:9999"
		getSystemInfo []byte
	)
	if tplinkClient, err = net.DialTimeout("tcp", hs300Location, time.Duration(10)*time.Second); err != nil {
		return
	}
	defer closer(tplinkClient)
	getSystemInfo = encrypt(KasaCommand)
	if _, err = tplinkClient.Write(getSystemInfo); err != nil {
		return
	}
	var bb = new(myBuff)
	var bytesRead int64
	bytesRead, err = bb.readFrom(tplinkClient)
	//_, err = bb.readFrom(tplinkClient)
	if err != nil {
		return
	}
	log.Printf("Bytes Read: %d\n", bytesRead)
	if bb.Len() >= 4 {
		return decrypt(bb.buf[4:]), nil
	}
	return
}

//So most of everything below this is pulled whole hog out of bytes.Buffer, but I needed a custom read to prevent us
//from hanging until the connection timed out.
func (b *myBuff) readFrom(r io.Reader) (n int64, err error) {
	const MinRead = 512
	var (
		errNegativeRead      = errors.New("reader returned negative count from Read")
		contentSizeExtracted bool
		SizeToRead           int32
		i, bytesRead         int
	)

	for {
		i = b.grow(MinRead)
		b.buf = b.buf[:i]
		bytesRead, err = r.Read(b.buf[i:cap(b.buf)])
		if bytesRead < 0 {
			panic(errNegativeRead)
		} else if !contentSizeExtracted {
			var bytesToRead = bytes.NewBuffer(b.buf[0:4])
			if err = binary.Read(bytesToRead, binary.BigEndian, &SizeToRead); err != nil {
				panic(err)
			}
			SizeToRead = SizeToRead + 4
			//log.Printf("Think I need to read %d bytes\n", SizeToRead)
			contentSizeExtracted = true
			//} else {
			//	log.Printf("contentSizeExtracted: %t, SizeToRead: %d\n", contentSizeExtracted, SizeToRead)
		}

		b.buf = b.buf[:i+bytesRead]
		n += int64(bytesRead)
		if contentSizeExtracted && n >= int64(SizeToRead) {
			return
		}
		if err == io.EOF {
			return n, nil // err is EOF, so return nil explicitly
		}
		if err != nil {
			return n, err
		}
	}
}

type myBuff struct {
	buf       []byte   // contents are the bytes buf[off : len(buf)]
	off       int      // read at &buf[off], write at &buf[len(buf)]
	bootstrap [64]byte // memory to hold first slice; helps small buffers avoid allocation.
}

func (b *myBuff) grow(n int) int {
	const maxInt = int(^uint(0) >> 1)

	m := b.Len()
	// If buffer is empty, reset to recover space.
	if m == 0 && b.off != 0 {
		b.Reset()
	}
	// Try to grow by means of a reslice.
	if i, ok := b.tryGrowByReslice(n); ok {
		return i
	}
	// Check if we can make use of bootstrap array.
	if b.buf == nil && n <= len(b.bootstrap) {
		b.buf = b.bootstrap[:n]
		return 0
	}
	c := cap(b.buf)
	if n <= c/2-m {
		// We can slide things down instead of allocating a new
		// slice. We only need m+n <= c to slide, but
		// we instead let capacity get twice as large so we
		// don't spend all our time copying.
		copy(b.buf, b.buf[b.off:])
	} else if c > maxInt-c-n {
		panic(errTooLarge)
	} else {
		// Not enough space anywhere, we need to allocate.
		buf := makeSlice(2*c + n)
		copy(buf, b.buf[b.off:])
		b.buf = buf
	}
	// Restore b.off and len(b.buf).
	b.off = 0
	b.buf = b.buf[:m+n]
	return m
}

// Len returns the number of bytes of the unread portion of the buffer;
// b.Len() == len(b.Bytes()).
func (b *myBuff) Len() int { return len(b.buf) - b.off }

// Reset resets the buffer to be empty,
// but it retains the underlying storage for use by future writes.
// Reset is the same as Truncate(0).
func (b *myBuff) Reset() {
	b.buf = b.buf[:0]
	b.off = 0
}

// tryGrowByReslice is a inlineable version of grow for the fast-case where the
// internal buffer only needs to be resliced.
// It returns the index where bytes should be written and whether it succeeded.
func (b *myBuff) tryGrowByReslice(n int) (int, bool) {
	if l := len(b.buf); n <= cap(b.buf)-l {
		b.buf = b.buf[:l+n]
		return l, true
	}
	return 0, false
}

// makeSlice allocates a slice of size n. If the allocation fails, it panics
// with errTooLarge.
func makeSlice(n int) []byte {
	// If the make fails, give a known error.
	defer func() {
		if recover() != nil {
			panic(errTooLarge)
		}
	}()
	return make([]byte, n)
}
