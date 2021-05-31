package parser

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

const parentContextKey = "__PARENT"

func SetParent(parent Atom, ctx context.Context) context.Context {
	if reflect.TypeOf(parent).Kind() != reflect.Ptr {
		panic("BUG: SetParent requires a pointer")
	}
	return context.WithValue(ctx, parentContextKey, parent)
}
func GetParent(ctx context.Context) Atom {
	if v, ok := ctx.Value(parentContextKey).(Atom); ok {
		return v
	}
	return nil
}

type AlphaConsumer struct{}

func (AlphaConsumer) Name() string     { return "ALPHA" }
func (a AlphaConsumer) String() string { return a.Name() }
func (AlphaConsumer) TryConsume(ctx context.Context, c *Cursor) (Atom, error) {
	v := c.Peek()
	if v != 0x00 && (v >= 0x41 && v <= 0x5A) || (v >= 0x61 && v <= 0x7A) {
		c.Consume()
		return Alpha{value: string(v), parent: GetParent(ctx)}, nil
	}
	return nil, Error(c, "Expected alpha character between a-z or A-Z. Found %q", v)
}
func (AlphaConsumer) Weight() int { return 0 }

type BitConsumer struct{}

func (BitConsumer) Name() string     { return "BIT" }
func (b BitConsumer) String() string { return b.Name() }
func (BitConsumer) Weight() int      { return 0 }
func (BitConsumer) TryConsume(ctx context.Context, c *Cursor) (Atom, error) {
	v := c.Peek()
	if v == '0' || v == '1' {
		c.Consume()
		return Bit{value: string(v), parent: GetParent(ctx)}, nil
	}
	return nil, Error(c, "Expected a bit (0-1). Found %q", v)
}

type CharConsumer struct{}

func (CharConsumer) Name() string     { return "CHAR" }
func (c CharConsumer) String() string { return c.Name() }
func (CharConsumer) Weight() int      { return 0 }
func (CharConsumer) TryConsume(ctx context.Context, c *Cursor) (Atom, error) {
	v := c.Peek()
	if v == 0x01 || v >= 0x7F {
		c.Consume()
		return Char{value: string(v), parent: GetParent(ctx)}, nil
	}
	return nil, Error(c, "Expected a value equal to 0x01 or greater than 0x7E. Found %q (0x%02x)", v, v)
}

type LFConsumer struct{}

func (LFConsumer) Name() string     { return "LF" }
func (l LFConsumer) String() string { return l.Name() }
func (LFConsumer) Weight() int      { return 0 }
func (LFConsumer) TryConsume(ctx context.Context, c *Cursor) (Atom, error) {
	v := c.Peek()
	if v == 0x0A {
		c.Consume()
		return LFVal{parent: GetParent(ctx)}, nil
	}
	return nil, Error(c, "Expected a linefeed. Found %q", v)
}

type CRConsumer struct{}

func (CRConsumer) Name() string     { return "CR" }
func (c CRConsumer) String() string { return c.Name() }
func (CRConsumer) Weight() int      { return 0 }
func (CRConsumer) TryConsume(ctx context.Context, c *Cursor) (Atom, error) {
	v := c.Peek()
	if v == 0x0D {
		c.Consume()
		return CRVal{parent: GetParent(ctx)}, nil
	}
	return nil, Error(c, "expected a carriage return. Found %q", v)
}

type CtlConsumer struct{}

func (CtlConsumer) Name() string     { return "CTL" }
func (c CtlConsumer) String() string { return c.Name() }
func (CtlConsumer) Weight() int      { return 0 }
func (CtlConsumer) TryConsume(ctx context.Context, c *Cursor) (Atom, error) {
	v := c.Peek()
	if v <= 0x1f || v == 0x7f {
		c.Consume()
		return Ctl{value: string(v), parent: GetParent(ctx)}, nil
	}
	return nil, Error(c, "Expected a control character (0x7F, or <= 0x1F). Found %q (0x%2x)", v, v)
}

type DigitConsumer struct{}

func (DigitConsumer) Name() string     { return "DIGIT" }
func (d DigitConsumer) String() string { return d.Name() }
func (DigitConsumer) Weight() int      { return 0 }
func (DigitConsumer) TryConsume(ctx context.Context, c *Cursor) (Atom, error) {
	v := c.Peek()
	if v >= 0x30 && v <= 0x39 {
		c.Consume()
		return Digit{value: string(v), parent: GetParent(ctx)}, nil
	}
	return nil, Error(c, "Expected a digit (0-9). Found %q", v)
}

