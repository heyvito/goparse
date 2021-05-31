package gen

import (
	"fmt"
	"strings"

	"github.com/heyvito/goparse/abnf"
	"github.com/heyvito/goparse/parser"
)

func Generate(list *abnf.RuleList) string {
	sb := strings.Builder{}
	sb.WriteString("var parser = map[string]p.Consumer{\n")
	for _, r := range list.Rules {
		sb.WriteRune('"')
		sb.WriteString(r.Name.Name)
		sb.WriteString(`": `)
		WriteElement(r.Elements, &sb)
		sb.WriteString(",\n")
	}
	sb.WriteRune('}')

	return sb.String()
}

func escapeLit(lit string) string {
	if lit == "'" {
		return "\\'"
	}
	if lit == "\\" {
		return "\\\\"
	}
	return lit
}

func escapeString(str string) string {
	return strings.ReplaceAll(str, "\"", "\\\"")
}

func WriteElement(element interface{}, sb *strings.Builder) {
	switch el := element.(type) {
	case abnf.Elements:
		WriteElement(el.Alternation, sb)
	case abnf.Alternation:
		if len(el.Elements) == 1 {
			WriteElement(el.Elements[0], sb)
			return
		}
		sb.WriteString("p.Alt(")
		for _, v := range el.Elements {
			WriteElement(v, sb)
			sb.WriteString(",")
		}
		sb.WriteString(")")
	case abnf.Concatenation:
		if len(el.Elements) == 1 {
			WriteElement(el.Elements[0], sb)
			return
		}
		sb.WriteString("p.Cat(")
		for _, v := range el.Elements {
			WriteElement(v, sb)
			sb.WriteString(",")
		}
		sb.WriteString(")")
	case abnf.Repetition:
		if el.Meta == nil {
			WriteElement(el.Element, sb)
			return
		}
		if el.Meta.Min == 0 && el.Meta.Max == 0 {
			sb.WriteString("p.Star(")
		} else if el.Meta.Min == 1 && el.Meta.Max == 0 {
			sb.WriteString("p.Plus(")
		} else {
			sb.WriteString(fmt.Sprintf("p.Repeat(%d, %d, ", el.Meta.Min, el.Meta.Max))
		}

		WriteElement(el.Element, sb)
		sb.WriteRune(')')
	case abnf.Group:
		WriteElement(el.Elements, sb)
	case abnf.Element:
		WriteElement(el.Inner, sb)
	case abnf.RuleName:
		if _, ok := parser.CoreConsumers[strings.ToLower(el.Name)]; ok {
			sb.WriteString("p.")
			sb.WriteString(strings.ToUpper(el.Name))
			return
		}
		sb.WriteString("p.Ref(\"")
		sb.WriteString(el.Name)
		sb.WriteString("\")")
	case abnf.Option:
		sb.WriteString("p.Opt(")
		WriteElement(el.Elements, sb)
		sb.WriteString(")")
	case abnf.CharVal:
		if len(el.Value) == 1 {
			sb.WriteString("p.Lit('")
			sb.WriteString(escapeLit(el.Value))
			sb.WriteString("')")
			return
		}
		sb.WriteString("p.Str(\"")
		sb.WriteString(escapeString(el.Value))
		sb.WriteString("\")")
	case abnf.HexVal:
		switch el.Mode {
		case abnf.NumericModeRange:
			sb.WriteString("p.HexRange(")
			sb.WriteString(fmt.Sprintf("0x%02x, 0x%02x", el.Range.From, el.Range.To))
			sb.WriteString(")")
		case abnf.NumericModeSingle:
			sb.WriteString("p.Hex(")
			sb.WriteString(fmt.Sprintf("0x%02x", el.Single))
			sb.WriteString(")")
		case abnf.NumericModeSequence:
			sb.WriteString("HexSeq(")
			for _, v := range el.Sequence {
				sb.WriteString(fmt.Sprintf("0x%2x,", v))
			}
			sb.WriteString(")")
		default:
			panic("Unimplemented")
		}

	case abnf.DecVal:
		switch el.Mode {
		case abnf.NumericModeRange:
			sb.WriteString("p.DecRange(")
			sb.WriteString(fmt.Sprintf("%d, %d", el.Range.From, el.Range.To))
			sb.WriteString(")")
		case abnf.NumericModeSingle:
			sb.WriteString("p.Dec(")
			sb.WriteString(fmt.Sprintf("%d", el.Single))
			sb.WriteString(")")
		case abnf.NumericModeSequence:
			sb.WriteString("p.DecSeq(")
			for _, v := range el.Sequence {
				sb.WriteString(fmt.Sprintf("%d,", v))
			}
			sb.WriteString(")")
		default:
			panic("Unimplemented")
		}
	default:
		fmt.Printf("Unimplemented: %T\n", el)
	}
}
