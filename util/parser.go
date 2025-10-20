package util

import (
	"fmt"
	"strconv"
	"strings"
)

func ToUint16(raw string) (uint16, error) {
	val, err := strconv.ParseUint(raw, 10, 16) // 限制位宽为16位
	if err != nil {
		return 0, fmt.Errorf("cannot convert %s to uint16, err:%v", raw, err.Error())
	}
	return uint16(val), nil
}

func ToUint32(raw string) (uint32, error) {
	val, err := strconv.ParseUint(raw, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("cannot convert %s to uint32, err:%v", raw, err.Error())
	}
	return uint32(val), nil
}

func ToInt64(raw string) (int64, error) {
	val, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot convert %s to uint32, err:%v", raw, err.Error())
	}
	return int64(val), nil
}

func ExtractLast(ref string) string {
	parts := strings.Split(ref, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ref
}

func MergeUint16ToUint32(high16, low16 uint16) uint32 {
	return uint32(high16)<<16 | uint32(low16)
}

func MergeHexUint16ToUint32(svcHex, eventHex string) (uint16, uint32, error) {
	s1, err := strconv.ParseUint(svcHex, 0, 16)
	if err != nil {
		return 0, 0, err
	}
	svcID := uint16(s1)

	e1, err := strconv.ParseUint(eventHex, 0, 16)
	if err != nil {
		return 0, 0, err
	}
	eventID := uint16(e1)
	return svcID, MergeUint16ToUint32(svcID, eventID), nil
}
