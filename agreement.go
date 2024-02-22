package xinvoice

import "github.com/invopop/gobl/bill"

// Agreement defines the structure of the ApplicableHeaderTradeAgreement of the CII standard
type Agreement struct {
	BuyerReference string  `xml:"ram:BuyerReference,omitempty"`
	Seller         *Seller `xml:"ram:SellerTradeParty"`
	Buyer          *Buyer  `xml:"ram:BuyerTradeParty"`
}

// Seller defines the structure of the SellerTradeParty of the CII standard
type Seller struct {
	Name                      string                     `xml:"ram:Name"`
	LegalOrganization         *LegalOrganization         `xml:"ram:SpecifiedLegalOrganization,omitempty"`
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
		BuyerReference: inv.Customer.TaxID.String(),
		Seller: &Seller{
			Name: inv.Supplier.Name,
			Contact: &Contact{
				Name:  inv.Supplier.People[0].Name.Given,
				Phone: inv.Supplier.Telephones[0].Number,
				Email: inv.Supplier.Emails[0].Address,
			},
			PostalTradeAddress: &PostalTradeAddress{
				Postcode:  inv.Supplier.Addresses[0].Code,
				LineOne:   inv.Supplier.Addresses[0].Street,
				City:      inv.Supplier.Addresses[0].Locality,
				CountryID: string(inv.Supplier.Addresses[0].Country),
			},
			URIUniversalCommunication: &URIUniversalCommunication{
				URIID:    inv.Supplier.Emails[0].Address,
				SchemeID: "EM",
			},
			SpecifiedTaxRegistration: &SpecifiedTaxRegistration{
				ID:       inv.Supplier.TaxID.String(),
				SchemeID: "VA",
			},
		},
		Buyer: &Buyer{
			ID:   inv.Customer.TaxID.String(),
			Name: inv.Customer.Name,
			PostalTradeAddress: &PostalTradeAddress{
				Postcode:  inv.Customer.Addresses[0].Code,
				LineOne:   inv.Customer.Addresses[0].Street,
				City:      inv.Customer.Addresses[0].Locality,
				CountryID: string(inv.Supplier.Addresses[0].Country),
			},
			URIUniversalCommunication: &URIUniversalCommunication{
				URIID:    inv.Customer.Emails[0].Address,
				SchemeID: "EM",
			},
		},
	}
}
