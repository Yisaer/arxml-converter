package ast

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"

	"github.com/yisaer/arxml-converter/mod"
)

type DataTypesParser struct {
	*Parser

	implementationDataTypesArPackage *etree.Element
	applicationDatatypeArPackage     *etree.Element

	applicationDataTypes    map[string]*mod.DataType
	implementationDataTypes map[string]*mod.DataType
}

func NewDataTypesParser(parser *Parser) *DataTypesParser {
	return &DataTypesParser{
		Parser:                  parser,
		applicationDataTypes:    make(map[string]*mod.DataType),
		implementationDataTypes: make(map[string]*mod.DataType),
	}
}

func (dp *DataTypesParser) parseApplicationDatatypes(root *etree.Element) error {
	elements, err := dp.getElements(root)
	if err != nil {
		return err
	}
	for index, apdt := range elements.SelectElements("APPLICATION-PRIMITIVE-DATA-TYPE") {
		if err := dp.ParseApplicationDataType(apdt); err != nil {
			return fmt.Errorf("parse index %v APPLICATION-PRIMITIVE-DATA-TYPE failed, err:%v", index, err.Error())
		}
	}
	for index, aadt := range elements.SelectElements("APPLICATION-ARRAY-DATA-TYPE") {
		if err := dp.ParseApplicationDataType(aadt); err != nil {
			return fmt.Errorf("parse index %v APPLICATION-ARRAY-DATA-TYPE failed, err:%v", index, err.Error())
		}
	}
	for index, ardt := range elements.SelectElements("APPLICATION-RECORD-DATA-TYPE") {
		if err := dp.ParseApplicationDataType(ardt); err != nil {
			return fmt.Errorf("parse index %v APPLICATION-RECORD-DATA-TYPE failed, err:%v", index, err.Error())
		}
	}

	return nil
}

func (dp *DataTypesParser) ParseApplicationDataType(root *etree.Element) (err error) {
	sn, err := dp.Parser.getShortname(root)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = fmt.Errorf("shortname:%v ,err:%v", sn, err.Error())
		}
	}()

	category, err := dp.Parser.getCategory(root)
	if err != nil {
		return fmt.Errorf("parse category failed err:%v", err.Error())
	}
	switch category {
	case "STRING":
		sddpc, err := dp.Parser.getSWDataDefPropsConditional(root)
		if err != nil {
			return err
		}
		stp := sddpc.SelectElement("SW-TEXT-PROPS")
		if stp == nil {
			return fmt.Errorf("no SW-TEXT-PROPS found")
		}
		isDynamicString, err := dp.Parser.getArraySizeSemantics(stp)
		if err != nil {
			return err
		}
		if !isDynamicString {
			return fmt.Errorf("fixed length string not supported now")
		}
		btr := stp.SelectElement("BASE-TYPE-REF")
		if btr == nil {
			return fmt.Errorf("no BASE-TYPE-REF found")
		}
		btrRaw := btr.Text()
		if !strings.Contains(strings.ToUpper(btrRaw), "UTF_8") {
			return fmt.Errorf("BASE-TYPE ref should be UTF_8, got:%v", btrRaw)
		}
		dp.applicationDataTypes[sn] = mod.NewStringDataType(sn, category, 0)
	case "VALUE":
		idtrKey, ok := dp.dataTypeMappings[sn]
		if !ok {
			return fmt.Errorf("failed to find mapping for applicationDataType:%v", sn)
		}
		dt, ok := dp.implementationDataTypes[idtrKey]
		if !ok {
			return fmt.Errorf("failed to implementationDataType:%v for application key:%v", idtrKey, sn)
		}
		dt.ShorName = sn
		dp.applicationDataTypes[sn] = dt
	case "ARRAY":
		element := root.SelectElement("ELEMENT")
		if element == nil {
			return fmt.Errorf("no ELEMENT found")
		}
		typeRef := element.SelectElement("TYPE-TREF")
		if typeRef == nil {
			return fmt.Errorf("no TYPE-TREF found for sn %v", sn)
		}
		arrayRef := strings.TrimPrefix(typeRef.Text(), appDataTypePrefix)
		isDynamicArray, err := dp.Parser.getArraySizeSemantics(element)
		if err != nil {
			return err
		}
		if !isDynamicArray {
			return fmt.Errorf("fixed length array not supported now")
		}
		dp.applicationDataTypes[sn] = mod.NewArrayDataType(sn, category, arrayRef, 0)
	case "STRUCTURE":
		elements := root.SelectElement("ELEMENTS")
		if elements == nil {
			return fmt.Errorf("no ELEMENTS found")
		}
		s := &mod.Structure{
			STRList: make([]*mod.StructureTypRef, 0),
		}
		records := elements.SelectElements("APPLICATION-RECORD-ELEMENT")
		for _, record := range records {
			ref := &mod.StructureTypRef{}
			recordSN, err := dp.Parser.getShortname(record)
			if err != nil {
				return err
			}
			ref.ShorName = recordSN
			typeRef := record.SelectElement("TYPE-TREF")
			if typeRef == nil {
				return fmt.Errorf("no TYPE-REF found for sn %v", recordSN)
			}
			ref.Ref = strings.TrimPrefix(typeRef.Text(), appDataTypePrefix)
			s.STRList = append(s.STRList, ref)
		}
		dp.applicationDataTypes[sn] = mod.NewStructureDataType(sn, category, s)
	default:
		return fmt.Errorf("unknown category:%v", category)
	}
	return nil
}
