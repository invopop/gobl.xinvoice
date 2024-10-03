package to_gobl

import (
	"time"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
)

const (
	PaymentMeansCash               = "10"
	PaymentMeansCheque             = "20"
	PaymentMeansBankTransfer       = "30"
	PaymentMeansBankAccount        = "42"
	PaymentMeansCard               = "48"
	PaymentMeansDirectDebit        = "49"
	PaymentMeansStandingOrder      = "57"
	PaymentMeansSEPACreditTransfer = "58"
	PaymentMeansSEPADirectDebit    = "59"
	PaymentMeansReport             = "97"
)

const (
	StandardSalesTax  = "S"
	ZeroRatedGoodsTax = "Z"
	TaxExempt         = "E"
)

// Convert a date string to a cal.Date
func ParseDate(date string) cal.Date {
	t, err := time.Parse("20060102", date)
	if err != nil {
		return cal.Date{}
	}

	return cal.MakeDate(t.Year(), t.Month(), t.Day())
}

// Map X-Rechnung rate to GOBL equivalent
func FindTaxKey(taxType string) cbc.Key {
	switch taxType {
	case StandardSalesTax:
		return tax.RateStandard
	case ZeroRatedGoodsTax:
		return tax.RateZero
	case TaxExempt:
		return tax.RateExempt
	}
	return tax.RateStandard
}

// Map X-Rechnung invoice type to GOBL equivalent
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

// Map UN/ECE code to GOBL equivalent
func UnitFromUNECE(unece cbc.Code) org.Unit {
	for _, def := range org.UnitDefinitions {
		if def.UNECE == unece {
			return def.Unit
		}
	}
	// If no match is found, return the original UN/ECE code as a Unit
	return org.Unit(unece)
}

// Map X-Rechnung payment means to GOBL equivalent
func PaymentMeansTypeCodeParse(typeCode string) cbc.Key {
	switch typeCode {
	case PaymentMeansCash:
		return pay.MeansKeyCash
	case PaymentMeansCheque:
		return pay.MeansKeyCheque
	case PaymentMeansBankTransfer, PaymentMeansBankAccount:
		return pay.MeansKeyDebitTransfer
	case PaymentMeansCard:
		return pay.MeansKeyCard
	case PaymentMeansSEPACreditTransfer:
		return pay.MeansKeyCreditTransfer
	case PaymentMeansSEPADirectDebit, PaymentMeansDirectDebit:
		return pay.MeansKeyDirectDebit
	default:
		return pay.MeansKeyOther
	}
}
