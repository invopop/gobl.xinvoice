package to_gobl

import (
	"github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
)

// NewDocument converts a XRechnung document into a GOBL envelope
func NewDocumentGOBL(doc *XMLDoc) (*gobl.Envelope, error) {

	inv := &bill.Invoice{
		Code:      cbc.Code(doc.ExchangedDocument.ID),
		Type:      TypeCodeParse(doc.ExchangedDocument.TypeCode),
		IssueDate: parseDate(doc.ExchangedDocument.IssueDateTime.DateTimeString.Value),
		Currency:  doc.SupplyChainTradeTransaction.ApplicableHeaderTradeSettlement.InvoiceCurrencyCode,
		Supplier:  parseParty(&doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.SellerTradeParty),
		Customer:  parseParty(&doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.BuyerTradeParty),
		Lines:     parseLines(&doc.SupplyChainTradeTransaction),
		// All 1..1
		Payment: parsePayment(&doc.SupplyChainTradeTransaction.ApplicableHeaderTradeSettlement),
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
		inv.Ordering = &bill.Ordering{
			Code: cbc.Code(*doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.BuyerReference),
		}
	}

	env, err := gobl.Envelop(inv)
	if err != nil {
		return nil, err
	}
	return env, nil
}