type DQuoteConsumer struct{}

func (DQuoteConsumer) Name() string     { return "DQUOTE" }
func (d DQuoteConsumer) String() string { return d.Name() }
func (DQuoteConsumer) Weight() int      { return 0 }
func (DQuoteConsumer) TryConsume(ctx context.Context, c *Cursor) (Atom, error) {
	v := c.Peek()
	if v == 0x22 {
		c.Consume()
		return DQuote{parent: GetParent(ctx)}, nil
	}
	return nil, Error(c, "Expected a double-quote. Found %q", v)
}

type HTabConsumer struct{}

func (HTabConsumer) Name() string     { return "HTAB" }
func (h HTabConsumer) String() string { return h.Name() }
func (HTabConsumer) Weight() int      { return 0 }
func (HTabConsumer) TryConsume(ctx context.Context, c *Cursor) (Atom, error) {
	v := c.Peek()
	if v == 0x09 {
		c.Consume()
		return HTab{parent: GetParent(ctx)}, nil
	}
	return nil, Error(c, "Expected a horizontal tab. Found %q", v)
}

type OctetConsumer struct{}

func (OctetConsumer) Name() string     { return "OCTET" }
func (o OctetConsumer) String() string { return o.Name() }
func (OctetConsumer) Weight() int      { return 0 }
func (OctetConsumer) TryConsume(ctx context.Context, c *Cursor) (Atom, error) {
	if ok, v := c.TryPeek(); ok {
		c.Consume()
		return Octet{value: v, parent: GetParent(ctx)}, nil
	}
	return nil, Error(c, "Expected octet, found EOF")
}

type SPConsumer struct{}

func (SPConsumer) Name() string     { return "SP" }
func (s SPConsumer) String() string { return s.Name() }
func (SPConsumer) Weight() int      { return 0 }
func (SPConsumer) TryConsume(ctx context.Context, c *Cursor) (Atom, error) {
	v := c.Peek()
	if v == 0x20 {
		c.Consume()
		return SPVal{parent: GetParent(ctx)}, nil
	}
	return nil, Error(c, "Expected a space, found %q", v)
}

type VCharConsumer struct{}

func (VCharConsumer) Name() string     { return "VCHAR" }
func (v VCharConsumer) String() string { return v.Name() }
func (VCharConsumer) Weight() int      { return 0 }
func (VCharConsumer) TryConsume(ctx context.Context, c *Cursor) (Atom, error) {
	v := c.Peek()
	if v >= 0x21 && v <= 0x7E {
		c.Consume()
		return VChar{value: string(v), parent: GetParent(ctx)}, nil
	}
	return nil, Error(c, "Expected a visible character, found %q (0x%02x) instead", v, v)
}

type ConcatenationConsumer struct {
	cons []Consumer
}

func (c ConcatenationConsumer) Name() string { return "CONCATENATION" }
func (ConcatenationConsumer) Weight() int    { return 1 }
func (c ConcatenationConsumer) String() string {
	str := make([]string, len(c.cons))
	for i, c := range c.cons {
		str[i] = c.String()
	}
	return fmt.Sprintf("( %s )", strings.Join(str, " "))
}
func (c ConcatenationConsumer) TryConsume(ctx context.Context, cur *Cursor) (Atom, error) {
	var results []Atom
	ret := AtomList{parent: GetParent(ctx)}
	cd := cur.dup()
	for _, v := range c.cons {
		if res, err := v.TryConsume(SetParent(&ret, ctx), &cd); err == nil {
			if res != nil {
				results = append(results, res)
			}
		} else {
			return nil, err
		}
	}

	ret.value = results
	cur.Merge(cd)
	return ret, nil
}

type OptionalConsumer struct {
	con Consumer
}

func (o OptionalConsumer) Name() string   { return fmt.Sprintf("OPT(%s)", o.con.String()) }
func (o OptionalConsumer) String() string { return fmt.Sprintf("[ %s ]", o.con.String()) }
func (OptionalConsumer) Weight() int      { return 0 }
func (o OptionalConsumer) TryConsume(ctx context.Context, c *Cursor) (Atom, error) {
	cd := c.dup()
	ret := OptionVal{parent: GetParent(ctx)}

	if v, err := o.con.TryConsume(SetParent(&ret, ctx), &cd); err == nil {
		c.Merge(cd)
		ret.Valid = true
		ret.value = v
	}
	return ret, nil
}

