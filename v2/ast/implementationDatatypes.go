package ast

import (
	"fmt"

	"github.com/beevik/etree"

	"github.com/yisaer/arxml-converter/mod"
)

func (p *Parser) parseImplementationDataTypes(root *etree.Element) error {
	elements := root.SelectElement("ELEMENTS")
	if elements == nil {
		return fmt.Errorf("no elements found")
	}
	idts := elements.SelectElements("IMPLEMENTATION-DATA-TYPE")
	if len(idts) < 1 {
		return fmt.Errorf("no IMPLEMENTATION-DATA-TYPE found")
	}
	for index, idt := range idts {
		if err := p.parseImplementationDataType(idt); err != nil {
			return fmt.Errorf("parse %v ImplementationDataType failed, err:%v", index, err.Error())
		}
	}
	return nil
}

func (p *Parser) parseImplementationDataType(root *etree.Element) (err error) {
	sn, err := p.getShortname(root)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = fmt.Errorf("shortname:%v,err:%v", sn, err.Error())
		}
	}()

	category, err := p.getCategory(root)
	if err != nil {
		return err
	}
	if category != "VALUE" {
		return nil
	}
	sddpc, err := p.getSWDataDefPropsConditional(root)
	if err != nil {
		return err
	}
	byr := sddpc.SelectElement("BASE-TYPE-REF")
	if byr == nil {
		return fmt.Errorf("no BASE-TYPE-REF found")
	}
	ref := byr.Text()
	if err := p.validBasicType(ref); err != nil {
		return err
	}
	p.implementationDataTypes[sn] = mod.NewBasicDataType(sn, category, ref)
	return nil
}
