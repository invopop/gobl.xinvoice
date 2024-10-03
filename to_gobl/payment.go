package to_gobl

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
)

// Parses the XML information for a Payment object
func ParsePayment(settlement *ApplicableHeaderTradeSettlement) *bill.Payment {
	payment := &bill.Payment{}

	if settlement.PayeeTradeParty != nil {
		payee := &org.Party{Name: settlement.PayeeTradeParty.Name}
		if settlement.PayeeTradeParty.PostalTradeAddress.LineOne != "" {
			payee.Addresses = []*org.Address{
				{
					Street:   settlement.PayeeTradeParty.PostalTradeAddress.LineOne,
					Locality: settlement.PayeeTradeParty.PostalTradeAddress.CityName,
					Code:     settlement.PayeeTradeParty.PostalTradeAddress.PostcodeCode,
					Country:  settlement.PayeeTradeParty.PostalTradeAddress.CountryID,
				},
			}
		}
		payment.Payee = payee
	}
	if len(settlement.SpecifiedTradePaymentTerms) > 0 {
		if settlement.SpecifiedTradePaymentTerms[0].DueDateDateTime != nil {
			payment.Terms = parsePaymentTerms(settlement)
		}
	}

	if len(settlement.SpecifiedTradeSettlementPaymentMeans) > 0 && settlement.SpecifiedTradeSettlementPaymentMeans[0].TypeCode != "1" {
		payment.Instructions = parsePaymentMeans(settlement)
	}

	if settlement.SpecifiedAdvancePayment.FormattedReceivedDateTime.DateTimeString != "" {
		advancePaymentReceivedDateTime := ParseDate(settlement.SpecifiedAdvancePayment.FormattedReceivedDateTime.DateTimeString)
		advance := &pay.Advance{
			Amount: num.AmountFromFloat64(settlement.SpecifiedAdvancePayment.PaidAmount, 0),
			Date:   &advancePaymentReceivedDateTime,
		}
		payment.Advances = []*pay.Advance{advance}
	}

	return payment
}

func parsePaymentTerms(settlement *ApplicableHeaderTradeSettlement) *pay.Terms {
	terms := &pay.Terms{}
	var dueDates []*pay.DueDate

	for _, paymentTerm := range settlement.SpecifiedTradePaymentTerms {
		if terms.Detail == "" {
			terms.Detail = paymentTerm.Description
		}

		if paymentTerm.DueDateDateTime.DateTimeString != "" {
			paymentTermsDueDateDateTime := ParseDate(paymentTerm.DueDateDateTime.DateTimeString)
			dueDate := &pay.DueDate{
				Date: &paymentTermsDueDateDateTime,
			}
			if paymentTerm.PartialPaymentAmount != nil {
				dueDate.Amount, _ = num.AmountFromString(*paymentTerm.PartialPaymentAmount)
			}
			dueDates = append(dueDates, dueDate)
		}
	}
	terms.DueDates = dueDates
	return terms
}

func parsePaymentMeans(settlement *ApplicableHeaderTradeSettlement) *pay.Instructions {
	paymentMeans := settlement.SpecifiedTradeSettlementPaymentMeans[0]
	instructions := &pay.Instructions{
		Key: PaymentMeansTypeCodeParse(paymentMeans.TypeCode),
	}

	if paymentMeans.Information != nil {
		instructions.Detail = *paymentMeans.Information
	}

	if paymentMeans.ApplicableTradeSettlementFinancialCard != nil {
		if paymentMeans.ApplicableTradeSettlementFinancialCard != nil {
			card := paymentMeans.ApplicableTradeSettlementFinancialCard
			instructions.Card = &pay.Card{
				Last4: card.ID[len(card.ID)-4:],
			}
			if card.CardholderName != "" {
				instructions.Card.Holder = card.CardholderName
			}
		}
	}

	if paymentMeans.PayeePartyCreditorFinancialAccount != nil {
		account := paymentMeans.PayeePartyCreditorFinancialAccount
		if account.IBANID != "" {
			instructions.CreditTransfer = []*pay.CreditTransfer{
				{
					IBAN: account.IBANID,
				},
			}
		}
		if account.AccountName != "" {
			//No issue because X-Rechnung only supports one credit transfer per instruction
			instructions.CreditTransfer[0].Name = account.AccountName
		}
	}
	return instructions
}
