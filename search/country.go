package search

import (
	"github.com/freepk/hlc18r4/proto"
	//"github.com/freepk/hlc18r4/tokens"
	"github.com/freepk/inverted"
)

func countryProc(doc *inverted.Document, acc *proto.Account) {
	if acc.Country > 0 {
		doc.Parts = append(doc.Parts, int(acc.Country))
	}
	for i := range acc.Interests {
		if acc.Interests[i] == 0 {
			break
		}
		doc.Fields[InterestField] = append(doc.Fields[InterestField], int(acc.Interests[i]))
	}
}
