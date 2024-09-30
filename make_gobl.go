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
	"github.com/invopop/gobl/org"
)

// Document is a pseudo-model for containing the XML document being created
type XMLDoc struct {
	XMLName                     xml.Name                    `xml:"CrossIndustryInvoice"`
	BusinessProcessContext      string                      `xml:"ExchangedDocumentContext>BusinessProcessSpecifiedDocumentContextParameter>ID"`
	GuidelineContext            string                      `xml:"ExchangedDocumentContext>GuidelineSpecifiedDocumentContextParameter>ID"`
	ExchangedDocument           ExchangedDocument           `xml:"ExchangedDocument"`
	SupplyChainTradeTransaction SupplyChainTradeTransaction `xml:"SupplyChainTradeTransaction"`
}

type ExchangedDocument struct {
	ID            string `xml:"ID"`
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
		LineID string `xml:"LineID"`
	} `xml:"AssociatedDocumentLineDocument"`
	SpecifiedTradeProduct struct {
		Name string `xml:"Name"`
	} `xml:"SpecifiedTradeProduct"`
	SpecifiedLineTradeAgreement struct {
		NetPriceProductTradePrice struct {
			ChargeAmount string `xml:"ChargeAmount"`
		} `xml:"NetPriceProductTradePrice"`
	} `xml:"SpecifiedLineTradeAgreement"`
	SpecifiedLineTradeDelivery struct {
		BilledQuantity struct {
			Value    string `xml:",chardata"`
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
			LineTotalAmount string `xml:"LineTotalAmount"`
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
			Value    string `xml:",chardata"`
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
	SpecifiedTradePaymentTerms struct {
		Description string `xml:"Description"`
	} `xml:"SpecifiedTradePaymentTerms"`
	SpecifiedTradeSettlementHeaderMonetarySummation struct {
		LineTotalAmount     string `xml:"LineTotalAmount"`
		TaxBasisTotalAmount string `xml:"TaxBasisTotalAmount"`
		TaxTotalAmount      struct {
			Value      string `xml:",chardata"`
			CurrencyID string `xml:"currencyID,attr"`
		} `xml:"TaxTotalAmount"`
		GrandTotalAmount string `xml:"GrandTotalAmount"`
		DuePayableAmount string `xml:"DuePayableAmount"`
	} `xml:"SpecifiedTradeSettlementHeaderMonetarySummation"`
}

// NewDocument converts a XRechnung document into a GOBL envelope
func NewDocumentGOBL(doc *XMLDoc) (*gobl.Envelope, error) {

	issueDate, err := parseDate(doc.ExchangedDocument.IssueDateTime.DateTimeString.Value)
	if err != nil {
		return nil, err
	}

	inv := &bill.Invoice{
		Code:      cbc.Code(doc.ExchangedDocument.ID),
		Type:      "standard",
		IssueDate: issueDate,
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
			// TaxID: extractTaxID(doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.SellerTradeParty.SpecifiedTaxRegistration),
		},
		// Customer: &org.Party{
		// 	Name: doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.BuyerTradeParty.Name,
		// 	Addresses: []*org.Address{
		// 		{
		// 			Street:  doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.BuyerTradeParty.PostalTradeAddress.LineOne,
		// 			City:    doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.BuyerTradeParty.PostalTradeAddress.CityName,
		// 			Code:    doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.BuyerTradeParty.PostalTradeAddress.PostcodeCode,
		// 			Country: doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.BuyerTradeParty.PostalTradeAddress.CountryID,
		// 		},
		// 	},
		// },
		// Lines:   parseLines(doc.SupplyChainTradeTransaction.IncludedSupplyChainTradeLineItem),
		// Payment: parsePayment(doc.SupplyChainTradeTransaction.ApplicableHeaderTradeSettlement),
		// Totals:  parseTotals(doc.SupplyChainTradeTransaction.ApplicableHeaderTradeSettlement),
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

// Parse XML to JSON
// func parseDocument(doc *Document) ()

func parseDate(date string) (cal.Date, error) {
	t, err := time.Parse("20060102", date)
	if err != nil {
		return cal.Date{}, err
	}
	return cal.MakeDate(t.Day(), t.Month(), t.Year()), nil
}
