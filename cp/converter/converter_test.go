package converter

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yisaer/idl-parser/converter"
)

func TestGetDataType(t *testing.T) {
	c, err := NewArxmlCPConverter("../../test/s1_cp_test.xml", converter.IDlConverterConfig{IsLittleEndian: false, LengthFieldLength: 4, PaddingLength: 4})
	require.NoError(t, err)
	tr, err := c.GetDataTypeByID(33282, 2181169157)
	require.NoError(t, err)
	fmt.Println(tr.TypeName())
}
