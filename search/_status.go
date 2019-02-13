package search

import (
	"github.com/freepk/hlc18r4/proto"
	"github.com/freepk/hlc18r4/tokens"
	"github.com/freepk/inverted"
)

func statusProc(doc *inverted.Document, acc *proto.Account) {
	switch acc.Status {
	case tokens.SingleStatus:
		doc.Parts = append(doc.Parts, tokens.SingleStatus, tokens.NotInRelStatus, tokens.NotComplStatus)
	case tokens.InRelStatus:
		doc.Parts = append(doc.Parts, tokens.InRelStatus, tokens.NotSingleStatus, tokens.NotComplStatus)
	case tokens.ComplStatus:
		doc.Parts = append(doc.Parts, tokens.ComplStatus, tokens.NotSingleStatus, tokens.NotInRelStatus)
	}
	switch acc.Sex {
	case tokens.MaleSex:
		doc.Fields[SexField] = append(doc.Fields[SexField], tokens.MaleSex)
	case tokens.FemaleSex:
		doc.Fields[SexField] = append(doc.Fields[SexField], tokens.FemaleSex)
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
