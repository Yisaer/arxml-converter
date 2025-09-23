package converter

import (
	"fmt"
	"strings"

	idlAst "github.com/yisaer/idl-parser/ast"
	"github.com/yisaer/idl-parser/ast/typeref"
	"github.com/yisaer/idl-parser/converter"

	"arxml-converter/ast"
)

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

type ArXMLConverter struct {
	Parser            *ast.Parser
	convertedTypeRefs map[string]typeref.TypeRef
	idlModule         *idlAst.Module
	idlConverter      *converter.IDLConverter
}

func NewConverter(path string, config converter.IDlConverterConfig) (*ArXMLConverter, error) {
	parser, err := ast.NewParser(path)
	if err != nil {
		return nil, err
	}
	c := &ArXMLConverter{
		Parser:            parser,
		convertedTypeRefs: make(map[string]typeref.TypeRef),
	}
	if err := c.Parser.Parse(); err != nil {
		return nil, err
	}
	for _, dt := range c.Parser.DataTypes {
		typeRef, err := c.convertDataTypeToTypeRef(dt)
		if err != nil {
			return nil, fmt.Errorf("failed to convert datatype %s: %w", dt.ShorName, err)
		}
		c.convertedTypeRefs[strings.ToLower(dt.ShorName)] = typeRef
	}

	c.idlModule, err = c.TransformToIDLModule()
	if err != nil {
		return nil, err
	}
	c.idlConverter, err = converter.NewIDLConverterWithModule(config, *c.idlModule)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *ArXMLConverter) GetTypeByID(serviceID, eventID int) (string, error) {
	svc, ok := c.Parser.Services[serviceID]
	if !ok {
		return "", fmt.Errorf("service %v not found", serviceID)
	}
	event, ok := svc.Events[eventID]
	if !ok {
		return "", fmt.Errorf("event %v not found", eventID)
	}
	t, ok := c.convertedTypeRefs[strings.ToLower(event.ShortName)]
	if !ok {
		return "", fmt.Errorf("not found type ref %v for event id %v", event.ShortName, event.EventID)
	}
	return t.TypeName(), nil
}

// GetConvertedTypeRef 获取转换后的 TypeRef
func (c *ArXMLConverter) GetConvertedTypeRef(key string) (typeref.TypeRef, bool) {
	typeRef, exists := c.convertedTypeRefs[key]
	return typeRef, exists
}

// GetAllConvertedTypeRefs 获取所有转换后的 TypeRef
func (c *ArXMLConverter) GetAllConvertedTypeRefs() map[string]typeref.TypeRef {
	return c.convertedTypeRefs
}

// convertDataTypeToTypeRef 将 ArXML DataType 转换为 typeref.TypeRef
func (c *ArXMLConverter) convertDataTypeToTypeRef(dt *ast.DataType) (typeref.TypeRef, error) {
	switch dt.Category {
	case "TYPE_REFERENCE":
		return c.convertTypReference(dt.TypReference)
	case "ARRAY":
		return c.convertArray(dt.Array)
	case "STRUCTURE":
		return c.convertStructure(dt.Structure, dt.ShorName)
	default:
		return nil, fmt.Errorf("unsupported category: %s", dt.Category)
	}
}

// convertTypReference 转换 TypReference 为 TypeRef
func (c *ArXMLConverter) convertTypReference(tr *ast.TypReference) (typeref.TypeRef, error) {
	if tr == nil {
		return nil, fmt.Errorf("typReference is nil")
	}
	if strings.Contains(strings.ToLower(tr.Ref), "string") {
		if tr.StringSize > 0 {
			return typeref.NewFixedLengthStringType(int(tr.StringSize)), nil
		}
		return typeref.NewStringType(), nil
	}
	if basicType := c.getBasicTypeFromRef(tr.Ref); basicType != nil {
		return basicType, nil
	}
	return nil, fmt.Errorf("unknown type reference: %s", tr.Ref)
}

// convertArray 转换 Array 为 TypeRef
func (c *ArXMLConverter) convertArray(arr *ast.Array) (typeref.TypeRef, error) {
	if arr == nil {
		return nil, fmt.Errorf("array is nil")
	}
	innerType, err := c.convertRefToTypeRef(arr.RefType)
	if err != nil {
		return nil, fmt.Errorf("failed to convert array inner type: %w", err)
	}
	if arr.ArraySize > 0 {
		return typeref.NewArrayType(innerType, int(arr.ArraySize)), nil
	}
	return typeref.NewSequence(innerType), nil
}

// convertStructure 转换 Structure 为 TypeRef
func (c *ArXMLConverter) convertStructure(structData *ast.Structure, structName string) (typeref.TypeRef, error) {
	if structData == nil {
		return nil, fmt.Errorf("structure is nil")
	}

	return typeref.NewTypeName(structName), nil
}

// convertRefToTypeRef 根据引用字符串转换 TypeRef
func (c *ArXMLConverter) convertRefToTypeRef(ref string) (typeref.TypeRef, error) {
	// 检查是否是已知的 dataType 引用
	key := fmt.Sprintf("/dataTypes/%s", c.extractTypeNameFromRef(ref))
	if typeRef, exists := c.convertedTypeRefs[key]; exists {
		return typeRef, nil
	}

	// 检查是否是基础类型
	if basicType := c.getBasicTypeFromRef(ref); basicType != nil {
		return basicType, nil
	}

	// 检查是否是字符串类型
	if strings.Contains(strings.ToLower(ref), "string") {
		return typeref.NewStringType(), nil
	}

	// 默认作为自定义类型名
	typeName := c.extractTypeNameFromRef(ref)
	return typeref.NewTypeName(typeName), nil
}

// getBasicTypeFromRef 从引用中获取基础类型
func (c *ArXMLConverter) getBasicTypeFromRef(ref string) typeref.TypeRef {
	typeName := c.extractTypeNameFromRef(ref)
	lowerName := strings.ToLower(typeName)

	switch {
	case strings.Contains(lowerName, "uint8"):
		return typeref.NewOctetType()
	case strings.Contains(lowerName, "uint16"):
		return typeref.NewUnsignedShortType()
	case strings.Contains(lowerName, "uint32"):
		return typeref.NewUnsignedLong()
	case strings.Contains(lowerName, "uint64"):
		return typeref.NewUnsignedLongLong()
	case strings.Contains(lowerName, "bool"):
		return typeref.NewBooleanType()
	case strings.Contains(lowerName, "int8"):
		return typeref.NewOctetType()
	case strings.Contains(lowerName, "int16"):
		return typeref.NewShortType()
	case strings.Contains(lowerName, "int32"):
		return typeref.NewLongType()
	case strings.Contains(lowerName, "int64"):
		return typeref.NewLongLongType()
	case strings.Contains(lowerName, "float"):
		return typeref.NewFloatType()
	case strings.Contains(lowerName, "double"):
		return typeref.NewDoubleType()
	}
	return nil
}

// extractTypeNameFromRef 从引用路径中提取类型名
func (c *ArXMLConverter) extractTypeNameFromRef(ref string) string {
	// 移除路径前缀，只保留最后的类型名
	parts := strings.Split(ref, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ref
}
