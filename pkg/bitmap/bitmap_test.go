package bitmap

import "testing"

func TestBitmap_Set(t *testing.T) {
	b := NewBitmap(1)

	b.Set("hello")
	b.Set("world")
	b.Set("happy")
	b.Set("unhappy")
	b.Set("nihao")
	b.Set("zhongyu")
	b.Set("yuzhong")
	b.Set("haode")

	for _, bit := range b.bits {
		t.Logf("%b , %v ", bit, bit)
	}
}
