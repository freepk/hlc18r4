package inverted

const reverseMaxID = uint32(10000000)

type ReverseIter struct {
	a []uint32
	n int
	i int
}

func NewReverseIter(a []uint32) *ReverseIter {
	n := len(a)
	return &ReverseIter{a: a, n: n, i: 0}
}

func (it *ReverseIter) Reset() {
	it.i = 0
}

func (it *ReverseIter) Next() (int, bool) {
	i := it.i
	if i < it.n {
		it.i++
		reverseID := it.a[i]
		return int(reverseMaxID - reverseID), true
	}
	return 0, false
}
