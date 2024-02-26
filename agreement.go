package xinvoice

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
)

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

// NewAgreement creates the ApplicableHeaderTradeAgreement part of a EN 16931 compliant invoice
func NewAgreement(inv *bill.Invoice) (*Agreement, error) {
	if inv.Customer == nil {
		return nil, fmt.Errorf("Customer not found")
	}
	customer := inv.Customer

	if inv.Supplier == nil {
		return nil, fmt.Errorf("Supplier not found")
	}
	supplier := inv.Supplier

	if customer.TaxID == nil {
		return nil, fmt.Errorf("Customer TaxID not found")
	}
	ref := customer.TaxID.String()

	buyer, err := newBuyer(customer)
	if err != nil {
		return nil, err
	}

	seller, err := newSeller(supplier)
	if err != nil {
		return nil, err
	}

	agreement := &Agreement{
		BuyerReference: ref,
		Seller:         seller,
		Buyer:          buyer,
	}

	return agreement, nil
}

func newBuyer(customer *org.Party) (*Buyer, error) {
	if customer.TaxID == nil {
		return nil, fmt.Errorf("Customer TaxID not found")
	}
	ref := customer.TaxID.String()

	address, err := newPostalTradeAddress(customer.Addresses)
	if err != nil {
		return nil, err
	}

	email, err := newEmail(customer.Emails)
	if err != nil {
		return nil, err
	}

	buyer := &Buyer{
		ID:                        ref,
		Name:                      customer.Name,
		PostalTradeAddress:        address,
		URIUniversalCommunication: email,
	}

	return buyer, nil
}

func newSeller(supplier *org.Party) (*Seller, error) {
	if supplier.TaxID == nil {
		return nil, fmt.Errorf("Supplier TaxID not found")
	}
	taxID := supplier.TaxID.String()

	contact, err := newContact(supplier)
	if err != nil {
		return nil, err
	}

	address, err := newPostalTradeAddress(supplier.Addresses)
	if err != nil {
		return nil, err
	}

	email, err := newEmail(supplier.Emails)
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
	name := supplier.People[0].Name.Given

	if len(supplier.Telephones) == 0 {
		return nil, fmt.Errorf("Supplier Telephones not found")
	}
	phone := supplier.Telephones[0].Number

	if len(supplier.Emails) == 0 {
		return nil, fmt.Errorf("Supplier Emails not found")
	}
	email := supplier.Emails[0].Address

	contact := &Contact{
		Name:  name,
		Phone: phone,
		Email: email,
	}

	return contact, nil
}

func newPostalTradeAddress(addresses []*org.Address) (*PostalTradeAddress, error) {
	if len(addresses) == 0 {
		return nil, fmt.Errorf("No addresses found")
	}
	address := addresses[0]

	postalTradeAddress := &PostalTradeAddress{
		Postcode:  address.Code,
		LineOne:   address.Street,
		City:      address.Locality,
		CountryID: string(address.Country),
	}

	return postalTradeAddress, nil
}

func newEmail(emails []*org.Email) (*URIUniversalCommunication, error) {
	if len(emails) == 0 {
		return nil, fmt.Errorf("No emails found")
	}

	email := &URIUniversalCommunication{
		URIID:    emails[0].Address,
		SchemeID: "EM",
	}

	return email, nil
}
