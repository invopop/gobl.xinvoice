package xinvoice_test

import (
	"testing"

	xtog "github.com/invopop/gobl.xinvoice/internal/xtog"
	"github.com/invopop/gobl.xinvoice/test"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsePayment(t *testing.T) {
	// Read the XML file
	doc, err := test.LoadTestXMLDoc("invoice-test-4.xml")
	require.NoError(t, err)

	payment := xtog.ParsePayment(&doc.SupplyChainTradeTransaction.ApplicableHeaderTradeSettlement)

	assert.NotNil(t, payment)

	assert.Equal(t, "Ilona Schidt", payment.Payee.Name)
	assert.Len(t, payment.Payee.Addresses, 1)
	assert.Equal(t, "Auf der Wacht 45", payment.Payee.Addresses[0].Street)
	assert.Equal(t, "Kaisersesch", payment.Payee.Addresses[0].Locality)
	assert.Equal(t, "56799", payment.Payee.Addresses[0].Code)
	assert.Equal(t, l10n.ISOCountryCode("DE"), payment.Payee.Addresses[0].Country)

	assert.NotNil(t, payment.Terms)
	assert.Equal(t, "Partial Payment", payment.Terms.Detail)
	assert.Len(t, payment.Terms.DueDates, 1)
	assert.Equal(t, "2024-10-01", payment.Terms.DueDates[0].Date.String())
	expectedAmount, _ := num.AmountFromString("20.00")
	assert.Equal(t, expectedAmount, payment.Terms.DueDates[0].Amount)

	assert.NotNil(t, payment.Instructions)
	assert.Equal(t, pay.MeansKeyDebitTransfer, payment.Instructions.Key)
	assert.Equal(t, "Barzahlung", payment.Instructions.Detail)

	assert.NotNil(t, payment.Instructions.Card)
	assert.Equal(t, "3456", payment.Instructions.Card.Last4)
	assert.Equal(t, "Ilona Schidt", payment.Instructions.Card.Holder)

	assert.Len(t, payment.Instructions.CreditTransfer, 1)
	assert.Equal(t, "123456789012345678", payment.Instructions.CreditTransfer[0].IBAN)
}
