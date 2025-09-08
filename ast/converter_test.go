package ast

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArXMLConverter_Decode_RoadIndex(t *testing.T) {
	// Create converter with big-endian (大端) for testing
	converter, err := NewConverter("example.xml", false) // false = big-endian
	require.NoError(t, err)
	require.NotNil(t, converter)

	// Test data for roadIndex structure
	// roadIndex contains:
	// - SegmentIndex (uint32_t) - 4 bytes
	// - LinkIndex (uint32_t) - 4 bytes
	// Total: 8 bytes

	// Create test data: SegmentIndex = 12345, LinkIndex = 67890
	// Big-endian representation
	testData := make([]byte, 8)
	binary.BigEndian.PutUint32(testData[0:4], 12345) // SegmentIndex
	binary.BigEndian.PutUint32(testData[4:8], 67890) // LinkIndex

	t.Run("decode_roadindex_big_endian", func(t *testing.T) {
		result, err := converter.Decode(testData, "roadIndex")

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, map[string]interface{}{
			"SegmentIndex": int64(12345),
			"LinkIndex":    int64(67890),
		}, result)
	})

	t.Run("decode_roadindex_little_endian", func(t *testing.T) {
		// Create converter with little-endian for comparison
		converterLE, err := NewConverter("example.xml", true) // true = little-endian
		require.NoError(t, err)

		// Create test data with little-endian representation
		testDataLE := make([]byte, 8)
		binary.LittleEndian.PutUint32(testDataLE[0:4], 12345) // SegmentIndex
		binary.LittleEndian.PutUint32(testDataLE[4:8], 67890) // LinkIndex

		result, err := converterLE.Decode(testDataLE, "roadIndex")

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, map[string]interface{}{
			"SegmentIndex": int64(12345),
			"LinkIndex":    int64(67890),
		}, result)
	})
}
