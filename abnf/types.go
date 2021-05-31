package abnf

type NodeKind int

const (
	NodeKindRulelist NodeKind = iota + 1
	NodeKindRule
	NodeKindRulename
	NodeKindDefinedAs
	NodeKindElements
	NodeKindCWsp
	NodeKindCNl
	NodeKindComment
	NodeKindAlternation
	NodeKindConcatenation
	NodeKindRepetition
	NodeKindRepeat
	NodeKindElement
	NodeKindGroup
	NodeKindOption
	NodeKindCharVal
	NodeKindBinVal
	NodeKindDecVal
	NodeKindHexVal
	NodeKindProseVal
)

type Node interface {
	Kind() NodeKind
}

type NumericMode int

const (
	NumericModeSingle NumericMode = iota
	NumericModeSequence
	NumericModeRange
)

type Range struct {
	From, To int
}

type ProseVal struct {
	Value string
}

type Numeric struct {
	Mode     NumericMode
	Single   int
	Sequence []int
	Range    Range
}

type BinVal struct{ Numeric }

func (b BinVal) Kind() NodeKind {
	return NodeKindBinVal
}

type DecVal struct{ Numeric }

func (d DecVal) Kind() NodeKind {
	return NodeKindDecVal
}

type HexVal struct{ Numeric }

func (h HexVal) Kind() NodeKind {
	return NodeKindHexVal
}

type CharVal struct {
	Value string
}

func (c CharVal) Kind() NodeKind {
	return NodeKindCharVal
}

type Option struct {
	Elements Alternation
}

func (o Option) Kind() NodeKind {
	return NodeKindOption
}

type Group struct {
	Elements Alternation
}

func (g Group) Kind() NodeKind {
	return NodeKindGroup
}

type Element struct {
	Inner Node
}

type Repeat struct {
	Min int
	Max int
}

type Repetition struct {
	Meta    *Repeat
	Element Element
}

type Concatenation struct {
	Elements []Repetition
}

type Alternation struct {
	Elements []Concatenation
}

type Comment struct {
	Value string
}

type CNl struct {
	Comment *Comment
}

type CWsp struct {
	CNl CNl
}

type Elements struct {
	Alternation Alternation
}

type DefinedAs struct {
	Value string
}

type RuleName struct {
	Name string
}

func (r RuleName) Kind() NodeKind {
	return NodeKindRulename
}

type Rule struct {
	Name      RuleName
	DefinedAs DefinedAs
	Elements  Elements
}

type RuleList struct {
	Rules []Rule
}
