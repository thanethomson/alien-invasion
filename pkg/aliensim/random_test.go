package aliensim

import "testing"

func TestSequenceGenerator(t *testing.T) {
	g := NewSequenceGenerator()
	if g.next != 0 {
		t.Error("Expected g.next == 0, but got ", g.next)
	}
	for i := uint32(0); i <= 100; i++ {
		if v := g.Uint32(); v != i {
			t.Error("Expected g.Uint32() == ", i, " but got ", v)
		}
	}
}
