package xinvoice

import (
	"encoding/xml"
	// "encoding/json"

	"time"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
)

// Model for XML excluding namespaces
type XMLDoc struct {
	XMLName                     xml.Name                    `xml:"CrossIndustryInvoice"`
	BusinessProcessContext      string                      `xml:"ExchangedDocumentContext>BusinessProcessSpecifiedDocumentContextParameter>ID"`
	GuidelineContext            string                      `xml:"ExchangedDocumentContext>GuidelineSpecifiedDocumentContextParameter>ID"`
	ExchangedDocument           ExchangedDocument           `xml:"ExchangedDocument"`
	SupplyChainTradeTransaction SupplyChainTradeTransaction `xml:"SupplyChainTradeTransaction"`
}

type ExchangedDocument struct {
	ID            string `xml:"ID"`
	Name          string `xml:"Name"`
	TypeCode      string `xml:"TypeCode"`
	IssueDateTime struct {
		DateTimeString struct {
			Value  string `xml:",chardata"`
			Format string `xml:"format,attr"`
		} `xml:"DateTimeString"`
	} `xml:"IssueDateTime"`
	IncludedNote []IncludedNote `xml:"IncludedNote"`
}

type SupplyChainTradeTransaction struct {
	IncludedSupplyChainTradeLineItem []IncludedSupplyChainTradeLineItem `xml:"IncludedSupplyChainTradeLineItem"`
	ApplicableHeaderTradeAgreement   ApplicableHeaderTradeAgreement     `xml:"ApplicableHeaderTradeAgreement"`
	ApplicableHeaderTradeDelivery    struct{}                           `xml:"ApplicableHeaderTradeDelivery"`
	ApplicableHeaderTradeSettlement  ApplicableHeaderTradeSettlement    `xml:"ApplicableHeaderTradeSettlement"`
}

type IncludedSupplyChainTradeLineItem struct {
	AssociatedDocumentLineDocument struct {
		LineID       int            `xml:"LineID"`
		IncludedNote []IncludedNote `xml:"IncludedNote"`
	} `xml:"AssociatedDocumentLineDocument"`
	SpecifiedTradeProduct struct {
		Name string `xml:"Name"`
	} `xml:"SpecifiedTradeProduct"`
	SpecifiedLineTradeAgreement struct {
		NetPriceProductTradePrice struct {
			ChargeAmount float64 `xml:"ChargeAmount"`
		} `xml:"NetPriceProductTradePrice"`
	} `xml:"SpecifiedLineTradeAgreement"`
	SpecifiedLineTradeDelivery *struct {
		BilledQuantity struct {
			Value    int64  `xml:",chardata"`
			UnitCode string `xml:"unitCode,attr"`
		} `xml:"BilledQuantity"`
	} `xml:"SpecifiedLineTradeDelivery"`
	SpecifiedLineTradeSettlement struct {
		ApplicableTradeTax struct {
			TypeCode              string `xml:"TypeCode"`
			CategoryCode          string `xml:"CategoryCode"`
			RateApplicablePercent string `xml:"RateApplicablePercent"`
		} `xml:"ApplicableTradeTax"`
		SpecifiedTradeSettlementLineMonetarySummation struct {
			LineTotalAmount float64 `xml:"LineTotalAmount"`
		} `xml:"SpecifiedTradeSettlementLineMonetarySummation"`
	} `xml:"SpecifiedLineTradeSettlement"`
}

type ApplicableHeaderTradeAgreement struct {
	BuyerReference   *string    `xml:"BuyerReference"`
	SellerTradeParty TradeParty `xml:"SellerTradeParty"`
	BuyerTradeParty  TradeParty `xml:"BuyerTradeParty"`
}

type TradeParty struct {
	ID                  string `xml:"ID,omitempty"`
	Name                string `xml:"Name"`
	DefinedTradeContact *struct {
		PersonName                      string `xml:"PersonName"`
		TelephoneUniversalCommunication *struct {
			CompleteNumber string `xml:"CompleteNumber"`
		} `xml:"TelephoneUniversalCommunication"`
		EmailURIUniversalCommunication *struct {
			URIID string `xml:"URIID"`
		} `xml:"EmailURIUniversalCommunication"`
	} `xml:"DefinedTradeContact,omitempty"`
	PostalTradeAddress *struct {
		PostcodeCode string              `xml:"PostcodeCode"`
		LineOne      string              `xml:"LineOne"`
		CityName     string              `xml:"CityName"`
		CountryID    l10n.ISOCountryCode `xml:"CountryID"`
	} `xml:"PostalTradeAddress"`
	URIUniversalCommunication struct {
		URIID struct {
			Value    string `xml:",chardata"`
			SchemeID string `xml:"schemeID,attr"`
		} `xml:"URIID"`
	} `xml:"URIUniversalCommunication"`
	SpecifiedTaxRegistration *struct {
		ID *struct {
			Value string `xml:",chardata"`
			//VA used for VAT-ID used in B2B, FC for tax number (Steuernummer).
			SchemeID string `xml:"schemeID,attr"`
		} `xml:"ID"`
	} `xml:"SpecifiedTaxRegistration,omitempty"`
}

