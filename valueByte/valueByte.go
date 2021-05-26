package valueByte

type ValueByte struct {
	B []byte
}

// 该缓存的值都需实现Len方法，实现了Value接口
func (v ValueByte) Len() int{
	return len(v.B)
}

// 因为b只读，需要值只能进行拷贝查看，防止被外部修改
func (v ValueByte)Clone() []byte{
	return CloneBytes(v.B)
}

// 为了外部转化需要给一个ToString 方法
func (v ValueByte)ToString() string{
	return string(v.B)
}

func CloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}


