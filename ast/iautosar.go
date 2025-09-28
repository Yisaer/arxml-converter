package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/beevik/etree"
)

func (p *Parser) parseIautoSar() error {
	eles := p.iautoSarElement.SelectElement("ELEMENTS")
	if eles == nil {
		return fmt.Errorf("no ELEMENTS in iautosar arpackage")
	}
	serviceInterfaces := eles.SelectElements("SOMEIP-SERVICE-INTERFACE-DEPLOYMENT")
	if len(serviceInterfaces) < 1 {
		return fmt.Errorf("no someip services in iautosar arpackage")
	}
	for _, si := range serviceInterfaces {
		service, err := p.ParseServiceInterface(si)
		if err != nil {
			return fmt.Errorf("parsing service interface: %w", err)
		}
		p.Services[service.ServiceID] = service
	}
	return nil
}

func (p *Parser) ParseServiceInterface(si *etree.Element) (*Service, error) {
	s := &Service{
		Events:      make(map[int]Event),
		FieldNotify: make(map[int]FieldNotify),
	}
	sn := si.SelectElement("SHORT-NAME")
	if sn == nil {
		return nil, fmt.Errorf("no SHORT-name in service")
	}
	s.ShortName = sn.Text()
	sid := si.SelectElement("SERVICE-INTERFACE-ID")
	if sid == nil {
		return nil, fmt.Errorf("no SERVICE-INTERFACE ID in service %v", s.ShortName)
	}
	sidI, err := strconv.Atoi(sid.Text())
	if err != nil {
		return nil, fmt.Errorf("invalid SERVICE-INTERFACE ID in service %v", s.ShortName)
	}
	s.ServiceID = sidI
	refRaw := si.SelectElement("SERVICE-INTERFACE-REF")
	if refRaw == nil {
		return nil, fmt.Errorf("no SERVICE-INTERFACE REF in service %v", s.ShortName)
	}
	s.ServiceInterfaceRef = refRaw.Text()
	if s.ServiceInterfaceRef == "" {
		return nil, fmt.Errorf("no SERVICE-INTERFACE REF in service %v", s.ShortName)
	}

	eds := si.SelectElement("EVENT-DEPLOYMENTS")
	if eds != nil {
		eventdeployments := eds.SelectElements("SOMEIP-EVENT-DEPLOYMENT")
		for _, ed := range eventdeployments {
			event := Event{}
			esn := ed.SelectElement("SHORT-NAME")
			if esn == nil {
				return nil, fmt.Errorf("no event SHORT-NAME in service %v", sn.Text())
			}
			event.ShortName = esn.Text()
			eidraw := ed.SelectElement("EVENT-ID")
			if eidraw == nil {
				return nil, fmt.Errorf("no EVENT-ID in service %v, event %v", sn.Text(), s.ShortName)
			}
			eid, err := strconv.Atoi(eidraw.Text())
			if err != nil {
				return nil, fmt.Errorf("invalid EVENT-ID in service %v, event %v", sn.Text(), s.ShortName)
			}
			refRaw := ed.SelectElement("EVENT-REF")
			if refRaw == nil {
				return nil, fmt.Errorf("no EVENT-REF in service %v, event %v", s.ShortName, event.ShortName)
			}
			event.EventRef = refRaw.Text()
			event.EventID = eid
			s.Events[eid] = event
		}
	}
	fds := si.SelectElement("FIELD-DEPLOYMENTS")
	if fds != nil {
		fieldDeployments := fds.SelectElements("SOMEIP-FIELD-DEPLOYMENT")
		for _, fd := range fieldDeployments {
			fieldNotify := FieldNotify{}
			esn := fd.SelectElement("SHORT-NAME")
			if esn == nil {
				return nil, fmt.Errorf("no SHORT-NAME in service %v field deployment", s.ShortName)
			}
			fieldNotify.ShortName = esn.Text()
			refRaw := fd.SelectElement("FIELD-REF")
			if refRaw == nil || refRaw.Text() == "" {
				return nil, fmt.Errorf("no FIELD-REF in service %v field deployment %v", s.ShortName, fieldNotify.ShortName)
			}
			fieldNotify.FieldRef = refRaw.Text()
			notifier := fd.SelectElement("NOTIFIER")
			if notifier != nil {
				eidraw := notifier.SelectElement("EVENT-ID")
				if eidraw == nil {
					return nil, fmt.Errorf("no EVENT-ID in service %v, event %v", s.ShortName, fieldNotify.ShortName)
				}
				eid, err := strconv.Atoi(eidraw.Text())
				if err != nil {
					return nil, fmt.Errorf("invalid EVENT-ID in service %v, fieldNotify %v", s.ShortName, fieldNotify.ShortName)
				}
				fieldNotify.EventID = eid
				s.FieldNotify[eid] = fieldNotify
			}

		}
	}
	if len(s.Events) < 1 && len(s.FieldNotify) < 1 {
		return nil, fmt.Errorf("no events/FieldsNotify in service %v", sn.Text())
	}
	return s, s.Validate()
}

func (s *Service) Validate() error {
	if len(s.ServiceInterfaceRef) < 1 {
		return fmt.Errorf("no service interface ref in service %v", s.ShortName)
	}
	serviceRef := strings.ToLower(s.ServiceInterfaceRef)
	for _, event := range s.Events {
		if !strings.HasPrefix(strings.ToLower(event.EventRef), serviceRef) {
			return fmt.Errorf("invalid event ref in service %v event %v, serviceRef:%v, eventRef:%v", s.ShortName, event.EventRef, s.ServiceInterfaceRef, event.EventRef)
		}
	}
	for _, field := range s.FieldNotify {
		if !strings.HasPrefix(strings.ToLower(field.FieldRef), serviceRef) {
			return fmt.Errorf("invalid field ref in service %v field %v, serviceRef:%v, fieldRef:%v", s.ShortName, field.FieldRef, s.ServiceInterfaceRef, field.FieldRef)
		}
	}
	return nil
}

type Service struct {
	ShortName           string
	ServiceInterfaceRef string

	ServiceID   int
	Events      map[int]Event
	FieldNotify map[int]FieldNotify
}

type Event struct {
	EventID   int
	ShortName string
	EventRef  string
}

type FieldNotify struct {
	EventID   int
	ShortName string
	FieldRef  string
}
