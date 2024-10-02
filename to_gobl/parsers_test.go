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

func TestParseLines(t *testing.T) {
	doc, err := LoadTestXMLDoc("invoice-test-1.xml")
	require.NoError(t, err)

	lines := to_gobl.ParseLines(&doc.SupplyChainTradeTransaction)
	require.Len(t, lines, 2)

	assert.Equal(t, "2h Beschaffung + Aufbau des neuen Tisches a 25€/h netto + 7% MwSt.", lines[0].Item.Name)
	assert.Equal(t, "5350.00", lines[0].Item.Price.String())
	assert.Equal(t, "1", lines[0].Quantity.String())
	assert.Equal(t, "VAT", string(lines[0].Taxes[0].Category))
	assert.Equal(t, "7.00", lines[0].Taxes[0].Percent.String())

	assert.Equal(t, "1x Couchtisch inklusive 19% MwSt.", lines[1].Item.Name)
	assert.Equal(t, "2000.00", lines[1].Item.Price.String())
	assert.Equal(t, "1", lines[1].Quantity.String())
	assert.Equal(t, "VAT", string(lines[1].Taxes[0].Category))
	assert.Equal(t, "19.00", lines[1].Taxes[0].Percent.String())
}
