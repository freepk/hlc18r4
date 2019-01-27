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
	acc *proto.Account
	doc *inverted.Document
	rep *repo.AccountsRepo
}

func newDefaultIndexer(rep *repo.AccountsRepo) *defaultIndexer {
	acc := &proto.Account{}
	doc := &inverted.Document{ID: 0, Parts: make([]int, 1), Fields: make([][]int, 12)}
	return &defaultIndexer{pos: 0, acc: acc, doc: doc, rep: rep}
}

func (ix *defaultIndexer) Reset() {
	ix.pos = 0
}

func (ix *defaultIndexer) Next() (*inverted.Document, bool) {
	n := ix.rep.Len()
	for i := ix.pos; i < n; i++ {
		id := n - i - 1
		*ix.acc = *ix.rep.Get(id)
		if ix.acc.Email.Len > 0 {
			ix.pos = i + 1
			return ix.processDocument(id, ix.acc), true
		}
	}
	return nil, false
}

func (ix *defaultIndexer) resetDocument() *inverted.Document {
	doc := ix.doc
	doc.ID = 0
	doc.Parts = doc.Parts[:0]
	for field := range doc.Fields {
		doc.Fields[field] = doc.Fields[field][:0]
	}
	return doc
}

func (ix *defaultIndexer) processDocument(id int, acc *proto.Account) *inverted.Document {
	doc := ix.resetDocument()
	doc.ID = 2000000 - id
	doc.Parts = append(doc.Parts, defaultPartition)
	switch acc.Sex {
	case proto.MaleSex:
		doc.Fields[sexField] = append(doc.Fields[sexField], MaleToken)
	case proto.FemaleSex:
		doc.Fields[sexField] = append(doc.Fields[sexField], FemaleToken)
	}
	switch acc.Status {
	case proto.SingleStatus:
		doc.Fields[statusField] = append(doc.Fields[statusField], SingleToken, NotInRelToken, NotComplToken)
	case proto.InRelStatus:
		doc.Fields[statusField] = append(doc.Fields[statusField], InRelToken, NotSingleToken, NotComplToken)
	case proto.ComplStatus:
		doc.Fields[statusField] = append(doc.Fields[statusField], ComplToken, NotSingleToken, NotInRelToken)
	}
	if acc.Fname > 0 {
		doc.Fields[fnameField] = append(doc.Fields[fnameField], NotNullToken, int(acc.Fname))
	} else {
		doc.Fields[fnameField] = append(doc.Fields[fnameField], NullToken)
	}
	if acc.Sname > 0 {
		doc.Fields[snameField] = append(doc.Fields[snameField], NotNullToken, int(acc.Sname))
	} else {
		doc.Fields[snameField] = append(doc.Fields[snameField], NullToken)
	}
	if acc.Country > 0 {
		doc.Fields[countryField] = append(doc.Fields[countryField], NotNullToken, int(acc.Country))
	} else {
		doc.Fields[countryField] = append(doc.Fields[countryField], NullToken)
	}
	if acc.City > 0 {
		doc.Fields[cityField] = append(doc.Fields[cityField], NotNullToken, int(acc.City))
	} else {
		doc.Fields[cityField] = append(doc.Fields[cityField], NullToken)
	}
	for i := range acc.Interests {
		if acc.Interests[i] == 0 {
			break
		}
		doc.Fields[interestField] = append(doc.Fields[interestField], int(acc.Interests[i]))
	}
	doc.Fields[birthYearField] = append(doc.Fields[birthYearField], birthYearTokenTS(int(acc.BirthTS)))
	if acc.PremiumFinish[0] > 0 {
		doc.Fields[premiumField] = append(doc.Fields[premiumField], NotNullToken)
	} else {
		doc.Fields[premiumField] = append(doc.Fields[premiumField], NullToken)
	}
	if premiumNow(acc.PremiumFinish[:]) {
		doc.Fields[premiumField] = append(doc.Fields[premiumField], PremiumNowToken)
	}
	if acc.Phone[0] > 0 {
		doc.Fields[phoneCodeField] = append(doc.Fields[phoneCodeField], NotNullToken)
	} else {
		doc.Fields[phoneCodeField] = append(doc.Fields[phoneCodeField], NullToken)
	}
	if code, ok := phoneCode(acc.Phone[:]); ok {
		doc.Fields[phoneCodeField] = append(doc.Fields[phoneCodeField], phoneCodeToken(code))
	}
	if domain, ok := emailDomain(acc.Email.Buf[:acc.Email.Len]); ok {
		doc.Fields[emailDomainField] = append(doc.Fields[emailDomainField], emailDomainToken(domain))
	}
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

func (idx *DefaultIndex) Sex(token int) *inverted.TokenIter {
	return idx.inv.Iterator(defaultPartition, sexField, token)
}

func (idx *DefaultIndex) Status(token int) *inverted.TokenIter {
	return idx.inv.Iterator(defaultPartition, statusField, token)
}

func (idx *DefaultIndex) Fname(token int) *inverted.TokenIter {
	return idx.inv.Iterator(defaultPartition, fnameField, token)
}

func (idx *DefaultIndex) Sname(token int) *inverted.TokenIter {
	return idx.inv.Iterator(defaultPartition, snameField, token)
}

func (idx *DefaultIndex) Country(token int) *inverted.TokenIter {
	return idx.inv.Iterator(defaultPartition, countryField, token)
}

func (idx *DefaultIndex) City(token int) *inverted.TokenIter {
	return idx.inv.Iterator(defaultPartition, cityField, token)
}

func (idx *DefaultIndex) Interest(token int) *inverted.TokenIter {
	return idx.inv.Iterator(defaultPartition, interestField, token)
}

func (idx *DefaultIndex) BirthYear(token int) *inverted.TokenIter {
	return idx.inv.Iterator(defaultPartition, birthYearField, token)
}

func (idx *DefaultIndex) Premium(token int) *inverted.TokenIter {
	return idx.inv.Iterator(defaultPartition, premiumField, token)
}

func (idx *DefaultIndex) PhoneCode(token int) *inverted.TokenIter {
	return idx.inv.Iterator(defaultPartition, phoneCodeField, token)
}

func (idx *DefaultIndex) EmailDomain(token int) *inverted.TokenIter {
	return idx.inv.Iterator(defaultPartition, emailDomainField, token)
}
