package ast

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
)

func (p *Parser) parseApplicationDatatypes(root *etree.Element) error {
	elements, err := p.getElements(root)
	if err != nil {
		return err
	}
	apdts := elements.SelectElements("APPLICATION-PRIMITIVE-DATA-TYPE")
	for index, apdt := range apdts {
		if err := p.ParseApplicationPrimitiveDataType(apdt); err != nil {
			return fmt.Errorf("parse index %v APPLICATION-PRIMITIVE-DATA-TYPE failed, err:%v", index, err.Error())
		}
	}
	return nil
}

func (p *Parser) ParseApplicationPrimitiveDataType(root *etree.Element) (err error) {
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
		return fmt.Errorf("parse category failed err:%v", err.Error())
	}
	switch category {
	case "STRING":
		sddpc, err := p.getSWDataDefPropsConditional(root)
		if err != nil {
			return err
		}
		stp := sddpc.SelectElement("SW-TEXT-PROPS")
		if stp == nil {
			return fmt.Errorf("no SW-TEXT-PROPS found")
		}
		isDynamicString, err := p.getArraySizeSemantics(stp)
		if err != nil {
			return err
		}
		if !isDynamicString {
			return fmt.Errorf("fixed length string not supported now")
		}
		btr := stp.SelectElement("BASE-TYPE-REF")
		if btr == nil {
			return fmt.Errorf("no BASE-TYPE-REF found")
		}
		if !strings.Contains(btr.Text(), "UTF_8") {
			return fmt.Errorf("BASE-TYPE ref should be UTF_8")
		}

	case "VALUE":
	default:
		return fmt.Errorf("unknown category:%v", category)
	}
	return nil
}
