package blob

import (
	"bytes"
	"strconv"
)

type Blob struct {
	data []byte
}

func New(data []byte) Blob {
	return Blob{data}
}

func (b Blob) Type() string {
	return "blob"
}

func (b Blob) Serialize() []byte {
	var buf bytes.Buffer
	buf.WriteString(b.Type())
	buf.WriteRune(' ')
	buf.WriteString(strconv.Itoa(len(b.data)))
	buf.WriteRune('\x00')
	buf.Write(b.data)
	return buf.Bytes()
}
