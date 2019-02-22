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
	if acc.Fname > 0 {
		doc.Fields[FnameField] = append(doc.Fields[FnameField], int(acc.Fname))
	}
	if acc.Sname > 0 {
		doc.Fields[SnameField] = append(doc.Fields[SnameField], int(acc.Sname))
	}
	if acc.Country > 0 {
		doc.Fields[CountryField] = append(doc.Fields[CountryField], int(acc.Country))
	}
	if acc.City > 0 {
		doc.Fields[CityField] = append(doc.Fields[CityField], int(acc.City))
	}
	for i := range acc.Interests {
		if acc.Interests[i] == 0 {
			break
		}
		doc.Fields[InterestField] = append(doc.Fields[InterestField], int(acc.Interests[i]))
	}
	birthYear := time.Unix(int64(acc.BirthTS), 0).UTC().Year()
	doc.Fields[BirthYearField] = append(doc.Fields[BirthYearField], birthYear-tokens.EpochYear)
	if code, ok := phoneCode(acc.Phone[:]); ok {
		doc.Fields[PhoneCodeField] = append(doc.Fields[PhoneCodeField], tokens.AddPhoneCode(code))
	}
	if domain, ok := emailDomain(acc.Email.Buf[:acc.Email.Len]); ok {
		doc.Fields[EmailDomainField] = append(doc.Fields[EmailDomainField], tokens.AddEmailDomain(domain))
	}
}
