package converter

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yisaer/idl-parser/converter"

	"github.com/yisaer/arxml-converter/util"
)

func TestGetDataType(t *testing.T) {
	c, err := NewArxmlCPConverter("../../test/s1_cp_test.xml", converter.IDlConverterConfig{IsLittleEndian: false, LengthFieldLength: 4, PaddingLength: 4})
	require.NoError(t, err)
	key, tr, err := c.GetDataTypeByID(33282, 2181169157)
	require.NoError(t, err)
	require.Equal(t, "string", tr.TypeName())
	require.Equal(t, "adt_WiFiApName", key)
	testData := []byte{
		0x00, 0x00, 0x00, 0x08, // 长度字段 (8字节)
		0xEF, 0xBB, 0xBF, // UTF-8 BOM
		0x54, 0x65, 0x73, 0x74, // "Test"
		0x00, // UTF-8 终止符
	}
	key, v, err := c.Convert(33282, 2181169157, testData)
	require.NoError(t, err)
	require.Equal(t, "Test", v)
}

func TestGetDataType2(t *testing.T) {
	//t.Skip()
	svcID, headerID, err := util.MergeHexUint16ToUint32("0xab02", "0x8015")
	require.NoError(t, err)
	c, err := NewArxmlCPConverter("../../tmp/example2.xml", converter.IDlConverterConfig{IsLittleEndian: false, LengthFieldLength: 4, PaddingLength: 4})
	require.NoError(t, err)
	key, tr, err := c.GetDataTypeByID(svcID, headerID)
	require.NoError(t, err)
	fmt.Println(key)
	fmt.Println(tr.TypeName())
	data, err := hex.DecodeString("000000000000000202000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)
	k, _, err := c.Convert(svcID, headerID, data)
	require.NoError(t, err)
	fmt.Println(k)
}

func TestGetDataType3(t *testing.T) {
	svcID, headerID, err := util.MergeHexUint16ToUint32("0xab04", "0x8009")
	require.NoError(t, err)
	fmt.Println(svcID)
	fmt.Println(headerID)
	c, err := NewArxmlCPConverter("../../tmp/example2.xml", converter.IDlConverterConfig{IsLittleEndian: false, LengthFieldLength: 4, PaddingLength: 4})
	require.NoError(t, err)
	key, tr, err := c.GetDataTypeByID(svcID, headerID)
	require.NoError(t, err)
	fmt.Println(key)
	fmt.Println(tr.TypeName())
}
