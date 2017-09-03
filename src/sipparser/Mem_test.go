package sipparser

import (
	//"fmt"
	"testing"
)

//*
func TestMemAllocatorAlloc(t *testing.T) {

	testdata := []struct {
		memSize   int32
		allocSize int32
		ok        bool
	}{

		{200, 1, true},
		{200, 199, true},

		{100, 0, false},
		{100, -1, false},
		{100, -100, false},
	}

	prefix := FuncName()
	used := int32(0)

	for i, v := range testdata {

		allocator := NewMemAllocator(v.memSize)
		if allocator == nil {
			t.Errorf("%s[%d] failed: NewMemAllocator failed\n", prefix, i)
			break
		}

		if allocator.Capacity() != v.memSize {
			t.Errorf("%s[%d] failed: wrong capacity =%d, wanted = %d\n", prefix, i, allocator.Capacity(), v.memSize)
			break
		}

		mem, addr := allocator.Alloc(v.allocSize)

		if v.ok {
			used += v.allocSize
			used = RoundToAlign(used, SIP_MEM_ALIGN)
		}

		if mem == nil && v.ok {
			t.Errorf("%s[%d] failed: mem should not be nil\n", prefix, i)
			continue
		}

		if mem != nil && !v.ok {
			t.Errorf("%s[%d] failed: mem should be nil\n", prefix, i)
			continue
		}

		if addr == ABNF_PTR_NIL && v.ok {
			t.Errorf("%s[%d] failed: addr should not be nil\n", prefix, i)
			continue
		}

		if addr != ABNF_PTR_NIL && !v.ok {
			t.Errorf("%s[%d] failed: addr should be nil\n", prefix, i)
			continue
		}
	}
}

func TestMemAllocatorAllocEx(t *testing.T) {

	testdata := []struct {
		memSize      int32
		allocSize    int32
		memAllocSize int32
		ok           bool
	}{

		{200, 1, SIP_MEM_ALIGN, true},
		{200, SIP_MEM_ALIGN + 1, 2 * SIP_MEM_ALIGN, true},
	}

	prefix := FuncName()
	used := int32(0)

	for i, v := range testdata {

		allocator := NewMemAllocator(v.memSize)
		if allocator == nil {
			t.Errorf("%s[%d] failed: NewMemAllocator failed\n", prefix, i)
			break
		}

		if allocator.Capacity() != v.memSize {
			t.Errorf("%s[%d] failed: wrong capacity =%d, wanted = %d\n", prefix, i, allocator.Capacity(), v.memSize)
			break
		}

		addr, alloc := allocator.AllocEx(v.allocSize)

		if v.ok {
			used += v.allocSize
			used = RoundToAlign(used, SIP_MEM_ALIGN)
		}

		if addr == ABNF_PTR_NIL && v.ok {
			t.Errorf("%s[%d] failed: addr should not be nil\n", prefix, i)
			continue
		}

		if addr != ABNF_PTR_NIL && !v.ok {
			t.Errorf("%s[%d] failed: addr should be nil\n", prefix, i)
			continue
		}

		if alloc != v.memAllocSize {
			t.Errorf("%s[%d] failed: wong mem alloc size = %d, wanted = %d\n", prefix, i, alloc, v.memAllocSize)
			continue
		}
	}
}

func TestMemAllocatorUsed(t *testing.T) {
	testdata := []struct {
		allocSize int32
		ok        bool
	}{
		{101, true},
		{203, true},
		{-1, false},
		{0, false},
		{21, true},
		{1, true},
	}

	allocator := NewMemAllocator(1000)
	prefix := FuncName()
	used := int32(0)
	allocNum := int32(0)
	allocNumOk := int32(0)

	for i, v := range testdata {
		allocNum++
		if v.ok {
			allocNumOk++
			used += v.allocSize
			used = RoundToAlign(used, SIP_MEM_ALIGN)
		}

		mem, _ := allocator.Alloc(v.allocSize)
		if mem == nil && v.ok {
			t.Errorf("%s[%d] failed: mem should not be nil\n", prefix, i)
			continue
		}

		if mem != nil && !v.ok {
			t.Errorf("%s[%d] failed: mem should be nil\n", prefix, i)
			continue
		}

		if allocator.Used() != used {
			t.Errorf("%s[%d] failed: wong used = %d, wanted = %d\n", prefix, i, allocator.Used(), used)
			continue
		}

		if allocator.AllocNum() != allocNum {
			t.Errorf("%s[%d] failed: wong AllocNum = %d, wanted = %d\n", prefix, i, allocator.AllocNum(), used)
			continue
		}

		if allocator.AllocNumOk() != allocNumOk {
			t.Errorf("%s[%d] failed: wong AllocNumOk = %d, wanted = %d\n", prefix, i, allocator.AllocNumOk(), used)
			continue
		}

		if allocator.FreeAllNum() != 0 {
			t.Errorf("%s[%d] failed: wong AllocNumOk = %d, wanted = 0\n", prefix, i, allocator.FreeAllNum())
			continue
		}

		if allocator.FreePartNum() != 0 {
			t.Errorf("%s[%d] failed: wong FreePartNum = %d, wanted = 0\n", prefix, i, allocator.FreePartNum())
			continue
		}
	}

	buf := []byte("%as\r\xe8")
	PrintByteSliceHex(0, len(buf), buf)
	PrintByteSliceHex(0, 10, nil)
	allocator.String(0, 128)
	allocator.String(-1, 128)
	allocator.String(10, 1)
	allocator.String(0, int(allocator.Capacity()+1))
	allocator.String(int(allocator.Capacity()+2), int(allocator.Capacity()+1))

	allocator.FreePart(100)
	if allocator.Used() != 100 {
		t.Errorf("%s failed: wong used = %d, wanted = 100\n", prefix, allocator.Used())
	}

	allocator.FreePart(200)
	if allocator.Used() != 100 {
		t.Errorf("%s failed: wong used = %d, wanted = 100\n", prefix, allocator.Used())
	}

	allocator.FreePart(-1)
	if allocator.Used() != 0 {
		t.Errorf("%s failed: wong used = %d, wanted = 100\n", prefix, allocator.Used())
	}

	allocator.FreeAll()

	if allocator.Used() != 0 {
		t.Errorf("%s failed: wong used = %d, wanted = %d\n", prefix, allocator.Used(), 0)
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
func BenchmarkMemAllocEx(b *testing.B) {
	b.StopTimer()
	allocator := NewMemAllocator(1024 * 128)
	b.ReportAllocs()
	b.SetBytes(2)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		allocator.FreeAll()
		_, _ = allocator.AllocEx(1000)
	}
	//fmt.Printf("uri = %s\n", uri.String())
	//fmt.Printf("")
}
