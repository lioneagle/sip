package sipparser3

import (
	"reflect"
	"unsafe"
)

type Element struct {
	next, prev *Element
	list       *SipList
	Value      SipUriParam
}

func (e *Element) Next() *Element {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

func (e *Element) Prev() *Element {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

func allocElement(allocator *MemAllocator) *Element {
	bytes := allocator.Alloc(int(unsafe.Sizeof(Element{})))
	return (*Element)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&bytes)).Data))
	//return nil
}

// SipList represents a doubly linked list.
// The zero value for SipList is an empty list ready to use.
type SipList struct {
	root      Element // sentinel list element, only &root, root.prev, and root.next are used
	len       int     // current list length excluding (this) sentinel element
	allocator *MemAllocator
}

func (this *SipList) allocElement(value SipUriParam) *Element {
	e := allocElement(this.allocator)
	e.Value = value
	return e
}

func NewSipList(allocator *MemAllocator) *SipList {
	return new(SipList).Init(allocator)
}

func (this *SipList) Init(allocator *MemAllocator) *SipList {
	this.root.next = &this.root
	this.root.prev = &this.root
	this.allocator = allocator
	this.len = 0
	return this
}

func (this *SipList) Len() int { return this.len }

// Front returns the first element of list this or nil.
func (this *SipList) Front() *Element {
	if this.len == 0 {
		return nil
	}
	return this.root.next
}

// Back returns the last element of list l or nil.
func (this *SipList) Back() *Element {
	if this.len == 0 {
		return nil
	}
	return this.root.prev
}

// lazyInit lazily initializes a zero List value.
func (this *SipList) lazyInit() {
	if this.root.next == nil {
		this.Init(this.allocator)
	}
}

// insert inserts e after at, increments this.len, and returns e.
func (this *SipList) insert(e, at *Element) *Element {
	n := at.next
	at.next = e
	e.prev = at
	e.next = n
	n.prev = e
	e.list = this
	this.len++
	return e
}

// insertValue is a convenience wrapper for insert(&Element{Value: v}, at).
func (this *SipList) insertValue(v SipUriParam, at *Element) *Element {
	return this.insert(this.allocElement(v), at)
}

// remove removes e from its list, decrements this.len, and returns e.
func (this *SipList) remove(e *Element) *Element {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	this.len--
	return e
}

// Remove removes e from this if e is an element of list this.
// It returns the element value e.Value.
func (this *SipList) Remove(e *Element) interface{} {
	if e.list == this {
		// if e.list == this, this must have been initialized when e was inserted
		// in this or this == nil (e is a zero Element) and this.remove will crash
		this.remove(e)
	}
	return e.Value
}

// PushFront inserts a new element e with value v at the front of list this and returns e.
func (this *SipList) PushFront(v SipUriParam) *Element {
	this.lazyInit()
	return this.insertValue(v, &this.root)
}

// PushBack inserts a new element e with value v at the back of list this and returns e.
func (this *SipList) PushBack(v SipUriParam) *Element {
	this.lazyInit()
	return this.insertValue(v, this.root.prev)
}

// InsertBefore inserts a new element e with value v immediately before mark and returns e.
// If mark is not an element of this, the list is not modified.
func (this *SipList) InsertBefore(v SipUriParam, mark *Element) *Element {
	if mark.list != this {
		return nil
	}
	// see comment in List.Remove about initialization of this
	return this.insertValue(v, mark.prev)
}

// InsertAfter inserts a new element e with value v immediately after mark and returns e.
// If mark is not an element of this, the list is not modified.
func (this *SipList) InsertAfter(v SipUriParam, mark *Element) *Element {
	if mark.list != this {
		return nil
	}
	// see comment in List.Remove about initialization of this
	return this.insertValue(v, mark)
}

// MoveToFront moves element e to the front of list this.
// If e is not an element of this, the list is not modified.
func (this *SipList) MoveToFront(e *Element) {
	if e.list != this || this.root.next == e {
		return
	}
	// see comment in List.Remove about initialization of this
	this.insert(this.remove(e), &this.root)
}

// MoveToBack moves element e to the back of list this.
// If e is not an element of this, the list is not modified.
func (this *SipList) MoveToBack(e *Element) {
	if e.list != this || this.root.prev == e {
		return
	}
	// see comment in List.Remove about initialization of this
	this.insert(this.remove(e), this.root.prev)
}

// MoveBefore moves element e to its new position before mark.
// If e or mark is not an element of this, or e == mark, the list is not modified.
func (this *SipList) MoveBefore(e, mark *Element) {
	if e.list != this || e == mark || mark.list != this {
		return
	}
	this.insert(this.remove(e), mark.prev)
}

// MoveAfter moves element e to its new position after mark.
// If e or mark is not an element of this, or e == mark, the list is not modified.
func (this *SipList) MoveAfter(e, mark *Element) {
	if e.list != this || e == mark || mark.list != this {
		return
	}
	this.insert(this.remove(e), mark)
}

// PushBackList inserts a copy of an other list at the back of list this.
// The lists this and other may be the same.
func (this *SipList) PushBackList(other *SipList) {
	this.lazyInit()
	for i, e := other.Len(), other.Front(); i > 0; i, e = i-1, e.Next() {
		this.insertValue(e.Value, this.root.prev)
	}
}

// PushFrontList inserts a copy of an other list at the front of list this.
// The lists this and other may be the same.
func (this *SipList) PushFrontList(other *SipList) {
	this.lazyInit()
	for i, e := other.Len(), other.Back(); i > 0; i, e = i-1, e.Prev() {
		this.insertValue(e.Value, &this.root)
	}
}

func (this *SipList) RemoveAll() {
	var n *Element
	for e := this.Front(); e != nil; e = n {
		n = e.Next()
		this.Remove(e)
	}
}
