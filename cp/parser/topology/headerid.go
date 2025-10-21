package topology

import (
	"fmt"

	"github.com/beevik/etree"

	"github.com/yisaer/arxml-converter/util"
)

func (tp *TopoLogyParser) parseSOCKETCONNECTIONBUNDLE(node *etree.Element) (err error) {
	sn, err := util.GetShortname(node)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = fmt.Errorf("parse %v SOCKETCONNECTIONBUNDLE error: %v", sn, err)
		}
	}()
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
