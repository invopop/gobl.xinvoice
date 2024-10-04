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
func NewSeller(supplier *org.Party) *Seller {
	if supplier == nil {
		return nil
	}
	seller := &Seller{
		Name:                      supplier.Name,
		Contact:                   newContact(supplier),
		PostalTradeAddress:        NewPostalTradeAddress(supplier.Addresses),
		URIUniversalCommunication: NewEmail(supplier.Emails),
	}

	if supplier.TaxID != nil {
		seller.SpecifiedTaxRegistration = &SpecifiedTaxRegistration{
			ID:       supplier.TaxID.String(),
			SchemeID: "VA",
		}
	}

	return seller
}

func newContact(supplier *org.Party) *Contact {
	if len(supplier.People) == 0 || len(supplier.Telephones) == 0 || len(supplier.Emails) == 0 {
		return nil
	}

	contact := &Contact{
		PersonName: contactName(supplier.People[0].Name),
		Phone:      supplier.Telephones[0].Number,
		Email:      supplier.Emails[0].Address,
	}

	return contact
}

func contactName(personName *org.Name) string {
	given := personName.Given
	surname := personName.Surname

	if given == "" && surname == "" {
		return ""
	}

	if given == "" {
		return surname
	}

	if surname == "" {
		return given
	}

	return fmt.Sprintf("%s %s", given, surname)
}
