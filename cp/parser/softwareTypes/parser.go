package softwareTypes

import (
	"fmt"

	"github.com/beevik/etree"

	"github.com/yisaer/arxml-converter/util"
)

type SoftwareTypesParser struct {
	interfacesElement                *etree.Element
	interfaceRefMap                  map[string]string
	directionInArgumentDataPrototype *etree.Element
}

func NewSoftwareTypesParser() *SoftwareTypesParser {
	return &SoftwareTypesParser{
		interfaceRefMap: make(map[string]string),
	}
}

func (sp *SoftwareTypesParser) GetInterfaceRefMap() map[string]string {
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
		if err := sp.searchClientServerInterface(clientServerInterfaceElement); err != nil {
			return fmt.Errorf("parsing %v client server interface : %w", index, err)
		}
	}
	return nil
}

func (sp *SoftwareTypesParser) searchClientServerInterface(node *etree.Element) (err error) {
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
	clientServerOperation := operationsElement.SelectElement("CLIENT-SERVER-OPERATION")
	if clientServerOperation == nil {
		return nil
	}
	return sp.parseClientServerOperation(sn, clientServerOperation)
}

func (sp *SoftwareTypesParser) parseClientServerOperation(sn string, node *etree.Element) (err error) {
	argumentsElement := node.SelectElement("ARGUMENTS")
	if argumentsElement == nil {
		return nil
	}
	argumentDataProtoTypeList := argumentsElement.SelectElements("ARGUMENT-DATA-PROTOTYPE")
	if err := sp.searchDirectionINArgument(argumentDataProtoTypeList); err != nil {
		return err
	}
	typeRefElement := sp.directionInArgumentDataPrototype.SelectElement("TYPE-TREF")
	if typeRefElement == nil {
		return fmt.Errorf("no type reference element found for %v in argument direction in", sn)
	}
	sp.interfaceRefMap[sn] = typeRefElement.Text()
	return nil
}

func (sp *SoftwareTypesParser) searchDirectionINArgument(argumentDataProtoTypeList []*etree.Element) error {
	for _, argumentDataProtoTypeElement := range argumentDataProtoTypeList {
		DIRECTIONElement := argumentDataProtoTypeElement.SelectElement("DIRECTION")
		if DIRECTIONElement == nil {
			return nil
		}
		if DIRECTIONElement.Text() == "IN" {
			sp.directionInArgumentDataPrototype = argumentDataProtoTypeElement
			return nil
		}
	}
	return fmt.Errorf("no DIRECTION IN argument found")
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
