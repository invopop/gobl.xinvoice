package xinvoice

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
)

// SchemeIDEmail represents the Scheme ID for email addresses
const SchemeIDEmail = "EM"

// Agreement defines the structure of the ApplicableHeaderTradeAgreement of the CII standard
type Agreement struct {
	BuyerReference string  `xml:"ram:BuyerReference,omitempty"`
	Seller         *Seller `xml:"ram:SellerTradeParty"`
	Buyer          *Buyer  `xml:"ram:BuyerTradeParty"`
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

// NewAgreement creates the ApplicableHeaderTradeAgreement part of a EN 16931 compliant invoice
func NewAgreement(inv *bill.Invoice) (*Agreement, error) {
	agreement := new(Agreement)

	if inv.Ordering != nil {
		agreement.BuyerReference = inv.Ordering.Code
	}

	if supplier := inv.Supplier; supplier != nil {
		agreement.Seller = NewSeller(supplier)
	}

	if customer := inv.Customer; customer != nil {
		agreement.Buyer = NewBuyer(customer)
	}

	return agreement, nil
}

// NewPostalTradeAddress creates the PostalTradeAddress part of a EN 16931 compliant invoice
func NewPostalTradeAddress(addresses []*org.Address) *PostalTradeAddress {
	if len(addresses) == 0 {
		return nil
	}
	address := addresses[0]

	postalTradeAddress := &PostalTradeAddress{
		Postcode:  address.Code,
		LineOne:   address.Street,
		City:      address.Locality,
		CountryID: string(address.Country),
	}

	return postalTradeAddress
}

// NewEmail creates the URIUniversalCommunication part of a EN 16931 compliant invoice
func NewEmail(emails []*org.Email) *URIUniversalCommunication {
	if len(emails) == 0 {
		return nil
	}

	email := &URIUniversalCommunication{
		URIID:    emails[0].Address,
		SchemeID: SchemeIDEmail,
	}

	return email
}
