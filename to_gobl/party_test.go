package to_gobl_test

import (
	"testing"

	"github.com/invopop/gobl.xinvoice/to_gobl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseParty(t *testing.T) {
	doc, err := LoadTestXMLDoc("invoice-test-1.xml")
	require.NoError(t, err)

	seller := to_gobl.ParseParty(&doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.SellerTradeParty)
	require.NotNil(t, seller)

	assert.Equal(t, "Lieferant GmbH", seller.Name)
	assert.Equal(t, "DE", seller.TaxID.Country)
	assert.Equal(t, "123456789", seller.TaxID.Code)

	buyer := to_gobl.ParseParty(&doc.SupplyChainTradeTransaction.ApplicableHeaderTradeAgreement.BuyerTradeParty)
	require.NotNil(t, buyer)

	assert.Equal(t, "Kunden AG Mitte", buyer.Name)
	assert.Equal(t, "DE", buyer.TaxID.Country)
	assert.Equal(t, "987654321", buyer.TaxID.Code)
}
