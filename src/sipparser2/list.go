package sipparser2

import (
	"container/list"
)

type SipList struct {
	list.List
}

func (this *SipList) RemoveAll() {
	var n *list.Element
	for e := this.Front(); e != nil; e = n {
		n = e.Next()
		this.Remove(e)
	}
}
