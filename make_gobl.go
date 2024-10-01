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
}

type SupplyChainTradeTransaction struct {
	IncludedSupplyChainTradeLineItem []IncludedSupplyChainTradeLineItem `xml:"IncludedSupplyChainTradeLineItem"`
	ApplicableHeaderTradeAgreement   ApplicableHeaderTradeAgreement     `xml:"ApplicableHeaderTradeAgreement"`
	ApplicableHeaderTradeDelivery    struct{}                           `xml:"ApplicableHeaderTradeDelivery"`
	ApplicableHeaderTradeSettlement  ApplicableHeaderTradeSettlement    `xml:"ApplicableHeaderTradeSettlement"`
}

type IncludedSupplyChainTradeLineItem struct {
	AssociatedDocumentLineDocument struct {
		LineID int `xml:"LineID"`
	} `xml:"AssociatedDocumentLineDocument"`
	SpecifiedTradeProduct struct {
		Name string `xml:"Name"`
	} `xml:"SpecifiedTradeProduct"`
	SpecifiedLineTradeAgreement struct {
		NetPriceProductTradePrice struct {
			ChargeAmount float64 `xml:"ChargeAmount"`
		} `xml:"NetPriceProductTradePrice"`
	} `xml:"SpecifiedLineTradeAgreement"`
	SpecifiedLineTradeDelivery struct {
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
	BuyerReference   string     `xml:"BuyerReference"`
	SellerTradeParty TradeParty `xml:"SellerTradeParty"`
	BuyerTradeParty  TradeParty `xml:"BuyerTradeParty"`
}

type TradeParty struct {
	ID                  string `xml:"ID,omitempty"`
	Name                string `xml:"Name"`
	DefinedTradeContact struct {
		PersonName                      string `xml:"PersonName"`
		TelephoneUniversalCommunication struct {
			CompleteNumber string `xml:"CompleteNumber"`
		} `xml:"TelephoneUniversalCommunication"`
		EmailURIUniversalCommunication struct {
			URIID string `xml:"URIID"`
		} `xml:"EmailURIUniversalCommunication"`
	} `xml:"DefinedTradeContact,omitempty"`
	PostalTradeAddress struct {
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
	SpecifiedTaxRegistration struct {
		ID struct {
			Value string `xml:",chardata"`
			//VA used for VAT-ID used in B2B, FC for tax number.
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
	PayeeTradeParty            TradeParty `xml:"PayeeTradeParty"`
	SpecifiedTradePaymentTerms struct {
		Description     string `xml:"Description"`
		DueDateDateTime struct {
			DateTimeString string `xml:"DateTimeString"`
			Format         string `xml:"format,attr"`
		} `xml:"DueDateDateTime"`
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

// NewDocument converts a XRechnung document into a GOBL envelope
func NewDocumentGOBL(doc *XMLDoc) (*gobl.Envelope, error) {

	PaymentTermsDueDateDateTime := parseDate(doc.SupplyChainTradeTransaction.ApplicableHeaderTradeSettlement.SpecifiedTradePaymentTerms.DueDateDateTime.DateTimeString)
	AdvancePaymentReceivedDateTime := parseDate(doc.SupplyChainTradeTransaction.ApplicableHeaderTradeSettlement.SpecifiedAdvancePayment.FormattedReceivedDateTime.DateTimeString)
	inv := &bill.Invoice{
		Code:      cbc.Code(doc.ExchangedDocument.ID),
		Type:      "standard",
		IssueDate: parseDate(doc.ExchangedDocument.IssueDateTime.DateTimeString.Value),
		Currency:  doc.SupplyChainTradeTransaction.ApplicableHeaderTradeSettlement.InvoiceCurrencyCode,
		Supplier: &org.Party{
			Name: doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.SellerTradeParty.Name,
			Addresses: []*org.Address{
				{
					Street:   doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.SellerTradeParty.PostalTradeAddress.LineOne,
					Locality: doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.SellerTradeParty.PostalTradeAddress.CityName,
					Code:     doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.SellerTradeParty.PostalTradeAddress.PostcodeCode,
					Country:  doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.SellerTradeParty.PostalTradeAddress.CountryID,
				},
			},
			TaxID: &tax.Identity{
				Country: l10n.TaxCountryCode(doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.SellerTradeParty.PostalTradeAddress.CountryID),
				Code:    cbc.Code(doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.SellerTradeParty.SpecifiedTaxRegistration.ID.Value),
			},
		},
		Customer: &org.Party{
			Name: doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.BuyerTradeParty.Name,
			Addresses: []*org.Address{
				{
					Street:   doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.BuyerTradeParty.PostalTradeAddress.LineOne,
					Locality: doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.BuyerTradeParty.PostalTradeAddress.CityName,
					Code:     doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.BuyerTradeParty.PostalTradeAddress.PostcodeCode,
					Country:  doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.BuyerTradeParty.PostalTradeAddress.CountryID,
				},
			},
		},
		Lines: parseLines(&doc.SupplyChainTradeTransaction),
		Payment: &bill.Payment{
			Payee: &org.Party{
				Name: doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.BuyerTradeParty.Name,
				Addresses: []*org.Address{
					{
						Street:   doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.BuyerTradeParty.PostalTradeAddress.LineOne,
						Locality: doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.BuyerTradeParty.PostalTradeAddress.CityName,
						Code:     doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.BuyerTradeParty.PostalTradeAddress.PostcodeCode,
						Country:  doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.BuyerTradeParty.PostalTradeAddress.CountryID,
					},
				},
			},
			Terms: &pay.Terms{
				Detail: doc.SupplyChainTradeTransaction.ApplicableHeaderTradeSettlement.SpecifiedTradePaymentTerms.Description,
				DueDates: []*pay.DueDate{
					{
						Date: &PaymentTermsDueDateDateTime,
					},
				},
			},
			Advances: []*pay.Advance{
				{
					Amount: num.AmountFromFloat64(doc.SupplyChainTradeTransaction.ApplicableHeaderTradeSettlement.SpecifiedAdvancePayment.PaidAmount, 0),
					Date:   &AdvancePaymentReceivedDateTime,
				},
			},
			Instructions: nil,
		},
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

func parseDate(date string) cal.Date {
	t, err := time.Parse("20060102", date)
	if err != nil {
		return cal.Date{}
	}
	return cal.MakeDate(t.Day(), t.Month(), t.Year())
}

func parseLines(transaction *SupplyChainTradeTransaction) []*bill.Line {
	items := transaction.IncludedSupplyChainTradeLineItem
	lines := make([]*bill.Line, 0, len(transaction.IncludedSupplyChainTradeLineItem))

	for _, item := range items {
		quantity := num.MakeAmount(item.SpecifiedLineTradeDelivery.BilledQuantity.Value, 0)
		price := num.AmountFromFloat64(item.SpecifiedLineTradeAgreement.NetPriceProductTradePrice.ChargeAmount, 0)
		percent, _ := num.PercentageFromString(item.SpecifiedLineTradeSettlement.ApplicableTradeTax.RateApplicablePercent)
		// total := num.MakeAmount(item.SpecifiedLineTradeSettlement.SpecifiedTradeSettlementLineMonetarySummation.LineTotalAmount, 0)
		//discount := num.MakePercent(0, 0)

		line := &bill.Line{
			Index:    item.AssociatedDocumentLineDocument.LineID,
			Quantity: quantity,
			Item: &org.Item{
				Name:  item.SpecifiedTradeProduct.Name,
				Price: price,
			},
			// Sum: total,
			// Total: total,
			Taxes: tax.Set{
				{
					Category: cbc.Code(item.SpecifiedLineTradeSettlement.ApplicableTradeTax.CategoryCode),
					Rate:     findTaxKey(item.SpecifiedLineTradeSettlement.ApplicableTradeTax.TypeCode),
					Percent:  &percent,
				},
			},
			// Notes: []*cbc.Note{
			// 	{
			// 		Content: cbc.Code(item.SpecifiedLineTradeSettlement.ApplicableTradeTax.TypeCode),
			// 	},
			// },
		}

		// Set the unit if available
		if item.SpecifiedLineTradeDelivery.BilledQuantity.UnitCode != "" {
			line.Item.Unit = org.Unit(item.SpecifiedLineTradeDelivery.BilledQuantity.UnitCode)
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
