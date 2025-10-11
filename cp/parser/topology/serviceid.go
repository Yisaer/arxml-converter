package topology

import (
	"fmt"

	"github.com/beevik/etree"

	"github.com/yisaer/arxml-converter/util"
)

func (tp *TopoLogyParser) parseSOCKETADDRESS(node *etree.Element) (err error) {
	applicationEndpointElement := node.SelectElement("APPLICATION-ENDPOINT")
	if applicationEndpointElement == nil {
		return nil
	}
	proServiceInstancesElement := applicationEndpointElement.SelectElement("PROVIDED-SERVICE-INSTANCES")
	if proServiceInstancesElement == nil {
		return nil
	}
	providedServiceInstanceList := proServiceInstancesElement.SelectElements("PROVIDED-SERVICE-INSTANCE")
	for index, providedInstance := range providedServiceInstanceList {
		if err := tp.parseProvidedServiceInstance(providedInstance); err != nil {
			return fmt.Errorf("parse %v providedServiceInstance err: %v", index, err)
		}
	}
	return nil
}

func (tp *TopoLogyParser) parseProvidedServiceInstance(node *etree.Element) (err error) {
	sn, err := util.GetShortname(node)
	if err != nil {
		return err
	}
	serviceIDElement := node.SelectElement("SERVICE-IDENTIFIER")
	if serviceIDElement == nil {
		return nil
	}
	serviceID, err := util.ToUint16(serviceIDElement.Text())
	if err != nil {
		return err
	}
	tp.serviceIDMap[serviceID] = sn
	return nil
}
