package xinvoice

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/tax"
)

const (
	// StandardSalesTax is the tax type for sales tax applying at the standard rate
	StandardSalesTax = "S"

	// ZeroRatedGoodsTax is the tax type for goods taxable at the zero rate
	ZeroRatedGoodsTax = "Z"

	// TaxExempt indicates the tax type for exempt transactions
	TaxExempt = "E"

	// ReverseCharge indicates the tax type for reverse charge transactions
	ReverseCharge = "AE"

	// NoSalesTaxIntraCommunity indicates no sales tax is shown for intra-community deliveries
	NoSalesTaxIntraCommunity = "K"

	// TaxNotChargedExportOutsideEU indicates tax is not charged due to export outside the EU
	TaxNotChargedExportOutsideEU = "G"

	// OutOfTaxScope indicates transactions outside the tax scope
	OutOfTaxScope = "O"

	// IGIC represents IGIC (Canary Islands)
	IGIC = "L"

	// IPSI represents IPSI (Ceuta/Melilla)
	IPSI = "M"
)

// FindTaxCode finds the tax code for the provided bill line.
// It returns the found tax code.
//
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
func FindTaxCode(line *bill.Line) string {
	if len(line.Taxes) == 0 {
		return StandardSalesTax
	}
	t := line.Taxes[0]

	switch t.Rate {
	case tax.RateStandard:
		return StandardSalesTax
	case tax.RateZero:
		return ZeroRatedGoodsTax
	case tax.RateExempt:
		return TaxExempt
	}

	switch t.Category {
	case "IGIC":
		return IGIC
	case "IPSI":
		return IPSI
	}

	return StandardSalesTax
}
