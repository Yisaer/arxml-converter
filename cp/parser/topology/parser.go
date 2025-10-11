package topology

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"

	"github.com/yisaer/arxml-converter/util"
)

type TopoLogyParser struct {
	clusterArPackage *etree.Element
	serviceIDMap     map[uint16]string
	headerIdRef      map[uint32]string
}

func NewTopoLogyParser() *TopoLogyParser {
	return &TopoLogyParser{
		serviceIDMap: make(map[uint16]string),
		headerIdRef:  make(map[uint32]string),
	}
}

func (tp *TopoLogyParser) GetServiceIDMap() map[uint16]string {
	return tp.serviceIDMap
}

func (tp *TopoLogyParser) GetHeaderRef() map[uint32]string {
	return tp.headerIdRef
}

func (tp *TopoLogyParser) ParseTopoLogy(node *etree.Element) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("ParseTopoLogy err: %v", err)
		}
	}()
	arpackagesElement := node.SelectElement("AR-PACKAGES")
	if arpackagesElement == nil {
		return fmt.Errorf("AR-PACKAGES not found")
	}
	arpackagesList := arpackagesElement.SelectElements("AR-PACKAGE")
	if err := tp.searchCluster(arpackagesList); err != nil {
		return err
	}
	if err := tp.parseCluster(); err != nil {
		return fmt.Errorf("parse cluster err: %v", err)
	}
	return nil
}

func (tp *TopoLogyParser) parseCluster() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("parseCluster err: %v", err)
		}
	}()
	elements, err := util.GetElements(tp.clusterArPackage)
	if err != nil {
		return err
	}
	ethClusterElement := elements.SelectElement("ETHERNET-CLUSTER")
	if ethClusterElement == nil {
		return fmt.Errorf("ETHERNET-CLUSTER not found")
	}
	ethClusterVar := ethClusterElement.SelectElement("ETHERNET-CLUSTER-VARIANTS")
	if ethClusterVar == nil {
		return fmt.Errorf("ETHERNET-CLUSTER-VARIANTS not found")
	}
	ethClusterCondition := ethClusterVar.SelectElement("ETHERNET-CLUSTER-CONDITIONAL")
	if ethClusterCondition == nil {
		return fmt.Errorf("ETHERNET-CLUSTER-CONDITIONAL not found")
	}
	phyChannelsElement := ethClusterCondition.SelectElement("PHYSICAL-CHANNELS")
	if phyChannelsElement == nil {
		return fmt.Errorf("PHYSICAL-CHANNELS not found")
	}
	ethPhyChannels := phyChannelsElement.SelectElements("ETHERNET-PHYSICAL-CHANNEL")
	for index, ethPhyChannel := range ethPhyChannels {
		if err := tp.parseETHERNETPHYSICALCHANNEL(ethPhyChannel); err != nil {
			return fmt.Errorf("parse %v ETHERNET-PHYSICAL-CHANNEL err: %v", index, err)
		}
	}
	return nil
}

func (tp *TopoLogyParser) parseETHERNETPHYSICALCHANNEL(node *etree.Element) (err error) {
	soAdConfigElement := node.SelectElement("SO-AD-CONFIG")
	if soAdConfigElement == nil {
		return fmt.Errorf("SO-AD-CONFIG not found")
	}
	// parse service id
	socketAddresssElement := soAdConfigElement.SelectElement("SOCKET-ADDRESSS")
	if socketAddresssElement == nil {
		return fmt.Errorf("SOCKET-ADDRESSS not found")
	}
	socketAddressList := socketAddresssElement.SelectElements("SOCKET-ADDRESS")
	for index, socketAddress := range socketAddressList {
		if err := tp.parseSOCKETADDRESS(socketAddress); err != nil {
			return fmt.Errorf("parse %v SOCKET-ADDRESS err: %v", index, err)
		}
	}

	// parse header id
	connectionBundlesElement := soAdConfigElement.SelectElement("CONNECTION-BUNDLES")
	if connectionBundlesElement == nil {
		return fmt.Errorf("CONNECTION-BUNDLES not found")
	}
	socketConnectionBundleList := connectionBundlesElement.SelectElements("SOCKET-CONNECTION-BUNDLE")
	for index, socketConnectionBundle := range socketConnectionBundleList {
		if err := tp.parseSOCKETCONNECTIONBUNDLE(socketConnectionBundle); err != nil {
			return fmt.Errorf("parse %v SOCKET-CONNECTION-BUNDLE err :%v", index, err)
		}
	}

	return nil
}

func (tp *TopoLogyParser) parseSOCKETCONNECTIONBUNDLE(node *etree.Element) (err error) {
	bundleConnectionsElement := node.SelectElement("BUNDLED-CONNECTIONS")
	if bundleConnectionsElement == nil {
		return nil
	}
	socketConnectionElement := bundleConnectionsElement.SelectElement("SOCKET-CONNECTION")
	if socketConnectionElement == nil {
		return nil
	}
	pdusElement := socketConnectionElement.SelectElement("PDUS")
	if pdusElement == nil {
		return nil
	}
	socketConnectionIPDUIdentifierList := pdusElement.SelectElements("SOCKET-CONNECTION-IPDU-IDENTIFIER")
	for index, scipdui := range socketConnectionIPDUIdentifierList {
		if err := tp.parseSOCKETCONNECTIONIPDUIDENTIFIER(scipdui); err != nil {
			return fmt.Errorf("parse %v SOCKET-CONNECTION-IPDU-IDENTIFIER err: %v", index, err)
		}
	}
	return nil
}

func (tp *TopoLogyParser) parseSOCKETCONNECTIONIPDUIDENTIFIER(node *etree.Element) (err error) {

	headerIDElement := node.SelectElement("HEADER-ID")
	if headerIDElement == nil {
		return fmt.Errorf("HEADER-ID not found")
	}
	headerID, err := util.ToUint32(headerIDElement.Text())
	if err != nil {
		return fmt.Errorf("parse HEADER-ID err: %v", err)
	}
	pduTriggeringRefElement := node.SelectElement("PDU-TRIGGERING-REF")
	pduTriggeringRefElementRaw := pduTriggeringRefElement.Text()
	if strings.Contains(pduTriggeringRefElementRaw, "call") {
		tp.headerIdRef[headerID] = pduTriggeringRefElementRaw
	}
	return nil
}

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

func (tp *TopoLogyParser) searchCluster(arpackagesList []*etree.Element) error {
	for _, arPackage := range arpackagesList {
		sn, err := util.GetShortname(arPackage)
		if err != nil {
			return err
		}
		if sn == "Clusters" {
			tp.clusterArPackage = arPackage
			return nil
		}
	}
	return fmt.Errorf("AR-PACKAGES Cluster not found")
}
