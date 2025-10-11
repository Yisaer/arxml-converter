package util

import (
	"fmt"
	"strconv"
)

var (
	AppDataTypePrefix       = "/DataTypes/ApplicationDataType/"
	ImplementDataTypePrefix = "/DataTypes/ImplementationDataTypes/"
)

func ToUint16(raw string) (uint16, error) {
	val, err := strconv.ParseUint(raw, 10, 16) // 限制位宽为16位
	if err != nil {
		return 0, fmt.Errorf("cannot convert %s to uint16, err:%v", raw, err.Error())
	}
	return uint16(val), nil
}
