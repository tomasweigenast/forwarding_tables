package packets

type IPv4Payload interface {
	Encode() []byte
	Len() int64
	Protocol() IPv4Protocol
}

type IPv4Protocol byte
