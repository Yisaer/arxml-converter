package converter

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
	"github.com/yisaer/idl-parser/converter"

	apconverter "github.com/yisaer/arxml-converter/ap/converter"
	cpconverter "github.com/yisaer/arxml-converter/cp/converter"
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
	switch c.version {
	case AUTOSAR_00048Version:
		c.apArxmlConverter, err = apconverter.NewConverterWithDoc(doc, config)
		if err != nil {
			return nil, err
		}
	case AUTOSAR_4_2_2Version:
		c.cpArxmlConverter, err = cpconverter.NewArxmlCPConverterWithDoc(doc, config)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
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
