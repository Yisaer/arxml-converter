package parser

import (
	"fmt"

	"github.com/beevik/etree"

	"github.com/yisaer/arxml-converter/util"
)

func (p *Parser) parseDataTypeMappingSets(node *etree.Element) error {
	elements, err := util.GetElements(node)
	if err != nil {
		return err
	}
	dtms := elements.SelectElement("DATA-TYPE-MAPPING-SET")
	if dtms == nil {
		return fmt.Errorf("no DATA-TYPE-MAPPING-SET found")
	}
	sn, err := util.GetShortname(dtms)
	if err != nil {
		return err
	}
	if sn != "Data_Type_Mappings" {
		return fmt.Errorf("no Data_Type_Mappings found")
	}
	dtm := dtms.SelectElement("DATA-TYPE-MAPS")
	if dtm == nil {
		return fmt.Errorf("no DATA-TYPE-MAPS found")
	}
	subdtms := dtm.SelectElements("DATA-TYPE-MAP")
	for index, subdtm := range subdtms {
		if err := p.parseSubDtm(subdtm); err != nil {
			return fmt.Errorf("parse %v DATA-TYPE-MAP failed: %w", index, err)
		}
	}
	return nil
}

func (p *Parser) parseSubDtm(subdtm *etree.Element) error {
	adtr := subdtm.SelectElement("APPLICATION-DATA-TYPE-REF")
	if adtr == nil {
		return fmt.Errorf("no APPLICATION-DATA-TYPE-REF found")
	}
	adtrKey := extractLast(adtr.Text())
	idtr := subdtm.SelectElement("IMPLEMENTATION-DATA-TYPE-REF")
	if idtr == nil {
		return fmt.Errorf("no IMPLEMENTATION-DATA-TYPE-REF found")
	}
	idtrKey := extractLast(idtr.Text())
	p.dataTypeMappings[adtrKey] = idtrKey
	return nil
}
