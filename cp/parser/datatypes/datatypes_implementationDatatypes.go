package datatypes

import (
	"fmt"

	"github.com/beevik/etree"

	"github.com/yisaer/arxml-converter/ast"
	"github.com/yisaer/arxml-converter/util"
)

func (dp *DataTypesParser) parseImplementationDataTypes(node *etree.Element) error {
	for index, idt := range node.FindElements("//IMPLEMENTATION-DATA-TYPE") {
		if err := dp.parseImplementationValueDataType(idt); err != nil {
			return fmt.Errorf("parse %v ImplementationDataType failed, err:%v", index, err.Error())
		}
	}
	return nil
}

func (dp *DataTypesParser) parseImplementationValueDataType(root *etree.Element) (err error) {
	sn, err := util.GetShortname(root)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = fmt.Errorf("shortname:%v,err:%v", sn, err.Error())
		}
	}()

	category, err := util.GetCategory(root)
	if err != nil {
		return err
	}
	if category != "VALUE" {
		return nil
	}
	sddpc, err := util.GetSWDataDefPropsConditional(root)
	if err != nil {
		return err
	}
	byr := sddpc.SelectElement("BASE-TYPE-REF")
	if byr == nil {
		return fmt.Errorf("no BASE-TYPE-REF found")
	}
	ref := byr.Text()
	if err := util.ValidBasicType(ref); err != nil {
		return err
	}
	dp.implementationDataTypes[sn] = ast.NewBasicDataType(sn, category, ref)
	return nil
}
