package xinvoice

import (
	"github.com/invopop/gobl/l10n"
)

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
		PersonName                      *string `xml:"PersonName"`
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
	SpecifiedTaxRegistration *[]struct {
		ID *struct {
			Value string `xml:",chardata"`
			//VA used for VAT-ID used in B2B, FC for tax number (Steuernummer).
			SchemeID string `xml:"schemeID,attr"`
		} `xml:"ID"`
	} `xml:"SpecifiedTaxRegistration,omitempty"`
}

type ApplicableHeaderTradeSettlement struct {
	InvoiceCurrencyCode                  string `xml:"InvoiceCurrencyCode"`
	SpecifiedTradeSettlementPaymentMeans []struct {
		TypeCode                               string  `xml:"TypeCode"`
		Information                            *string `xml:"Information"`
		ApplicableTradeSettlementFinancialCard *struct {
			ID             string `xml:"ID"`
			CardholderName string `xml:"CardholderName"`
		} `xml:"ApplicableTradeSettlementFinancialCard"`
		PayeePartyCreditorFinancialAccount *struct {
			IBANID      string `xml:"IBANID"`
			AccountName string `xml:"AccountName"`
		} `xml:"PayeePartyCreditorFinancialAccount"`
		PayerPartyDebtorFinancialAccount *struct {
			IBANID string `xml:"IBANID"`
		} `xml:"PayerPartyDebtorFinancialAccount"`
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
		DueDateDateTime *struct {
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
