package topology

import (
	"fmt"

	"github.com/beevik/etree"

	"github.com/yisaer/arxml-converter/util"
)

func (tp *TopoLogyParser) parsePDUTRIGGERINGS(node *etree.Element) error {
	pduTriggeringList := node.SelectElements("PDU-TRIGGERING")
	for index, pduTriggeringElement := range pduTriggeringList {
		if err := tp.parsePDUTRIGGERING(pduTriggeringElement); err != nil {
			return fmt.Errorf("parse %v pdu trigger err:%v", index, err)
		}
	}
	return nil
}

func (tp *TopoLogyParser) parsePDUTRIGGERING(node *etree.Element) error {
	sn, err := util.GetShortname(node)
	if err != nil {
		return err
	}
	iPDURefElement := node.SelectElement("I-PDU-REF")
	if iPDURefElement == nil {
		return fmt.Errorf("I-PDU-REF element not found")
	}
	tp.pduTriggeringRef[sn] = iPDURefElement.Text()
	return nil
}
