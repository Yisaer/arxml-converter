package parser

import (
	"fmt"
	"strings"

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

var (
	appDataTypePrefix       = "/DataTypes/ApplicationDataType/"
	implementDataTypePrefix = "/DataTypes/ImplementationDataTypes/"
)

func (p *Parser) parseSubDtm(subdtm *etree.Element) error {
	adtr := subdtm.SelectElement("APPLICATION-DATA-TYPE-REF")
	if adtr == nil {
		return fmt.Errorf("no APPLICATION-DATA-TYPE-REF found")
	}
	if !strings.HasPrefix(adtr.Text(), appDataTypePrefix) {
		return fmt.Errorf("invalid APPLICATION-DATA-TYPE-REF:%v", adtr.Text())
	}
	adtrKey := strings.TrimPrefix(adtr.Text(), appDataTypePrefix)
	idtr := subdtm.SelectElement("IMPLEMENTATION-DATA-TYPE-REF")
	if idtr == nil {
		return fmt.Errorf("no IMPLEMENTATION-DATA-TYPE-REF found")
	}
	if !strings.HasPrefix(idtr.Text(), implementDataTypePrefix) {
		return fmt.Errorf("invalid IMPLEMENTATION-DATA-TYPE-REF:%v", idtr.Text())
	}
	idtrKey := strings.TrimPrefix(idtr.Text(), implementDataTypePrefix)
	p.dataTypeMappings[adtrKey] = idtrKey
	return nil
}
