package converter

import (
	"fmt"
	"strings"

	idlAst "github.com/yisaer/idl-parser/ast"
	"github.com/yisaer/idl-parser/ast/struct_type"
	"github.com/yisaer/idl-parser/ast/typeref"

	"arxml-converter/ast"
)

// ArXMLToIDLConverter 将 ArXML AST 转换为 idlAst.Module
type ArXMLToIDLConverter struct {
	parser *ast.Parser
}

// NewArXMLToIDLConverter 创建新的转换器
func NewArXMLToIDLConverter(parser *ast.Parser) *ArXMLToIDLConverter {
	return &ArXMLToIDLConverter{
		parser: parser,
	}
}

// ConvertToIDLModule 将 ArXML parser 的结果转换为 idlAst.Module
func (c *ArXMLToIDLConverter) ConvertToIDLModule() (*idlAst.Module, error) {
	// 创建模块内容列表
	var content []idlAst.ModuleContent

	// 遍历所有数据类型
	for _, dt := range c.parser.DtList {
		// 只转换结构体类型
		if dt.Category == "STRUCTURE" && dt.Structure != nil {
			structContent, err := c.convertStructure(dt)
			if err != nil {
				return nil, fmt.Errorf("failed to convert structure %s: %w", dt.ShorName, err)
			}
			content = append(content, structContent)
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
func (c *ArXMLToIDLConverter) convertStructure(dt *ast.DataType) (*struct_type.Struct, error) {
	if dt.Structure == nil {
		return nil, fmt.Errorf("structure is nil for %s", dt.ShorName)
	}
	var fields []struct_type.Field
	// 转换结构体字段
	for _, strField := range dt.Structure.STRList {
		field, err := c.convertField(strField)
		if err != nil {
			return nil, fmt.Errorf("failed to convert field %s: %w", strField.ShorName, err)
		}
		fields = append(fields, *field)
	}
	return &struct_type.Struct{
		Name:   dt.ShorName,
		Fields: fields,
		Type:   "Struct",
	}, nil
}

// convertField 将 ArXML StructureTypRef 转换为 idlAst Field
func (c *ArXMLToIDLConverter) convertField(strField *ast.StructureTypRef) (*struct_type.Field, error) {
	// 根据引用类型创建对应的 TypeRef
	typeRef, err := c.createTypeRef(strField.Ref)
	if err != nil {
		return nil, fmt.Errorf("failed to create type ref for %s: %w", strField.Ref, err)
	}

	return &struct_type.Field{
		Type: typeRef,
		Name: strField.ShorName,
	}, nil
}

// createTypeRef 根据引用字符串创建对应的 TypeRef
func (c *ArXMLToIDLConverter) createTypeRef(ref string) (typeref.TypeRef, error) {
	// 清理引用路径，提取类型名
	typeName := c.extractTypeName(ref)

	// 检查是否是基础类型
	if basicType := c.getBasicType(typeName); basicType != nil {
		return basicType, nil
	}

	// 检查是否是数组类型
	if arrayType := c.getArrayType(ref); arrayType != nil {
		return arrayType, nil
	}

	// 检查是否是字符串类型
	if stringType := c.getStringType(ref); stringType != nil {
		return stringType, nil
	}

	// 默认作为自定义类型
	return typeref.NewTypeName(typeName), nil
}

// extractTypeName 从引用路径中提取类型名
func (c *ArXMLToIDLConverter) extractTypeName(ref string) string {
	// 移除路径前缀，只保留最后的类型名
	parts := strings.Split(ref, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ref
}

// getBasicType 获取基础类型
func (c *ArXMLToIDLConverter) getBasicType(typeName string) typeref.TypeRef {
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

// getArrayType 检查并创建数组类型
func (c *ArXMLToIDLConverter) getArrayType(ref string) typeref.TypeRef {
	// 在 ArXML 中，数组类型通过 Array 结构体表示
	// 这里我们需要检查是否有对应的 Array 定义
	for _, dt := range c.parser.DtList {
		if dt.Category == "ARRAY" && dt.Array != nil {
			// 检查引用是否匹配
			arrayRef := fmt.Sprintf("/dataTypes/%s", dt.ShorName)
			if ref == arrayRef {
				innerType, err := c.createTypeRef(dt.Array.RefType)
				if err != nil {
					continue
				}
				if dt.Array.ArraySize > 0 {
					return typeref.NewArrayType(innerType, int(dt.Array.ArraySize))
				}
				return typeref.NewSequence(innerType)
			}
		}
	}
	return nil
}

// getStringType 检查并创建字符串类型
func (c *ArXMLToIDLConverter) getStringType(ref string) typeref.TypeRef {
	lowerRef := strings.ToLower(ref)
	if strings.Contains(lowerRef, "string") {
		// 检查是否有固定长度
		for _, dt := range c.parser.DtList {
			if dt.TypReference != nil && dt.TypReference.Ref == ref && dt.TypReference.StringSize > 0 {
				return typeref.NewFixedLengthStringType(int(dt.TypReference.StringSize))
			}
		}
		return typeref.NewStringType()
	}
	return nil
}
