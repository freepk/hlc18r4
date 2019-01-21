package inverted

type TokenIter struct {
	a []uint32
	n int
	i int
}

func NewTokenIter(a []uint32) *TokenIter {
	n := len(a)
	return &TokenIter{a: a, n: n, i: 0}
}

func (it *TokenIter) Reset() {
	it.i = 0
}

func (it *TokenIter) Next() (int, bool) {
	i := it.i
	if i < it.n {
		it.i++
		return int(it.a[i]), true
	}
	return 0, false
}
