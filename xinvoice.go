// Package xinvoice helps convert GOBL into XRechnung and Factur-X documents.
package xinvoice

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
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
	XMLName                xml.Name           `xml:"rsm:CrossIndustryInvoice"`
	RSMNamespace           string             `xml:"xmlns:rsm,attr"`
	RAMNamespace           string             `xml:"xmlns:ram,attr"`
	QDTNamespace           string             `xml:"xmlns:qdt,attr"`
	UDTNamespace           string             `xml:"xmlns:udt,attr"`
	BusinessProcessContext string             `xml:"rsm:ExchangedDocumentContext>ram:BusinessProcessSpecifiedDocumentContextParameter>ram:ID"`
	GuidelineContext       string             `xml:"rsm:ExchangedDocumentContext>ram:GuidelineSpecifiedDocumentContextParameter>ram:ID"`
	ExchangedDocument      *ExchangedDocument `xml:"rsm:ExchangedDocument"`
	Transaction            *Transaction       `xml:"rsm:SupplyChainTradeTransaction"`
}

// ExchangedDocument a collection of data for a Cross Industry Invoice Header that is exchanged between two or more parties in written, printed or electronic form.
type ExchangedDocument struct {
	ID           string `xml:"ram:ID"`
	TypeCode     string `xml:"ram:TypeCode"`
	IssueDate    *Date  `xml:"ram:IssueDateTime>udt:DateTimeString"`
	IncludedNote *Note  `xml:"ram:IncludedNote"`
}

// Date defines date in the UDT structure
type Date struct {
	Date   string `xml:",chardata"`
	Format string `xml:"format,attr,omitempty"`
}

// Note defines note in the RAM structure
type Note struct {
	Content     string `xml:"ram:Content"`
	SubjectCode string `xml:"ram:SubjectCode"`
}

// NewDocument converts a GOBL envelope into a XRechnung and Factur-X document
func NewDocument(env *gobl.Envelope) (*Document, error) {
	inv, ok := env.Extract().(*bill.Invoice)
	if !ok {
		return nil, fmt.Errorf("invalid type %T", env.Document)
	}

	doc := Document{
		RSMNamespace:           RSM,
		RAMNamespace:           RAM,
		QDTNamespace:           QDT,
		UDTNamespace:           UDT,
		BusinessProcessContext: BusinessProcess,
		GuidelineContext:       GuidelineContext,
		ExchangedDocument:      newHeader(inv),
		Transaction:            NewTransaction(),
	}
	return &doc, nil
}

func newHeader(inv *bill.Invoice) *ExchangedDocument {
	typeCode := invoiceTypeCode(inv.Type)
	date := formatIssueDate(inv.IssueDate)

	return &ExchangedDocument{
		ID:       "123456XX",
		TypeCode: typeCode,
		IssueDate: &Date{
			Date:   date,
			Format: "102",
		},
		IncludedNote: &Note{
			Content:     "Es gelten unsere Allgem. Geschäftsbedingungen, die Sie unter […] finden.",
			SubjectCode: "ADU",
		},
	}
}

func formatIssueDate(date cal.Date) string {
	t := time.Date(date.Year, date.Month, date.Day, 0, 0, 0, 0, time.UTC)
	return t.Format("20060102")
}

// For German suppliers, the element "Invoice type code" (BT-3) should only contain the
// following values from code list UNTDID 1001:
// - 326 (Partial invoice)
// - 380 (Commercial invoice)
// - 384 (Corrected invoice)
// - 389 (Self-billed invoice)
// - 381 (Credit note)
func invoiceTypeCode(t cbc.Key) string {
	hash := map[string]string{
		"standard":    "380",
		"corrective":  "384",
		"credit-note": "381",
	}
	return hash[t.String()]
}

// Bytes returns the XML representation of the document in bytes
func (d *Document) Bytes() ([]byte, error) {
	bytes, err := xml.MarshalIndent(d, "", "  ")
	if err != nil {
		return nil, err
	}

	return append([]byte(xml.Header), bytes...), nil
}
