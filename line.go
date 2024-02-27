package xinvoice

import (
	"fmt"

	"github.com/invopop/gobl/bill"
)

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
	ApplicableTradeTax []*ApplicableTradeTax `xml:"ram:ApplicableTradeTax"`
	Sum                string                `xml:"ram:SpecifiedTradeSettlementLineMonetarySummation>ram:LineTotalAmount"`
}

// ApplicableTradeTax defines the structure of ApplicableTradeTax of the CII standard
type ApplicableTradeTax struct {
	TaxType        string `xml:"ram:TypeCode"`
	TaxCode        string `xml:"ram:CategoryCode"`
	TaxRatePercent string `xml:"ram:RateApplicablePercent"`
}

func newLine(line *bill.Line) (*Line, error) {
	if line.Item == nil {
		return nil, fmt.Errorf("No item provided in line")
	}
	item := line.Item

	settlement, err := newTradeSettlement(line)
	if err != nil {
		return nil, err
	}

	lineItem := &Line{
		ID:       item.Name,
		Name:     item.Name,
		NetPrice: item.Price.String(),
		TradeDelivery: &Quantity{
			Amount:   line.Quantity.String(),
			UnitCode: string(item.Unit.UNECE()),
		},
		TradeSettlement: settlement,
	}

	return lineItem, nil
}

func newTradeSettlement(line *bill.Line) (*TradeSettlement, error) {
	var applicableTradeTax []*ApplicableTradeTax
	for _, tax := range line.Taxes {
		if tax.Percent == nil {
			return nil, fmt.Errorf("No Tax percent provided for item")
		}

		tradeTax := &ApplicableTradeTax{
			TaxType: tax.Category.String(),
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
			TaxRatePercent: tax.Percent.StringWithoutSymbol(),
		}

		applicableTradeTax = append(applicableTradeTax, tradeTax)
	}

	settlement := &TradeSettlement{
		ApplicableTradeTax: applicableTradeTax,
		Sum:                line.Total.String(),
	}

	return settlement, nil
}

// NewLines generates lines for XInvoice
func NewLines(lines []*bill.Line) ([]*Line, error) {
	var Lines []*Line

	for _, line := range lines {
		lineItem, err := newLine(line)
		if err != nil {
			return nil, err
		}

		Lines = append(Lines, lineItem)
	}

	return Lines, nil
}
