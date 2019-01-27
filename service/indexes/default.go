package indexes

import (
	"gitlab.com/freepk/hlc18r4/inverted"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

const (
	defaultPartition = 0
)

const (
	sexField = iota
	statusField
	fnameField
	snameField
	snamePrefixField
	countryField
	cityField
	interestField
	birthYearField
	premiumField
	phoneCodeField
	emailDomainField
)

type defaultIndexer struct {
	pos int
	doc *inverted.Document
	rep *repo.AccountsRepo
}

func newDefaultIndexer(rep *repo.AccountsRepo) *defaultIndexer {
	doc := &inverted.Document{ID: 0, Parts: make([]int, 1), Tokens: make([][]int, 12)}
	return &defaultIndexer{pos: 0, doc: doc, rep: rep}
}

func (ix *defaultIndexer) Reset() {
	ix.pos = 0
}

func (ix *defaultIndexer) Next() (*inverted.Document, bool) {
	if id, ok := ix.next(); ok {
		return ix.processDocument(id), true
	}
	return nil, false
}

func (ix *defaultIndexer) next() (int, bool) {
	n := ix.rep.Len()
	acc := proto.Account{}
	for i := ix.pos; i < n; i++ {
		id := n - i - 1
		acc = *ix.rep.Get(id)
		if acc.Email.Len > 0 {
			ix.pos = i + 1
			return id, true
		}
	}
	return 0, false
}

func (ix *defaultIndexer) resetDocument() *inverted.Document {
	doc := ix.doc
	doc.ID = 0
	doc.Parts = doc.Parts[:0]
	for field := range doc.Tokens {
		doc.Tokens[field] = doc.Tokens[field][:0]
	}
	return doc
}

func (ix *defaultIndexer) processDocument(id int) *inverted.Document {
	acc := *ix.rep.Get(id)
	doc := ix.resetDocument()
	doc.ID = 2000000 - id
	doc.Parts = append(doc.Parts, defaultPartition)
	switch acc.Sex {
	case proto.MaleSex:
		doc.Tokens[sexField] = append(doc.Tokens[sexField], MaleToken)
	case proto.FemaleSex:
		doc.Tokens[sexField] = append(doc.Tokens[sexField], FemaleToken)
	}
	switch acc.Status {
	case proto.SingleStatus:
		doc.Tokens[statusField] = append(doc.Tokens[statusField], SingleToken, NotInRelToken, NotComplToken)
	case proto.InRelStatus:
		doc.Tokens[statusField] = append(doc.Tokens[statusField], InRelToken, NotSingleToken, NotComplToken)
	case proto.ComplStatus:
		doc.Tokens[statusField] = append(doc.Tokens[statusField], ComplToken, NotSingleToken, NotInRelToken)
	}
	if acc.Fname > 0 {
		doc.Tokens[fnameField] = append(doc.Tokens[fnameField], NotNullToken, int(acc.Fname))
	} else {
		doc.Tokens[fnameField] = append(doc.Tokens[fnameField], NullToken)
	}
	if acc.Sname > 0 {
		doc.Tokens[snameField] = append(doc.Tokens[snameField], NotNullToken, int(acc.Sname))
	} else {
		doc.Tokens[snameField] = append(doc.Tokens[snameField], NullToken)
	}
	if acc.Country > 0 {
		doc.Tokens[countryField] = append(doc.Tokens[countryField], NotNullToken, int(acc.Country))
	} else {
		doc.Tokens[countryField] = append(doc.Tokens[countryField], NullToken)
	}
	if acc.City > 0 {
		doc.Tokens[cityField] = append(doc.Tokens[cityField], NotNullToken, int(acc.City))
	} else {
		doc.Tokens[cityField] = append(doc.Tokens[cityField], NullToken)
	}
	for i := range acc.Interests {
		if acc.Interests[i] == 0 {
			break
		}
		doc.Tokens[interestField] = append(doc.Tokens[interestField], int(acc.Interests[i]))
	}
	doc.Tokens[birthYearField] = append(doc.Tokens[birthYearField], birthYearTokenTS(int(acc.BirthTS)))
	if acc.PremiumFinish[0] > 0 {
		doc.Tokens[premiumField] = append(doc.Tokens[premiumField], NotNullToken)
	} else {
		doc.Tokens[premiumField] = append(doc.Tokens[premiumField], NullToken)
	}
	if premiumNow(acc.PremiumFinish[:]) {
		doc.Tokens[premiumField] = append(doc.Tokens[premiumField], PremiumNowToken)
	}
	if acc.Phone[0] > 0 {
		doc.Tokens[phoneCodeField] = append(doc.Tokens[phoneCodeField], NotNullToken)
	} else {
		doc.Tokens[phoneCodeField] = append(doc.Tokens[phoneCodeField], NullToken)
	}
	//if code, ok := phoneCode(acc.Phone); ok {
	//	doc.Tokens[phoneCodeField] = append(doc.Tokens[phoneCodeField], phoneCodeToken(code))
	//}
	//if domain, ok := emailDomain(acc.Email.Buf[:acc.Email.Len]); ok {
	//	doc.Tokens[emailDomainField] = append(doc.Tokens[emailDomainField], emailDomainToken(domain))
	//}
	return doc
}

type DefaultIndex struct {
	inv *inverted.Inverted
}

func NewDefaultIndex(rep *repo.AccountsRepo) *DefaultIndex {
	src := newDefaultIndexer(rep)
	inv := inverted.NewInverted(src)
	return &DefaultIndex{inv: inv}
}

func (idx *DefaultIndex) Rebuild() {
	idx.inv.Rebuild()
}

func (idx *DefaultIndex) Sex(token int) *inverted.ArrayIter {
	return idx.inv.Iterator(defaultPartition, sexField, token)
}

func (idx *DefaultIndex) Status(token int) *inverted.ArrayIter {
	return idx.inv.Iterator(defaultPartition, statusField, token)
}

func (idx *DefaultIndex) Fname(token int) *inverted.ArrayIter {
	return idx.inv.Iterator(defaultPartition, fnameField, token)
}

func (idx *DefaultIndex) Sname(token int) *inverted.ArrayIter {
	return idx.inv.Iterator(defaultPartition, snameField, token)
}

func (idx *DefaultIndex) Country(token int) *inverted.ArrayIter {
	return idx.inv.Iterator(defaultPartition, countryField, token)
}

func (idx *DefaultIndex) City(token int) *inverted.ArrayIter {
	return idx.inv.Iterator(defaultPartition, cityField, token)
}

func (idx *DefaultIndex) Interest(token int) *inverted.ArrayIter {
	return idx.inv.Iterator(defaultPartition, interestField, token)
}

func (idx *DefaultIndex) BirthYear(token int) *inverted.ArrayIter {
	return idx.inv.Iterator(defaultPartition, birthYearField, token)
}

func (idx *DefaultIndex) Premium(token int) *inverted.ArrayIter {
	return idx.inv.Iterator(defaultPartition, premiumField, token)
}

func (idx *DefaultIndex) PhoneCode(token int) *inverted.ArrayIter {
	return idx.inv.Iterator(defaultPartition, phoneCodeField, token)
}

func (idx *DefaultIndex) EmailDomain(token int) *inverted.ArrayIter {
	return idx.inv.Iterator(defaultPartition, emailDomainField, token)
}
