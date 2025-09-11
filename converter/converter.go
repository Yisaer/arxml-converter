package converter

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"

	"arxml-converter/ast"
)

type ParseAble interface {
	Parse([]byte) (interface{}, error)
}

type FieldRefType int

const (
	StringType FieldRefType = iota
	BoolType
	FloatType
	DoubleType
	Int8Type
	Int16Type
	Int32Type
	Int64Type
	Uint8Type
	Uint16Type
	Uint32Type
	Uint64Type
)

func (ft FieldRefType) String() string {
	switch ft {
	case StringType:
		return "String"
	case BoolType:
		return "Bool"
	case FloatType:
		return "Float"
	case DoubleType:
		return "Double"
	case Int8Type:
		return "Int8"
	case Int16Type:
		return "Int16"
	case Int32Type:
		return "Int32"
	case Int64Type:
		return "Int64"
	case Uint8Type:
		return "Uint8"
	case Uint16Type:
		return "Uint16"
	case Uint32Type:
		return "Uint32"
	case Uint64Type:
		return "Uint64"
	default:
		return "Unknown"
	}
}

func RefTypeToFieldType(refType string) (FieldRefType, error) {
	lower := strings.ToLower(refType)
	switch {
	case strings.Contains(lower, "string"):
		return StringType, nil
	case strings.Contains(lower, "bool"):
		return BoolType, nil
	case strings.Contains(lower, "float"):
		return FloatType, nil
	case strings.Contains(lower, "double"):
		return DoubleType, nil
	case strings.Contains(lower, "uint8"):
		return Uint8Type, nil
	case strings.Contains(lower, "uint16"):
		return Uint16Type, nil
	case strings.Contains(lower, "uint32"):
		return Uint32Type, nil
	case strings.Contains(lower, "uint64"):
		return Uint64Type, nil
	case strings.Contains(lower, "int8"):
		return Int8Type, nil
	case strings.Contains(lower, "int16"):
		return Int16Type, nil
	case strings.Contains(lower, "int32"):
		return Int32Type, nil
	case strings.Contains(lower, "int64"):
		return Int64Type, nil
	}
	return StringType, fmt.Errorf("unknown field type: %s", refType)
}

type ArXMLConverter struct {
	Parser         *ast.Parser
	IsLittleEndian bool
	typeRefs       map[string]*ast.TypReference
	arrRefS        map[string]*ast.Array
	structRefs     map[string]*ast.Structure
}

func NewConverter(path string, IsLittleEndian bool) (*ArXMLConverter, error) {
	parser, err := ast.NewParser(path)
	if err != nil {
		return nil, err
	}
	c := &ArXMLConverter{
		Parser:         parser,
		IsLittleEndian: IsLittleEndian,
		typeRefs:       make(map[string]*ast.TypReference),
		arrRefS:        make(map[string]*ast.Array),
		structRefs:     make(map[string]*ast.Structure),
	}
	if err := c.Parser.Parse(); err != nil {
		return nil, err
	}
	for _, dt := range c.Parser.DtList {
		key := fmt.Sprintf("/dataTypes/%s", dt.ShorName)
		switch {
		case dt.TypReference != nil:
			c.typeRefs[key] = dt.TypReference
		case dt.Array != nil:
			c.arrRefS[key] = dt.Array
		case dt.Structure != nil:
			c.structRefs[key] = dt.Structure
		default:
			return nil, fmt.Errorf("unknown dt: %s, not found: %v", dt.ShorName, key)
		}
	}
	return c, nil
}

func (c *ArXMLConverter) Decode(data []byte, stName string) (interface{}, error) {
	sts, ok := c.structRefs[fmt.Sprintf("/dataTypes/%s", stName)]
	if !ok {
		return nil, fmt.Errorf("struct ref not found: %s", stName)
	}
	v, _, err := c.ParseStructure(data, sts)
	return v, err
}

