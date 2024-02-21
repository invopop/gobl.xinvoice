package xinvoice

import "github.com/invopop/gobl/bill"

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

// PostalTradeAddress defines the structure of the PostalTradeAddress of the CII standard
type PostalTradeAddress struct {
	Postcode  string `xml:"ram:PostcodeCode"`
	LineOne   string `xml:"ram:LineOne"`
	City      string `xml:"ram:CityName"`
	CountryID string `xml:"ram:CountryID"`
}

// Buyer defines the structure of the BuyerTradeParty of the CII standard
type Buyer struct {
	ID                        string                     `xml:"ram:ID"`
	Name                      string                     `xml:"ram:Name"`
	PostalTradeAddress        *PostalTradeAddress        `xml:"ram:PostalTradeAddress"`
	URIUniversalCommunication *URIUniversalCommunication `xml:"ram:URIUniversalCommunication>ram:URIID"`
}

// URIUniversalCommunication defines the structure of URIUniversalCommunication of the CII standard
type URIUniversalCommunication struct {
	URIID    string `xml:",chardata"`
	SchemeID string `xml:"schemeID,attr"`
}

func NewAgreement(inv *bill.Invoice) *Agreement {
	return &Agreement{
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
}
