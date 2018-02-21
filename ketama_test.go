package consistent

import (
	"math/rand"
	"testing"
)

func TestKetama(t *testing.T) {
	k := Ketama{}
	t.Log(k.Get("a"))
	k.Reset(128 * 3)
	k.Add("a", 128)
	k.Add("b", 128)
	k.Add("c", 128)
	k.Build()
	t.Log(k.Get("a"))
	r := rand.New(rand.NewSource(0))
	cnt := map[string]int{}
	for i := 0; i < 900; i++ {
		v, _ := k.Get32(r.Uint32())
		cnt[v]++
	}
	t.Log(cnt)
}
