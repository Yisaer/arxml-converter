package converter

import (
	"fmt"

	"github.com/beevik/etree"
	"github.com/yisaer/idl-parser/converter"
	cpconverter "github.com/yisaer/arxml-converter/cp/converter"
	apconverter "github.com/yisaer/arxml-converter/ap/converter"
)

type ArxmlConverter struct {
	path             string
	config           converter.IDlConverterConfig
	cpArxmlConverter *cpconverter.ArxmlCPConverter
	apArxmlConverter *apconverter.ArXMLConverter
	doc              *etree.Document
}

func NewConverter(path string, config converter.IDlConverterConfig) (*ArxmlConverter, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(path); err != nil {
		return nil, err
	}
	doc.SelectElement("AUTOSAR")
	return nil, nil
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
	return nil
}
