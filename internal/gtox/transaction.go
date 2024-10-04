package xinvoice

import (
	"github.com/invopop/gobl/bill"
)

// Transaction defines the structure of the transaction in the CII standard
type Transaction struct {
	Lines      []*Line     `xml:"ram:IncludedSupplyChainTradeLineItem"`
	Agreement  *Agreement  `xml:"ram:ApplicableHeaderTradeAgreement"`
	Delivery   *Delivery   `xml:"ram:ApplicableHeaderTradeDelivery"`
	Settlement *Settlement `xml:"ram:ApplicableHeaderTradeSettlement"`
}

// Delivery defines the structure of ApplicableHeaderTradeDelivery of the CII standard
type Delivery struct {
	Event *Date `xml:"ram:ActualDeliverySupplyChainEvent>ram:OccurrenceDateTime>udt:DateTimeString,omitempty"`
}

// NewTransaction creates the transaction part of a EN 16931 compliant invoice
func NewTransaction(inv *bill.Invoice) (*Transaction, error) {
	agreement, err := NewAgreement(inv)
	if err != nil {
		return nil, err
	}

	transaction := &Transaction{
		Lines:      NewLines(inv.Lines),
		Agreement:  agreement,
		Delivery:   &Delivery{},
		Settlement: NewSettlement(inv),
	}

	return transaction, nil
}
