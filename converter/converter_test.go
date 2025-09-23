package converter

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yisaer/idl-parser/converter"
)

func TestParser(t *testing.T) {
	p, err := NewConverter("../test/example.xml", converter.IDlConverterConfig{IsLittleEndian: false, LengthFieldLength: 2, PaddingLength: 4})
	require.NoError(t, err)
	typ, err := p.GetTypeByID(22530, 36865)
	require.NoError(t, err)
	fmt.Println(typ)
}
