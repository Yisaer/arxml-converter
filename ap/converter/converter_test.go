package converter

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yisaer/idl-parser/converter"
)

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
