package inverted

type ArrayIter struct {
	a []uint32
	i int
}

func NewArrayIter(a []uint32) *ArrayIter {
	return &ArrayIter{a: a, i: 0}
}

func (it *ArrayIter) Reset() {
	it.i = 0
}

func (it *ArrayIter) Next() (int, bool) {
	i := it.i
	if i < len(it.a) {
		it.i++
		return int(it.a[i]), true
	}
	return 0, false
}
