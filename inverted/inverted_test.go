package inverted

import (
	"testing"

	"gitlab.com/freepk/hlc18r4/backup"
)

func TestRebuild(t *testing.T) {
	rep, err := backup.Restore("../tmp/data/data.zip")
	if err != nil {
		t.Fatal(err)
	}
	interests := NewInvertedIndex(rep, DefaultParts, InterestsTokens)
	fnames := NewInvertedIndex(rep, DefaultParts, FnameTokens)
	snames := NewInvertedIndex(rep, DefaultParts, SnameTokens)
	countries := NewInvertedIndex(rep, DefaultParts, CountryTokens)
	cities := NewInvertedIndex(rep, DefaultParts, CityTokens)

	total, grow := 0, 0

	t.Log("Frist pass")

	total, grow = interests.Rebuild()
	t.Log("Interests", total, grow)
	total, grow = fnames.Rebuild()
	t.Log("Fnames", total, grow)
	total, grow = snames.Rebuild()
	t.Log("Snames", total, grow)
	total, grow = countries.Rebuild()
	t.Log("Countries", total, grow)
	total, grow = cities.Rebuild()
	t.Log("Cities", total, grow)

	t.Log("Second pass")

	total, grow = interests.Rebuild()
	t.Log("Interests", total, grow)
	total, grow = fnames.Rebuild()
	t.Log("Fnames", total, grow)
	total, grow = snames.Rebuild()
	t.Log("Snames", total, grow)
	total, grow = countries.Rebuild()
	t.Log("Countries", total, grow)
	total, grow = cities.Rebuild()
	t.Log("Cities", total, grow)
}
