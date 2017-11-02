package sipparser

import (
	//"fmt"
	//"reflect"
	"unsafe"
)

type AbnfListNode struct {
	//AbnfPtr
	next  AbnfPtr
	prev  AbnfPtr
	Value AbnfPtr
}

func NewAbnfListNode(context *ParseContext) AbnfPtr {
	addr := context.allocator.Alloc(uint32(unsafe.Sizeof(AbnfListNode{})))
	if addr == ABNF_PTR_NIL {
		return ABNF_PTR_NIL
	}

	addr.GetAbnfListNode(context).Init()
	return addr
}

func (this *AbnfListNode) Init() {
	this.next = ABNF_PTR_NIL
	this.prev = ABNF_PTR_NIL
	this.Value = ABNF_PTR_NIL
}

func (this *AbnfListNode) Next(context *ParseContext) *AbnfListNode {
	if this.next == ABNF_PTR_NIL {
		return nil
	}
	return this.next.GetAbnfListNode(context)
}

func (this *AbnfListNode) Prev(context *ParseContext) *AbnfListNode {
	if this.prev == ABNF_PTR_NIL {
		return nil
	}
	return this.prev.GetAbnfListNode(context)
}

type AbnfList struct {
	head AbnfPtr
	tail AbnfPtr
	size int32
}

func NewAbnfList(context *ParseContext) AbnfPtr {
	addr := context.allocator.Alloc(uint32(unsafe.Sizeof(AbnfList{})))
	if addr == ABNF_PTR_NIL {
		return ABNF_PTR_NIL
	}
	addr.GetAbnfList(context).Init()
	return addr
}

func (this *AbnfList) Init() {
	this.head = ABNF_PTR_NIL
	this.tail = ABNF_PTR_NIL
	this.size = 0
}

func (this *AbnfList) Len() int32 { return this.size }

func (this *AbnfList) Front(context *ParseContext) *AbnfListNode {
	if this.size == 0 {
		return nil
	}
	return this.head.GetAbnfListNode(context)
}

func (this *AbnfList) Back(context *ParseContext) *AbnfListNode {
	if this.size == 0 {
		return nil
	}
	return this.tail.GetAbnfListNode(context)
}

func (this *AbnfList) PushBack(context *ParseContext, value AbnfPtr) (node *AbnfListNode) {
	addr := NewAbnfListNode(context)
	if addr == ABNF_PTR_NIL {
		return nil
	}

	node = addr.GetAbnfListNode(context)
	node.prev = this.tail
	node.Value = value

	if this.size == 0 {
		this.head = addr
	} else {
		this.tail.GetAbnfListNode(context).next = addr
	}
	this.tail = addr
	this.size++
	return node
}

func (this *AbnfList) PushFront(context *ParseContext, value AbnfPtr) (node *AbnfListNode) {
	addr := NewAbnfListNode(context)
	if addr == ABNF_PTR_NIL {
		return nil
	}

	node = addr.GetAbnfListNode(context)
	node.next = this.head
	node.Value = value

	if this.size == 0 {
		this.tail = addr
	} else {
		this.head.GetAbnfListNode(context).prev = addr
	}

	this.head = addr
	this.size++
	return node
}

func (this *AbnfList) PopBack(context *ParseContext) *AbnfListNode {
	if this.size == 0 {
		return nil
	}

	tail := this.tail.GetAbnfListNode(context)

	if this.size == 1 {
		this.head = ABNF_PTR_NIL
		this.tail = ABNF_PTR_NIL
	} else {
		this.tail = tail.prev
		this.tail.GetAbnfListNode(context).next = ABNF_PTR_NIL
	}

	tail.next = ABNF_PTR_NIL
	tail.prev = ABNF_PTR_NIL
	this.size--
	return nil
}

func (this *AbnfList) PopFront(context *ParseContext) *AbnfListNode {
	if this.size == 0 {
		return nil
	}

	head := this.head.GetAbnfListNode(context)

	if this.size == 1 {
		this.head = ABNF_PTR_NIL
		this.tail = ABNF_PTR_NIL
	} else {
		this.head = head.next
		this.head.GetAbnfListNode(context).prev = ABNF_PTR_NIL
	}

	head.next = ABNF_PTR_NIL
	head.prev = ABNF_PTR_NIL
	this.size--

	if this.head == ABNF_PTR_NIL {
		return nil
	}
	return this.head.GetAbnfListNode(context)
}

func (this *AbnfList) Remove(context *ParseContext, e *AbnfListNode) *AbnfListNode {
	if this.size == 0 {
		return nil
	}

	if e.prev == ABNF_PTR_NIL {
		return this.PopFront(context)
	}

	if e.next == ABNF_PTR_NIL {
		return this.PopBack(context)
	}

	e.prev.GetAbnfListNode(context).next = e.next
	next := e.next.GetAbnfListNode(context)
	next.prev = e.prev

	e.next = ABNF_PTR_NIL
	e.prev = ABNF_PTR_NIL
	this.size--
	return next
}

func (this *AbnfList) RemoveAll(context *ParseContext) {
	this.Init()
}
