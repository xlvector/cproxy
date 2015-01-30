package cproxy

const (
	IMMUTABLE_UNKNOWN = 0
	IMMUTABLE_YES     = 1
	IMMUTABLE_NO      = 2
)

type Body struct {
	data      []byte
	sn        int64
	immutable int
	hit       int
}

func time33(b []byte) int64 {
	ret := int64(0)
	for _, e := range b {
		ret *= 33
		ret += int64(e)
	}
	return ret
}
