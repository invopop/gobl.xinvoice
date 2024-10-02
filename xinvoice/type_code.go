package xinvoice

import (
	"strconv"

	"github.com/invopop/gobl/bill"
)

// TypeCodeInstrumentNotDefined is the Type Code for an undefined Payment Instrument
const TypeCodeInstrumentNotDefined = "1"

// Acceptable codes:
// 1   - Instrument not defined
// 2   - Automated clearing house credit
// 3   - Automated clearing house debit
// 4   - ACH demand debit reversal
// 5   - ACH demand credit reversal
// 6   - ACH demand credit
// 7   - ACH demand debit
// 8   - Hold
// 9   - National or regional clearing
// 10  - In cash
// 11  - ACH savings credit reversal
// 12  - ACH savings debit reversal
// 13  - ACH savings credit
// 14  - ACH savings debit
// 15  - Bookentry credit
// 16  - Bookentry debit
// 17  - ACH demand cash concentration/disbursement (CCD) credit
// 18  - ACH demand cash concentration/disbursement (CCD) debit
// 19  - ACH demand corporate trade payment (CTP) credit
// 20  - Cheque
// 21  - Banker's draft
// 22  - Certified banker's draft
// 23  - Bank cheque (issued by a banking or similar establishment)
// 24  - Bill of exchange awaiting acceptance
// 25  - Certified cheque
// 26  - Local cheque
// 27  - ACH demand corporate trade payment (CTP) debit
// 28  - ACH demand corporate trade exchange (CTX) credit
// 29  - ACH demand corporate trade exchange (CTX) debit
// 30  - Credit transfer
// 31  - Debit transfer
// 32  - ACH demand cash concentration/disbursement plus (CCD+) credit
// 33  - ACH demand cash concentration/disbursement plus (CCD+) debit
// 34  - ACH prearranged payment and deposit (PPD)
// 35  - ACH savings cash concentration/disbursement (CCD) credit
// 36  - ACH savings cash concentration/disbursement (CCD) debit
// 37  - ACH savings corporate trade payment (CTP) credit
// 38  - ACH savings corporate trade payment (CTP) debit
// 39  - ACH savings corporate trade exchange (CTX) credit
// 40  - ACH savings corporate trade exchange (CTX) debit
// 41  - ACH savings cash concentration/disbursement plus (CCD+) credit
// 42  - Payment to bank account
// 43  - ACH savings cash concentration/disbursement plus (CCD+) debit
// 44  - Accepted bill of exchange
// 45  - Referenced home-banking credit transfer
// 46  - Interbank debit transfer
// 47  - Home-banking debit transfer
// 48  - Bank card
// 49  - Direct debit
// 50  - Payment by postgiro
// 51  - FR, norme 6 97-Telereglement CFONB (French Organisation for Banking Standards)  - Option A
// 52  - Urgent commercial payment
// 53  - Urgent Treasury Payment
// 54  - Credit card
// 55  - Debit card
// 56  - Bankgiro
// 57  - Standing agreement
// 58  - SEPA credit transfer
// 59  - SEPA direct debit
// 60  - Promissory note
// 61  - Promissory note signed by the debtor
// 62  - Promissory note signed by the debtor and endorsed by a bank
// 63  - Promissory note signed by the debtor and endorsed by a third party
// 64  - Promissory note signed by a bank
// 65  - Promissory note signed by a bank and endorsed by another bank
// 66  - Promissory note signed by a third party
// 67  - Promissory note signed by a third party and endorsed by a bank
// 68  - Online payment service
// 69  - Transfer Advice
// 70  - Bill drawn by the creditor on the debtor
// 74  - Bill drawn by the creditor on a bank
// 75  - Bill drawn by the creditor, endorsed by another bank
// 76  - Bill drawn by the creditor on a bank and endorsed by a third party
// 77  - Bill drawn by the creditor on a third party
// 78  - Bill drawn by creditor on third party, accepted and endorsed by bank
// 91  - Not transferable banker's draft
// 92  - Not transferable local cheque
// 93  - Reference giro
// 94  - Urgent giro
// 95  - Free format giro
// 96  - Requested method for payment was not used
// 97  - Clearing between partners
// ZZZ - Mutually defined
//
// source of information:
// https://portal3.gefeg.com/projectdata/invoice/deliverables/installed/publishingproject/xrechnung%202.2.0%20-%20(ab%2001.08.2022)/xrechnung_cii_v2.2.0_01.08.2022.scm/html/DE/02216.htm
func isValidTypeCode(code string) bool {
	if code == "ZZZ" {
		return true
	}

	num, err := strconv.Atoi(code)
	if err != nil {
		return false
	}
	return ((num >= 1 && num <= 70) || (num >= 74 && num <= 78) || (num >= 91 && num <= 97))
}

// FindTypeCode finds the type code for invoice.
// The default return value is TypeCodeInstrumentNotDefined
func FindTypeCode(inv *bill.Invoice) string {
	code := TypeCodeInstrumentNotDefined
	if inv.Payment == nil || inv.Payment.Instructions == nil {
		return TypeCodeInstrumentNotDefined
	}
	if !isValidTypeCode(code) {
		return TypeCodeInstrumentNotDefined
	}
	return TypeCodeInstrumentNotDefined
}
