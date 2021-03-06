package test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/heyvito/goparse/abnf"
	"github.com/heyvito/goparse/abnf1"
)

func TestParseProgressive(t *testing.T) {
	data, err := os.ReadFile("../grammars/abnf.abnf")
	require.NoError(t, err)
	n := time.Now()
	v, err := abnf1.ParseABNFNaive(string(data))
	diff := time.Since(n)
	fmt.Printf("Parse took %s\n", diff.String())
	require.NoError(t, err)
	fmt.Println(abnf.Generate(v))
}
