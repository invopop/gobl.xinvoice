package to_gobl_test

import (
	"testing"

	"github.com/invopop/gobl.xinvoice/to_gobl"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Define tests for the ParseParty function
func TestParseParty(t *testing.T) {
	doc, err := LoadTestXMLDoc("invoice-test-1.xml")
	require.NoError(t, err)

	seller := to_gobl.ParseParty(&doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.SellerTradeParty)
	require.NotNil(t, seller)

	assert.Equal(t, "TimeOut Immo GmbH", seller.Name)
	assert.Equal(t, l10n.TaxCountryCode("DE"), seller.TaxID.Country)
	assert.Equal(t, cbc.Code("363188747"), seller.TaxID.Code)

	buyer := to_gobl.ParseParty(&doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.BuyerTradeParty)
	require.NotNil(t, buyer)

	assert.Equal(t, "mih GmbH", buyer.Name)
	assert.Equal(t, "An der Wurth 2 – 4", buyer.Addresses[0].Street)
	assert.Equal(t, "Horstmar", buyer.Addresses[0].Locality)
	assert.Equal(t, "48612", buyer.Addresses[0].Code)
	assert.Equal(t, l10n.ISOCountryCode("DE"), buyer.Addresses[0].Country)
}
