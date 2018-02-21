package consistent

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"sort"
)

type vnode struct {
	h     uint32
	value string
}

type sortable []vnode

func (s *sortable) Len() int {
	return len(*s)
}

func (s *sortable) Swap(i, j int) {
	v := (*s)[i]
	(*s)[i] = (*s)[j]
	(*s)[j] = v
}

func (s *sortable) Less(i, j int) bool {
	return (*s)[i].h < (*s)[j].h
}

type Ketama struct {
	nodes []vnode
}

func (k *Ketama) Reset(capacity int) {
	if capacity == 0 {
		k.nodes = nil
	} else {
		if capacity > cap(k.nodes) {
			k.nodes = make([]vnode, 0, capacity)
		} else {
			k.nodes = k.nodes[:0:capacity]
		}
	}
}

func (k *Ketama) Add(value string, n int) {
	const perhash = 4
	n = (n + (perhash - 1)) / perhash
	for i := 0; i < n; i++ {
		sum := md5.Sum([]byte(fmt.Sprintf("%s+%d", value, i)))
		for j := 0; j < perhash; j++ {
			k.nodes = append(k.nodes, vnode{
				h:     binary.LittleEndian.Uint32(sum[j*4:]),
				value: value,
			})
		}
	}
}

func (k *Ketama) Build() {
	sort.Sort((*sortable)(&k.nodes))
}

func (k *Ketama) Get32(key uint32) (string, bool) {
	if len(k.nodes) == 0 {
		return "", false
	}
	idx := sort.Search(len(k.nodes), func(i int) bool {
		return key <= k.nodes[i].h
	})
	if idx == len(k.nodes) {
		idx = 0
	}
	return k.nodes[idx].value, true
}

func (k *Ketama) Get(key string) (string, bool) {
	return k.Get32(crc32.ChecksumIEEE([]byte(key)))
}
