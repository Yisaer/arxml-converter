package datatypes

import (
	"fmt"

	"github.com/beevik/etree"
)

func (dp *DataTypesParser) ParseDataTypes(root *etree.Element) error {
	arpackagesElement := root.SelectElement("AR-PACKAGES")
	if arpackagesElement == nil {
		return fmt.Errorf("AR-PACKAGES element not found")
	}
	if err := dp.parseImplementationDataTypes(arpackagesElement); err != nil {
		return fmt.Errorf("parse ImplementationDataTypes failed, err:%v", err.Error())
	}
	if err := dp.parseApplicationDatatypes(arpackagesElement); err != nil {
		return fmt.Errorf("parse application data types: %w", err)
	}
	return nil
}
