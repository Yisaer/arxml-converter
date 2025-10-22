package tpConfig

import "github.com/beevik/etree"

type TpConfigParser struct {
	pduMap map[string]string
}

func NewTpConfigParser() *TpConfigParser {
	return &TpConfigParser{pduMap: make(map[string]string)}
}

func (p *TpConfigParser) ParseTpConfig(node *etree.Element) error {
	for _, element := range node.FindElements("//SOMEIP-TP-CONNECTION") {
		p.parseSOMEIPTPCONNECTION(element)
	}
	return nil
}

func (p *TpConfigParser) parseSOMEIPTPCONNECTION(node *etree.Element) error {
	tpSDUREFElement := node.SelectElement("TP-SDU-REF")
	if tpSDUREFElement == nil {
		return nil
	}
	TRANSPORTPDUREFElement := node.SelectElement("TRANSPORT-PDU-REF")
	if TRANSPORTPDUREFElement == nil {
		return nil
	}
	p.pduMap[TRANSPORTPDUREFElement.Text()] = tpSDUREFElement.Text()
	return nil
}

func (p *TpConfigParser) GetTpConfigPDUMap() map[string]string {
	return p.pduMap
}
