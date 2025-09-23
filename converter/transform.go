package converter

import (
	"fmt"
	"strings"

	idlAst "github.com/yisaer/idl-parser/ast"
	"github.com/yisaer/idl-parser/ast/struct_type"
	"github.com/yisaer/idl-parser/ast/typeref"

	"arxml-converter/ast"
)

func (c *ArXMLConverter) TransformToIDLModule() (*idlAst.Module, error) {
	var content []idlAst.ModuleContent
	for _, dt := range c.Parser.DataTypes {
		if dt.Category == "STRUCTURE" && dt.Structure != nil {
			structContent, err := c.transformStructure(dt)
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
func (c *ArXMLConverter) transformStructure(dt *ast.DataType) (*struct_type.Struct, error) {
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
func (c *ArXMLConverter) transformField(strField *ast.StructureTypRef) (typeref.TypeRef, error) {
	typeRef, err := c.createTypeRef(strField.Ref)
	if err != nil {
		return nil, fmt.Errorf("failed to create type ref for %s: %w", strField.Ref, err)
	}
	return typeRef, nil
}

// createTypeRef 根据引用字符串创建对应的 TypeRef
func (c *ArXMLConverter) createTypeRef(ref string) (typeref.TypeRef, error) {
	if strings.Contains(ref, "AppSrv_Inputkey") {
		fmt.Println("here")
	}
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

// getArrayType 检查并创建数组类型
func (c *ArXMLConverter) getArrayType(ref string) typeref.TypeRef {
	// 在 ArXML 中，数组类型通过 Array 结构体表示
	// 这里我们需要检查是否有对应的 Array 定义
	for _, dt := range c.Parser.DataTypes {
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
func (c *ArXMLConverter) getStringType(ref string) typeref.TypeRef {
	lowerRef := strings.ToLower(ref)
	if strings.Contains(lowerRef, "string") {
		for _, dt := range c.Parser.DataTypes {
			if dt.TypReference != nil && dt.TypReference.Ref == ref && dt.TypReference.StringSize > 0 {
				return typeref.NewFixedLengthStringType(int(dt.TypReference.StringSize))
			}
		}
		return typeref.NewStringType()
	}
	return nil
}
