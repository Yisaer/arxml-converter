package topology

import (
	"fmt"

	"github.com/beevik/etree"
)

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
