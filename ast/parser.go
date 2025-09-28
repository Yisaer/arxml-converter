package ast

import (
	"fmt"

	"github.com/beevik/etree"

	"arxml-converter/mod"
)

type Parser struct {
	Path              string
	Doc               *etree.Document
	iautoSarElement   *etree.Element
	dataTypesElement  *etree.Element
	interfacesElement *etree.Element

	Interfaces map[string]*ServiceInterface
	DataTypes  map[string]*mod.DataType
	Services   map[int]*Service
}

func NewParser(path string) (*Parser, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(path); err != nil {
		return nil, err
	}
	p := &Parser{Path: path, Doc: doc}
	p.Interfaces = make(map[string]*ServiceInterface)
	p.DataTypes = make(map[string]*mod.DataType)
	p.Services = make(map[int]*Service)
	return p, nil
}

func (p *Parser) Parse() error {
	autoSar := p.Doc.SelectElement("AUTOSAR")
	if autoSar == nil {
		return fmt.Errorf("no autosar")
	}
	if err := p.search(autoSar); err != nil {
		return err
	}
	if err := p.parseInterfaces(); err != nil {
		return err
	}
	if err := p.parseDataTypes(); err != nil {
		return err
	}
	if err := p.parseIautoSar(); err != nil {
		return err
	}
	return nil
}

func (p *Parser) search(autoSar *etree.Element) error {
	arPackages := autoSar.SelectElement("AR-PACKAGES")
	if arPackages == nil {
		return fmt.Errorf("no ar-packages")
	}
	arPackageList := arPackages.SelectElements("AR-PACKAGE")
	if len(arPackageList) < 1 {
		return fmt.Errorf("no ar-packages")
	}
	if err := p.searchDataTypes(arPackageList); err != nil {
		return err
	}
	if err := p.searchIAutoSar(arPackageList); err != nil {
		return err
	}
	if err := p.searchInterfaces(arPackageList); err != nil {
		return err
	}
	return nil
}

func (p *Parser) searchIAutoSar(arPackageElements []*etree.Element) error {
	for _, arPackage := range arPackageElements {
		s := arPackage.SelectElement("SHORT-NAME")
		if s == nil {
			continue
		}
		if s.Text() == "IAUTOSAR" {
			p.iautoSarElement = arPackage
			return nil
		}
	}
	return fmt.Errorf("no IAUTOSAR find in ar package")
}

func (p *Parser) searchInterfaces(arPackageElements []*etree.Element) error {
	for _, arPackage := range arPackageElements {
		s := arPackage.SelectElement("SHORT-NAME")
		if s == nil {
			continue
		}
		if s.Text() == "interfaces" {
			p.interfacesElement = arPackage
			return nil
		}
	}
	return fmt.Errorf("no interfaces find in ar package")
}
