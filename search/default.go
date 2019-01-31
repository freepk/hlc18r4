package search

import (
	//"fmt"

	"github.com/freepk/inverted"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
	"gitlab.com/freepk/hlc18r4/tokens"
)

const (
	DefaultPartition = 0
)

const (
	SexField = iota
	StatusField
	CountryField
	CityField
	InterestField
	FnameField
	SnameField
	BirthYearField
	PremiumField
	PhoneCodeField
	EmailDomainField
)

type DefaultIndex struct {
	inv *inverted.Inverted
}

func NewDefaultIndex(rep *repo.AccountsRepo) *DefaultIndex {
	proc := NewAccountsProcessor(rep, DefaultProc)
	inv := inverted.NewInverted(proc)
	return &DefaultIndex{inv: inv}
}

func (idx *DefaultIndex) Rebuild() {
	idx.inv.Rebuild()
}

func DefaultProc(doc *inverted.Document, acc *proto.Account) {
	//fmt.Println(doc.ID)
	return
	doc.Parts = append(doc.Parts, DefaultPartition)
	switch acc.Sex {
	case tokens.MaleSex:
		doc.Fields[SexField] = append(doc.Fields[SexField], tokens.MaleSex)
	case tokens.FemaleSex:
		doc.Fields[SexField] = append(doc.Fields[SexField], tokens.FemaleSex)
	}
	/*
		switch acc.Status {
		case proto.SingleStatus:
			doc.Fields[StatusField] = append(doc.Fields[StatusField], SingleToken, NotInRelToken, NotComplToken)
		case proto.InRelStatus:
			doc.Fields[StatusField] = append(doc.Fields[StatusField], InRelToken, NotSingleToken, NotComplToken)
		case proto.ComplStatus:
			doc.Fields[StatusField] = append(doc.Fields[StatusField], ComplToken, NotSingleToken, NotInRelToken)
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
	*/
}

/*
func (idx *DefaultIndex) part() *inverted.Part {
	return idx.inv.Part(defaultPartition)
}

func (idx *DefaultIndex) Sex(t int) *inverted.TokenIter {
	return idx.part().Field(sexField).Iterator(t)
}

func (idx *DefaultIndex) Status(t int) *inverted.TokenIter {
	return idx.part().Field(statusField).Iterator(t)
}

func (idx *DefaultIndex) EmailDomain(t int) *inverted.TokenIter {
	return idx.part().Field(emailDomainField).Iterator(t)
}

func (idx *DefaultIndex) Fname(t int) *inverted.TokenIter {
	return idx.part().Field(fnameField).Iterator(t)
}

func (idx *DefaultIndex) Sname(t int) *inverted.TokenIter {
	return idx.part().Field(snameField).Iterator(t)
}

func (idx *DefaultIndex) PhoneCode(t int) *inverted.TokenIter {
	return idx.part().Field(phoneCodeField).Iterator(t)
}
func (idx *DefaultIndex) Country(t int) *inverted.TokenIter {
	return idx.part().Field(countryField).Iterator(t)
}

func (idx *DefaultIndex) City(t int) *inverted.TokenIter {
	return idx.part().Field(cityField).Iterator(t)
}
func (idx *DefaultIndex) BirthYear(t int) *inverted.TokenIter {
	return idx.part().Field(birthYearField).Iterator(t)
}

func (idx *DefaultIndex) Interest(t int) *inverted.TokenIter {
	return idx.part().Field(interestField).Iterator(t)
}

func (idx *DefaultIndex) Premium(t int) *inverted.TokenIter {
	return idx.part().Field(premiumField).Iterator(t)
}
*/
