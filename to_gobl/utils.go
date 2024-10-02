package to_gobl

import (
	"time"

	xinvoice "github.com/invopop/gobl.xinvoice/xinvoice"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

func ParseDate(date string) cal.Date {
	t, err := time.Parse("20060102", date)
	if err != nil {
		return cal.Date{}
	}

	return cal.MakeDate(t.Year(), t.Month(), t.Day())
}

func FindTaxKey(taxType string) cbc.Key {
	switch taxType {
	case xinvoice.StandardSalesTax:
		return tax.RateStandard
	case xinvoice.ZeroRatedGoodsTax:
		return tax.RateZero
	case xinvoice.TaxExempt:
		return tax.RateExempt
	}
	return tax.RateStandard
}

func TypeCodeParse(typeCode string) cbc.Key {
	switch typeCode {
	case "380":
		return bill.InvoiceTypeStandard
	case "381":
		return bill.InvoiceTypeCreditNote
	case "384":
		return bill.InvoiceTypeCorrective
	case "325":
		return bill.InvoiceTypeProforma
	case "383":
		return bill.InvoiceTypeDebitNote
	}
	return bill.InvoiceTypeOther
}

func UnitFromUNECE(unece cbc.Code) org.Unit {
	for _, def := range org.UnitDefinitions {
		if def.UNECE == unece {
			return def.Unit
		}
	}
	// If no match is found, return the original UN/ECE code as a Unit
	return org.Unit(unece)
}
