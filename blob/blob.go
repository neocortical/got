package blob

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
	return b.data
}
