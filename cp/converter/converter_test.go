package converter

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yisaer/idl-parser/converter"
)

func TestGetDataType(t *testing.T) {
	c, err := NewArxmlCPConverter("../../test/s1_cp_test.xml", converter.IDlConverterConfig{IsLittleEndian: false, LengthFieldLength: 4, PaddingLength: 4})
	require.NoError(t, err)
	tr, err := c.GetDataTypeByID(33282, 2181169157)
	require.NoError(t, err)
	require.Equal(t, "string", tr.TypeName())
	testData := []byte{
		0x00, 0x00, 0x00, 0x08, // 长度字段 (8字节)
		0xEF, 0xBB, 0xBF, // UTF-8 BOM
		0x54, 0x65, 0x73, 0x74, // "Test"
		0x00, // UTF-8 终止符
	}
	v, err := c.Convert(33282, 2181169157, testData)
	require.NoError(t, err)
	require.Equal(t, "Test", v)
}
