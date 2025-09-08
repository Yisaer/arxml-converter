package ast

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	p, err := NewParser("example.xml")
	require.NoError(t, err)
	require.NoError(t, p.Parse())
	require.Len(t, p.ArPackageElements, 4)
	allTypes := make(map[string]struct{})
	for _, dt := range p.dtList {
		if dt.TypReference != nil {
			allTypes[dt.TypReference.Ref] = struct{}{}
		}
	}
	fmt.Println(allTypes)
}
