package ast

import (
	"fmt"

	"github.com/beevik/etree"
)

func (p *Parser) parseDataTypes(root *etree.Element) error {
	arpackages := root.SelectElements("AR-PACKAGES")
	if len(arpackages) < 1 {
		return fmt.Errorf("no ar-packages found")
	}
	for index, arpkg := range arpackages {
		shortname, err := p.getShortname(arpkg)
		if err != nil {
			return fmt.Errorf("could not get short name for ar-packages[%d]", index)
		}
		switch shortname {
		case "ImplementationDataTypes":
			p.implementationDataTypesArPackage = arpkg
		case "ApplicationDataTypes":
			p.applicationDatatypeArPackage = arpkg
		case "DataTypeMappingSets":
			p.dataTypeMappingSetsArPackage = arpkg
		}
	}
	if err := p.parseDataTypeMappingSets(p.dataTypeMappingSetsArPackage); err != nil {
		return fmt.Errorf("parse dataTypeMappingSets failed, err:%v", err.Error())
	}
	if err := p.parseImplementationDataTypes(p.implementationDataTypesArPackage); err != nil {
		return fmt.Errorf("parse ImplementationDataTypes failed, err:%v", err.Error())
	}
	if err := p.parseApplicationDatatypes(p.applicationDatatypeArPackage); err != nil {
		return fmt.Errorf("parse application data types: %w", err)
	}
	return nil
}
