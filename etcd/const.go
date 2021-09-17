package etcd

import (
	"sync"
)

type node2Url struct {
  sync.RWMutex
  urlList []string
}

func(n *node2Url) Len() int {
	n.RLock()
	defer n.RUnlock()
	return len(n.urlList)
}

func(n *node2Url) GetUrl(index uint64) string {
	n.RLock()
	defer n.RUnlock()
	l := n.Len()
	if l > 0 {
		return n.urlList[index%uint64(l)]
	}
	return ""
}
