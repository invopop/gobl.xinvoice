package xinvoice

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// Parses the XML information for a Lines object
func ParseXMLLines(transaction *SupplyChainTradeTransaction) []*bill.Line {
	items := transaction.IncludedSupplyChainTradeLineItem
	lines := make([]*bill.Line, 0, len(transaction.IncludedSupplyChainTradeLineItem))

	for _, item := range items {
		price := num.AmountFromFloat64(item.SpecifiedLineTradeAgreement.NetPriceProductTradePrice.ChargeAmount, 0)

		line := &bill.Line{
			// Index:    item.AssociatedDocumentLineDocument.LineID, //generated field
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  item.SpecifiedTradeProduct.Name,
				Price: price,
			},
			Taxes: tax.Set{
				{
					Rate:     FindTaxKey(item.SpecifiedLineTradeSettlement.ApplicableTradeTax.CategoryCode),
					Category: cbc.Code(item.SpecifiedLineTradeSettlement.ApplicableTradeTax.TypeCode),
				},
			},
		}

		if item.SpecifiedLineTradeDelivery != nil {
			line.Quantity = num.MakeAmount(item.SpecifiedLineTradeDelivery.BilledQuantity.Value, 0)
		}

		if len(item.AssociatedDocumentLineDocument.IncludedNote) > 0 {
			line.Notes = make([]*cbc.Note, 0, len(item.AssociatedDocumentLineDocument.IncludedNote))
			for _, note := range item.AssociatedDocumentLineDocument.IncludedNote {
				n := &cbc.Note{}
				if note.Content != "" {
					n.Text = note.Content
				}
				if note.ContentCode != "" {
					n.Code = note.ContentCode
				}
				line.Notes = append(line.Notes, n)
			}
		}

		if item.SpecifiedLineTradeSettlement.ApplicableTradeTax.RateApplicablePercent != "" {
			percent, _ := num.PercentageFromString(item.SpecifiedLineTradeSettlement.ApplicableTradeTax.RateApplicablePercent)
			line.Taxes[0].Percent = &percent
		}

		if item.SpecifiedLineTradeDelivery.BilledQuantity.UnitCode != "" {
			line.Item.Unit = UnitFromUNECE(cbc.Code(item.SpecifiedLineTradeDelivery.BilledQuantity.UnitCode))
		}

		lines = append(lines, line)
	}

	return lines
}
