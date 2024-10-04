package xinvoice_test

import (
	"testing"

	xtog "github.com/invopop/gobl.xinvoice/internal/xtog"
	"github.com/invopop/gobl.xinvoice/test"
	"github.com/invopop/gobl/num"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Define tests for the ParseXMLLines function
func TestParseLines(t *testing.T) {
	doc, err := test.LoadTestXMLDoc("invoice-test-1.xml")
	require.NoError(t, err)

	lines := xtog.ParseXMLLines(&doc.SupplyChainTradeTransaction)
	require.Len(t, lines, 2)

	assert.Equal(t, "2h Beschaffung + Aufbau des neuen Tisches a 25€/h netto + 7% MwSt.", lines[0].Item.Name)
	assert.Equal(t, num.MakeAmount(5350, 0), lines[0].Item.Price)
	assert.Equal(t, num.MakeAmount(1, 0), lines[0].Quantity)
	assert.Equal(t, "VAT", string(lines[0].Taxes[0].Category))
	assert.Equal(t, num.MakePercentage(7, 0), *lines[0].Taxes[0].Percent)

	assert.Equal(t, "1x Couchtisch inklusive 19% MwSt.", lines[1].Item.Name)
	assert.Equal(t, num.MakeAmount(149, 0), lines[1].Item.Price)
	assert.Equal(t, num.MakeAmount(1, 0), lines[1].Quantity)
	assert.Equal(t, "VAT", string(lines[1].Taxes[0].Category))
	assert.Equal(t, num.MakePercentage(19, 0), *lines[1].Taxes[0].Percent)
}
