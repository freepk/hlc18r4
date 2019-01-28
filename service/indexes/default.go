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
	countryField
	cityField
	interestField
	fnameField
	snameField
	birthYearField
	premiumField
	phoneCodeField
	emailDomainField
)

type defaultIter struct {
	accountsIter
}

func newDefaultIter(rep *repo.AccountsRepo) *defaultIter {
	return &defaultIter{accountsIter: *newAccountsIter(rep, 12)}
}

func (it *defaultIter) Next() (*inverted.Document, bool) {
	if it.next() {
		return it.processDocument(), true
	}
	return nil, false
}

func (it *defaultIter) processDocument() *inverted.Document {
	it.resetDocument()
	acc := it.acc
	doc := it.doc
	doc.ID = 2000000 - it.id()
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
	doc.Fields[birthYearField] = append(doc.Fields[birthYearField], yearTokenTS(int(acc.BirthTS)))
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
	src := newDefaultIter(rep)
	inv := inverted.NewInverted(src)
	return &DefaultIndex{inv: inv}
}

func (idx *DefaultIndex) Rebuild() {
	idx.inv.Rebuild()
}

func (idx *DefaultIndex) SexEq() {
}

func (idx *DefaultIndex) StatusEq() {
}

func (idx *DefaultIndex) StatusNeq() {
}

func (idx *DefaultIndex) EmailDomain() {
}

func (idx *DefaultIndex) FnameEq() {
}

func (idx *DefaultIndex) FnameAny() {
}

func (idx *DefaultIndex) FnameNull() {
}

func (idx *DefaultIndex) SnameEq() {
}

func (idx *DefaultIndex) SnameNull() {
}

func (idx *DefaultIndex) PhoneCode() {
}

func (idx *DefaultIndex) PhoneNull() {
}

func (idx *DefaultIndex) CountryEq() {
}

func (idx *DefaultIndex) CountryNull() {
}

func (idx *DefaultIndex) CityEq() {
}

func (idx *DefaultIndex) CityAny() {
}

func (idx *DefaultIndex) CityNull() {
}

func (idx *DefaultIndex) BirthYear() {
}

func (idx *DefaultIndex) InterestsAny() {
}

func (idx *DefaultIndex) InterestsContains() {
}
