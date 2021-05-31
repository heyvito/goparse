package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type Atom interface {
	Parent() Atom
	Value() interface{}
	Kind() AtomKind
}

type Alpha struct {
	value  string
	parent Atom
}

func (a Alpha) Parent() Atom       { return a.parent }
func (a Alpha) Value() interface{} { return a.value }
func (a Alpha) Kind() AtomKind     { return KindAlpha }

type Bit struct {
	value  string
	parent Atom
}

func (b Bit) Parent() Atom       { return b.parent }
func (b Bit) Value() interface{} { return b.value }
func (b Bit) Kind() AtomKind     { return KindBit }

type Char struct {
	value  string
	parent Atom
}

func (c Char) Parent() Atom       { return c.parent }
func (c Char) Value() interface{} { return c.value }
func (c Char) Kind() AtomKind     { return KindChar }

type CRVal struct{ parent Atom }

func (c CRVal) Parent() Atom       { return c.parent }
func (c CRVal) Value() interface{} { return "\r" }
func (c CRVal) Kind() AtomKind     { return KindCR }

type LFVal struct{ parent Atom }

func (l LFVal) Parent() Atom       { return l.parent }
func (l LFVal) Value() interface{} { return "\n" }
func (l LFVal) Kind() AtomKind     { return KindLF }

type Ctl struct {
	parent Atom
	value  string
}

func (c Ctl) Parent() Atom       { return c.parent }
func (c Ctl) Value() interface{} { return c.value }
func (c Ctl) Kind() AtomKind     { return KindCtl }

type Digit struct {
	parent Atom
	value  string
}

func (d Digit) Parent() Atom       { return d.parent }
func (d Digit) Value() interface{} { return d.value }
func (d Digit) Kind() AtomKind     { return KindDigit }

type DQuote struct{ parent Atom }

func (d DQuote) Parent() Atom       { return d.parent }
func (d DQuote) Value() interface{} { return "\"" }
func (d DQuote) Kind() AtomKind     { return KindDQuote }

type HTab struct{ parent Atom }

func (h HTab) Parent() Atom       { return h.parent }
func (h HTab) Value() interface{} { return "\t" }
func (h HTab) Kind() AtomKind     { return KindHTab }

type Octet struct {
	parent Atom
	value  rune
}

func (o Octet) Parent() Atom       { return o.parent }
func (o Octet) Value() interface{} { return o.value }
func (o Octet) Kind() AtomKind     { return KindOctet }

type SPVal struct{ parent Atom }

func (s SPVal) Parent() Atom       { return s.parent }
func (s SPVal) Value() interface{} { return " " }
func (s SPVal) Kind() AtomKind     { return KindSP }

type VChar struct {
	parent Atom
	value  string
}

func (v VChar) Parent() Atom       { return v.parent }
func (v VChar) Value() interface{} { return v.value }
func (v VChar) Kind() AtomKind     { return KindVChar }

type OptionVal struct {
	parent Atom
	Valid  bool
	value  Atom
}

func (o OptionVal) Parent() Atom       { return o.parent }
func (o OptionVal) Value() interface{} { return o.value }
func (o OptionVal) Kind() AtomKind     { return KindOption }

type AtomList struct {
	parent Atom
	value  []Atom
}

func (a AtomList) Len() int     { return len(a.value) }
func (a AtomList) Parent() Atom { return a.parent }
func (a AtomList) Value() interface{} {
	if a.value == nil {
		return nil
	}
	return a.value
}
func (a AtomList) Kind() AtomKind { return KindAtomList }
func (a AtomList) AllTerminals() bool {
	for _, v := range a.value {
		switch i := v.(type) {
		case AtomList:
			if !i.AllTerminals() {
				return false
			}
		case RefResult:
			return false
		case OptionVal:
			return false
		}
	}
	return true
}
func (a AtomList) ReduceAsString() string {
	str := strings.Builder{}
	for _, v := range a.value {
		switch inst := v.(type) {
		case Alpha:
			str.WriteString(inst.value)
		case Char:
			str.WriteString(inst.value)
		case Digit:
			str.WriteString(inst.value)
		case DQuote:
			str.WriteString("\"")
		case SPVal:
			str.WriteString(" ")
		case LFVal:
			str.WriteString("\n")
		case VChar:
			str.WriteString(inst.value)
		case AtomList:
			str.WriteString(inst.ReduceAsString())
		default:
			panic(fmt.Sprintf("Cannot reduce %T as string", v))
		}
	}
	return str.String()
}
func (a AtomList) ReduceAsInt() (bool, int) {
	str := a.ReduceAsString()
	if str == "" {
		return false, 0
	}
	i, err := strconv.Atoi(str)
	if err != nil {
		panic(fmt.Sprintf("ReduceAsInt failed: %s", err))
	}
	return true, i
}
func (a AtomList) ReduceAsHex() (bool, int) {
	str := a.ReduceAsString()
	if str == "" {
		return false, 0
	}
	i, err := strconv.ParseUint(str, 16, 8)
	if err != nil {
		panic(fmt.Sprintf("ReduceAsHex failed: %s", err))
	}
	return true, int(i)
}
func (a AtomList) First() Atom {
	return a.Nth(0)
}
func (a AtomList) Nth(i int) Atom {
	return a.value[i]
}

type RefResult struct {
	Name   string
	value  Atom
	parent Atom
}

func (r RefResult) Parent() Atom       { return r.parent }
func (r RefResult) Value() interface{} { return r.value }
func (r RefResult) Kind() AtomKind     { return KindRefResult }