type ApplicableHeaderTradeSettlement struct {
	InvoiceCurrencyCode                  currency.Code `xml:"InvoiceCurrencyCode"`
	SpecifiedTradeSettlementPaymentMeans struct {
		TypeCode string `xml:"TypeCode"`
	} `xml:"SpecifiedTradeSettlementPaymentMeans"`
	ApplicableTradeTax []struct {
		CalculatedAmount      string `xml:"CalculatedAmount"`
		TypeCode              string `xml:"TypeCode"`
		BasisAmount           string `xml:"BasisAmount"`
		CategoryCode          string `xml:"CategoryCode"`
		RateApplicablePercent string `xml:"RateApplicablePercent"`
	} `xml:"ApplicableTradeTax"`
	SpecifiedTradeSettlementHeaderMonetarySummation struct {
		LineTotalAmount     float64 `xml:"LineTotalAmount"`
		TaxBasisTotalAmount string  `xml:"TaxBasisTotalAmount"`
		TaxTotalAmount      struct {
			Value      string `xml:",chardata"`
			CurrencyID string `xml:"currencyID,attr"`
		} `xml:"TaxTotalAmount"`
		GrandTotalAmount string `xml:"GrandTotalAmount"`
		DuePayableAmount string `xml:"DuePayableAmount"`
	} `xml:"SpecifiedTradeSettlementHeaderMonetarySummation"`
	PayeeTradeParty            *TradeParty `xml:"PayeeTradeParty"`
	SpecifiedTradePaymentTerms []struct {
		Description     string `xml:"Description"`
		DueDateDateTime struct {
			DateTimeString string `xml:"DateTimeString"`
			Format         string `xml:"format,attr"`
		} `xml:"DueDateDateTime"`
		PartialPaymentAmount               *string `xml:"PartialPaymentAmount"`
		ApplicableTradePaymentPenaltyTerms struct {
			BasisAmount string `xml:"BasisAmount"`
		} `xml:"ApplicableTradePaymentPenaltyTerms"`
		ApplicableTradePaymentDiscountTerms struct {
			BasisAmount string `xml:"BasisAmount"`
		} `xml:"ApplicableTradePaymentDiscountTerms"`
	} `xml:"SpecifiedTradePaymentTerms"`
	SpecifiedAdvancePayment struct {
		PaidAmount                float64          `xml:"PaidAmount"`
		FormattedReceivedDateTime DateTimeFormat   `xml:"FormattedReceivedDateTime"`
		IncludedTradeTax          IncludedTradeTax `xml:"IncludedTradeTax"`
	} `xml:"SpecifiedAdvancePayment"`
}

type DateTimeFormat struct {
	DateTimeString string `xml:"DateTimeString"`
	Format         string `xml:"format,attr"`
}

type IncludedTradeTax struct {
	CalculatedAmount      string `xml:"CalculatedAmount"`
	TypeCode              string `xml:"TypeCode"`
	BasisAmount           string `xml:"BasisAmount"`
	CategoryCode          string `xml:"CategoryCode"`
	RateApplicablePercent string `xml:"RateApplicablePercent"`
}

type IncludedNote struct {
	ContentCode string `xml:"ContentCode"`
	Content     string `xml:"Content"`
	SubjectCode string `xml:"SubjectCode"`
}

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

// Bytes returns the XML representation of the document in bytes
func (d *Document) BytesGOBL() ([]byte, error) {
	bytes, err := xml.MarshalIndent(d, "", "  ")
	if err != nil {
		return nil, err
	}

	return append([]byte(xml.Header), bytes...), nil
}

