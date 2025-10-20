package converter

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
	"github.com/yisaer/idl-parser/ast/typeref"
	"github.com/yisaer/idl-parser/converter"

	apconverter "github.com/yisaer/arxml-converter/ap/converter"
	cpconverter "github.com/yisaer/arxml-converter/cp/converter"
	"github.com/yisaer/arxml-converter/util"
)

type AutosarXsdVersion int

const (
	AUTOSAR_00048Version AutosarXsdVersion = iota
	AUTOSAR_4_2_2Version
)

type ArxmlConverter struct {
	path             string
	config           converter.IDlConverterConfig
	cpArxmlConverter *cpconverter.ArxmlCPConverter
	apArxmlConverter *apconverter.ArXMLConverter
	doc              *etree.Document
	version          AutosarXsdVersion
}

func NewConverter(path string, config converter.IDlConverterConfig) (*ArxmlConverter, error) {
	var err error
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(path); err != nil {
		return nil, err
	}
	doc.SelectElement("AUTOSAR")
	c := &ArxmlConverter{
		path:   path,
		config: config,
	}
	if err := c.parseXML(path); err != nil {
		return nil, err
	}
	isCp, err1 := c.IsCP()
	if err1 != nil {
		return nil, err1
	}
	isAp, err2 := c.IsAP()
	if err2 != nil {
		return nil, err2
	}
	if isCp {
		c.cpArxmlConverter, err = cpconverter.NewArxmlCPConverterWithDoc(doc, config)
		if err != nil {
			return nil, err
		}
		return c, nil
	}
	if isAp {
		c.apArxmlConverter, err = apconverter.NewConverterWithDoc(doc, config)
		if err != nil {
			return nil, err
		}
		return c, nil
	}

	return nil, fmt.Errorf("target arxml isn't cp or ap")
}

func (c *ArxmlConverter) GetDataTypeByID(serviceID uint16, eventID uint16) (string, typeref.TypeRef, error) {
	if c.apArxmlConverter != nil {
		return c.apArxmlConverter.GetTypeByID(int(serviceID), int(eventID))
	}
	if c.cpArxmlConverter != nil {
		return c.cpArxmlConverter.GetDataTypeByID(serviceID, MergeUint16ToUint32(serviceID, eventID))
	}
	return "", nil, fmt.Errorf("target arxml isn't cp or ap")
}

func (c *ArxmlConverter) Decode(serviceID uint16, eventID uint16, data []byte) (string, interface{}, error) {
	if c.apArxmlConverter != nil {
		return c.apArxmlConverter.DecodeWithID(int(serviceID), int(eventID), data)
	}
	if c.cpArxmlConverter != nil {
		return c.cpArxmlConverter.Convert(serviceID, MergeUint16ToUint32(serviceID, eventID), data)
	}
	return "", nil, fmt.Errorf("no converter found")
}

func MergeUint16ToUint32(high16, low16 uint16) uint32 {
	return uint32(high16)<<16 | uint32(low16)
}

func (c *ArxmlConverter) parseXML(path string) error {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(path); err != nil {
		return err
	}
	autosarElement := doc.SelectElement("AUTOSAR")
	if autosarElement == nil {
		return fmt.Errorf("autosar element not found")
	}
	c.doc = doc

	// 检查 xsi:schemaLocation 属性
	schemaLocation := autosarElement.SelectAttr("xsi:schemaLocation")
	if schemaLocation == nil {
		return fmt.Errorf("xsi:schemaLocation attribute not found")
	}
	schemaLocationValue := schemaLocation.Value
	switch {
	case strings.Contains(schemaLocationValue, "AUTOSAR_00048.xsd"):
		c.version = AUTOSAR_00048Version
	case strings.Contains(schemaLocationValue, "AUTOSAR_4-2-2.xsd"):
		c.version = AUTOSAR_4_2_2Version
	default:
		return fmt.Errorf("unkown schema location %s", schemaLocationValue)
	}
	return nil
}

func (c *ArxmlConverter) IsAP() (bool, error) {
	autosarElement := c.doc.SelectElement("AUTOSAR")
	if autosarElement == nil {
		return false, nil
	}
	var interfacesElement *etree.Element
	var datatypes *etree.Element
	var IAUTOSAR *etree.Element
	arpackagesElement := autosarElement.SelectElement("AR-PACKAGES")
	if arpackagesElement == nil {
		return false, nil
	}
	arpList := arpackagesElement.SelectElements("AR-PACKAGE")
	for _, arg := range arpList {
		sn, err := util.GetShortname(arg)
		if err != nil {
			return false, err
		}
		switch sn {
		case "interfaces":
			interfacesElement = arg
		case "datatypes":
			datatypes = arg
		case "AUTOSAR":
			IAUTOSAR = arg
		}
	}
	if interfacesElement == nil {
		return false, nil
	}
	if datatypes == nil {
		return false, nil
	}
	if IAUTOSAR == nil {
		return false, nil
	}
	return true, nil
}

func (c *ArxmlConverter) IsCP() (bool, error) {
	autosarElement := c.doc.SelectElement("AUTOSAR")
	if autosarElement == nil {
		return false, nil
	}
	var DataTypes *etree.Element
	var Communication *etree.Element
	var SoftwareTypes *etree.Element
	var SoAdRoutingGroups *etree.Element
	var System *etree.Element
	var Topology *etree.Element
	var DataTypeMappingSets *etree.Element
	arpackagesElement := autosarElement.SelectElement("AR-PACKAGES")
	if arpackagesElement == nil {
		return false, nil
	}
	arpList := arpackagesElement.SelectElements("AR-PACKAGE")
	for _, arg := range arpList {
		sn, err := util.GetShortname(arg)
		if err != nil {
			continue
		}
		switch sn {
		case "DataTypes":
			DataTypes = arg
		case "Communication":
			Communication = arg
		case "SoftwareTypes":
			SoftwareTypes = arg
		case "SoAdRoutingGroups":
			SoAdRoutingGroups = arg
		case "System":
			System = arg
		case "Topology":
			Topology = arg
		case "DataTypeMappingSets":
			DataTypeMappingSets = arg

		}
	}
	if DataTypes == nil {
		return false, nil
	}
	if Communication == nil {
		return false, nil
	}
	if SoftwareTypes == nil {
		return false, nil
	}
	if SoAdRoutingGroups == nil {
		return false, nil
	}
	if System == nil {
		return false, nil
	}
	if Topology == nil {
		return false, nil
	}
	if DataTypeMappingSets == nil {
		return false, nil
	}
	return true, nil
}
