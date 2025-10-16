package parser

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
	idlAst "github.com/yisaer/idl-parser/ast"
	"github.com/yisaer/idl-parser/ast/typeref"

	"github.com/yisaer/arxml-converter/ast"
	"github.com/yisaer/arxml-converter/cp/parser/communication"
	"github.com/yisaer/arxml-converter/cp/parser/datatypes"
	"github.com/yisaer/arxml-converter/cp/parser/softwareTypes"
	"github.com/yisaer/arxml-converter/cp/parser/system"
	"github.com/yisaer/arxml-converter/cp/parser/topology"
)

type Parser struct {
	Path                       string
	Doc                        *etree.Document
	dataTypesElement           *etree.Element
	dataTypeMappingSetsElement *etree.Element
	topologyElement            *etree.Element
	communicationElement       *etree.Element
	systemElement              *etree.Element
	softwareTypesElement       *etree.Element

	dataTypesParser     *datatypes.DataTypesParser
	topologyParser      *topology.TopoLogyParser
	communicationParser *communication.CommunicationParser
	systemParser        *system.SystemParser
	softwareTypesParser *softwareTypes.SoftwareTypesParser

	dataTypeMappings map[string]string

	transformer *ast.TransformHelper
	idlModule   *idlAst.Module
}

func NewParser(path string) (*Parser, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(path); err != nil {
		return nil, err
	}
	p := &Parser{
		Path: path, Doc: doc,
		dataTypeMappings: make(map[string]string),
	}
	return p, nil
}

func (p *Parser) Parse() error {
	autosar := p.Doc.SelectElement("AUTOSAR")
	if autosar == nil {
		return fmt.Errorf("no autosar")
	}
	arPackages := autosar.SelectElement("AR-PACKAGES")
	if arPackages == nil {
		return fmt.Errorf("no AR-PACKAGES found")
	}

	if err := p.search(arPackages); err != nil {
		return err
	}
	if err := p.parse(); err != nil {
		return err
	}
	p.transformer = ast.NewTransformHelper(p.dataTypesParser.GetApplicationDataTypes())
	m, err := p.transformer.TransformIntoModule()
	if err != nil {
		return fmt.Errorf("transform error: %s", err)
	}
	p.idlModule = m
	return nil
}

func (p *Parser) parse() error {
	if err := p.parseDataTypeMappingSets(p.dataTypeMappingSetsElement); err != nil {
		return fmt.Errorf("parse dataTypeMappingSets: %w", err)
	}
	p.dataTypesParser = datatypes.NewDataTypesParser(p.dataTypeMappings)
	if err := p.dataTypesParser.ParseDataTypes(p.dataTypesElement); err != nil {
		return fmt.Errorf("parse dataTypes: %w", err)
	}
	p.topologyParser = topology.NewTopoLogyParser()
	if err := p.topologyParser.ParseTopoLogy(p.topologyElement); err != nil {
		return fmt.Errorf("parse topology: %w", err)
	}
	p.communicationParser = communication.NewCommunicationParser()
	if err := p.communicationParser.ParseCommunication(p.communicationElement); err != nil {
		return fmt.Errorf("parse communication: %w", err)
	}
	p.systemParser = system.NewSystemParser()
	if err := p.systemParser.ParseSystem(p.systemElement); err != nil {
		return fmt.Errorf("parse system: %w", err)
	}
	p.softwareTypesParser = softwareTypes.NewSoftwareTypesParser()
	if err := p.softwareTypesParser.ParseSoftwareTypes(p.softwareTypesElement); err != nil {
		return fmt.Errorf("parse softwareTypes: %w", err)
	}
	return nil
}

func (p *Parser) FindTypeRefByID(serviceID uint16, headerID uint32) (string, typeref.TypeRef, error) {
	serviceIDMap := p.topologyParser.GetServiceIDMap()
	_, ok := serviceIDMap[serviceID]
	if !ok {
		return "", nil, fmt.Errorf("no service found for %d", serviceID)
	}
	headerIDMap := p.topologyParser.GetHeaderRef()
	pduRef, ok := headerIDMap[headerID]
	if !ok {
		return "", nil, fmt.Errorf("no header ref for %d", headerID)
	}
	pduTriggeringRef := p.topologyParser.GetPDUTriggeringRef()
	pduTriggering, ok := pduTriggeringRef[extractLast(pduRef)]
	if !ok {
		return "", nil, fmt.Errorf("no pdu triggered for %v", pduRef)
	}
	communicationPDURefMap := p.communicationParser.GetPduRefMap()
	communicationPduRef, ok := communicationPDURefMap[extractLast(pduTriggering)]
	if !ok {
		return "", nil, fmt.Errorf("no pdu triggering ref for %v", pduTriggering)
	}
	systemSignalRef, ok := p.communicationParser.GetSignalRefMap()[extractLast(communicationPduRef)]
	if !ok {
		return "", nil, fmt.Errorf("no signal ref for %v", communicationPduRef)
	}

	find := false
	operationRef := ""
	for k, v := range p.systemParser.GetOperationRef() {
		if strings.HasSuffix(k, systemSignalRef) {
			find = true
			operationRef = v
		}
	}
	if !find {
		return "", nil, fmt.Errorf("no operation ref for %v", communicationPduRef)
	}
	InterfaceRefMap := p.softwareTypesParser.GetInterfaceRefMap()
	interfaceKey, err := extractLast2(operationRef)
	if err != nil {
		return "", nil, err
	}
	InterfaceRef, ok := InterfaceRefMap[interfaceKey]
	if !ok {
		return "", nil, fmt.Errorf("no interface ref for %v", operationRef)
	}
	tr, ok := p.transformer.GetConverterRef()[strings.ToLower(extractLast(InterfaceRef))]
	if !ok {
		return "", nil, fmt.Errorf("no converter ref for %v", InterfaceRef)
	}
	return extractLast(InterfaceRef), tr, nil
}

func (p *Parser) GetModule() *idlAst.Module {
	return p.idlModule
}

func extractLast(ref string) string {
	parts := strings.Split(ref, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ref
}

func extractLast2(ref string) (string, error) {
	parts := strings.Split(ref, "/")
	if len(parts) > 1 {
		return parts[len(parts)-2], nil
	}
	return "", fmt.Errorf("no last 2 element in %v", ref)
}
