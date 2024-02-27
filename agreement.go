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

	buyer, err := NewBuyer(customer)
	if err != nil {
		return nil, err
	}

	seller, err := NewSeller(supplier)
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

// NewPostalTradeAddress creates the PostalTradeAddress part of a EN 16931 compliant invoice
func NewPostalTradeAddress(addresses []*org.Address) (*PostalTradeAddress, error) {
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

// NewEmail creates the URIUniversalCommunication part of a EN 16931 compliant invoice
func NewEmail(emails []*org.Email) (*URIUniversalCommunication, error) {
	if len(emails) == 0 {
		return nil, fmt.Errorf("No emails found")
	}

	email := &URIUniversalCommunication{
		URIID:    emails[0].Address,
		SchemeID: "EM",
	}

	return email, nil
}
