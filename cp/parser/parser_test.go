package parser

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/yisaer/arxml-converter/ast"
)

func TestParser(t *testing.T) {
	p, err := NewParser("../../test/s1_cp_test.xml")
	require.NoError(t, err)
	require.NoError(t, p.Parse())
	transformer := ast.NewTransformHelper(p.dataTypesParser.GetApplicationDataTypes())
	m, err := transformer.TransformIntoModule()
	require.NoError(t, err)
	v, _ := json.Marshal(m)
	fmt.Println(string(v))
	fmt.Println(p.topologyParser.GetServiceIDMap())
	fmt.Println(p.topologyParser.GetHeaderRef())
	fmt.Println(p.topologyParser.GetPDUTriggeringRef())
	fmt.Println(p.communicationParser.GetPduRefMap())
	fmt.Println(p.systemParser.GetOperationRef())
	fmt.Println(p.softwareTypesParser.GetInterfaceRefMap())
}
