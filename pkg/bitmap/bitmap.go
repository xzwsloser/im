package bitmap

type Bitmap struct {
	bits []byte
	size int
}

func NewBitmap(size int) *Bitmap {
	if size == 0 {
		size = 250
	}

	return &Bitmap{
		bits: make([]byte, size),
		size: size * 8,
	}
}

func (b *Bitmap) Set(id string) {
	// 计算在那一个 bit
	idx := hash(id) % b.size
	// 计算在哪一个 byte
	byteIdx := idx / 8
	// 记录在与一个byte中的位置
	bitIdx := idx % 8

	// 此时只有 bitIdx 的位置为 1 , 其他都是 0
	b.bits[byteIdx] |= 1 << bitIdx
}

func (b *Bitmap) IsSet(id string) bool {
	idx := hash(id) % b.size
	byteIdx := idx / 8
	bitIdx := idx % 8

	return (b.bits[byteIdx])&(1<<bitIdx) != 0
}

func (b *Bitmap) Export() []byte {
	return b.bits
}

// 0 0 0 0 0 0 0 1
func Load(bits []byte) *Bitmap {
	if len(bits) == 0 {
		return NewBitmap(0)
	}

	return &Bitmap{
		bits: bits,
		size: len(bits) * 8,
	}
}

// @Description: String -> int(还是可能会有重复情况)
func hash(id string) int {
	seed := 131313
	hash := 0

	for _, c := range id {
		hash = hash*seed + int(c)
	}
	return hash & 0x7FFFFFFF
}
