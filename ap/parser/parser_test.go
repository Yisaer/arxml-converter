package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	p, err := NewParser("../../test/s1_ap_test.xml")
	require.NoError(t, err)
	require.NoError(t, p.Parse())
}
