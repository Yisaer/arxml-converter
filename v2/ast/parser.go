package ast

import (
	"fmt"

	"github.com/beevik/etree"
)

type Parser struct {
	Path                       string
	Doc                        *etree.Document
	dataTypesElement           *etree.Element
	dataTypeMappingSetsElement *etree.Element

	dataTypesParser *DataTypesParser

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
	p.dataTypesParser = NewDataTypesParser(p)
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

func (p *Parser) searchDataTypes(arPackagesElement *etree.Element) error {
	arPackages := arPackagesElement.SelectElements("AR-PACKAGE")
	for _, arPackage := range arPackages {
		sn, err := p.getShortname(arPackage)
		if err != nil {
			return err
		}
		if sn == "DataTypes" {
			p.dataTypesElement = arPackage
			return nil
		}
	}
	return fmt.Errorf("no DataTypes found")
}

func (p *Parser) searchDataTypeMappingSets(arPackagesElement *etree.Element) error {
	arPackages := arPackagesElement.SelectElements("AR-PACKAGE")
	for _, arPackage := range arPackages {
		sn, err := p.getShortname(arPackage)
		if err != nil {
			return err
		}
		if sn == "DataTypeMappingSets" {
			p.dataTypeMappingSetsElement = arPackage
			return nil
		}
	}
	return fmt.Errorf("no DataTypes found")
}

func (p *Parser) search(arPackages *etree.Element) error {
	if err := p.searchDataTypes(arPackages); err != nil {
		return fmt.Errorf("search data types: %w", err)
	}
	if err := p.searchDataTypeMappingSets(arPackages); err != nil {
		return fmt.Errorf("search data types mappings: %w", err)
	}
	return nil
}

func (p *Parser) parse() error {
	if err := p.parseDataTypeMappingSets(p.dataTypeMappingSetsElement); err != nil {
		return fmt.Errorf("parse dataTypeMappingSets: %w", err)
	}
	if err := p.dataTypesParser.parseDataTypes(p.dataTypesElement); err != nil {
		return fmt.Errorf("parse dataTypes: %w", err)
	}
	return nil
}
