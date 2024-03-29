// Package xinvoice helps convert GOBL into XRechnung and Factur-X documents.
package xinvoice

import (
	"encoding/xml"
	"fmt"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
)

// CFDI schema constants
const (
	RSM              = "urn:un:unece:uncefact:data:standard:CrossIndustryInvoice:100"
	RAM              = "urn:un:unece:uncefact:data:standard:ReusableAggregateBusinessInformationEntity:100"
	QDT              = "urn:un:unece:uncefact:data:standard:QualifiedDataType:100"
	UDT              = "urn:un:unece:uncefact:data:standard:UnqualifiedDataType:100"
	BusinessProcess  = "urn:fdc:peppol.eu:2017:poacc:billing:01:1.0"
	GuidelineContext = "urn:cen.eu:en16931:2017#compliant#urn:xeinkauf.de:kosit:xrechnung_3.0"
)

// Document is a pseudo-model for containing the XML document being created
type Document struct {
	XMLName                xml.Name     `xml:"rsm:CrossIndustryInvoice"`
	RSMNamespace           string       `xml:"xmlns:rsm,attr"`
	RAMNamespace           string       `xml:"xmlns:ram,attr"`
	QDTNamespace           string       `xml:"xmlns:qdt,attr"`
	UDTNamespace           string       `xml:"xmlns:udt,attr"`
	BusinessProcessContext string       `xml:"rsm:ExchangedDocumentContext>ram:BusinessProcessSpecifiedDocumentContextParameter>ram:ID"`
	GuidelineContext       string       `xml:"rsm:ExchangedDocumentContext>ram:GuidelineSpecifiedDocumentContextParameter>ram:ID"`
	ExchangedDocument      *Header      `xml:"rsm:ExchangedDocument"`
	Transaction            *Transaction `xml:"rsm:SupplyChainTradeTransaction"`
}

// Date defines date in the UDT structure
type Date struct {
	Date   string `xml:",chardata"`
	Format string `xml:"format,attr,omitempty"`
}

// NewDocument converts a GOBL envelope into a XRechnung and Factur-X document
func NewDocument(env *gobl.Envelope) (*Document, error) {
	inv, ok := env.Extract().(*bill.Invoice)
	if !ok {
		return nil, fmt.Errorf("invalid type %T", env.Document)
	}

	transaction, err := NewTransaction(inv)
	if err != nil {
		return nil, err
	}

	doc := Document{
		RSMNamespace:           RSM,
		RAMNamespace:           RAM,
		QDTNamespace:           QDT,
		UDTNamespace:           UDT,
		BusinessProcessContext: BusinessProcess,
		GuidelineContext:       GuidelineContext,
		ExchangedDocument:      NewHeader(inv),
		Transaction:            transaction,
	}
	return &doc, nil
}

// Bytes returns the XML representation of the document in bytes
func (d *Document) Bytes() ([]byte, error) {
	bytes, err := xml.MarshalIndent(d, "", "  ")
	if err != nil {
		return nil, err
	}

	return append([]byte(xml.Header), bytes...), nil
}
