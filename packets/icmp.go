package packets

type ICMP struct {
	Type byte
	Code byte
}

const ICMP_Protocol IPv4Protocol = 1

func NewICMP(icmp_type, code byte) ICMP {
	return ICMP{icmp_type, code}
}

func DecodeICMP(data []byte) *ICMP {
	if len(data) != 2 {
		return nil
	}

	return &ICMP{
		Type: data[0],
		Code: data[1],
	}
}

func (p ICMP) Encode() []byte {
	buffer := make([]byte, 2)
	buffer[0] = p.Type
	buffer[1] = p.Code
	return buffer
}

func (p ICMP) Protocol() IPv4Protocol {
	return ICMP_Protocol
}

func (p ICMP) Len() int64 {
	return 2
}