func parseParty(party *TradeParty) *org.Party {
	p := &org.Party{
		Name: party.Name,
	}

	if party.DefinedTradeContact != nil && party.DefinedTradeContact.PersonName != "" {
		p.People = []*org.Person{
			{
				Name: &org.Name{
					Given: party.DefinedTradeContact.PersonName,
				},
			},
		}
	}

	if party.PostalTradeAddress != nil {
		p.Addresses = []*org.Address{
			{
				Street:   party.PostalTradeAddress.LineOne,
				Locality: party.PostalTradeAddress.CityName,
				Code:     party.PostalTradeAddress.PostcodeCode,
				Country:  party.PostalTradeAddress.CountryID,
			},
		}
	}

	if party.DefinedTradeContact != nil && party.DefinedTradeContact.TelephoneUniversalCommunication != nil {
		p.Telephones = []*org.Telephone{
			{
				Number: party.DefinedTradeContact.TelephoneUniversalCommunication.CompleteNumber,
			},
		}
	}

	if party.DefinedTradeContact != nil && party.DefinedTradeContact.EmailURIUniversalCommunication != nil {
		p.Emails = []*org.Email{
			{
				Address: party.DefinedTradeContact.EmailURIUniversalCommunication.URIID,
			},
		}
	}

	if party.SpecifiedTaxRegistration != nil && party.SpecifiedTaxRegistration.ID != nil {
		p.TaxID = &tax.Identity{
			Country: l10n.TaxCountryCode(party.PostalTradeAddress.CountryID),
			Code:    cbc.Code(party.SpecifiedTaxRegistration.ID.Value),
		}
	}

	return p
}

func parseDate(date string) cal.Date {
	t, err := time.Parse("20060102", date)
	if err != nil {
		return cal.Date{}
	}

	return cal.MakeDate(t.Year(), t.Month(), t.Day())
}

func parsePayment(settlement *ApplicableHeaderTradeSettlement) *bill.Payment {
	payment := &bill.Payment{}

	if settlement.PayeeTradeParty != nil {
		payee := &org.Party{Name: settlement.PayeeTradeParty.Name}
		if settlement.PayeeTradeParty.PostalTradeAddress.LineOne != "" {
			payee.Addresses = []*org.Address{
				{
					Street:   settlement.PayeeTradeParty.PostalTradeAddress.LineOne,
					Locality: settlement.PayeeTradeParty.PostalTradeAddress.CityName,
					Code:     settlement.PayeeTradeParty.PostalTradeAddress.PostcodeCode,
					Country:  settlement.PayeeTradeParty.PostalTradeAddress.CountryID,
				},
			}
		}
		payment.Payee = payee
	}
	if len(settlement.SpecifiedTradePaymentTerms) > 0 {
		terms := &pay.Terms{}
		var dueDates []*pay.DueDate

		for _, paymentTerm := range settlement.SpecifiedTradePaymentTerms {
			if terms.Detail == "" {
				terms.Detail = paymentTerm.Description
			}

			if paymentTerm.DueDateDateTime.DateTimeString != "" {
				paymentTermsDueDateDateTime := parseDate(paymentTerm.DueDateDateTime.DateTimeString)
				dueDate := &pay.DueDate{
					Date: &paymentTermsDueDateDateTime,
				}
				if paymentTerm.PartialPaymentAmount != nil {
					dueDate.Amount, _ = num.AmountFromString(*paymentTerm.PartialPaymentAmount)
				}
				dueDates = append(dueDates, dueDate)
			}
		}

		terms.DueDates = dueDates
		payment.Terms = terms
	}

	if settlement.SpecifiedAdvancePayment.FormattedReceivedDateTime.DateTimeString != "" {
		advancePaymentReceivedDateTime := parseDate(settlement.SpecifiedAdvancePayment.FormattedReceivedDateTime.DateTimeString)
		advance := &pay.Advance{
			Amount: num.AmountFromFloat64(settlement.SpecifiedAdvancePayment.PaidAmount, 0),
			Date:   &advancePaymentReceivedDateTime,
		}
		payment.Advances = []*pay.Advance{advance}
	}

	return payment
}

func parseLines(transaction *SupplyChainTradeTransaction) []*bill.Line {
	items := transaction.IncludedSupplyChainTradeLineItem
	lines := make([]*bill.Line, 0, len(transaction.IncludedSupplyChainTradeLineItem))

	for _, item := range items {
		price := num.AmountFromFloat64(item.SpecifiedLineTradeAgreement.NetPriceProductTradePrice.ChargeAmount, 0)

		line := &bill.Line{
			// Index:    item.AssociatedDocumentLineDocument.LineID,
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  item.SpecifiedTradeProduct.Name,
				Price: price,
				// Unit:  org.Unit(org.Unit(item.SpecifiedLineTradeDelivery.BilledQuantity.UnitCode).UNECE()),
			},
			Taxes: tax.Set{
				{
					Rate:     findTaxKey(item.SpecifiedLineTradeSettlement.ApplicableTradeTax.CategoryCode),
					Category: cbc.Code(item.SpecifiedLineTradeSettlement.ApplicableTradeTax.TypeCode),
					// Percent:  &percent,
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

		// Set the unit if available
		if item.SpecifiedLineTradeDelivery.BilledQuantity.UnitCode != "" {
			line.Item.Unit = UnitFromUNECE(cbc.Code(item.SpecifiedLineTradeDelivery.BilledQuantity.UnitCode))
		}

		lines = append(lines, line)
	}

	return lines
}

func findTaxKey(taxType string) cbc.Key {
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
