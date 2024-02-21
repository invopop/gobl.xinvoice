package xinvoice

import "github.com/invopop/gobl/bill"

// Transaction defines the structure of the transaction in the CII standard
type Transaction struct {
	Lines      []*Line     `xml:"ram:IncludedSupplyChainTradeLineItem"`
	Agreement  *Agreement  `xml:"ram:ApplicableHeaderTradeAgreement"`
	Delivery   *Delivery   `xml:"ram:ApplicableHeaderTradeDelivery"`
	Settlement *Settlement `xml:"ram:ApplicableHeaderTradeSettlement"`
}

// Agreement defines the structure of the ApplicableHeaderTradeAgreement of the CII standard
type Agreement struct {
	BuyerReference string  `xml:"ram:BuyerReference"`
	Seller         *Seller `xml:"ram:SellerTradeParty"`
	Buyer          *Buyer  `xml:"ram:BuyerTradeParty"`
}

// Seller defines the structure of the SellerTradeParty of the CII standard
type Seller struct {
	Name                      string                     `xml:"ram:Name"`
	LegalOrganization         *LegalOrganization         `xml:"ram:SpecifiedLegalOrganization"`
	Contact                   *Contact                   `xml:"ram:DefinedTradeContact"`
	PostalTradeAddress        *PostalTradeAddress        `xml:"ram:PostalTradeAddress"`
	URIUniversalCommunication *URIUniversalCommunication `xml:"ram:URIUniversalCommunication>ram:URIID"`
	SpecifiedTaxRegistration  *SpecifiedTaxRegistration  `xml:"ram:SpecifiedTaxRegistration>ram:ID"`
}

// Contact defines the structure of the DefinedTradeContact of the CII standard
type Contact struct {
	Name  string `xml:"ram:PersonName"`
	Phone string `xml:"ram:TelephoneUniversalCommunication>ram:CompleteNumber"`
	Email string `xml:"ram:EmailURIUniversalCommunication>ram:URIID"`
}

// SpecifiedTaxRegistration defines the structure of the SpecifiedTaxRegistration of the CII standard
type SpecifiedTaxRegistration struct {
	ID       string `xml:",chardata"`
	SchemeID string `xml:"schemeID,attr"`
}

// LegalOrganization defines the structure of the SpecifiedLegalOrganization of the CII standard
type LegalOrganization struct {
	ID   string `xml:"ram:ID"`
	Name string `xml:"ram:TradingBusinessName"`
}

// Buyer defines the structure of the BuyerTradeParty of the CII standard
type Buyer struct {
	ID                        string                     `xml:"ram:ID"`
	Name                      string                     `xml:"ram:Name"`
	PostalTradeAddress        *PostalTradeAddress        `xml:"ram:PostalTradeAddress"`
	URIUniversalCommunication *URIUniversalCommunication `xml:"ram:URIUniversalCommunication>ram:URIID"`
}

// PostalTradeAddress defines the structure of the PostalTradeAddress of the CII standard
type PostalTradeAddress struct {
	Postcode  string `xml:"ram:PostcodeCode"`
	LineOne   string `xml:"ram:LineOne"`
	City      string `xml:"ram:CityName"`
	CountryID string `xml:"ram:CountryID"`
}

// URIUniversalCommunication defines the structure of URIUniversalCommunication of the CII standard
type URIUniversalCommunication struct {
	URIID    string `xml:",chardata"`
	SchemeID string `xml:"schemeID,attr"`
}

// Delivery defines the structure of ApplicableHeaderTradeDelivery of the CII standard
type Delivery struct {
	Event *Date `xml:"ram:ActualDeliverySupplyChainEvent>ram:OccurrenceDateTime>udt:DateTimeString,omitempty"`
}

// Settlement defines the structure of ApplicableHeaderTradeSettlement of the CII standard
type Settlement struct {
	Currency              string   `xml:"ram:InvoiceCurrencyCode"`
	TypeCode              string   `xml:"ram:SpecifiedTradeSettlementPaymentMeans>ram:TypeCode"`
	PayeeFinancialAccount string   `xml:"ram:SpecifiedTradeSettlementPaymentMeans>ram:PayeePartyCreditorFinancialAccount>ram:IBANID"`
	Tax                   *Tax     `xml:"ram:ApplicableTradeTax"`
	Description           string   `xml:"ram:SpecifiedTradePaymentTerms>ram:Description"`
	Summary               *Summary `xml:"ram:SpecifiedTradeSettlementHeaderMonetarySummation"`
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

// NewTransaction creates the transaction part of a EN 16931 compliant invoice
func NewTransaction(inv *bill.Invoice) *Transaction {
	agreement := &Agreement{
		BuyerReference: "04011000-12345-03",
		Seller: &Seller{
			Name: "[Seller name]",
			LegalOrganization: &LegalOrganization{
				ID:   "[HRA-Eintrag]",
				Name: "[Seller trading name]",
			},
			Contact: &Contact{
				Name:  "nicht vorhanden",
				Phone: "+49 1234-5678",
				Email: "seller@email.de",
			},
			PostalTradeAddress: &PostalTradeAddress{
				Postcode:  "12345",
				LineOne:   "[Seller address line 1]",
				City:      "[Seller city]",
				CountryID: "DE",
			},
			URIUniversalCommunication: &URIUniversalCommunication{
				URIID:    "seller@email.de",
				SchemeID: "EM",
			},
			SpecifiedTaxRegistration: &SpecifiedTaxRegistration{
				ID:       "DE 123456789",
				SchemeID: "VA",
			},
		},
		Buyer: &Buyer{
			ID:   "[Buyer identifier]",
			Name: "[Buyer name]",
			PostalTradeAddress: &PostalTradeAddress{
				Postcode:  "12345",
				LineOne:   "[Buyer address line 1]",
				City:      "[Buyer city]",
				CountryID: "DE",
			},
			URIUniversalCommunication: &URIUniversalCommunication{
				URIID:    "buyer@info.de",
				SchemeID: "EM",
			},
		},
	}
	return &Transaction{
		Lines:     NewLines(inv.Lines),
		Agreement: agreement,
		Delivery: &Delivery{
			Event: &Date{
				Date:   "20160621",
				Format: "102",
			},
		},
		Settlement: &Settlement{
			Currency:              "EUR",
			TypeCode:              "58",
			PayeeFinancialAccount: "DE75512108001245126199",
			Description:           "Zahlbar sofort ohne Abzug.",
			Tax: &Tax{
				CalculatedAmount:      "22.04",
				TypeCode:              "VAT",
				BasisAmount:           "314.86",
				CategoryCode:          "S",
				RateApplicablePercent: "7",
			},
			Summary: &Summary{
				TotalAmount:         "314.86",
				TaxBasisTotalAmount: "314.86",
				GrandTotalAmount:    "336.9",
				DuePayableAmount:    "336.9",
				TaxTotalAmount: &TaxTotalAmount{
					Amount:   "22.04",
					Currency: "EUR",
				},
			},
		},
	}
}
