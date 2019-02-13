package search

import (
	"github.com/freepk/hlc18r4/proto"
	"github.com/freepk/hlc18r4/tokens"
	"github.com/freepk/inverted"
)

func sexProc(doc *inverted.Document, acc *proto.Account) {
	switch acc.Sex {
	case tokens.MaleSex:
		doc.Parts = append(doc.Parts, tokens.MaleSex)
	case tokens.FemaleSex:
		doc.Parts = append(doc.Parts, tokens.FemaleSex)
	}
	switch acc.Status {
	case tokens.SingleStatus:
		doc.Fields[StatusField] = append(doc.Fields[StatusField], tokens.SingleStatus, tokens.NotInRelStatus, tokens.NotComplStatus)
	case tokens.InRelStatus:
		doc.Fields[StatusField] = append(doc.Fields[StatusField], tokens.InRelStatus, tokens.NotSingleStatus, tokens.NotComplStatus)
	case tokens.ComplStatus:
		doc.Fields[StatusField] = append(doc.Fields[StatusField], tokens.ComplStatus, tokens.NotSingleStatus, tokens.NotInRelStatus)
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
}
