package ByteView

/**
ByteView 只是对value做进一步封装
*/

type ByteView struct {
	B []byte
}

func (b ByteView) Len() int {
	return len(b.B)
}

func CloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

func (b ByteView) ByteSlice() []byte {
	return CloneBytes(b.B)
}

func (b ByteView) String() string {
	return string(b.B)
}
