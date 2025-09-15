package converter

import (
	"fmt"

	idlAst "github.com/yisaer/idl-parser/ast"
	"github.com/yisaer/idl-parser/converter"

	"arxml-converter/ast"
)

type FieldRefType int

const (
	StringType FieldRefType = iota
	BoolType
	FloatType
	DoubleType
	Int8Type
	Int16Type
	Int32Type
	Int64Type
	Uint8Type
	Uint16Type
	Uint32Type
	Uint64Type
)

func (ft FieldRefType) String() string {
	switch ft {
	case StringType:
		return "String"
	case BoolType:
		return "Bool"
	case FloatType:
		return "Float"
	case DoubleType:
		return "Double"
	case Int8Type:
		return "Int8"
	case Int16Type:
		return "Int16"
	case Int32Type:
		return "Int32"
	case Int64Type:
		return "Int64"
	case Uint8Type:
		return "Uint8"
	case Uint16Type:
		return "Uint16"
	case Uint32Type:
		return "Uint32"
	case Uint64Type:
		return "Uint64"
	default:
		return "Unknown"
	}
}

type ArXMLConverter struct {
	Parser       *ast.Parser
	typeRefs     map[string]*ast.TypReference
	arrRefS      map[string]*ast.Array
	structRefs   map[string]*ast.Structure
	module       *idlAst.Module
	idlConverter *converter.IDLConverter
}

func NewConverter(path string, config converter.IDlConverterConfig) (*ArXMLConverter, error) {
	parser, err := ast.NewParser(path)
	if err != nil {
		return nil, err
	}
	c := &ArXMLConverter{
		Parser:     parser,
		typeRefs:   make(map[string]*ast.TypReference),
		arrRefS:    make(map[string]*ast.Array),
		structRefs: make(map[string]*ast.Structure),
	}
	if err := c.Parser.Parse(); err != nil {
		return nil, err
	}
	for _, dt := range c.Parser.DtList {
		key := fmt.Sprintf("/dataTypes/%s", dt.ShorName)
		switch {
		case dt.TypReference != nil:
			c.typeRefs[key] = dt.TypReference
		case dt.Array != nil:
			c.arrRefS[key] = dt.Array
		case dt.Structure != nil:
			c.structRefs[key] = dt.Structure
		default:
			return nil, fmt.Errorf("unknown dt: %s, not found: %v", dt.ShorName, key)
		}
	}
	module, err := c.ToIDLModule()
	if err != nil {
		return nil, err
	}
	c.module = module
	idlConverter, err := converter.NewIDLConverterWithModule(config, *module)
	if err != nil {
		return nil, err
	}
	c.idlConverter = idlConverter
	return c, nil
}

func (c *ArXMLConverter) Decode(stName string, data []byte) (interface{}, error) {
	return c.idlConverter.Decode(fmt.Sprintf("ArXMLDataTypes.%s", stName), data)
}

// ToIDLModule 将 ArXML parser 的结果转换为 idlAst.Module
func (c *ArXMLConverter) ToIDLModule() (*idlAst.Module, error) {
	converter := NewArXMLToIDLConverter(c.Parser)
	return converter.ConvertToIDLModule()
}
