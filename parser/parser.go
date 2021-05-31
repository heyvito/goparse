package parser

import (
	"context"
	"strings"
)

type AtomKind int

const (
	KindAlpha AtomKind = iota + 1
	KindBit
	KindChar
	KindCR
	KindLF
	KindCtl
	KindDigit
	KindDQuote
	KindHTab
	KindOctet
	KindSP
	KindVChar
	KindOption
	KindAtomList
	KindRefResult
)

var atomKindString = map[AtomKind]string{
	KindAlpha:     "KindAlpha",
	KindBit:       "KindBit",
	KindChar:      "KindChar",
	KindCR:        "KindCR",
	KindLF:        "KindLF",
	KindCtl:       "KindCtl",
	KindDigit:     "KindDigit",
	KindDQuote:    "KindDQuote",
	KindHTab:      "KindHTab",
	KindOctet:     "KindOctet",
	KindSP:        "KindSP",
	KindVChar:     "KindVChar",
	KindOption:    "KindOption",
	KindAtomList:  "KindAtomList",
	KindRefResult: "KindRefResult",
}

func (a AtomKind) String() string {
	return atomKindString[a]
}

type Cursor struct {
	buffer []rune
	bufLen int
	pos    int
}

func CursorFromString(data string) Cursor {
	return Cursor{
		buffer: []rune(data),
		bufLen: len(data),
		pos:    -1,
	}
}

func (c Cursor) dup() Cursor {
	return Cursor{
		buffer: c.buffer,
		bufLen: c.bufLen,
		pos:    c.pos,
	}
}

func (c *Cursor) Merge(other Cursor) {
	c.pos = other.pos
}

func (c Cursor) Peek() rune {
	if c.pos+1 >= c.bufLen {
		return 0x00
	}
	return c.buffer[c.pos+1]
}

func (c Cursor) TryPeek() (bool, rune) {
	if c.pos+1 >= c.bufLen {
		return false, 0x00
	}
	return true, c.buffer[c.pos+1]
}

func (c *Cursor) Consume() {
	c.pos++
}

type Consumer interface {
	TryConsume(ctx context.Context, c *Cursor) (Atom, error)
	String() string
	Name() string
	Weight() int
}

func Lit(lit rune) *LitConsumer                   { return &LitConsumer{lit: lit} }
func Cat(cons ...Consumer) *ConcatenationConsumer { return &ConcatenationConsumer{cons: cons} }
func Opt(con Consumer) *OptionalConsumer          { return &OptionalConsumer{con: con} }
func Alt(cons ...Consumer) *AlternationConsumer   { return &AlternationConsumer{cons: cons} }
func Ref(name string) *RefConsumer                { return &RefConsumer{name: strings.ToLower(name)} }
func HexRange(from, to rune) *HexRangeConsumer    { return &HexRangeConsumer{from: from, to: to} }
func B(con Consumer) *BlankConsumer               { return &BlankConsumer{con: con} }
func Str(val string) *ConcatenationConsumer {
	var cons []Consumer
	for _, v := range val {
		cons = append(cons, Lit(v))
	}
	return Cat(cons...)
}
func Dec(i int) *DecimalConsumer              { return &DecimalConsumer{i} }
func DecRange(from, to int) *DecRangeConsumer { return &DecRangeConsumer{from, to} }
func Star(con Consumer) *RepetitionConsumer {
	return &RepetitionConsumer{
		mode: RepeatStar,
		con:  con,
	}
}
func Plus(con Consumer) *RepetitionConsumer {
	return &RepetitionConsumer{
		mode: RepeatPlus,
		con:  con,
	}
}

func Repeat(min, max int, con Consumer) *RepetitionConsumer {
	if max != 0 {
		return &RepetitionConsumer{
			mode: RepeatMinMax,
			min:  min,
			max:  max,
			con:  con,
		}
	}
	return &RepetitionConsumer{
		mode: RepeatMin,
		min:  min,
		max:  max,
		con:  con,
	}
}

var (
	ALPHA  = &AlphaConsumer{}
	BIT    = &BitConsumer{}
	CHAR   = &CharConsumer{}
	CR     = &CRConsumer{}
	LF     = &LFConsumer{}
	CTL    = &CtlConsumer{}
	DIGIT  = &DigitConsumer{}
	DQUOTE = &DQuoteConsumer{}
	HTAB   = &HTabConsumer{}
	OCTET  = &OctetConsumer{}
	SP     = &SPConsumer{}
	VCHAR  = &VCharConsumer{}
	CRLF   = Cat(CR, LF)
	HEXDIG = Alt(DIGIT, Lit('A'), Lit('B'), Lit('C'), Lit('D'), Lit('E'), Lit('F'))
	WSP    = Alt(SP, HTAB)
	LWSP   = Star(Alt(WSP, Cat(CRLF, WSP)))
)

var CoreConsumers = map[string]Consumer{
	"alpha":  ALPHA,
	"bit":    BIT,
	"char":   CHAR,
	"lf":     LF,
	"cr":     CR,
	"crlf":   CRLF,
	"ctl":    CTL,
	"digit":  DIGIT,
	"dquote": DQUOTE,
	"htab":   HTAB,
	"octet":  OCTET,
	"sp":     SP,
	"vchar":  VCHAR,
	"hexdig": HEXDIG,
	"wsp":    WSP,
	"lwsp":   LWSP,
}

func MakeRules(m map[string]Consumer) map[string]Consumer {
	r := map[string]Consumer{}
	for k, v := range CoreConsumers {
		r[k] = v
	}
	for k, v := range m {
		r[k] = v
	}
	return r
}

const ruleMapKey = "__RULEMAP"

func KickoffParser(cur *Cursor, parser map[string]Consumer, initialRule string) (Atom, error) {
	startAt := Ref(initialRule)
	ctx := context.WithValue(context.Background(), ruleMapKey, parser)
	return startAt.TryConsume(ctx, cur)
}
