package ast

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
)

func (p *Parser) getShortname(node *etree.Element) (string, error) {
	sn := node.SelectElement("SHORT-NAME")
	if sn == nil {
		return "", fmt.Errorf("no short name found")
	}
	sname := sn.Text()
	if len(sname) < 1 {
		return "", fmt.Errorf("no short name found")
	}
	return sname, nil
}

func (p *Parser) getCategory(node *etree.Element) (string, error) {
	sn := node.SelectElement("CATEGORY")
	if sn == nil {
		return "", fmt.Errorf("no CATEGORY found")
	}
	sname := sn.Text()
	if len(sname) < 1 {
		return "", fmt.Errorf("no CATEGORY found")
	}
	return sname, nil
}

func (p *Parser) getElements(node *etree.Element) (*etree.Element, error) {
	es := node.SelectElement("ELEMENTS")
	if es == nil {
		return nil, fmt.Errorf("no ELEMENTS found")
	}
	return es, nil
}

func (p *Parser) getSWDataDefPropsConditional(node *etree.Element) (*etree.Element, error) {
	sddp := node.SelectElement("SW-DATA-DEF-PROPS")
	if sddp == nil {
		return nil, fmt.Errorf("no SW-DATA-DEF-PROPS found")
	}
	sddpv := sddp.SelectElement("SW-DATA-DEF-PROPS-VARIANTS")
	if sddpv == nil {
		return nil, fmt.Errorf("no SW-DATA-DEF-PROPS-VARIANTS found")
	}
	sddpc := sddpv.SelectElement("SW-DATA-DEF-PROPS-CONDITIONAL")
	if sddpc == nil {
		return nil, fmt.Errorf("no SW-DATA-DEF-PROPS-CONDITIONAL found")
	}
	return sddpc, nil
}

func (p *Parser) getArraySizeSemantics(node *etree.Element) (isDynamic bool, err error) {
	ass := node.SelectElement("ARRAY-SIZE-SEMANTICS")
	if ass == nil {
		return false, fmt.Errorf("no ARRAY-SIZE-SEMANTICS found")
	}
	switch ass.Text() {
	case "VARIABLE-SIZE":
		return true, nil
	case "FIXED-SIZE":
		return false, nil
	}
	return false, fmt.Errorf("invalid ARRAY-SIZE-SEMANTICS:%v", ass.Text())
}

func (p *Parser) validBasicType(ref string) error {
	typeName := extractTypeNameFromRef(ref)
	lowerName := strings.ToLower(typeName)
	switch {
	case strings.Contains(lowerName, "uint8"):
		return nil
	case strings.Contains(lowerName, "uint16"):
		return nil
	case strings.Contains(lowerName, "uint32"):
		return nil
	case strings.Contains(lowerName, "uint64"):
		return nil
	case strings.Contains(lowerName, "bool"):
		return nil
	case strings.Contains(lowerName, "int8"):
		return nil
	case strings.Contains(lowerName, "int16"):
		return nil
	case strings.Contains(lowerName, "int32"):
		return nil
	case strings.Contains(lowerName, "int64"):
		return nil
	case strings.Contains(lowerName, "float"):
		return nil
	case strings.Contains(lowerName, "double"):
		return nil
	}
	return fmt.Errorf("invalid basic type:%v", ref)
}

func extractTypeNameFromRef(ref string) string {
	// 移除路径前缀，只保留最后的类型名
	parts := strings.Split(ref, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ref
}
