package ast

import (
	"fmt"
	"strconv"

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
		Events: make(map[int]Event),
	}
	sn := si.SelectElement("SHORT-NAME")
	if sn == nil {
		return nil, fmt.Errorf("no SHORT-name in service")
	}
	s.ShortName = sn.Text()
	sid := si.SelectElement("SERVICE-INTERFACE-ID")
	if sid == nil {
		return nil, fmt.Errorf("no SERVICE-INTERFACE ID in service")
	}
	sidI, err := strconv.Atoi(sid.Text())
	if err != nil {
		return nil, fmt.Errorf("invalid SERVICE-INTERFACE ID")
	}
	s.ServiceID = sidI
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
				return nil, fmt.Errorf("no EVENT-ID in service %v, event %v", sn.Text(), event.ShortName)
			}
			eid, err := strconv.Atoi(eidraw.Text())
			if err != nil {
				return nil, fmt.Errorf("invalid EVENT-ID in service %v, event %v", sn.Text(), event.ShortName)
			}
			event.EventID = eid
			s.Events[eid] = event
		}
	}
	fds := si.SelectElement("FORM-DEPLOYMENTS")
	if fds != nil {
		eventdeployments := fds.SelectElements("SOMEIP-FIELD-DEPLOYMENT")
		for _, ed := range eventdeployments {
			event := Event{}
			esn := ed.SelectElement("SHORT-NAME")
			if esn == nil {
				return nil, fmt.Errorf("no event SHORT-NAME in service %v", sn.Text())
			}
			event.ShortName = esn.Text()
			notifier := ed.SelectElement("NOTIFIER")
			if notifier == nil {
				return nil, fmt.Errorf("no NOTIFIER in service %v event %v", sn.Text(), event.ShortName)
			}
			eidraw := notifier.SelectElement("EVENT-ID")
			if eidraw == nil {
				return nil, fmt.Errorf("no EVENT-ID in service %v, event %v", sn.Text(), event.ShortName)
			}
			eid, err := strconv.Atoi(eidraw.Text())
			if err != nil {
				return nil, fmt.Errorf("invalid EVENT-ID in service %v, event %v", sn.Text(), event.ShortName)
			}
			event.EventID = eid
			s.Events[eid] = event
		}
	}
	if len(s.Events) < 1 {
		return nil, fmt.Errorf("no events in service %v", sn.Text())
	}
	return s, nil
}

type Service struct {
	ShortName string
	ServiceID int
	Events    map[int]Event
}

type Event struct {
	EventID   int
	ShortName string
}
