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
