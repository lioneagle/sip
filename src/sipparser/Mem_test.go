package sipparser

import (
	"bytes"
	"fmt"
	"testing"
)

//*
func TestMemAlloc(t *testing.T) {
	allocator := NewMemAllocator(1000)

	if allocator.Capacity() != 1000 {
		t.Errorf("TestMemAlloc failed, wong cap = %d, wanted = %d\n", allocator.Capacity(), 1000)
		return
	}

	testdata := []struct {
		allocSize int32
		values    []byte
	}{
		{100, []byte{1, 2}},
		{200, []byte{11, 12}},
	}

	used := int32(0)

	for i, v := range testdata {
		x := allocator.AllocByteSlice(v.allocSize)
		used += v.allocSize + SLICE_HEADER_LEN
		used = RoundToAlign(used, SIP_MEM_ALIGN)

		if allocator.Used() != used {
			t.Errorf("TestMemAlloc[%d] failed, wong used = %d, wanted = %d\n", i, allocator.Used(), used)
			break
		}

		if allocator.AllocNum() != int32(i+1) {
			t.Errorf("TestMemAlloc[%d] failed, wong AllocNum = %d, wanted = %d\n", i, allocator.AllocNum(), i+1)
			break
		}

		if allocator.FreeAllNum() != 0 {
			t.Errorf("TestMemAlloc[%d] failed, wong FreeAllNum = %d, wanted = %d\n", i, allocator.FreeAllNum(), 0)
			break
		}

		x1 := append(x, v.values...)

		if !bytes.Equal(x1, v.values) {
			t.Errorf("TestMemAlloc[%d] failed, wong values = %v, wanted = %v\n", i, x1, v.values)
			break
		}

		/*
			if &x1 != &x {
				t.Errorf("TestMemAlloc[%d] failed, wong ptr = %p, wanted = %p\n", i, &x1, &x)
				fmt.Println("allocator.mem =", allocator.mem)
				fmt.Println("x1 =", x1)
				fmt.Println("x =", x)
				continue
			} //*/
	}

	allocator.FreeAll()

	if allocator.Used() != 0 {
		t.Errorf("TestMemAlloc failed, wong used = %d, wanted = %d\n", allocator.Used(), 0)
		return
	}

}

func BenchmarkMemAlloc(b *testing.B) {
	b.StopTimer()
	allocator := NewMemAllocator(1024 * 128)
	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		allocator.FreeAll()
		_, _ = allocator.Alloc(1000)
	}
	//fmt.Printf("uri = %s\n", uri.String())
	//fmt.Printf("")
}

func BenchmarkMemAllocByteSlice(b *testing.B) {
	b.StopTimer()
	allocator := NewMemAllocator(1024 * 128)
	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		allocator.FreeAll()
		_ = allocator.AllocByteSlice(1000)
	}
	//fmt.Printf("uri = %s\n", uri.String())
	fmt.Printf("")
} //*/
