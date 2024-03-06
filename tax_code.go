package xinvoice

import "github.com/invopop/gobl/bill"

// StandardSalesTax is the tax type for sales tax applying at the standard rate
const StandardSalesTax = "S"

// ZeroRatedGoodsTax is the tax type for goods taxable at the zero rate
const ZeroRatedGoodsTax = "Z"

// TaxExempt indicates the tax type for exempt transactions
const TaxExempt = "E"

// ReverseCharge indicates the tax type for reverse charge transactions
const ReverseCharge = "AE"

// NoSalesTaxIntraCommunity indicates no sales tax is shown for intra-community deliveries
const NoSalesTaxIntraCommunity = "K"

// TaxNotChargedExportOutsideEU indicates tax is not charged due to export outside the EU
const TaxNotChargedExportOutsideEU = "G"

// OutOfTaxScope indicates transactions outside the tax scope
const OutOfTaxScope = "O"

// IGIC represents IGIC (Canary Islands)
const IGIC = "L"

// IPSI represents IPSI (Ceuta/Melilla)
const IPSI = "M"

// FindTaxCode finds the tax code for the provided bill line.
// It returns the found tax code.
func FindTaxCode(line *bill.Line) string {
	if len(line.Taxes) == 0 {
		return StandardSalesTax
	}
	tax := line.Taxes[0]

	switch tax.Rate {
	case "standard":
		return StandardSalesTax
	case "zero":
		return ZeroRatedGoodsTax
	}

	switch tax.Category {
	case "IGIC":
		return IGIC
	case "IPSI":
		return IPSI
	}

	return StandardSalesTax
}
