package consistent

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"sort"
)

type vnode[T any] struct {
	h     uint32
	value T
}

type sortable[T any] []vnode[T]

func (s *sortable[T]) Len() int {
	return len(*s)
}

func (s *sortable[T]) Swap(i, j int) {
	v := (*s)[i]
	(*s)[i] = (*s)[j]
	(*s)[j] = v
}

func (s *sortable[T]) Less(i, j int) bool {
	return (*s)[i].h < (*s)[j].h
}

type Ketama[T comparable] struct {
	nodes []vnode[T]
}

func (k *Ketama[T]) Reset(capacity int) {
	k.nodes = make([]vnode[T], 0, capacity)
}

const perHash = 4

func (k *Ketama[T]) Add(hashTag string, vn int, value T) {
	for i := 0; i < vn; i += perHash {
		sum := md5.Sum([]byte(fmt.Sprintf("%s/%d", hashTag, i)))
		n := min(vn-i, perHash)
		for j := 0; j < n; j++ {
			k.nodes = append(k.nodes, vnode[T]{
				h:     binary.LittleEndian.Uint32(sum[j*4:]),
				value: value,
			})
		}
	}
}

func (k *Ketama[T]) Remove(value T) {
	j := 0
	for i := 0; i < len(k.nodes); i++ {
		if k.nodes[i].value == value {
			continue
		}
		if i != j {
			k.nodes[j] = k.nodes[i]
		}
		j++
	}
	k.nodes = k.nodes[:j]
}

func (k *Ketama[T]) Build() {
	sort.Sort((*sortable[T])(&k.nodes))
}

func (k *Ketama[T]) Get32(key uint32) (T, bool) {
	if len(k.nodes) == 0 {
		var zero T
		return zero, false
	}
	idx := sort.Search(len(k.nodes), func(i int) bool {
		return key <= k.nodes[i].h
	})
	if idx == len(k.nodes) {
		idx = 0
	}
	return k.nodes[idx].value, true
}

func (k *Ketama[T]) Get(key string) (T, bool) {
	return k.Get32(crc32.ChecksumIEEE([]byte(key)))
}
