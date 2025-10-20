package system

import (
	"fmt"

	"github.com/beevik/etree"

	"github.com/yisaer/arxml-converter/util"
)

type SystemParser struct {
	operationRef map[string]string
}

func NewSystemParser() *SystemParser {
	return &SystemParser{
		operationRef: make(map[string]string),
	}
}

func (sp *SystemParser) GetOperationRef() map[string]string {
	return sp.operationRef
}

func (sp *SystemParser) ParseSystem(node *etree.Element) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("parse system error: %s", err.Error())
		}
	}()
	elements, err := util.GetElements(node)
	if err != nil {
		return err
	}
	systemElement := elements.SelectElement("SYSTEM")
	if systemElement == nil {
		return fmt.Errorf("no system")
	}
	systemSN, err := util.GetShortname(systemElement)
	if err != nil {
		return err
	}
	if systemSN == "SystemDescription" {
		return sp.parseSystemMapping(systemElement)
	}
	categoryElement := systemElement.SelectElement("CATEGORY")
	if categoryElement != nil && categoryElement.Text() == "SYSTEM_DESCRIPTION" {
		return sp.parseSystemMapping(systemElement)
	}
	return fmt.Errorf("unsupported system type: %s", systemSN)
}

func (sp *SystemParser) parseSystemMapping(systemElement *etree.Element) (err error) {
	mappingsElement := systemElement.SelectElement("MAPPINGS")
	if mappingsElement == nil {
		return fmt.Errorf("no mappings")
	}
	systemMappingElement := mappingsElement.SelectElement("SYSTEM-MAPPING")
	if systemMappingElement == nil {
		return fmt.Errorf("no SYSTEM-MAPPING ")
	}
	dataMappingsElement := systemMappingElement.SelectElement("DATA-MAPPINGS")
	if dataMappingsElement == nil {
		return fmt.Errorf("no DATA-MAPPINGS")
	}

	clientServerToSignalMappingList := dataMappingsElement.SelectElements("CLIENT-SERVER-TO-SIGNAL-MAPPING")
	for index, clientServerToSignalMappingElement := range clientServerToSignalMappingList {
		if err := sp.parseCLIENTSERVERTOSIGNALMAPPING(clientServerToSignalMappingElement); err != nil {
			return fmt.Errorf("parse %v CLIENT-SERVER-TO-SIGNAL-MAPPING error: %s", index, err.Error())
		}
	}
	SENDERRECEIVERTOSIGNALMAPPINGList := dataMappingsElement.SelectElements("SENDER-RECEIVER-TO-SIGNAL-MAPPING")
	for index, SENDERRECEIVERTOSIGNALMAPPINGElement := range SENDERRECEIVERTOSIGNALMAPPINGList {
		if err := sp.paraseSENDERRECEIVERTOSIGNALMAPPING(SENDERRECEIVERTOSIGNALMAPPINGElement); err != nil {
			return fmt.Errorf("parse %v SENDER-RECEIVER-TO-SIGNAL-MAPPING error: %s", index, err.Error())
		}
	}
	return nil
}

func (sp *SystemParser) parseCLIENTSERVERTOSIGNALMAPPING(node *etree.Element) (err error) {
	callSignalRefElement := node.SelectElement("CALL-SIGNAL-REF")
	if callSignalRefElement == nil {
		return nil
	}
	clientServerOperationIRefElement := node.SelectElement("CLIENT-SERVER-OPERATION-IREF")
	if clientServerOperationIRefElement == nil {
		return nil
	}
	targetOperationRefElement := clientServerOperationIRefElement.SelectElement("TARGET-OPERATION-REF")
	if targetOperationRefElement == nil {
		return nil
	}
	a := targetOperationRefElement.Text()
	sp.operationRef[callSignalRefElement.Text()] = a
	return nil
}

func (sp *SystemParser) paraseSENDERRECEIVERTOSIGNALMAPPING(node *etree.Element) (err error) {
	srElement := node.SelectElement("SYSTEM-SIGNAL-REF")
	if srElement == nil {
		return nil
	}
	DATAELEMENTIREF := node.SelectElement("DATA-ELEMENT-IREF")
	if DATAELEMENTIREF == nil {
		return nil
	}
	TARGETDATAPROTOTYPEREF := DATAELEMENTIREF.SelectElement("TARGET-DATA-PROTOTYPE-REF")
	if TARGETDATAPROTOTYPEREF == nil {
		return nil
	}
	sp.operationRef[srElement.Text()] = TARGETDATAPROTOTYPEREF.Text()
	return nil
}
