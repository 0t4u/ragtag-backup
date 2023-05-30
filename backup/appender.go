package backup

import (
	"errors"
	"io"
	"os"
	"regexp"
)

// https://stackoverflow.com/a/47747702
var (
	tailCheckLen    int64 = 16
	arrayEndsObject       = regexp.MustCompile("(\\[\\s*)?](\\s*}\\s*)$")
	justArray             = regexp.MustCompile("(\\[\\s*)?](\\s*)$")
)

type Appender struct {
	f          *os.File
	needsComma bool
	tail       []byte
}

func (a Appender) Write(b []byte) (n int, err error) {
	bytes := 0

	// add a comma
	if a.needsComma {
		n, err = a.f.Write([]byte(","))
		bytes += n
		if err != nil {
			return bytes, err
		}
	}

	if b[len(b)-1] == 10 {
		b = b[:len(b)-1]
	}

	// write the json
	n, err = a.f.Write(b)
	bytes += n
	if err != nil {
		return bytes, err
	}

	return bytes, err
}

func (a Appender) Close() error {
	// add the tail
	if _, err := a.f.Write(a.tail); err != nil {
		defer a.f.Close()
		return err
	}
	return a.f.Close()
}

func ArrayAppender(file string) (io.WriteCloser, error) {
	f, err := os.OpenFile(file, os.O_RDWR, 0664)
	if err != nil {
		return nil, err
	}

	pos, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, err
	}

	if pos < tailCheckLen {
		pos = 0
	} else {
		pos -= tailCheckLen
	}
	_, err = f.Seek(pos, io.SeekStart)
	if err != nil {
		return nil, err
	}

	tail, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	needsComma := true

	if len(tail) == 0 {
		needsComma = false
		_, err = f.Write([]byte("["))
		if err != nil {
			return nil, err
		}
	} else {
		var g [][]byte
		if g = arrayEndsObject.FindSubmatch(tail); g != nil {
		} else if g = justArray.FindSubmatch(tail); g != nil {
		} else {
			return nil, errors.New("does not end with array")
		}

		_, err = f.Seek(-int64(len(g[2])+1), io.SeekEnd) // 1 for ]
		if err != nil {
			return nil, err
		}
		// tail = g[2]
	}

	tail = []byte("]")

	return Appender{f: f, needsComma: needsComma, tail: tail}, nil
}
