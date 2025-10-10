package ast

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/yisaer/arxml-converter/mod"
)

func TestParser(t *testing.T) {
	p, err := NewParser("../../test/s1_cp_test.xml")
	require.NoError(t, err)
	require.NoError(t, p.Parse())
	transformer := mod.NewTransformHelper(p.dataTypesParser.applicationDataTypes)
	m, err := transformer.TransformIntoModule()
	require.NoError(t, err)
	v, _ := json.Marshal(m)
	fmt.Println(string(v))
}
