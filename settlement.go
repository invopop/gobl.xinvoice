package xinvoice

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

// Settlement defines the structure of ApplicableHeaderTradeSettlement of the CII standard
type Settlement struct {
	Currency    string   `xml:"ram:InvoiceCurrencyCode"`
	TypeCode    string   `xml:"ram:SpecifiedTradeSettlementPaymentMeans>ram:TypeCode"`
	Tax         []*Tax   `xml:"ram:ApplicableTradeTax"`
	Description string   `xml:"ram:SpecifiedTradePaymentTerms>ram:Description,omitempty"`
	Summary     *Summary `xml:"ram:SpecifiedTradeSettlementHeaderMonetarySummation"`
}

// Tax defines the structure of ApplicableTradeTax of the CII standard
type Tax struct {
	CalculatedAmount      string `xml:"ram:CalculatedAmount"`
	TypeCode              string `xml:"ram:TypeCode"`
	BasisAmount           string `xml:"ram:BasisAmount"`
	CategoryCode          string `xml:"ram:CategoryCode"`
	RateApplicablePercent string `xml:"ram:RateApplicablePercent"`
}

// Summary defines the structure of SpecifiedTradeSettlementHeaderMonetarySummation of the CII standard
type Summary struct {
	TotalAmount         string          `xml:"ram:LineTotalAmount"`
	TaxBasisTotalAmount string          `xml:"ram:TaxBasisTotalAmount"`
	TaxTotalAmount      *TaxTotalAmount `xml:"ram:TaxTotalAmount"`
	GrandTotalAmount    string          `xml:"ram:GrandTotalAmount"`
	DuePayableAmount    string          `xml:"ram:DuePayableAmount"`
}

// TaxTotalAmount defines the structure of the TaxTotalAmount of the CII standard
type TaxTotalAmount struct {
	Amount   string `xml:",chardata"`
	Currency string `xml:"currencyID,attr"`
}

// NewSettlement creates the ApplicableHeaderTradeSettlement part of a EN 16931 compliant invoice
func NewSettlement(inv *bill.Invoice) (*Settlement, error) {
	if inv.Totals == nil {
		return nil, fmt.Errorf("Totals not provided")
	}

	taxes, err := newTaxes(inv.Totals.Taxes)
	if err != nil {
		return nil, err
	}

	settlement := &Settlement{
		Currency:    string(inv.Currency),
		TypeCode:    "1",
		Description: inv.Payment.Terms.Detail,
		Tax:         taxes,
		Summary: &Summary{
			TotalAmount:         inv.Totals.Total.String(),
			TaxBasisTotalAmount: inv.Totals.Total.String(),
			GrandTotalAmount:    inv.Totals.TotalWithTax.String(),
			DuePayableAmount:    inv.Totals.Payable.String(),
			TaxTotalAmount: &TaxTotalAmount{
				Amount:   inv.Totals.Tax.String(),
				Currency: string(inv.Currency),
			},
		},
	}

	return settlement, nil
}

func newTaxes(total *tax.Total) ([]*Tax, error) {
	var Taxes []*Tax

	if total == nil {
		return nil, fmt.Errorf("Total taxes not provided")
	}

	for _, category := range total.Categories {
		for _, rate := range category.Rates {
			tax, err := newTax(rate)
			if err != nil {
				return nil, err
			}

			Taxes = append(Taxes, tax)
		}
	}

	return Taxes, nil
}

func newTax(rate *tax.RateTotal) (*Tax, error) {
	if rate.Percent == nil {
		return nil, fmt.Errorf("No tax rate percent provided")
	}
	percent := rate.Percent.StringWithoutSymbol()

	tax := &Tax{
		CalculatedAmount:      rate.Amount.String(),
		TypeCode:              "VAT",
		BasisAmount:           rate.Base.String(),
		CategoryCode:          taxCategoryCode(rate.Key),
		RateApplicablePercent: percent,
	}

	return tax, nil
}

// AE - VAT Reverse Charge
// E  - Exempt from tax
// G  - Free export item, tax not charged
// K  - VAT exempt for EEA intra-community supply of goods and services
// L  - Canary Islands general indirect tax
// M  - Tax for production, services and importation in Ceuta and Melilla
// O  - Services outside scope of tax
// S  - Standard rate
// Z  - Zero rated goods
func taxCategoryCode(key cbc.Key) string {
	hash := map[cbc.Key]string{
		tax.RateStandard: "S",
		tax.RateReduced:  "S",
		tax.RateZero:     "Z",
	}

	return hash[key]
}
