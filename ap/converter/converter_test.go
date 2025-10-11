package converter

import (
	"encoding/binary"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yisaer/idl-parser/converter"
)

func TestParser(t *testing.T) {
	p, err := NewConverter("../../test/example.xml", converter.IDlConverterConfig{IsLittleEndian: false, LengthFieldLength: 2, PaddingLength: 4})
	require.NoError(t, err)
	_, result, err := p.DecodeWithID(22530, 36865, testdata())
	require.NoError(t, err)
	require.Equal(t, expectData(), result)
}

func testdata() []byte {
	uint32s := make([]uint32, 60)
	for i := range uint32s {
		uint32s[i] = 1
	}
	data := make([]byte, 4*len(uint32s))
	for i, v := range uint32s {
		start := i * 4
		binary.BigEndian.PutUint32(data[start:start+4], v)
	}
	return data
}

func expectData() interface{} {
	v := make([]interface{}, 0)
	for i := 0; i < 30; i++ {
		v = append(v, map[string]interface{}{
			"key":   uint32(1),
			"Value": uint32(1),
		})
	}
	return v
}

func TestS1APCase(t *testing.T) {
	c, err := NewConverter("../../test/s1_ap_test.xml", converter.IDlConverterConfig{IsLittleEndian: false, LengthFieldLength: 4, PaddingLength: 4})
	require.NoError(t, err)
	hexStr := "0000000200000090efbbbfe4b8ade69687205749464900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c00000022efbbbf456e676c697368205749464900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000380000004e"
	data, err := hex.DecodeString(hexStr)
	require.NoError(t, err)
	name, v, err := c.DecodeWithID(33282, 32769, data)
	require.NoError(t, err)
	require.Equal(t, map[string]interface{}{
		"wiFiApNum": int32(2),
		"wiFiApArray": []interface{}{
			map[string]interface{}{
				"wiFiApName":     "中文 WIFI",
				"wiFiStrength":   int32(12),
				"wiFiEncryption": int32(34),
			},
			map[string]interface{}{
				"wiFiApName":     "English WIFI",
				"wiFiStrength":   int32(56),
				"wiFiEncryption": int32(78),
			},
		},
	}, v)
	require.Equal(t, "reportWiFiApList", name)
}
