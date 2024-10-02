package to_gobl

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/de"
	"github.com/invopop/gobl/tax"
)

// Parses the XML information for a Party object
func ParseParty(party *TradeParty) *org.Party {
	p := &org.Party{
		Name: party.Name,
	}

	if party.DefinedTradeContact != nil && party.DefinedTradeContact.PersonName != nil {
		p.People = []*org.Person{
			{
				Name: &org.Name{
					Given: *party.DefinedTradeContact.PersonName,
				},
			},
		}
	}

	if party.PostalTradeAddress != nil {
		p.Addresses = []*org.Address{
			{
				Street:   party.PostalTradeAddress.LineOne,
				Locality: party.PostalTradeAddress.CityName,
				Code:     party.PostalTradeAddress.PostcodeCode,
				Country:  party.PostalTradeAddress.CountryID,
			},
		}
	}

	if party.DefinedTradeContact != nil && party.DefinedTradeContact.TelephoneUniversalCommunication != nil {
		p.Telephones = []*org.Telephone{
			{
				Number: party.DefinedTradeContact.TelephoneUniversalCommunication.CompleteNumber,
			},
		}
	}

	if party.DefinedTradeContact != nil && party.DefinedTradeContact.EmailURIUniversalCommunication != nil {
		p.Emails = []*org.Email{
			{
				Address: party.DefinedTradeContact.EmailURIUniversalCommunication.URIID,
			},
		}
	}
	if party.SpecifiedTaxRegistration != nil {
		for _, taxReg := range *party.SpecifiedTaxRegistration {
			if taxReg.ID != nil {
				switch taxReg.ID.SchemeID {
				case "VA":
					p.TaxID = &tax.Identity{
						Country: l10n.TaxCountryCode(party.PostalTradeAddress.CountryID),
						Code:    cbc.Code(taxReg.ID.Value),
					}
				case "FC":
					identity := &org.Identity{
						Key:  de.IdentityKeyTaxNumber,
						Code: cbc.Code(taxReg.ID.Value),
					}
					p.Identities = append(p.Identities, identity)
				}
			}
		}
	}

	return p
}

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
		payment.Terms = terms
	}

	if settlement.SpecifiedAdvancePayment.FormattedReceivedDateTime.DateTimeString != "" {
		advancePaymentReceivedDateTime := ParseDate(settlement.SpecifiedAdvancePayment.FormattedReceivedDateTime.DateTimeString)
		advance := &pay.Advance{
			Amount: num.AmountFromFloat64(settlement.SpecifiedAdvancePayment.PaidAmount, 0),
			Date:   &advancePaymentReceivedDateTime,
		}
		payment.Advances = []*pay.Advance{advance}
	}
	if len(settlement.SpecifiedTradeSettlementPaymentMeans) > 0 {
		paymentMeans := settlement.SpecifiedTradeSettlementPaymentMeans[0]
		instructions := &pay.Instructions{}

		// 10: Bargeld (Cash)
		// 20: Scheck (Check)
		// 30: Überweisung (Bank Transfer)
		// 42: Payment to bank account
		// 48: Kartenzahlung (Card Payment)
		// 49: Lastschrift (Direct Debit)
		// 57: Dauerauftrag (Standing Order)
		// 58: SEPA Credit Transfer
		// 59: SEPA Direct Debit
		// 97: Report (Reporting, potentially for financial reporting or statement purposes)

		if paymentMeans.TypeCode != "" {
			instructions.Key = cbc.Key(paymentMeans.TypeCode)
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
		}

		// if paymentMeans.PayerPartyDebtorFinancialAccount != nil {
		// 	account := paymentMeans.PayerPartyDebtorFinancialAccount
		// 	if account.IBANID != "" {
		// 		instructions.DebitTransfer = []*pay.DebitTransfer{
		// 			{
		// 				IBAN: account.IBANID,
		// 			},
		// 		}
		// 	}
		// }

		if instructions.Key != "" || instructions.CreditTransfer != nil {
			payment.Instructions = instructions
		}
	}

	return payment
}

// Parses the XML information for a Lines object
func ParseLines(transaction *SupplyChainTradeTransaction) []*bill.Line {
	items := transaction.IncludedSupplyChainTradeLineItem
	lines := make([]*bill.Line, 0, len(transaction.IncludedSupplyChainTradeLineItem))

	for _, item := range items {
		price := num.AmountFromFloat64(item.SpecifiedLineTradeAgreement.NetPriceProductTradePrice.ChargeAmount, 0)

		line := &bill.Line{
			// Index:    item.AssociatedDocumentLineDocument.LineID, //generated field
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  item.SpecifiedTradeProduct.Name,
				Price: price,
			},
			Taxes: tax.Set{
				{
					Rate:     FindTaxKey(item.SpecifiedLineTradeSettlement.ApplicableTradeTax.CategoryCode),
					Category: cbc.Code(item.SpecifiedLineTradeSettlement.ApplicableTradeTax.TypeCode),
				},
			},
		}

		if item.SpecifiedLineTradeDelivery != nil {
			line.Quantity = num.MakeAmount(item.SpecifiedLineTradeDelivery.BilledQuantity.Value, 0)
		}

		if len(item.AssociatedDocumentLineDocument.IncludedNote) > 0 {
			line.Notes = make([]*cbc.Note, 0, len(item.AssociatedDocumentLineDocument.IncludedNote))
			for _, note := range item.AssociatedDocumentLineDocument.IncludedNote {
				n := &cbc.Note{}
				if note.Content != "" {
					n.Text = note.Content
				}
				if note.ContentCode != "" {
					n.Code = note.ContentCode
				}
				line.Notes = append(line.Notes, n)
			}
		}

		if item.SpecifiedLineTradeSettlement.ApplicableTradeTax.RateApplicablePercent != "" {
			percent, _ := num.PercentageFromString(item.SpecifiedLineTradeSettlement.ApplicableTradeTax.RateApplicablePercent)
			line.Taxes[0].Percent = &percent
		}

		if item.SpecifiedLineTradeDelivery.BilledQuantity.UnitCode != "" {
			line.Item.Unit = UnitFromUNECE(cbc.Code(item.SpecifiedLineTradeDelivery.BilledQuantity.UnitCode))
		}

		lines = append(lines, line)
	}

	return lines
}
