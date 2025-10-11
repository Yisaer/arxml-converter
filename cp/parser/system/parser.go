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
	if systemSN != "SystemDescription" {
		return fmt.Errorf("system shortname is not SystemDescription, got %s", systemSN)
	}
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
			return fmt.Errorf("parse %v CLIENTSERVERTOSIGNALMAPPING error: %s", index, err.Error())
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
	sp.operationRef[callSignalRefElement.Text()] = targetOperationRefElement.Text()
	return nil
}
