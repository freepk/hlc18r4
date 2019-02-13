package search

import (
	"github.com/freepk/hlc18r4/proto"
	"github.com/freepk/hlc18r4/tokens"
	"github.com/freepk/inverted"
)

func countryProc(doc *inverted.Document, acc *proto.Account) {
	if acc.Country > 0 {
		doc.Parts = append(doc.Parts, tokens.NotNull, int(acc.Country))
	} else {
		doc.Parts = append(doc.Parts, tokens.Null)
	}
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
	for i := range acc.Interests {
		if acc.Interests[i] == 0 {
			break
		}
		doc.Fields[InterestField] = append(doc.Fields[InterestField], int(acc.Interests[i]))
	}
}
