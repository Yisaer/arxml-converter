package parser

import (
	"fmt"

	"github.com/beevik/etree"

	"github.com/yisaer/arxml-converter/cp/parser/datatypes"
	"github.com/yisaer/arxml-converter/cp/parser/topology"
)

type Parser struct {
	Path                       string
	Doc                        *etree.Document
	dataTypesElement           *etree.Element
	dataTypeMappingSetsElement *etree.Element
	topologyElement            *etree.Element

	dataTypesParser *datatypes.DataTypesParser
	topologyParser  *topology.TopoLogyParser

	dataTypeMappings map[string]string
}

func NewParser(path string) (*Parser, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(path); err != nil {
		return nil, err
	}
	p := &Parser{
		Path: path, Doc: doc,
		dataTypeMappings: make(map[string]string),
	}
	return p, nil
}

func (p *Parser) Parse() error {
	autosar := p.Doc.SelectElement("AUTOSAR")
	if autosar == nil {
		return fmt.Errorf("no autosar")
	}
	arPackages := autosar.SelectElement("AR-PACKAGES")
	if arPackages == nil {
		return fmt.Errorf("no AR-PACKAGES found")
	}

	if err := p.search(arPackages); err != nil {
		return err
	}
	if err := p.parse(); err != nil {
		return err
	}
	return nil
}

func (p *Parser) parse() error {
	if err := p.parseDataTypeMappingSets(p.dataTypeMappingSetsElement); err != nil {
		return fmt.Errorf("parse dataTypeMappingSets: %w", err)
	}
	p.dataTypesParser = datatypes.NewDataTypesParser(p.dataTypeMappings)
	if err := p.dataTypesParser.ParseDataTypes(p.dataTypesElement); err != nil {
		return fmt.Errorf("parse dataTypes: %w", err)
	}
	p.topologyParser = topology.NewTopoLogyParser()
	if err := p.topologyParser.ParseTopoLogy(p.topologyElement); err != nil {
		return fmt.Errorf("parse topology: %w", err)
	}
	return nil
}
