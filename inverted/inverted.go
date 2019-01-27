package inverted

type Document struct {
	ID     int
	Parts  []int
	Fields [][]int
}

type Iterator interface {
	Reset()
	Next() (*Document, bool)
}

type Ref uint32

type Token struct {
	refs  []Ref
	count int
}

func (t *Token) Reset() {
	t.refs = t.refs[:0]
	t.count = 0
}

func (t *Token) Len() int {
	return len(t.refs)
}

func (t *Token) Iterator() *TokenIter {
	return NewTokenIter(t.refs)
}

func (t *Token) Incr() {
	t.count++
}

func (t *Token) Count() int {
	return t.count
}

type Field struct {
	tokens []Token
}

func (f *Field) Reset() {
	for i := range f.tokens {
		f.tokens[i].Reset()
	}
	f.tokens = f.tokens[:0]
}

func (f *Field) token(token int) *Token {
	if grow := token + 1 - len(f.tokens); grow > 0 {
		f.tokens = append(f.tokens, make([]Token, grow)...)
	}
	return &f.tokens[token]
}

func (f *Field) Token(token int) *Token {
	if token < len(f.tokens) {
		return &f.tokens[token]
	}
	return nil
}

func (f *Field) Len() int {
	return len(f.tokens)
}

type Part struct {
	fields []Field
}

func (p *Part) Reset() {
	for i := range p.fields {
		p.fields[i].Reset()
	}
	p.fields = p.fields[:0]
}

func (p *Part) field(field int) *Field {
	if grow := field + 1 - len(p.fields); grow > 0 {
		p.fields = append(p.fields, make([]Field, grow)...)
	}
	return &p.fields[field]
}

func (p *Part) Field(field int) *Field {
	if field < len(p.fields) {
		return &p.fields[field]
	}
	return nil
}

func (p *Part) Len() int {
	return len(p.fields)
}

type Inverted struct {
	iter  Iterator
	parts []Part
}

func NewInverted(iter Iterator) *Inverted {
	return &Inverted{iter: iter}
}

func (inv *Inverted) Reset() {
	for i := range inv.parts {
		inv.parts[i].Reset()
	}
	inv.parts = inv.parts[:0]
}

func (inv *Inverted) part(part int) *Part {
	if grow := part + 1 - len(inv.parts); grow > 0 {
		inv.parts = append(inv.parts, make([]Part, grow)...)
	}
	return &inv.parts[part]
}

func (inv *Inverted) Prepare() {
	inv.Reset()
	inv.iter.Reset()
	for {
		doc, ok := inv.iter.Next()
		if !ok {
			break
		}
		for p := range doc.Parts {
			part := inv.part(doc.Parts[p])
			for f := range doc.Fields {
				field := part.field(f)
				for t := range doc.Fields[f] {
					token := field.token(doc.Fields[f][t])
					token.Incr()
				}
			}
		}
	}
}

func (inv *Inverted) Rebuild() {
	inv.Prepare()
	inv.iter.Reset()
	for {
		doc, ok := inv.iter.Next()
		if !ok {
			break
		}
		for p := range doc.Parts {
			part := &inv.parts[doc.Parts[p]]
			for f := range doc.Fields {
				field := &part.fields[f]
				for t := range doc.Fields[f] {
					token := &field.tokens[doc.Fields[f][t]]
					if token.count > cap(token.refs) {
						grow := token.count * 110 / 100
						token.refs = make([]Ref, 0, grow)
					}
					token.refs = append(token.refs, Ref(doc.ID))
				}
			}
		}
	}
}

func (inv *Inverted) Part(part int) *Part {
	if part < len(inv.parts) {
		return &inv.parts[part]
	}
	return nil
}

func (inv *Inverted) Len() int {
	return len(inv.parts)
}

func (inv *Inverted) Iterator(part, field, token int) *TokenIter {
	if part := inv.Part(part); part != nil {
		if field := part.Field(field); field != nil {
			if token := field.Token(token); token != nil {
				return token.Iterator()
			}
		}
	}
	return nil
}
