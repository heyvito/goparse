package parser

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestALPHA(t *testing.T) {
	c := CursorFromString("A")
	v, err := ALPHA.TryConsume(context.Background(), &c)
	require.NoError(t, err)
	fmt.Println(v)
}

func TestAlt(t *testing.T) {
	c := CursorFromString("A1")
	v, err := Alt(ALPHA, DIGIT).TryConsume(context.Background(), &c)
	require.NoError(t, err)
	fmt.Println(v)
}

func TestStar(t *testing.T) {
	c := CursorFromString("ABC")
	v, err := Star(ALPHA).TryConsume(context.Background(), &c)
	require.NoError(t, err)
	fmt.Println(v)
}
