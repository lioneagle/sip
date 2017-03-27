package sipparser2

import (
	"container/list"
)

type AbnfList struct {
	list.List
}

func (this *AbnfList) RemoveAll() {
	var n *list.Element
	for e := this.Front(); e != nil; e = n {
		n = e.Next()
		this.Remove(e)
	}
}
