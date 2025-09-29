package mod

import (
	"strings"

	"github.com/yisaer/idl-parser/ast/typeref"
)

func GetBasicTypeFromRef(tr *TypReference) typeref.TypeRef {
	ref := tr.Ref
	typeName := ExtractTypeNameFromRef(ref)
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
