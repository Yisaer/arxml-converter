package ast

import (
	"fmt"

	"github.com/beevik/etree"
)

func (dp *DataTypesParser) parseDataTypes(root *etree.Element) error {
	arpackagesElement := root.SelectElement("AR-PACKAGES")
	if arpackagesElement == nil {
		return fmt.Errorf("AR-PACKAGES element not found")
	}
	arpackages := arpackagesElement.SelectElements("AR-PACKAGE")
	for index, arpkg := range arpackages {
		shortname, err := dp.getShortname(arpkg)
		if err != nil {
			return fmt.Errorf("could not get short name for ar-packages[%d]", index)
		}
		switch shortname {
		case "ImplementationDataTypes":
			dp.implementationDataTypesArPackage = arpkg
		case "ApplicationDataType":
			dp.applicationDatatypeArPackage = arpkg
		}
	}
	if dp.implementationDataTypesArPackage == nil {
		return fmt.Errorf("no implementationDataTypes found in AR-PACKAGES")
	}
	if dp.applicationDatatypeArPackage == nil {
		return fmt.Errorf("no applicationDataTypes found in AR-PACKAGES")
	}

	if err := dp.parseImplementationDataTypes(dp.implementationDataTypesArPackage); err != nil {
		return fmt.Errorf("parse ImplementationDataTypes failed, err:%v", err.Error())
	}
	if err := dp.parseApplicationDatatypes(dp.applicationDatatypeArPackage); err != nil {
		return fmt.Errorf("parse application data types: %w", err)
	}
	return nil
}
