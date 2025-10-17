package converter

import (
	"fmt"

	"github.com/beevik/etree"
	"github.com/yisaer/idl-parser/ast/typeref"
	"github.com/yisaer/idl-parser/converter"

	"github.com/yisaer/arxml-converter/cp/parser"
)

type ArxmlCPConverter struct {
	path         string
	config       converter.IDlConverterConfig
	parser       *parser.Parser
	idlConverter *converter.IDLConverter
}

func NewArxmlCPConverterWithDoc(doc *etree.Document, config converter.IDlConverterConfig) (*ArxmlCPConverter, error) {
	p := parser.NewParserWithDoc(doc)
	if err := p.Parse(); err != nil {
		return nil, err
	}
	idlConverter, err := converter.NewIDLConverterWithModule(config, *p.GetModule())
	if err != nil {
		return nil, fmt.Errorf("error creating idlConverter: %v", err)
	}
	return &ArxmlCPConverter{
		idlConverter: idlConverter,
		parser:       p,
		config:       config,
	}, nil
}

func NewArxmlCPConverter(path string, config converter.IDlConverterConfig) (*ArxmlCPConverter, error) {
	p, err := parser.NewParser(path)
	if err != nil {
		return nil, err
	}
	if err := p.Parse(); err != nil {
		return nil, err
	}
	idlConverter, err := converter.NewIDLConverterWithModule(config, *p.GetModule())
	if err != nil {
		return nil, fmt.Errorf("error creating idlConverter: %v", err)
	}
	return &ArxmlCPConverter{
		idlConverter: idlConverter,
		parser:       p,
		path:         path,
		config:       config,
	}, nil
}

func (c *ArxmlCPConverter) Convert(serviceID uint16, headerID uint32, data []byte) (string, interface{}, error) {
	key, tr, err := c.parser.FindTypeRefByID(serviceID, headerID)
	if err != nil {
		return "", nil, err
	}
	got, _, err := c.idlConverter.ParseDataByType(data, tr, *c.parser.GetModule())
	return key, got, err
}

func (c *ArxmlCPConverter) GetDataTypeByID(serviceID uint16, headerID uint32) (string, typeref.TypeRef, error) {
	return c.parser.FindTypeRefByID(serviceID, headerID)
}
