package converter

import (
	"fmt"
	"strings"

	idlAst "github.com/yisaer/idl-parser/ast"
	"github.com/yisaer/idl-parser/ast/struct_type"
	"github.com/yisaer/idl-parser/ast/typeref"

	"arxml-converter/mod"
)

func (c *ArXMLConverter) TransformToIDLModule() (*idlAst.Module, error) {
	var content []idlAst.ModuleContent
	for _, dt := range c.Parser.DataTypes {
		if dt.Category == "STRUCTURE" && dt.Structure != nil {
			structContent, err := c.transformStructure(dt)
			if err != nil {
				return nil, fmt.Errorf("failed to convert structure %s: %w", dt.ShorName, err)
			}
			content = append(content, *structContent)
		}
	}

	// 创建 IDL 模块
	module := &idlAst.Module{
		Name:    "ArXMLDataTypes",
		Content: content,
		Type:    "Module",
	}

	return module, nil
}

// convertStructure 将 ArXML Structure 转换为 idlAst Struct
func (c *ArXMLConverter) transformStructure(dt *mod.DataType) (*struct_type.Struct, error) {
	if dt.Structure == nil {
		return nil, fmt.Errorf("structure is nil for %s", dt.ShorName)
	}
	var fields []struct_type.Field
	for _, strField := range dt.Structure.STRList {

		fieldType, err := c.transformField(strField)
		if err != nil {
			return nil, fmt.Errorf("failed to convert field %s: %w", strField.ShorName, err)
		}
		fields = append(fields, struct_type.Field{Type: fieldType, Name: strField.ShorName})
	}
	return &struct_type.Struct{
		Name:   dt.ShorName,
		Fields: fields,
		Type:   "Struct",
	}, nil
}

// convertField 将 ArXML StructureTypRef 转换为 idlAst Field
func (c *ArXMLConverter) transformField(strField *mod.StructureTypRef) (typeref.TypeRef, error) {
	typeRef, err := c.createTypeRef(strField.Ref)
	if err != nil {
		return nil, fmt.Errorf("failed to create type ref for %s: %w", strField.Ref, err)
	}
	return typeRef, nil
}

// createTypeRef 根据引用字符串创建对应的 TypeRef
func (c *ArXMLConverter) createTypeRef(ref string) (typeref.TypeRef, error) {
	typeName := c.extractTypeName(ref)

	t, ok := c.convertedTypeRefs[strings.ToLower(typeName)]
	if !ok {
		return nil, fmt.Errorf("failed to convert type ref %s to type %s", ref, typeName)
	}
	return t, nil
}

// extractTypeName 从引用路径中提取类型名
func (c *ArXMLConverter) extractTypeName(ref string) string {
	// 移除路径前缀，只保留最后的类型名
	parts := strings.Split(ref, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ref
}

// getBasicType 获取基础类型
func (c *ArXMLConverter) getBasicType(typeName string) typeref.TypeRef {
	lowerName := strings.ToLower(typeName)
	switch {
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
	case strings.Contains(lowerName, "uint8"):
		return typeref.NewOctetType()
	case strings.Contains(lowerName, "uint16"):
		return typeref.NewUnsignedShortType()
	case strings.Contains(lowerName, "uint32"):
		return typeref.NewUnsignedLong()
	case strings.Contains(lowerName, "uint64"):
		return typeref.NewUnsignedLongLong()
	case strings.Contains(lowerName, "float"):
		return typeref.NewFloatType()
	case strings.Contains(lowerName, "double"):
		return typeref.NewDoubleType()
	}
	return nil
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
func (c *ArXMLConverter) convertDataTypeToTypeRef(dt *mod.DataType) (typeref.TypeRef, error) {
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
func (c *ArXMLConverter) convertTypReference(tr *mod.TypReference) (typeref.TypeRef, error) {
	if tr == nil {
		return nil, fmt.Errorf("typReference is nil")
	}
	if basicType := c.getBasicTypeFromRef(tr); basicType != nil {
		return basicType, nil
	}
	return nil, fmt.Errorf("unknown type reference: %s", tr.Ref)
}

// convertArray 转换 Array 为 TypeRef
func (c *ArXMLConverter) convertArray(arr *mod.Array) (typeref.TypeRef, error) {
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
func (c *ArXMLConverter) convertStructure(structData *mod.Structure, structName string) (typeref.TypeRef, error) {
	if structData == nil {
		return nil, fmt.Errorf("structure is nil")
	}

	return typeref.NewTypeName(structName), nil
}

// convertRefToTypeRef 根据引用字符串转换 TypeRef
func (c *ArXMLConverter) convertRefToTypeRef(ref string) (typeref.TypeRef, error) {
	if typeRef, exists := c.convertedTypeRefs[strings.ToLower(extractTypeNameFromRef(ref))]; exists {
		return typeRef, nil
	}
	return nil, fmt.Errorf("unkown ref %v", ref)
}

// getBasicTypeFromRef 从引用中获取基础类型
func (c *ArXMLConverter) getBasicTypeFromRef(tr *mod.TypReference) typeref.TypeRef {
	ref := tr.Ref
	typeName := extractTypeNameFromRef(ref)
	lowerName := strings.ToLower(typeName)
	switch {
	case strings.Contains(lowerName, "string"):
		if tr.StringSize > 0 {
			return typeref.NewFixedLengthStringType(int(tr.StringSize))
		}
		return typeref.NewStringType()
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

func extractTypeNameFromRef(ref string) string {
	// 移除路径前缀，只保留最后的类型名
	parts := strings.Split(ref, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ref
}
