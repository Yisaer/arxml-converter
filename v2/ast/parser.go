package ast

import (
	"github.com/beevik/etree"

	"github.com/yisaer/arxml-converter/mod"
)

type Parser struct {
	Path             string
	Doc              *etree.Document
	dataTypesElement *etree.Element
	DataTypes        map[string]*mod.DataType

	implementationDataTypesArPackage *etree.Element
	applicationDatatypeArPackage     *etree.Element
	dataTypeMappingSetsArPackage     *etree.Element

	implementationDataTypes map[string]*mod.DataType
	dataTypeMappings        map[string]string
}

func NewParser(path string) (*Parser, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(path); err != nil {
		return nil, err
	}
	p := &Parser{Path: path, Doc: doc}
	p.DataTypes = make(map[string]*mod.DataType)
	return p, nil
}
