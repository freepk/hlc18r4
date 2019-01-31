package search

import (
	"github.com/freepk/inverted"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
	"gitlab.com/freepk/hlc18r4/tokens"
)

type CommonIndex struct {
	inv *inverted.Inverted
}

func NewCommonIndex(rep *repo.AccountsRepo) *CommonIndex {
	proc := NewAccountsProcessor(rep, CommonProc, 1, 12)
	inv := inverted.NewInverted(proc)
	return &CommonIndex{inv: inv}
}

func (idx *CommonIndex) Rebuild() {
	idx.inv.Rebuild()
}

func CommonProc(doc *inverted.Document, acc *proto.Account) {
	doc.Parts = append(doc.Parts, CommonPartition)
	switch acc.Sex {
	case tokens.MaleSex:
		doc.Fields[SexField] = append(doc.Fields[SexField], tokens.MaleSex)
	case tokens.FemaleSex:
		doc.Fields[SexField] = append(doc.Fields[SexField], tokens.FemaleSex)
	}
	switch acc.Status {
	case tokens.SingleStatus:
		doc.Fields[StatusField] = append(doc.Fields[StatusField], tokens.SingleStatus, tokens.NotInRelStatus, tokens.NotComplStatus)
	case tokens.InRelStatus:
		doc.Fields[StatusField] = append(doc.Fields[StatusField], tokens.InRelStatus, tokens.NotSingleStatus, tokens.NotComplStatus)
	case tokens.ComplStatus:
		doc.Fields[StatusField] = append(doc.Fields[StatusField], tokens.ComplStatus, tokens.NotSingleStatus, tokens.NotInRelStatus)
	}
	if acc.Fname > 0 {
		doc.Fields[FnameField] = append(doc.Fields[FnameField], tokens.NotNull, int(acc.Fname))
	} else {
		doc.Fields[FnameField] = append(doc.Fields[FnameField], tokens.Null)
	}
	if acc.Sname > 0 {
		doc.Fields[SnameField] = append(doc.Fields[SnameField], tokens.NotNull, int(acc.Sname))
	} else {
		doc.Fields[SnameField] = append(doc.Fields[SnameField], tokens.Null)
	}
	if acc.Country > 0 {
		doc.Fields[CountryField] = append(doc.Fields[CountryField], tokens.NotNull, int(acc.Country))
	} else {
		doc.Fields[CountryField] = append(doc.Fields[CountryField], tokens.Null)
	}
	if acc.City > 0 {
		doc.Fields[CityField] = append(doc.Fields[CityField], tokens.NotNull, int(acc.City))
	} else {
		doc.Fields[CityField] = append(doc.Fields[CityField], tokens.Null)
	}
	for i := range acc.Interests {
		if acc.Interests[i] == 0 {
			break
		}
		doc.Fields[InterestField] = append(doc.Fields[InterestField], int(acc.Interests[i]))
	}
	//if birthYear, ok := tokens.YearFromTS(int(acc.BirthTS)); ok {
	//	doc.Fields[BirthYearField] = append(doc.Fields[BirthYearField], birthYear)
	//}
	if acc.PremiumFinish[0] > 0 {
		doc.Fields[PremiumField] = append(doc.Fields[PremiumField], tokens.NotNull)
	} else {
		doc.Fields[PremiumField] = append(doc.Fields[PremiumField], tokens.Null)
	}
	//if premiumNow(acc.PremiumFinish[:]) {
	//	doc.Fields[premiumField] = append(doc.Fields[premiumField], PremiumNowToken)
	//}
	if acc.Phone[0] > 0 {
		doc.Fields[PhoneCodeField] = append(doc.Fields[PhoneCodeField], tokens.NotNull)
	} else {
		doc.Fields[PhoneCodeField] = append(doc.Fields[PhoneCodeField], tokens.Null)
	}
	//if code, ok := phoneCode(acc.Phone[:]); ok {
	//	doc.Fields[phoneCodeField] = append(doc.Fields[phoneCodeField], phoneCodeToken(code))
	//}
	//if domain, ok := emailDomain(acc.Email.Buf[:acc.Email.Len]); ok {
	//	doc.Fields[emailDomainField] = append(doc.Fields[emailDomainField], emailDomainToken(domain))
	//}
}

/*
func iterFromField(field *inverted.Field, token int)

func (idx *CommonIndex) part() *inverted.Part {
	return idx.inv.Part(CommonPartition)
}

func (idx *CommonIndex) Sex(t int) *inverted.TokenIter {
	if field := idx.part().Field(SexField); field != nil {
		if token := field.Token(t); token != nil {
			return token.Iterator()
		}
	}
	return nil
}

func (idx *CommonIndex) Status(t int) *inverted.TokenIter {
	return idx.part().Field(StatusField).Iterator(t)
}

//func (idx *CommonIndex) EmailDomain(t int) *inverted.TokenIter {
//	return idx.part().Field(emailDomainField).Iterator(t)
//}

func (idx *CommonIndex) Fname(t int) *inverted.TokenIter {
	return idx.part().Field(FnameField).Iterator(t)
}

func (idx *CommonIndex) Sname(t int) *inverted.TokenIter {
	return idx.part().Field(SnameField).Iterator(t)
}

//func (idx *CommonIndex) PhoneCode(t int) *inverted.TokenIter {
//	return idx.part().Field(phoneCodeField).Iterator(t)
//}

func (idx *CommonIndex) Country(t int) *inverted.TokenIter {
	return idx.part().Field(CountryField).Iterator(t)
}

func (idx *CommonIndex) City(t int) *inverted.TokenIter {
	return idx.part().Field(CityField).Iterator(t)
}

//func (idx *CommonIndex) BirthYear(t int) *inverted.TokenIter {
//	return idx.part().Field(birthYearField).Iterator(t)
//}

func (idx *CommonIndex) Interest(t int) *inverted.TokenIter {
	return idx.part().Field(InterestField).Iterator(t)
}

//func (idx *CommonIndex) Premium(t int) *inverted.TokenIter {
//	return idx.part().Field(premiumField).Iterator(t)
//}

*/
