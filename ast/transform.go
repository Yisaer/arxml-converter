package ast

import (
	"fmt"
	"strings"

	idlAst "github.com/yisaer/idl-parser/ast"
	"github.com/yisaer/idl-parser/ast/struct_type"
	"github.com/yisaer/idl-parser/ast/typeref"
)

type TransformHelper struct {
	DataTypes         map[string]*DataType
	convertedTypeRefs map[string]typeref.TypeRef
}

func NewTransformHelper(dataTypes map[string]*DataType) *TransformHelper {
	return &TransformHelper{
		DataTypes:         dataTypes,
		convertedTypeRefs: make(map[string]typeref.TypeRef),
	}
}

func (t *TransformHelper) GetConverterRef() map[string]typeref.TypeRef {
	return t.convertedTypeRefs
}

func (t *TransformHelper) TransformIntoModule() (*idlAst.Module, error) {
	for _, dt := range t.DataTypes {
		if dt.Category != "ARRAY" && dt.Category != "VECTOR" {
			typeRef, err := t.convertDataTypeToTypeRef(dt)
			if err != nil {
				return nil, fmt.Errorf("failed to convert datatype %s: %w", dt.ShorName, err)
			}
			t.convertedTypeRefs[strings.ToLower(dt.ShorName)] = typeRef
		}
	}
	for _, dt := range t.DataTypes {
		if dt.Category == "ARRAY" || dt.Category == "VECTOR" {
			typeRef, err := t.convertDataTypeToTypeRef(dt)
			if err != nil {
				return nil, fmt.Errorf("failed to convert datatype %s: %w", dt.ShorName, err)
			}
			t.convertedTypeRefs[strings.ToLower(dt.ShorName)] = typeRef
		}
	}
	var content []idlAst.ModuleContent
	for _, dt := range t.DataTypes {
		if dt.Category == "STRUCTURE" && dt.Structure != nil {
			structContent, err := t.transformStructure(dt)
			if err != nil {
				return nil, fmt.Errorf("failed to convert structure %s: %w", dt.ShorName, err)
			}
			content = append(content, *structContent)
		}
	}
	module := &idlAst.Module{
		Name:    "ArXMLDataTypes",
		Content: content,
		Type:    "Module",
	}
	return module, nil
}

// convertStructure 将 ArXML Structure 转换为 idlAst Struct
func (t *TransformHelper) transformStructure(dt *DataType) (*struct_type.Struct, error) {
	if dt.Structure == nil {
		return nil, fmt.Errorf("structure is nil for %s", dt.ShorName)
	}
	var fields []struct_type.Field
	for _, strField := range dt.Structure.STRList {

		fieldType, err := t.transformField(strField)
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
func (t *TransformHelper) transformField(strField *StructureTypRef) (typeref.TypeRef, error) {
	typeRef, err := t.createTypeRef(strField.Ref)
	if err != nil {
		return nil, fmt.Errorf("failed to create type ref for %s: %w", strField.Ref, err)
	}
	return typeRef, nil
}

// createTypeRef 根据引用字符串创建对应的 TypeRef
func (t *TransformHelper) createTypeRef(ref string) (typeref.TypeRef, error) {
	key := ExtractTypeNameFromRef(ref)
	targetType, ok := t.convertedTypeRefs[key]
	if !ok {
		return nil, fmt.Errorf("failed to convert type ref %s", ref)
	}
	return targetType, nil
}

// convertDataTypeToTypeRef 将 ArXML DataType 转换为 typeref.TypeRef
func (t *TransformHelper) convertDataTypeToTypeRef(dt *DataType) (typeref.TypeRef, error) {
	switch {
	case dt.Category == "TYPE_REFERENCE" || dt.TypReference != nil:
		return t.convertTypReference(dt.TypReference)
	case dt.Category == "ARRAY":
		return t.convertArray(dt.Array)
	case dt.Category == "VECTOR":
		return t.convertVector(dt.Vector)
	case dt.Category == "STRUCTURE":
		return t.convertStructure(dt.Structure, dt.ShorName)
	default:
		return nil, fmt.Errorf("unsupported category: %s", dt.Category)
	}
}

// convertTypReference 转换 TypReference 为 TypeRef
func (t *TransformHelper) convertTypReference(tr *TypReference) (typeref.TypeRef, error) {
	if tr == nil {
		return nil, fmt.Errorf("typReference is nil")
	}
	if basicType := GetBasicTypeFromRef(tr); basicType != nil {
		return basicType, nil
	}
	return nil, fmt.Errorf("unknown type reference: %s", tr.Ref)
}

// convertArray 转换 Array 为 TypeRef
func (t *TransformHelper) convertArray(arr *Array) (typeref.TypeRef, error) {
	if arr == nil {
		return nil, fmt.Errorf("array is nil")
	}
	innerType, err := t.convertRefToTypeRef(arr.RefType)
	if err != nil {
		return nil, fmt.Errorf("failed to convert array inner type: %w", err)
	}
	if arr.ArraySize > 0 {
		return typeref.NewArrayType(innerType, int(arr.ArraySize)), nil
	}
	return typeref.NewSequence(innerType), nil
}

func (t *TransformHelper) convertVector(v *Vector) (typeref.TypeRef, error) {
	if v == nil {
		return nil, fmt.Errorf("vector is nil")
	}
	innerType, err := t.convertRefToTypeRef(v.RefType)
	if err != nil {
		return nil, fmt.Errorf("failed to convert array inner type: %w", err)
	}
	return typeref.NewSequence(innerType), nil
}

// convertStructure 转换 Structure 为 TypeRef
func (t *TransformHelper) convertStructure(structData *Structure, structName string) (typeref.TypeRef, error) {
	if structData == nil {
		return nil, fmt.Errorf("structure is nil")
	}

	return typeref.NewTypeName(structName), nil
}

// convertRefToTypeRef 根据引用字符串转换 TypeRef
func (t *TransformHelper) convertRefToTypeRef(ref string) (typeref.TypeRef, error) {
	if typeRef, exists := t.convertedTypeRefs[ExtractTypeNameFromRef(ref)]; exists {
		return typeRef, nil
	}
	return nil, fmt.Errorf("unkown ref %v", ref)
}

func ExtractTypeNameFromRef(ref string) string {
	extractedRef := extractType(ref)
	return strings.ToLower(extractedRef)
}

func extractType(ref string) string {
	parts := strings.Split(ref, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ref
}
