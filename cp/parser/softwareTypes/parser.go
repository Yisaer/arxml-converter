package softwareTypes

import (
	"fmt"

	"github.com/beevik/etree"

	"github.com/yisaer/arxml-converter/util"
)

type SoftwareTypesParser struct {
	interfacesElement *etree.Element
	interfaceRefMap   map[string]map[string]string
}

func NewSoftwareTypesParser() *SoftwareTypesParser {
	return &SoftwareTypesParser{
		interfaceRefMap: make(map[string]map[string]string),
	}
}

func (sp *SoftwareTypesParser) GetInterfaceRefMap() map[string]map[string]string {
	return sp.interfaceRefMap
}

func (sp *SoftwareTypesParser) ParseSoftwareTypes(node *etree.Element) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("parsing software types: %w", err)
		}
	}()

	arpackagesElement, err := util.GetArPackagesElement(node)
	if err != nil {
		return err
	}
	arpackageList := arpackagesElement.SelectElements("AR-PACKAGE")
	if err := sp.searchInterfaces(arpackageList); err != nil {
		return err
	}
	if err := sp.parseInterfaces(sp.interfacesElement); err != nil {
		return err
	}
	return nil
}

func (sp *SoftwareTypesParser) parseInterfaces(node *etree.Element) error {
	elements, err := util.GetElements(node)
	if err != nil {
		return err
	}
	clientServerInterfaceList := elements.SelectElements("CLIENT-SERVER-INTERFACE")
	for index, clientServerInterfaceElement := range clientServerInterfaceList {
		if err := sp.parseClientServerInterface(clientServerInterfaceElement); err != nil {
			return fmt.Errorf("parsing %v client server interface : %w", index, err)
		}
	}
	return nil
}

func (sp *SoftwareTypesParser) parseClientServerInterface(node *etree.Element) (err error) {
	sn, err := util.GetShortname(node)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = fmt.Errorf("searching client server interface %v: %w", sn, err)
		}
	}()
	operationsElement := node.SelectElement("OPERATIONS")
	if operationsElement == nil {
		return nil
	}
	for index, cso := range operationsElement.SelectElements("CLIENT-SERVER-OPERATION") {
		csoShortName, tref, err := sp.parseClientServerOperation(cso)
		if err != nil {
			return fmt.Errorf("parsing %v client server operation: %w", index, err)
		}
		sp.addClientServerInterfaceMap(sn, csoShortName, tref)
	}
	return nil
}

func (sp *SoftwareTypesParser) addClientServerInterfaceMap(csiShortname, csoShortname, tref string) {
	csoMap, ok := sp.interfaceRefMap[csiShortname]
	if !ok {
		csoMap = make(map[string]string)
		sp.interfaceRefMap[csiShortname] = csoMap
	}
	csoMap[csoShortname] = tref
	sp.interfaceRefMap[csiShortname] = csoMap
}

func (sp *SoftwareTypesParser) parseClientServerOperation(node *etree.Element) (shortname, tref string, err error) {
	sn, err := util.GetShortname(node)
	if err != nil {
		return "", "", err
	}
	defer func() {
		if err != nil {
			err = fmt.Errorf("parsing client server operation %v: %w", sn, err)
		}
	}()

	argumentsElement := node.SelectElement("ARGUMENTS")
	if argumentsElement == nil {
		return "", "", nil
	}
	for index, argument := range argumentsElement.SelectElements("ARGUMENT-DATA-PROTOTYPE") {
		DIRECTIONElement := argument.SelectElement("DIRECTION")
		if DIRECTIONElement == nil {
			continue
		}
		if DIRECTIONElement.Text() == "IN" {
			typeRefElement := argument.SelectElement("TYPE-TREF")
			if typeRefElement == nil {
				return "", "", fmt.Errorf("parsing %v ARGUMENT-DATA-PROTOTYPE no TYPE-TREF", index)
			}
			return sn, typeRefElement.Text(), nil
		}
	}
	return "", "", nil
}

func (sp *SoftwareTypesParser) searchInterfaces(arpackageList []*etree.Element) error {
	for _, arpackage := range arpackageList {
		sn, err := util.GetShortname(arpackage)
		if err != nil {
			return err
		}
		if sn == "Interfaces" {
			sp.interfacesElement = arpackage
			return nil
		}
	}
	return fmt.Errorf("could not find interfaces")
}
