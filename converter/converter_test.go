package converter

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yisaer/idl-parser/converter"
)

func TestParser(t *testing.T) {
	p, err := NewConverter("../test/example.xml", converter.IDlConverterConfig{IsLittleEndian: false, LengthFieldLength: 2, PaddingLength: 4})
	require.NoError(t, err)
	result, err := p.DecodeWithID(22530, 36865, testdata())
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