type LitConsumer struct {
	lit rune
}

func (l LitConsumer) Name() string   { return fmt.Sprintf("LIT(%c)", l.lit) }
func (l LitConsumer) String() string { return fmt.Sprintf("'%c'", l.lit) }
func (LitConsumer) Weight() int      { return 0 }
func (l LitConsumer) TryConsume(ctx context.Context, c *Cursor) (Atom, error) {
	ok, v := c.TryPeek()
	if !ok {
		return nil, Error(c, "Expected a literal %q, found EOF", l.lit)
	}

	if v == l.lit {
		c.Consume()
		return Char{value: string(v), parent: GetParent(ctx)}, nil
	}
	return nil, Error(c, "Expected a literal %q, found %q instead", l.lit, v)
}

type HexRangeConsumer struct {
	from rune
	to   rune
}

func (HexRangeConsumer) Weight() int      { return 0 }
func (h HexRangeConsumer) Name() string   { return fmt.Sprintf("HEXRANGE(0x%2x, 0x%2x)", h.from, h.to) }
func (h HexRangeConsumer) String() string { return fmt.Sprintf("%%x%2x-%2x", h.from, h.to) }
func (h HexRangeConsumer) TryConsume(ctx context.Context, c *Cursor) (Atom, error) {
	ok, v := c.TryPeek()
	if !ok {
		return nil, Error(c, "Expected a hexadecimal within range 0x%2x >= x <= 0x%2x, but found EOF", h.from, h.to)
	}
	if v >= h.from && v <= h.to {
		c.Consume()
		return Char{value: string(v), parent: GetParent(ctx)}, nil
	}
	return nil, Error(c, "expected a hexadecimal within range 0x%2x >= x <= 0x%2x, but found %q (0x%2x) instead", h.from, h.to, v, v)
}

// FIXME: Is this even working?!

type BlankConsumer struct {
	con Consumer
}

func (b BlankConsumer) TryConsume(ctx context.Context, c *Cursor) (Atom, error) {
	_, err := b.con.TryConsume(ctx, c)
	return nil, err
}
func (b BlankConsumer) String() string { return fmt.Sprintf("<%s>", b.con.String()) }
func (b BlankConsumer) Name() string   { return fmt.Sprintf("Blank(%s)", b.con.Name()) }
func (b BlankConsumer) Weight() int    { return b.con.Weight() }

type DecimalConsumer struct{ v int }

func (d DecimalConsumer) TryConsume(ctx context.Context, c *Cursor) (Atom, error) {
	ok, v := c.TryPeek()
	if !ok {
		return nil, Error(c, "Expected a decimal %d, found EOF", d.v)
	}

	if int(v) == d.v {
		c.Consume()
		return Char{value: string(v), parent: GetParent(ctx)}, nil
	}
	return nil, Error(c, "Expected a decimal %d, found %d instead", d.v, int(v))
}
func (d DecimalConsumer) String() string { return fmt.Sprintf("Dec(%d)", d.v) }
func (d DecimalConsumer) Name() string   { return "DEC" }
func (d DecimalConsumer) Weight() int    { return 0 }

type DecRangeConsumer struct {
	from int
	to   int
}

func (DecRangeConsumer) Weight() int      { return 0 }
func (d DecRangeConsumer) Name() string   { return fmt.Sprintf("DECRANGE(%d, %d)", d.from, d.to) }
func (d DecRangeConsumer) String() string { return fmt.Sprintf("%%d%d-%d", d.from, d.to) }
func (d DecRangeConsumer) TryConsume(ctx context.Context, c *Cursor) (Atom, error) {
	ok, v := c.TryPeek()
	if !ok {
		return nil, Error(c, "Expected a decimal within range %d >= x <= %d, but found EOF", d.from, d.to)
	}
	if int(v) >= d.from && int(v) <= d.to {
		c.Consume()
		return Char{value: string(v), parent: GetParent(ctx)}, nil
	}
	return nil, Error(c, "expected a decimal within range %d >= x <= %d, but found %q (%d) instead", d.from, d.to, v, int(v))
}
