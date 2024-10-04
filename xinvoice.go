// Package xinvoice helps convert GOBL into XRechnung and Factur-X documents.
package xinvoice

import (
	"encoding/xml"
	"fmt"

	"github.com/invopop/gobl"
	gtox "github.com/invopop/gobl.xinvoice/internal/gtox"
	xtog "github.com/invopop/gobl.xinvoice/internal/xtog"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
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
	XMLName                xml.Name          `xml:"rsm:CrossIndustryInvoice"`
	RSMNamespace           string            `xml:"xmlns:rsm,attr"`
	RAMNamespace           string            `xml:"xmlns:ram,attr"`
	QDTNamespace           string            `xml:"xmlns:qdt,attr"`
	UDTNamespace           string            `xml:"xmlns:udt,attr"`
	BusinessProcessContext string            `xml:"rsm:ExchangedDocumentContext>ram:BusinessProcessSpecifiedDocumentContextParameter>ram:ID"`
	GuidelineContext       string            `xml:"rsm:ExchangedDocumentContext>ram:GuidelineSpecifiedDocumentContextParameter>ram:ID"`
	ExchangedDocument      *gtox.Header      `xml:"rsm:ExchangedDocument"`
	Transaction            *gtox.Transaction `xml:"rsm:SupplyChainTradeTransaction"`
}

// Model for XML excluding namespaces
type XMLDoc struct {
	XMLName                     xml.Name                         `xml:"CrossIndustryInvoice"`
	BusinessProcessContext      string                           `xml:"ExchangedDocumentContext>BusinessProcessSpecifiedDocumentContextParameter>ID"`
	GuidelineContext            string                           `xml:"ExchangedDocumentContext>GuidelineSpecifiedDocumentContextParameter>ID"`
	ExchangedDocument           xtog.ExchangedDocument           `xml:"ExchangedDocument"`
	SupplyChainTradeTransaction xtog.SupplyChainTradeTransaction `xml:"SupplyChainTradeTransaction"`
}

// NewDocument converts a GOBL envelope into a XRechnung and Factur-X document
func NewDocument(env *gobl.Envelope) (*Document, error) {
	inv, ok := env.Extract().(*bill.Invoice)
	if !ok {
		return nil, fmt.Errorf("invalid type %T", env.Document)
	}

	transaction, err := gtox.NewTransaction(inv)
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
		ExchangedDocument:      gtox.NewHeader(inv),
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

// NewDocument converts a XRechnung document into a GOBL envelope
func NewDocumentGOBL(doc *XMLDoc) (*gobl.Envelope, error) {

	inv := &bill.Invoice{
		Code:      cbc.Code(doc.ExchangedDocument.ID),
		Type:      xtog.TypeCodeParse(doc.ExchangedDocument.TypeCode),
		IssueDate: xtog.ParseDate(doc.ExchangedDocument.IssueDateTime.DateTimeString.Value),
		Currency:  currency.Code(doc.SupplyChainTradeTransaction.ApplicableHeaderTradeSettlement.InvoiceCurrencyCode),
		Supplier:  xtog.ParseParty(&doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.SellerTradeParty),
		Customer:  xtog.ParseParty(&doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.BuyerTradeParty),
		Lines:     xtog.ParseXMLLines(&doc.SupplyChainTradeTransaction),
	}

	// Payment comprised of terms, means and payee. Check tehre is relevant info in at least one of them to create a payment
	if doc.SupplyChainTradeTransaction.ApplicableHeaderTradeSettlement.PayeeTradeParty != nil ||
		(len(doc.SupplyChainTradeTransaction.ApplicableHeaderTradeSettlement.SpecifiedTradePaymentTerms) > 0 &&
			doc.SupplyChainTradeTransaction.ApplicableHeaderTradeSettlement.SpecifiedTradePaymentTerms[0].DueDateDateTime != nil) ||
		(len(doc.SupplyChainTradeTransaction.ApplicableHeaderTradeSettlement.SpecifiedTradeSettlementPaymentMeans) > 0 &&
			doc.SupplyChainTradeTransaction.ApplicableHeaderTradeSettlement.SpecifiedTradeSettlementPaymentMeans[0].TypeCode != "1") {
		inv.Payment = xtog.ParsePayment(&doc.SupplyChainTradeTransaction.ApplicableHeaderTradeSettlement)
	}

	if len(doc.ExchangedDocument.IncludedNote) > 0 {
		inv.Notes = make([]*cbc.Note, 0, len(doc.ExchangedDocument.IncludedNote))
		for _, note := range doc.ExchangedDocument.IncludedNote {
			n := &cbc.Note{}
			if note.Content != "" {
				n.Text = note.Content
			}
			if note.ContentCode != "" {
				n.Code = note.ContentCode
			}
			inv.Notes = append(inv.Notes, n)
		}
	}

	if doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.BuyerReference != nil {
		if *doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.BuyerReference != "N/A" {
			inv.Ordering = &bill.Ordering{
				Code: cbc.Code(*doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.BuyerReference),
			}
		}
	}

	env, err := gobl.Envelop(inv)
	if err != nil {
		return nil, err
	}
	return env, nil
}
