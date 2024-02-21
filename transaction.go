package xinvoice

import "github.com/invopop/gobl/bill"

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

// Settlement defines the structure of ApplicableHeaderTradeSettlement of the CII standard
type Settlement struct {
	Currency              string   `xml:"ram:InvoiceCurrencyCode"`
	TypeCode              string   `xml:"ram:SpecifiedTradeSettlementPaymentMeans>ram:TypeCode"`
	PayeeFinancialAccount string   `xml:"ram:SpecifiedTradeSettlementPaymentMeans>ram:PayeePartyCreditorFinancialAccount>ram:IBANID"`
	Tax                   *Tax     `xml:"ram:ApplicableTradeTax"`
	Description           string   `xml:"ram:SpecifiedTradePaymentTerms>ram:Description"`
	Summary               *Summary `xml:"ram:SpecifiedTradeSettlementHeaderMonetarySummation"`
}

// Tax defines the structure of ApplicableTradeTax of the CII standard
type Tax struct {
	CalculatedAmount      string `xml:"ram:CalculatedAmount"`
	TypeCode              string `xml:"ram:TypeCode"`
	BasisAmount           string `xml:"ram:BasisAmount"`
	CategoryCode          string `xml:"ram:CategoryCode"`
	RateApplicablePercent string `xml:"ram:RateApplicablePercent"`
}

// Summary defines the structure of SpecifiedTradeSettlementHeaderMonetarySummation of the CII standard
type Summary struct {
	TotalAmount         string          `xml:"ram:LineTotalAmount"`
	TaxBasisTotalAmount string          `xml:"ram:TaxBasisTotalAmount"`
	TaxTotalAmount      *TaxTotalAmount `xml:"ram:TaxTotalAmount"`
	GrandTotalAmount    string          `xml:"ram:GrandTotalAmount"`
	DuePayableAmount    string          `xml:"ram:DuePayableAmount"`
}

// TaxTotalAmount defines the structure of the TaxTotalAmount of the CII standard
type TaxTotalAmount struct {
	Amount   string `xml:",chardata"`
	Currency string `xml:"currencyID,attr"`
}

// NewTransaction creates the transaction part of a EN 16931 compliant invoice
func NewTransaction(inv *bill.Invoice) *Transaction {
	return &Transaction{
		Lines:     NewLines(inv.Lines),
		Agreement: NewAgreement(inv),
		Delivery: &Delivery{
			Event: &Date{
				Date:   "20160621",
				Format: "102",
			},
		},
		Settlement: &Settlement{
			Currency:              "EUR",
			TypeCode:              "58",
			PayeeFinancialAccount: "DE75512108001245126199",
			Description:           "Zahlbar sofort ohne Abzug.",
			Tax: &Tax{
				CalculatedAmount:      "22.04",
				TypeCode:              "VAT",
				BasisAmount:           "314.86",
				CategoryCode:          "S",
				RateApplicablePercent: "7",
			},
			Summary: &Summary{
				TotalAmount:         "314.86",
				TaxBasisTotalAmount: "314.86",
				GrandTotalAmount:    "336.9",
				DuePayableAmount:    "336.9",
				TaxTotalAmount: &TaxTotalAmount{
					Amount:   "22.04",
					Currency: "EUR",
				},
			},
		},
	}
}
