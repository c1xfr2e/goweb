package codec

type Decoder interface {
	Decode(data []byte) ([]byte, error)
}

type Encoder interface {
	Encode(data []byte) ([]byte, error)
}
