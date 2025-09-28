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
		if dt.Category != "ARRAY" {
			typeRef, err := c.convertDataTypeToTypeRef(dt)
			if err != nil {
				return nil, fmt.Errorf("failed to convert datatype %s: %w", dt.ShorName, err)
			}
			c.convertedTypeRefs[strings.ToLower(dt.ShorName)] = typeRef
		}
	}
	for _, dt := range c.Parser.DataTypes {
		if dt.Category == "ARRAY" {
			typeRef, err := c.convertDataTypeToTypeRef(dt)
			if err != nil {
				return nil, fmt.Errorf("failed to convert datatype %s: %w", dt.ShorName, err)
			}
			c.convertedTypeRefs[strings.ToLower(dt.ShorName)] = typeRef
		}
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

func (c *ArXMLConverter) DecodeWithID(serviceID, eventID int, data []byte) (interface{}, error) {
	t, err := c.GetTypeByID(serviceID, eventID)
	if err != nil {
		return nil, err
	}
	result, _, err := c.idlConverter.ParseDataByType(data, t, *c.idlModule)
	return result, err
}

func (c *ArXMLConverter) GetTypeByID(serviceID, eventID int) (typeref.TypeRef, error) {
	svc, ok := c.Parser.Services[serviceID]
	if !ok {
		return nil, fmt.Errorf("service %v not found", serviceID)
	}
	interfaceRef := strings.ToLower(extractTypeNameFromRef(svc.ServiceInterfaceRef))
	targetInterface, ok := c.Parser.Interfaces[interfaceRef]
	if !ok {
		return nil, fmt.Errorf("interface %v not found for serviceID %v", interfaceRef, serviceID)
	}
	event, ok := svc.Events[eventID]
	if ok {
		eventRef := strings.ToLower(extractTypeNameFromRef(event.EventRef))
		targetEvent, ok := targetInterface.Events[eventRef]
		if !ok {
			return nil, fmt.Errorf("event %v not found in interface %v", eventRef, targetInterface.Shortname)
		}
		typeRef := strings.ToLower(extractTypeNameFromRef(targetEvent.TypeRef))
		targetTypRef, ok := c.convertedTypeRefs[typeRef]
		if !ok {
			return nil, fmt.Errorf("type %v not found in interface %v event %v", typeRef, interfaceRef, eventRef)
		}
		return targetTypRef, nil
	}
	fieldNotify, ok := svc.FieldNotify[eventID]
	if ok {
		fieldNotifyRef := strings.ToLower(extractTypeNameFromRef(fieldNotify.FieldRef))
		targetField, ok := targetInterface.Fields[fieldNotifyRef]
		if !ok {
			return nil, fmt.Errorf("field %v not found in interface %v", fieldNotifyRef, targetInterface.Shortname)
		}
		typeRef := strings.ToLower(extractTypeNameFromRef(targetField.TypeRef))
		targetTypRef, ok := c.convertedTypeRefs[typeRef]
		if !ok {
			return nil, fmt.Errorf("type %v not found in interface %v field %v", typeRef, interfaceRef, fieldNotifyRef)
		}
		return targetTypRef, nil
	}
	return nil, fmt.Errorf("unknown eventID:%v in serviceID:%v", serviceID, eventID)
}
