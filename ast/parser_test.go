package ast

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	p, err := NewParser("example.xml")
	require.NoError(t, err)
	require.NoError(t, p.Parse())
}
