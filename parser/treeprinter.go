package parser

import (
	"fmt"
	"strings"
)

func PrintTree(root Atom) string {
	return printTreeFn(root, 0)
}

func printTreeFn(root Atom, level int) string {
	var str strings.Builder
	str.WriteString(strings.Repeat(" ", 2*level))
	switch v := root.(type) {
	case Alpha:
		str.WriteString("- A: ")
		str.WriteString(v.value)
	case Bit:
		str.WriteString("- B: ")
		str.WriteString(v.value)
	case Char:
		str.WriteString("- C: ")
		str.WriteString(v.value)
	case CRVal:
		str.WriteString("- CR")
	case LFVal:
		str.WriteString("- LF")
	case Ctl:
		str.WriteString("- T: ")
		str.WriteString(fmt.Sprintf("0x%02x", v.value[0]))
	case Digit:
		str.WriteString("- D: ")
		str.WriteString(v.value)
	case DQuote:
		str.WriteString("- DQuote")
	case HTab:
		str.WriteString("- HTab")
	case Octet:
		str.WriteString("- O: ")
		str.WriteString(fmt.Sprintf("0x%02x", v.value))
	case SPVal:
		str.WriteString("- SP")
	case VChar:
		str.WriteString("- VChar: ")
		str.WriteString(fmt.Sprintf("%q", v.value))
	case OptionVal:
		if !v.Valid {
			str.WriteString("- Empty Opt")
		} else {
			str.WriteString("- Opt: \n")
			str.WriteString(printTreeFn(v.value, level+1))
		}
	case AtomList:
		if len(v.value) == 0 {
			str.WriteString("- Empty List")
		} else {
			str.WriteString("- List: \n")
			for _, v := range v.value {
				str.WriteString(printTreeFn(v, level+1))
			}
		}
	case RefResult:
		str.WriteString(fmt.Sprintf("- Rule %s:\n", v.Name))
		str.WriteString(printTreeFn(v.value, level+1))
	}

	str.WriteString("\n")
	return str.String()
}