func (c *ArXMLConverter) ParseStructure(data []byte, sts *ast.Structure) (interface{}, []byte, error) {
	m := make(map[string]interface{})
	var v interface{}
	var err error
	remainData := data
	for _, col := range sts.STRList {
		targetTyp, ok1 := c.typeRefs[col.Ref]
		targetArray, ok2 := c.arrRefS[col.Ref]
		targetSts, ok3 := c.structRefs[col.Ref]
		switch {
		case ok1:
			v, remainData, err = c.ParseTypReference(remainData, targetTyp)
			if err != nil {
				return nil, nil, err
			}
			m[col.ShorName] = v
		case ok2:
			v, remainData, err = c.ParseArray(remainData, targetArray)
			if err != nil {
				return nil, nil, err
			}
			m[col.ShorName] = v
		case ok3:
			v, remainData, err = c.ParseStructure(remainData, targetSts)
			if err != nil {
				return nil, nil, err
			}
			m[col.ShorName] = v
		default:
			return nil, nil, fmt.Errorf("unknown field type: %s", col.Ref)
		}
	}
	return m, remainData, nil
}

func (c *ArXMLConverter) ParseArray(data []byte, array *ast.Array) (interface{}, []byte, error) {
	targets := make([]interface{}, 0)
	var v interface{}
	var err error
	remainData := data
	targetTyp, ok1 := c.typeRefs[array.RefType]
	targetArray, ok2 := c.arrRefS[array.RefType]
	targetSts, ok3 := c.structRefs[array.RefType]
	switch {
	case ok1:
		for i := int64(0); i < array.ArraySize; i++ {
			v, remainData, err = c.ParseTypReference(remainData, targetTyp)
			if err != nil {
				return nil, nil, err
			}
			targets = append(targets, v)
		}
	case ok2:
		for i := int64(0); i < array.ArraySize; i++ {
			v, remainData, err = c.ParseArray(remainData, targetArray)
			if err != nil {
				return nil, nil, err
			}
			targets = append(targets, v)
		}
	case ok3:
		for i := int64(0); i < array.ArraySize; i++ {
			v, remainData, err = c.ParseStructure(remainData, targetSts)
			if err != nil {
				return nil, nil, err
			}
			targets = append(targets, v)
		}
	default:
		return nil, nil, fmt.Errorf("unknown array type: %s", array.RefType)
	}
	return targets, remainData, nil
}

func (c *ArXMLConverter) ParseTypReference(data []byte, t *ast.TypReference) (interface{}, []byte, error) {
	ft, err := RefTypeToFieldType(t.Ref)
	if err != nil {
		return nil, nil, err
	}
	switch ft {
	case StringType:
		if t.StringSize > 0 {
			return parseBytesToFixedLengthString(data, int(t.StringSize))
		}
		return ParseBytesToDynamicString(data)
	case BoolType:
		return c.parseBytesToBoolean(data)
	case FloatType:
		return c.parseBytesToFloat(data)
	case DoubleType:
		return c.parseBytesToDouble(data)
	case Int8Type:
		return c.parseBytesToInt8(data)
	case Int16Type:
		return c.parseBytesToInt16(data)
	case Int32Type:
		return c.parseBytesToInt32(data)
	case Int64Type:
		return c.parseBytesToInt64(data)
	case Uint8Type:
		return c.parseBytesToUint8(data)
	case Uint16Type:
		return c.parseBytesToUint16(data)
	case Uint32Type:
		return c.parseBytesToUint32(data)
	case Uint64Type:
		return c.parseBytesToUint64(data)
	}
	return nil, nil, fmt.Errorf("unknown field type: %s", ft.String())
}

func (c *ArXMLConverter) parseBytesToBoolean(data []byte) (bool, []byte, error) {
	if len(data) < 1 {
		return false, nil, fmt.Errorf("expect data len %v got len %v", 1, len(data))
	}
	return data[0] != 0x00, data[1:], nil
}

