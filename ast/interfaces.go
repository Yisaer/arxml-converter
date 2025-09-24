package ast

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
)

func (p *Parser) parseInterfaces() error {
	eles := p.interfacesElement.SelectElement("ELEMENTS")
	if eles == nil {
		return fmt.Errorf("no ELEMENTS in interfaces arpackage")
	}
	serviceInterfaces := eles.SelectElements("SERVICE-INTERFACE")
	if len(serviceInterfaces) < 1 {
		return fmt.Errorf("no SERVICE-INTERFACE in elements interfaces arpackage")
	}

	for index, serviceInterface := range serviceInterfaces {
		si, err := p.parseServiceInterface(serviceInterface)
		if err != nil {
			return fmt.Errorf("parsing index %v service interface failed, err:%v", index, err)
		}
		p.Interfaces[strings.ToLower(si.Shortname)] = si
	}
	return nil
}

func (p *Parser) parseServiceInterface(e *etree.Element) (*ServiceInterface, error) {
	sn := e.SelectElement("SHORT-NAME")
	if sn == nil {
		return nil, fmt.Errorf("no SHORT-NAME in service interface arpackage")
	}
	si := &ServiceInterface{
		Shortname: sn.Text(),
		Events:    make(map[string]ServiceInterfaceEvent),
		Fields:    make(map[string]ServiceInterfaceField),
	}
	es := e.SelectElement("EVENTS")
	if es != nil {
		vdps := es.SelectElements("VARIABLE-DATA-PROTOTYPE")
		for _, vdp := range vdps {
			sn := vdp.SelectElement("SHORT-NAME")
			if sn == nil {
				return nil, fmt.Errorf("no SHORT-NAME in service interface %v event vdp", si.Shortname)
			}
			eventShortname := sn.Text()
			typref := vdp.SelectElement("TYPE-TREF")
			if typref == nil {
				return nil, fmt.Errorf("no TYPE-TREF in serviceInterface %v event vdp %v", si.Shortname, eventShortname)
			}
			si.Events[strings.ToLower(eventShortname)] = ServiceInterfaceEvent{
				ShortName: eventShortname,
				TypeRef:   typref.Text(),
			}
		}
	}

	fss := e.SelectElement("FIELDS")
	if fss != nil {
		fs := fss.SelectElements("FIELD")
		for _, field := range fs {
			sn := field.SelectElement("SHORT-NAME")
			if sn == nil {
				return nil, fmt.Errorf("no SHORT-NAME in service interface %v field", si.Shortname)
			}
			fieldShortname := sn.Text()
			typref := field.SelectElement("TYPE-TREF")
			if typref == nil {
				return nil, fmt.Errorf("no TYPE-TREF in service interface %v field %v", si.Shortname, fieldShortname)
			}
			si.Fields[strings.ToLower(fieldShortname)] = ServiceInterfaceField{
				ShortName: fieldShortname,
				TypeRef:   typref.Text(),
			}
		}
	}
	if len(si.Events) < 1 && len(si.Fields) < 1 {
		return nil, fmt.Errorf("no EVENTS/Fields found in service interface %v", si.Shortname)
	}
	return si, nil
}

type ServiceInterface struct {
	Shortname string
	Events    map[string]ServiceInterfaceEvent
	Fields    map[string]ServiceInterfaceField
}

type ServiceInterfaceEvent struct {
	ShortName string
	TypeRef   string
}

type ServiceInterfaceField struct {
	ShortName string
	TypeRef   string
}
