package search

import (
	"time"

	"github.com/freepk/hlc18r4/proto"
	"github.com/freepk/hlc18r4/tokens"
	"github.com/freepk/inverted"
	"github.com/freepk/parse"
)

func phoneCode(b []byte) ([]byte, bool) {
	if len(b) > 5 && b[1] == '(' && b[5] == ')' {
		return b[2:5], true
	}
	return nil, false
}

func emailDomain(b []byte) ([]byte, bool) {
	if domain, _ := parse.ScanSymbol(b, 0x40); len(domain) > 0 {
		return domain, true
	}
	return nil, false
}

func commonProc(doc *inverted.Document, acc *proto.Account) {
	doc.Parts = append(doc.Parts, CommonPart)
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
	birthYear := time.Unix(int64(acc.BirthTS), 0).UTC().Year()
	doc.Fields[BirthYearField] = append(doc.Fields[BirthYearField], birthYear-tokens.EpochYear)
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
	if code, ok := phoneCode(acc.Phone[:]); ok {
		doc.Fields[PhoneCodeField] = append(doc.Fields[PhoneCodeField], tokens.AddPhoneCode(code))
	}
	if domain, ok := emailDomain(acc.Email.Buf[:acc.Email.Len]); ok {
		doc.Fields[EmailDomainField] = append(doc.Fields[EmailDomainField], tokens.AddEmailDomain(domain))
	}
}