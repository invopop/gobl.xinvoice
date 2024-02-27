package xinvoice

import (
	"fmt"

	"github.com/invopop/gobl/org"
)

// Buyer defines the structure of the BuyerTradeParty of the CII standard
type Buyer struct {
	ID                        string                     `xml:"ram:ID"`
	Name                      string                     `xml:"ram:Name"`
	PostalTradeAddress        *PostalTradeAddress        `xml:"ram:PostalTradeAddress"`
	URIUniversalCommunication *URIUniversalCommunication `xml:"ram:URIUniversalCommunication>ram:URIID"`
}

// NewBuyer creates the BuyerTradeParty part of a EN 16931 compliant invoice
func NewBuyer(customer *org.Party) (*Buyer, error) {
	if customer.TaxID == nil {
		return nil, fmt.Errorf("Customer TaxID not found")
	}
	ref := customer.TaxID.String()

	address, err := NewPostalTradeAddress(customer.Addresses)
	if err != nil {
		return nil, err
	}

	email, err := NewEmail(customer.Emails)
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
