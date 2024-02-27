package xinvoice

import (
	"fmt"

	"github.com/invopop/gobl/org"
)

// Seller defines the structure of the SellerTradeParty of the CII standard
type Seller struct {
	Name                      string                     `xml:"ram:Name"`
	LegalOrganization         *LegalOrganization         `xml:"ram:SpecifiedLegalOrganization,omitempty"`
	Contact                   *Contact                   `xml:"ram:DefinedTradeContact"`
	PostalTradeAddress        *PostalTradeAddress        `xml:"ram:PostalTradeAddress"`
	URIUniversalCommunication *URIUniversalCommunication `xml:"ram:URIUniversalCommunication>ram:URIID"`
	SpecifiedTaxRegistration  *SpecifiedTaxRegistration  `xml:"ram:SpecifiedTaxRegistration>ram:ID"`
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

// Contact defines the structure of the DefinedTradeContact of the CII standard
type Contact struct {
	PersonName string `xml:"ram:PersonName"`
	Phone      string `xml:"ram:TelephoneUniversalCommunication>ram:CompleteNumber"`
	Email      string `xml:"ram:EmailURIUniversalCommunication>ram:URIID"`
}

// NewSeller creates the SellerTradeParty part of a EN 16931 compliant invoice
func NewSeller(supplier *org.Party) (*Seller, error) {
	if supplier.TaxID == nil {
		return nil, fmt.Errorf("Supplier TaxID not found")
	}
	taxID := supplier.TaxID.String()

	contact, err := newContact(supplier)
	if err != nil {
		return nil, err
	}

	address, err := NewPostalTradeAddress(supplier.Addresses)
	if err != nil {
		return nil, err
	}

	email, err := NewEmail(supplier.Emails)
	if err != nil {
		return nil, err
	}

	seller := &Seller{
		Name:                      supplier.Name,
		Contact:                   contact,
		PostalTradeAddress:        address,
		URIUniversalCommunication: email,
		SpecifiedTaxRegistration: &SpecifiedTaxRegistration{
			ID:       taxID,
			SchemeID: "VA",
		},
	}

	return seller, nil
}

func newContact(supplier *org.Party) (*Contact, error) {
	if len(supplier.People) == 0 {
		return nil, fmt.Errorf("Supplier People not found")
	}
	name := contactName(&supplier.People[0].Name)

	if len(supplier.Telephones) == 0 {
		return nil, fmt.Errorf("Supplier Telephones not found")
	}
	phone := supplier.Telephones[0].Number

	if len(supplier.Emails) == 0 {
		return nil, fmt.Errorf("Supplier Emails not found")
	}
	email := supplier.Emails[0].Address

	contact := &Contact{
		PersonName: name,
		Phone:      phone,
		Email:      email,
	}

	return contact, nil
}

func contactName(personName *org.Name) string {
	return fmt.Sprintf("%s %s", personName.Given, personName.Surname)
}
