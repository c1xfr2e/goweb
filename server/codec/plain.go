package codec

type PlainCodec struct {
}

func NewPlainCodec() *PlainCodec {
	return &PlainCodec{}
}

func (c *PlainCodec) Encode(data []byte) ([]byte, error) {
	return data, nil
}

func (c *PlainCodec) Decode(data []byte) ([]byte, error) {
	return data, nil
}
