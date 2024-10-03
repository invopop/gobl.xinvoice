package to_gobl

import (
	"github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
)

// NewDocument converts a XRechnung document into a GOBL envelope
func NewDocumentGOBL(doc *XMLDoc) (*gobl.Envelope, error) {

	inv := &bill.Invoice{
		Code:      cbc.Code(doc.ExchangedDocument.ID),
		Type:      TypeCodeParse(doc.ExchangedDocument.TypeCode),
		IssueDate: ParseDate(doc.ExchangedDocument.IssueDateTime.DateTimeString.Value),
		Currency:  currency.Code(doc.SupplyChainTradeTransaction.ApplicableHeaderTradeSettlement.InvoiceCurrencyCode),
		Supplier:  ParseParty(&doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.SellerTradeParty),
		Customer:  ParseParty(&doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.BuyerTradeParty),
		Lines:     ParseXMLLines(&doc.SupplyChainTradeTransaction),
	}

	// Payment comprised of terms, means and payee. Check tehre is relevant info in at least one of them to create a payment
	if doc.SupplyChainTradeTransaction.ApplicableHeaderTradeSettlement.PayeeTradeParty != nil ||
		(len(doc.SupplyChainTradeTransaction.ApplicableHeaderTradeSettlement.SpecifiedTradePaymentTerms) > 0 &&
			doc.SupplyChainTradeTransaction.ApplicableHeaderTradeSettlement.SpecifiedTradePaymentTerms[0].DueDateDateTime != nil) ||
		(len(doc.SupplyChainTradeTransaction.ApplicableHeaderTradeSettlement.SpecifiedTradeSettlementPaymentMeans) > 0 &&
			doc.SupplyChainTradeTransaction.ApplicableHeaderTradeSettlement.SpecifiedTradeSettlementPaymentMeans[0].TypeCode != "1") {
		inv.Payment = ParsePayment(&doc.SupplyChainTradeTransaction.ApplicableHeaderTradeSettlement)
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
