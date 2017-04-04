package sipparser

import (
	//"bytes"
	"fmt"
	"testing"
)

func checkEmptyList(t *testing.T, prefix string, list *AbnfList) bool {
	if list.head != ABNF_PTR_NIL {
		t.Errorf("%s failed: wrong head = %d, wanted = %s\n", prefix, list.head, ABNF_PTR_NIL.String())
		return false
	}

	if list.tail != ABNF_PTR_NIL {
		t.Errorf("%s failed: wrong tail = %d, wanted = %s\n", prefix, list.tail, ABNF_PTR_NIL.String())
		return false
	}

	if list.size != 0 {
		t.Errorf("%s failed: wrong size = %d, wanted = 0\n", prefix, list.size)
		return false
	}

	return true
}

func checkNodeValue(t *testing.T, context *ParseContext, prefix string, node *AbnfListNode, v AbnfPtr) bool {
	if node.Value != v {
		t.Errorf("%s failed: wrong Value = %d, wanted = %d\n", prefix, node.Value, v)
		return false
	}
	return true
}

func checkListNode(t *testing.T, context *ParseContext, prefix string, node1 *AbnfListNode, node2 *AbnfListNode) bool {
	if node1.next != node2.next {
		t.Errorf("%s failed: wrong next = %d, wanted = %d\n", prefix, node1.next, node2.next)
		return false
	}
	if node1.prev != node2.prev {
		t.Errorf("%s failed: wrong prev = %d, wanted = %d\n", prefix, node1.prev, node2.prev)
		return false
	}
	if node1.Value != node2.Value {
		t.Errorf("%s failed: wrong Value = %d, wanted = %d\n", prefix, node1.Value, node2.Value)
		return false
	}
	return true
}

func checkList(t *testing.T, context *ParseContext, prefix string, list *AbnfList, nodes []*AbnfListNode) bool {
	if list.Len() != int32(len(nodes)) {
		t.Errorf("%s failed: wrong size = %d, wanted = %d\n", prefix, list.Len(), len(nodes))
		return false
	}

	if len(nodes) == 0 {
		return checkEmptyList(t, prefix, list)
	}

	var iter *AbnfListNode

	iter = list.Front(context)
	for i, node := range nodes {
		new_prefix := fmt.Sprintf("%s[%d]", prefix, i)
		if !checkListNode(t, context, new_prefix, iter, node) {
			return false
		}
		iter = iter.Next(context)
	}

	return true
}

func TestAbnfListNew(t *testing.T) {
	context := NewParseContext()
	context.allocator = NewMemAllocator(1000)
	prefix := FuncName()

	list, addr := NewAbnfList(context)
	if list == nil {
		t.Errorf("%s failed: should be ok\n", prefix)
		return
	}

	if addr != 0 {
		t.Errorf("%s failed: wrong addr = %d, wanted = 0\n", prefix, addr)
		return
	}

	if !checkEmptyList(t, prefix, list) {
		return
	}

	context.allocator = NewMemAllocator(1)
	list, _ = NewAbnfList(context)
	if list != nil {
		t.Errorf("%s failed: should not be ok\n", prefix)
		return
	}
}

func TestAbnfListAdd(t *testing.T) {
	context := NewParseContext()
	context.allocator = NewMemAllocator(1000)
	list, _ := NewAbnfList(context)
	prefix := FuncName()

	if list.Front(context) != nil {
		t.Errorf("%s failed: Front should be nil\n", prefix)
	}

	if list.Back(context) != nil {
		t.Errorf("%s failed: Back should be nil\n", prefix)
	}

	node1 := list.PushBack(context, 1)

	if node1.Prev(context) != nil {
		t.Errorf("%s failed: Prev should be nil\n", prefix)
	}

	if node1.Next(context) != nil {
		t.Errorf("%s failed: Next should be nil\n", prefix)
	}

	checkListNode(t, context, prefix, list.Front(context), node1)
	checkListNode(t, context, prefix, list.Back(context), node1)
	checkList(t, context, prefix, list, []*AbnfListNode{node1})

	list.PopBack(context)
	checkList(t, context, prefix, list, []*AbnfListNode{})

	node1 = list.PushFront(context, 1)
	checkNodeValue(t, context, prefix, node1, 1)
	checkList(t, context, prefix, list, []*AbnfListNode{node1})
	node2 := list.PushBack(context, 2)
	checkNodeValue(t, context, prefix, node2, 2)
	checkListNode(t, context, prefix, node1, node2.Prev(context))
	checkListNode(t, context, prefix, node2, node1.Next(context))
	checkList(t, context, prefix, list, []*AbnfListNode{node1, node2})
	node3 := list.PushFront(context, 3)
	checkList(t, context, prefix, list, []*AbnfListNode{node3, node1, node2})

}

func TestAbnfListPushBack(t *testing.T) {
	context := NewParseContext()
	context.allocator = NewMemAllocator(1000)
	list, _ := NewAbnfList(context)
	prefix := FuncName()

	if list.PopBack(context) != nil {
		t.Errorf("%s failed: empty list PopBack should be nil\n", prefix)
	}

	list.PushBack(context, 1)
	list.PopBack(context)
	checkList(t, context, prefix, list, []*AbnfListNode{})

	node1 := list.PushBack(context, 1)
	list.PushBack(context, 2)
	list.PopBack(context)
	checkList(t, context, prefix, list, []*AbnfListNode{node1})
	list.PopBack(context)
	checkList(t, context, prefix, list, []*AbnfListNode{})
}

func TestAbnfListPushFront(t *testing.T) {
	context := NewParseContext()
	context.allocator = NewMemAllocator(1000)
	list, _ := NewAbnfList(context)
	prefix := FuncName()

	if list.PopFront(context) != nil {
		t.Errorf("%s failed: empty list PopFront should be nil\n", prefix)
	}

	list.PushFront(context, 1)
	list.PopFront(context)
	checkList(t, context, prefix, list, []*AbnfListNode{})

	node1 := list.PushFront(context, 1)
	list.PushFront(context, 2)
	list.PopFront(context)
	checkList(t, context, prefix, list, []*AbnfListNode{node1})
	list.PopFront(context)
	checkList(t, context, prefix, list, []*AbnfListNode{})

}

func TestAbnfListRemove(t *testing.T) {
	context := NewParseContext()
	context.allocator = NewMemAllocator(1000)
	list, _ := NewAbnfList(context)
	prefix := FuncName()

	if list.Remove(context, nil) != nil {
		t.Errorf("%s failed: empty list Remove should be nil\n", prefix)
	}

	node1 := list.PushBack(context, 1)
	list.Remove(context, node1)
	checkList(t, context, prefix, list, []*AbnfListNode{})

	node1 = list.PushFront(context, 1)
	list.Remove(context, node1)
	checkList(t, context, prefix, list, []*AbnfListNode{})

	node1 = list.PushBack(context, 1)
	node2 := list.PushBack(context, 2)
	list.Remove(context, node2)
	checkList(t, context, prefix, list, []*AbnfListNode{node1})
	list.PopFront(context)
	checkList(t, context, prefix, list, []*AbnfListNode{})

	node1 = list.PushBack(context, 1)
	node2 = list.PushBack(context, 2)
	node3 := list.PushBack(context, 3)
	checkNodeValue(t, context, prefix, node1, 1)
	checkNodeValue(t, context, prefix, node2, 2)
	checkNodeValue(t, context, prefix, node3, 3)
	list.Remove(context, node2)
	checkList(t, context, prefix, list, []*AbnfListNode{node1, node3})
	list.RemoveAll(context)
	checkList(t, context, prefix, list, []*AbnfListNode{})
}
