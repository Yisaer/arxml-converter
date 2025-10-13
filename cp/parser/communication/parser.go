package communication

import (
	"fmt"

	"github.com/beevik/etree"

	"github.com/yisaer/arxml-converter/util"
)

type CommunicationParser struct {
	pdusElement    *etree.Element
	signalsElement *etree.Element
	pduRefMap      map[string]string
	signalRef      map[string]string
}

func NewCommunicationParser() *CommunicationParser {
	return &CommunicationParser{
		pduRefMap: make(map[string]string),
		signalRef: make(map[string]string),
	}
}

func (p *CommunicationParser) GetPduRefMap() map[string]string {
	return p.pduRefMap
}

func (p *CommunicationParser) GetSignalRefMap() map[string]string {
	return p.signalRef
}

func (p *CommunicationParser) ParseCommunication(node *etree.Element) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("prase communication err: %v", err)
		}
	}()

	arpackagesElement, err := util.GetArPackagesElement(node)
	if err != nil {
		return err
	}
	arpackageList := arpackagesElement.SelectElements("AR-PACKAGE")
	if err := p.searchPDUsElement(arpackageList); err != nil {
		return err
	}
	if err := p.searchSignals(arpackageList); err != nil {
		return err
	}
	if err := p.parsePDUs(p.pdusElement); err != nil {
		return err
	}
	if err := p.parseSignals(p.signalsElement); err != nil {
		return err
	}
	return nil
}

func (p *CommunicationParser) searchPDUsElement(arpackageList []*etree.Element) error {
	for _, arpackage := range arpackageList {
		sn, err := util.GetShortname(arpackage)
		if err != nil {
			return err
		}
		if sn == "PDUs" {
			p.pdusElement = arpackage
			return nil
		}
	}
	return fmt.Errorf("no PDUs found")
}

func (p *CommunicationParser) searchSignals(arpackageList []*etree.Element) error {
	for _, arpackage := range arpackageList {
		sn, err := util.GetShortname(arpackage)
		if err != nil {
			return err
		}
		if sn == "Signals" {
			p.signalsElement = arpackage
			return nil
		}
	}
	return fmt.Errorf("no Signals found")
}

func (p *CommunicationParser) parsePDUs(node *etree.Element) error {
	elements, err := util.GetElements(node)
	if err != nil {
		return err
	}
	iSignalPDUList := elements.SelectElements("I-SIGNAL-I-PDU")
	for index, iSignalPDU := range iSignalPDUList {
		if err := p.parseiSignalPDU(iSignalPDU); err != nil {
			return fmt.Errorf("parse %v iSignalPDU err: %v", index, err)
		}
	}
	return nil
}

func (p *CommunicationParser) parseiSignalPDU(node *etree.Element) error {
	sn, err := util.GetShortname(node)
	if err != nil {
		return err
	}
	iSignalToPduMappingsElement := node.SelectElement("I-SIGNAL-TO-PDU-MAPPINGS")
	if iSignalToPduMappingsElement == nil {
		return nil
	}
	iSignalToIPDUMapping := iSignalToPduMappingsElement.SelectElement("I-SIGNAL-TO-I-PDU-MAPPING")
	if iSignalToIPDUMapping == nil {
		return nil
	}
	iSignalRef := iSignalToIPDUMapping.SelectElement("I-SIGNAL-REF")
	if iSignalRef == nil {
		return nil
	}
	p.pduRefMap[sn] = iSignalRef.Text()
	return nil
}

func (p *CommunicationParser) parseSignals(node *etree.Element) error {
	elements, err := util.GetElements(node)
	if err != nil {
		return err
	}
	iSignalList := elements.SelectElements("I-SIGNAL")
	for index, iSignal := range iSignalList {
		if err := p.parseISignal(iSignal); err != nil {
			return fmt.Errorf("parse %v iSignal err: %v", index, err)
		}
	}
	return nil
}

func (p *CommunicationParser) parseISignal(node *etree.Element) error {
	sn, err := util.GetShortname(node)
	if err != nil {
		return nil
	}
	systemSignalRefElement := node.SelectElement("SYSTEM-SIGNAL-REF")
	if systemSignalRefElement == nil {
		return nil
	}
	p.signalRef[sn] = systemSignalRefElement.Text()
	return nil
}
