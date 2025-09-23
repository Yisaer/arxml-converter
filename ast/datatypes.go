package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/beevik/etree"
)

func (p *Parser) searchDataTypes(arPackageElements []*etree.Element) error {
	for _, arPackage := range arPackageElements {
		s := arPackage.SelectElement("SHORT-NAME")
		if s == nil {
			continue
		}
		if s.Text() == "dataTypes" {
			p.dataTypesElement = arPackage
			return nil
		}
	}
	return fmt.Errorf("no dataTypes find in ar package")
}

func (p *Parser) parseDataTypes() error {
	eles := p.dataTypesElement.SelectElement("ELEMENTS")
	if eles == nil {
		return fmt.Errorf("no ELEMENTS in datatypes arpackage")
	}
	dataTypes := eles.SelectElements("STD-CPP-IMPLEMENTATION-DATA-TYPE")
	if len(dataTypes) < 1 {
		return fmt.Errorf("no STD-CPP-IMPLEMENTATION-DATA-TYPE in elements datatypes arpackage")
	}
	for index, dataType := range dataTypes {
		dt, err := p.parseDataType(dataType)
		if err != nil {
			return fmt.Errorf("index %d STD-CPP-IMPLEMENTATION-DATA-TYPE has err:%v", index, err.Error())
		}
		p.DataTypes[dt.ShorName] = dt
	}
	return nil
}

func (p *Parser) parseDataType(d *etree.Element) (*DataType, error) {
	dt := &DataType{}
	sn := d.SelectElement("SHORT-NAME")
	if sn == nil {
		return nil, fmt.Errorf("no SHORT-NAME in %v", d.Text())
	}
	dt.ShorName = sn.Text()
	category := d.SelectElement("CATEGORY")
	if category == nil {
		return nil, fmt.Errorf("no CATEGORY in %v", d.Text())
	}
	dt.Category = category.Text()
	switch dt.Category {
	case "TYPE_REFERENCE":
		ref := d.SelectElement("TYPE-REFERENCE-REF")
		if ref == nil {
			return nil, fmt.Errorf("no TYPE-REFERENCE-REF")
		}
		dt.TypReference = &TypReference{Ref: ref.Text()}
		if strings.Contains(strings.ToLower(dt.TypReference.Ref), "string") {
			stringSize := d.SelectElement("ARRAY-SIZE")
			if stringSize == nil {
				return nil, fmt.Errorf("no ARRAY-SIZE for TYPE-REFERENCE-REF string")
			}
			as, err := strconv.ParseInt(stringSize.Text(), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid ARRAY-SIZE: %s", stringSize.Text())
			}
			dt.StringSize = as
		}
	case "ARRAY":
		arraySize := d.SelectElement("ARRAY-SIZE")
		if arraySize == nil {
			return nil, fmt.Errorf("no ARRAY-SIZE")
		}
		as, err := strconv.ParseInt(arraySize.Text(), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid ARRAY-SIZE: %s", arraySize.Text())
		}
		args := d.SelectElement("TEMPLATE-ARGUMENTS")
		if args == nil {
			return nil, fmt.Errorf("no TEMPLATE-ARGUMENTS")
		}
		cppArgs := args.SelectElement("CPP-TEMPLATE-ARGUMENT")
		if cppArgs == nil {
			return nil, fmt.Errorf("no CPP-TEMPLATE-ARGUMENT in TEMPLATE-ARGUMENTS")
		}
		inPlace := cppArgs.SelectElement("INPLACE")
		if inPlace == nil {
			return nil, fmt.Errorf("no INPLACE")
		}
		ip, err := strconv.ParseBool(inPlace.Text())
		if err != nil {
			return nil, fmt.Errorf("invalid INPLACE: %s", inPlace.Text())
		}
		typRef := cppArgs.SelectElement("TEMPLATE-TYPE-REF")
		if typRef == nil {
			return nil, fmt.Errorf("no TEMPLATE-TYPE-REF in CPP-TEMPLATE-ARGUMENT")
		}
		dt.Array = &Array{
			ArraySize: as,
			Inplace:   ip,
			RefType:   typRef.Text(),
		}
	case "STRUCTURE":
		if err := p.ParseStructure(dt, d); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid category: %s", dt.Category)
	}
	return dt, nil
}

type DataType struct {
	ShorName string `json:"shor_name"`
	Category string `json:"category"`
	*TypReference
	*Array
	*Structure
}

type TypReference struct {
	Ref        string `json:"ref"`
	StringSize int64  `json:"string_size"`
}

type Array struct {
	ArraySize int64  `json:"array_size"`
	Inplace   bool   `json:"inplace"`
	RefType   string `json:"ref_type"`
}

type Structure struct {
	STRList []*StructureTypRef `json:"str_list"`
}

type StructureTypRef struct {
	InPlace  bool   `json:"in_place"`
	Ref      string `json:"ref"`
	ShorName string `json:"shor_name"`
}

func (p *Parser) ParseStructure(dt *DataType, d *etree.Element) error {
	subElements := d.SelectElement("SUB-ELEMENTS")
	if subElements == nil {
		return fmt.Errorf("no SUB-ELEMENTS")
	}
	cppElements := subElements.SelectElements("CPP-IMPLEMENTATION-DATA-TYPE-ELEMENT")
	dt.Structure = &Structure{
		STRList: make([]*StructureTypRef, 0),
	}
	for _, cppElement := range cppElements {
		str := &StructureTypRef{}
		sn := cppElement.SelectElement("SHORT-NAME")
		if sn == nil {
			return fmt.Errorf("no SHORT-NAME")
		}
		typRef := cppElement.SelectElement("TYPE-REFERENCE")
		if typRef == nil {
			return fmt.Errorf("no TYPE-REFERENCE")
		}
		Inplace := typRef.SelectElement("INPLACE")
		if Inplace == nil {
			return fmt.Errorf("no INPLACE")
		}
		ip, err := strconv.ParseBool(Inplace.Text())
		if err != nil {
			return fmt.Errorf("invalid bool: %s", Inplace.Text())
		}
		trd := typRef.SelectElement("TYPE-REFERENCE-REF")
		if trd == nil {
			return fmt.Errorf("no TYPE-REFERENCE-REF in TYPE-REFERENCE")
		}
		str.ShorName = sn.Text()
		str.InPlace = ip
		str.Ref = trd.Text()
		dt.Structure.STRList = append(dt.Structure.STRList, str)
	}
	return nil
}
