package xinvoice

import (
	"github.com/invopop/gobl/org"
)

// Buyer defines the structure of the BuyerTradeParty of the CII standard
type Buyer struct {
	ID                        string                     `xml:"ram:ID,omitempty"`
	Name                      string                     `xml:"ram:Name"`
	PostalTradeAddress        *PostalTradeAddress        `xml:"ram:PostalTradeAddress"`
	URIUniversalCommunication *URIUniversalCommunication `xml:"ram:URIUniversalCommunication>ram:URIID"`
}

// NewBuyer creates the BuyerTradeParty part of a EN 16931 compliant invoice
func NewBuyer(customer *org.Party) *Buyer {
	buyer := &Buyer{
		Name:                      customer.Name,
		PostalTradeAddress:        NewPostalTradeAddress(customer.Addresses),
		URIUniversalCommunication: NewEmail(customer.Emails),
	}

	if customer.TaxID != nil {
		buyer.ID = customer.TaxID.String()
	}

	return buyer
}
