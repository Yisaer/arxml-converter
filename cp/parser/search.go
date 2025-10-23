package parser

import (
	"fmt"

	"github.com/beevik/etree"

	"github.com/yisaer/arxml-converter/util"
)

func (p *Parser) searchDataTypes(arPackagesElement *etree.Element) error {
	arPackages := arPackagesElement.SelectElements("AR-PACKAGE")
	for _, arPackage := range arPackages {
		sn, err := util.GetShortname(arPackage)
		if err != nil {
			return err
		}
		if sn == "DataTypes" {
			p.dataTypesElement = arPackage
			return nil
		}
	}
	return fmt.Errorf("no DataTypes found")
}

func (p *Parser) searchDataTypeMappingSets(arPackagesElement *etree.Element) error {
	arPackages := arPackagesElement.SelectElements("AR-PACKAGE")
	for _, arPackage := range arPackages {
		sn, err := util.GetShortname(arPackage)
		if err != nil {
			return err
		}
		if sn == "DataTypeMappingSets" {
			p.dataTypeMappingSetsElement = arPackage
			return nil
		}
	}
	return fmt.Errorf("no DataTypeMappingSets found")
}

func (p *Parser) searchTopology(arPackagesElement *etree.Element) error {
	arPackages := arPackagesElement.SelectElements("AR-PACKAGE")
	for _, arPackage := range arPackages {
		sn, err := util.GetShortname(arPackage)
		if err != nil {
			return err
		}
		if sn == "Topology" {
			p.topologyElement = arPackage
			return nil
		}
	}
	return fmt.Errorf("no Topology found")
}

func (p *Parser) searchCommunication(arPackagesElement *etree.Element) error {
	arPackages := arPackagesElement.SelectElements("AR-PACKAGE")
	for _, arPackage := range arPackages {
		sn, err := util.GetShortname(arPackage)
		if err != nil {
			return err
		}
		if sn == "Communication" {
			p.communicationElement = arPackage
			return nil
		}
	}
	return fmt.Errorf("no Communication found")
}

func (p *Parser) searchSystem(arPackagesElement *etree.Element) error {
	arPackages := arPackagesElement.SelectElements("AR-PACKAGE")
	for _, arPackage := range arPackages {
		sn, err := util.GetShortname(arPackage)
		if err != nil {
			return err
		}
		if sn == "System" {
			p.systemElement = arPackage
			return nil
		}
	}
	return fmt.Errorf("no System found")
}

func (p *Parser) searchSoftwareTypes(arPackagesElement *etree.Element) error {
	arPackages := arPackagesElement.SelectElements("AR-PACKAGE")
	for _, arPackage := range arPackages {
		sn, err := util.GetShortname(arPackage)
		if err != nil {
			return err
		}
		if sn == "SoftwareTypes" {
			p.softwareTypesElement = arPackage
			return nil
		}
	}
	return fmt.Errorf("no SoftwareTypes found")
}

func (p *Parser) searchTpConfig(arPackagesElement *etree.Element) error {
	arPackages := arPackagesElement.SelectElements("AR-PACKAGE")
	for _, arPackage := range arPackages {
		sn, err := util.GetShortname(arPackage)
		if err != nil {
			return err
		}
		if sn == "TpConfig" {
			p.tpConfigElement = arPackage
			return nil
		}
	}
	return fmt.Errorf("no TpConfig found")
}

func (p *Parser) search(arPackages *etree.Element) error {
	if err := p.searchDataTypes(arPackages); err != nil {
		return fmt.Errorf("search data types: %w", err)
	}
	if err := p.searchDataTypeMappingSets(arPackages); err != nil {
		return fmt.Errorf("search data types mappings: %w", err)
	}
	if err := p.searchTopology(arPackages); err != nil {
		return fmt.Errorf("search topology: %w", err)
	}
	if err := p.searchCommunication(arPackages); err != nil {
		return fmt.Errorf("search communication: %w", err)
	}
	if err := p.searchSystem(arPackages); err != nil {
		return fmt.Errorf("search system: %w", err)
	}
	if err := p.searchSoftwareTypes(arPackages); err != nil {
		return fmt.Errorf("search software types: %w", err)
	}
	p.searchTpConfig(arPackages)
	return nil
}
