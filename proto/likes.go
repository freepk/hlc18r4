package proto

import (
	"github.com/freepk/parse"
)

type LikeEx struct {
	Liker uint32
	Likee uint32
	TS    uint32
}

type Likes struct {
	Likes []LikeEx
}

func (l *Likes) reset() {
	l.Likes = l.Likes[:0]
}

func (l *Likes) UnmarshalJSON(buf []byte) ([]byte, bool) {
	var tail []byte
	var ok bool

	l.reset()

	if tail, ok = parse.SkipSymbol(buf, '{'); !ok {
		return buf, false
	}
	for {
		tail = parse.SkipSpaces(tail)
		switch {
		case len(tail) > 8 && string(tail[:8]) == `"likes":`:
			if tail, ok = parse.SkipSymbol(tail[8:], '['); !ok {
				return buf, false
			}
			for {
				liker := 0
				likee := 0
				ts := 0
				if tail, ok = parse.SkipSymbol(tail, '{'); !ok {
					break
				}
				for {
					tail = parse.SkipSpaces(tail)
					switch {
					case len(tail) > 8 && string(tail[:8]) == `"liker":`:
						if tail, liker, ok = parse.ParseInt(tail[8:]); !ok {
							return buf, false
						}
					case len(tail) > 8 && string(tail[:8]) == `"likee":`:
						if tail, likee, ok = parse.ParseInt(tail[8:]); !ok {
							return buf, false
						}
					case len(tail) > 5 && string(tail[:5]) == `"ts":`:
						if tail, ts, ok = parse.ParseInt(tail[5:]); !ok {
							return buf, false
						}
					}
					if tail, ok = parse.SkipSymbol(tail, ','); !ok {
						break
					}
				}
				if tail, ok = parse.SkipSymbol(tail, '}'); !ok {
					return buf, false
				}
				l.Likes = append(l.Likes, LikeEx{Liker: uint32(liker), Likee: uint32(likee), TS: uint32(ts)})
				if tail, ok = parse.SkipSymbol(tail, ','); !ok {
					break
				}
			}
			if tail, ok = parse.SkipSymbol(tail, ']'); !ok {
				return buf, false
			}
		}
		if tail, ok = parse.SkipSymbol(tail, ','); !ok {
			break
		}
	}
	if tail, ok = parse.SkipSymbol(tail, '}'); !ok {
		return buf, false
	}
	return tail, true
}
