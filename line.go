package xinvoice

import "github.com/invopop/gobl/bill"

// Line defines the structure of the IncludedSupplyChainTradeLineItem in the CII standard
type Line struct {
	ID              string           `xml:"ram:AssociatedDocumentLineDocument>ram:LineID"`
	Name            string           `xml:"ram:SpecifiedTradeProduct>ram:Name"`
	NetPrice        string           `xml:"ram:SpecifiedLineTradeAgreement>ram:NetPriceProductTradePrice>ram:ChargeAmount"`
	TradeDelivery   *Quantity        `xml:"ram:SpecifiedLineTradeDelivery>ram:BilledQuantity"`
	TradeSettlement *TradeSettlement `xml:"ram:SpecifiedLineTradeSettlement"`
}

// Quantity defines the structure of the quantity with its attributes for the CII standard
type Quantity struct {
	Amount   string `xml:",chardata"`
	UnitCode string `xml:"unitCode,attr"`
}

// TradeSettlement defines the structure of the SpecifiedLineTradeSettlement of the CII standard
type TradeSettlement struct {
	TaxType        string `xml:"ram:ApplicableTradeTax>ram:TypeCode"`
	TaxCode        string `xml:"ram:ApplicableTradeTax>ram:CategoryCode"`
	TaxRatePercent string `xml:"ram:ApplicableTradeTax>ram:RateApplicablePercent"`
	Sum            string `xml:"ram:SpecifiedTradeSettlementLineMonetarySummation>ram:LineTotalAmount"`
}

func newLine(line *bill.Line) *Line {
	return &Line{
		ID:       line.Item.Name,
		Name:     line.Item.Name,
		NetPrice: line.Item.Price.String(),
		TradeDelivery: &Quantity{
			Amount:   line.Quantity.String(),
			UnitCode: string(line.Item.Unit.UNECE()),
		},
		TradeSettlement: &TradeSettlement{
			TaxType: "VAT",
			// The sales tax category codes are as follows:
			// - S = Sales tax applies at the standard rate
			// - Z = goods taxable at the zero rate
			// - E = Tax exempt
			// - AE = Reversal of tax liability
			// - K = No sales tax is shown for intra-community deliveries
			// - G = Tax not charged due to export outside the EU
			// - O = Outside the tax scope
			// - L = IGIC (Canary Islands)
			// - M = IPSI (Ceuta/Melilla)
			TaxCode:        "S",
			TaxRatePercent: line.Taxes[0].Percent.StringWithoutSymbol(),
			Sum:            line.Total.String(),
		},
	}
}

// NewLines generates lines for the KSeF invoice
func NewLines(lines []*bill.Line) []*Line {
	var Lines []*Line

	for _, line := range lines {
		Lines = append(Lines, newLine(line))
	}

	return Lines
}
