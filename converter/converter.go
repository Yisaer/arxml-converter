package converter

import (
	"fmt"

	idlAst "github.com/yisaer/idl-parser/ast"
	"github.com/yisaer/idl-parser/ast/typeref"
	"github.com/yisaer/idl-parser/converter"

	"github.com/yisaer/arxml-converter/ast"
	"github.com/yisaer/arxml-converter/mod"
)

type ArXMLConverter struct {
	Parser       *ast.Parser
	idlModule    *idlAst.Module
	idlConverter *converter.IDLConverter
	transformer  *mod.TransformHelper
}

func NewConverter(path string, config converter.IDlConverterConfig) (*ArXMLConverter, error) {
	parser, err := ast.NewParser(path)
	if err != nil {
		return nil, err
	}
	c := &ArXMLConverter{
		Parser: parser,
	}
	if err := c.Parser.Parse(); err != nil {
		return nil, err
	}
	transformerHelper := mod.NewTransformHelper(c.Parser.DataTypes)
	c.transformer = transformerHelper
	c.idlModule, err = c.transformer.TransformIntoModule()
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
	interfaceRef := mod.ExtractTypeNameFromRef(svc.ServiceInterfaceRef)
	targetInterface, ok := c.Parser.Interfaces[interfaceRef]
	if !ok {
		return nil, fmt.Errorf("interface %v not found for serviceID %v", interfaceRef, serviceID)
	}
	event, ok := svc.Events[eventID]
	if ok {
		eventRef := mod.ExtractTypeNameFromRef(event.EventRef)
		targetEvent, ok := targetInterface.Events[eventRef]
		if !ok {
			return nil, fmt.Errorf("event %v not found in interface %v", eventRef, targetInterface.Shortname)
		}
		typeRef := mod.ExtractTypeNameFromRef(targetEvent.TypeRef)
		targetTypRef, ok := c.transformer.GetConverterRef()[typeRef]
		if !ok {
			return nil, fmt.Errorf("type %v not found in interface %v event %v", typeRef, interfaceRef, eventRef)
		}
		return targetTypRef, nil
	}
	fieldNotify, ok := svc.FieldNotify[eventID]
	if ok {
		fieldNotifyRef := mod.ExtractTypeNameFromRef(fieldNotify.FieldRef)
		targetField, ok := targetInterface.Fields[fieldNotifyRef]
		if !ok {
			return nil, fmt.Errorf("field %v not found in interface %v", fieldNotifyRef, targetInterface.Shortname)
		}
		typeRef := mod.ExtractTypeNameFromRef(targetField.TypeRef)
		targetTypRef, ok := c.transformer.GetConverterRef()[typeRef]
		if !ok {
			return nil, fmt.Errorf("type %v not found in interface %v field %v", typeRef, interfaceRef, fieldNotifyRef)
		}
		return targetTypRef, nil
	}
	return nil, fmt.Errorf("unknown eventID:%v in serviceID:%v", serviceID, eventID)
}