func (c *ArXMLConverter) parseBytesToFloat(data []byte) (float64, []byte, error) {
	if len(data) < 4 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 4, len(data))
	}
	parseData, remainData := data[:4], data[4:]
	var value uint32
	if c.IsLittleEndian {
		value = binary.LittleEndian.Uint32(parseData)
	} else {
		value = binary.BigEndian.Uint32(parseData)
	}
	return float64(math.Float32frombits(value)), remainData, nil
}

func (c *ArXMLConverter) parseBytesToDouble(data []byte) (float64, []byte, error) {
	if len(data) < 8 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 8, len(data))
	}
	parseData, remainData := data[:8], data[8:]
	var value uint64
	if c.IsLittleEndian {
		value = binary.LittleEndian.Uint64(parseData)
	} else {
		value = binary.BigEndian.Uint64(parseData)
	}
	return math.Float64frombits(value), remainData, nil
}

func (c *ArXMLConverter) parseBytesToInt8(data []byte) (int64, []byte, error) {
	if len(data) < 1 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 1, len(data))
	}
	parseData, remainData := data[:1], data[1:]
	return int64(int8(parseData[0])), remainData, nil
}

func (c *ArXMLConverter) parseBytesToInt16(data []byte) (int64, []byte, error) {
	if len(data) < 2 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 2, len(data))
	}
	parseData, remainData := data[:2], data[2:]
	var value int16
	if c.IsLittleEndian {
		value = int16(binary.LittleEndian.Uint16(parseData))
	} else {
		value = int16(binary.BigEndian.Uint16(parseData))
	}
	return int64(value), remainData, nil
}

func (c *ArXMLConverter) parseBytesToInt32(data []byte) (int64, []byte, error) {
	if len(data) < 4 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 4, len(data))
	}
	parseData, remainData := data[:4], data[4:]
	var value int32
	if c.IsLittleEndian {
		value = int32(binary.LittleEndian.Uint32(parseData))
	} else {
		value = int32(binary.BigEndian.Uint32(parseData))
	}
	return int64(value), remainData, nil
}

func (c *ArXMLConverter) parseBytesToInt64(data []byte) (int64, []byte, error) {
	if len(data) < 8 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 8, len(data))
	}
	parseData, remainData := data[:8], data[8:]
	var value int64
	if c.IsLittleEndian {
		value = int64(binary.LittleEndian.Uint64(parseData))
	} else {
		value = int64(binary.BigEndian.Uint64(parseData))
	}
	return value, remainData, nil
}

func (c *ArXMLConverter) parseBytesToUint8(data []byte) (uint64, []byte, error) {
	if len(data) < 1 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 1, len(data))
	}
	parseData, remainData := data[:1], data[1:]
	return uint64(parseData[0]), remainData, nil
}

func (c *ArXMLConverter) parseBytesToUint16(data []byte) (uint64, []byte, error) {
	if len(data) < 2 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 2, len(data))
	}
	parseData, remainData := data[:2], data[2:]
	var value uint16
	if c.IsLittleEndian {
		value = binary.LittleEndian.Uint16(parseData)
	} else {
		value = binary.BigEndian.Uint16(parseData)
	}
	return uint64(value), remainData, nil
}

func (c *ArXMLConverter) parseBytesToUint32(data []byte) (uint64, []byte, error) {
	if len(data) < 4 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 4, len(data))
	}
	parseData, remainData := data[:4], data[4:]
	var value uint32
	if c.IsLittleEndian {
		value = binary.LittleEndian.Uint32(parseData)
	} else {
		value = binary.BigEndian.Uint32(parseData)
	}
	return uint64(value), remainData, nil
}

func (c *ArXMLConverter) parseBytesToUint64(data []byte) (uint64, []byte, error) {
	if len(data) < 8 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 8, len(data))
	}
	parseData, remainData := data[:8], data[8:]
	var value uint64
	if c.IsLittleEndian {
		value = binary.LittleEndian.Uint64(parseData)
	} else {
		value = binary.BigEndian.Uint64(parseData)
	}
	return value, remainData, nil
}
