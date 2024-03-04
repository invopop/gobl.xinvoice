package xinvoice_test

import (
	"testing"

	"github.com/invopop/gobl.xinvoice/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLines(t *testing.T) {
	t.Run("should contain the lines from invoice", func(t *testing.T) {
		doc, err := test.NewDocumentFrom("invoice-de-de.json")
		require.NoError(t, err)

		assert.Nil(t, err)
		assert.Equal(t, "Development services", doc.Transaction.Lines[0].ID)
		assert.Equal(t, "Development services", doc.Transaction.Lines[0].Name)
		assert.Equal(t, "90.00", doc.Transaction.Lines[0].NetPrice)
		assert.Equal(t, "20", doc.Transaction.Lines[0].TradeDelivery.Amount)
		assert.Equal(t, "HUR", doc.Transaction.Lines[0].TradeDelivery.UnitCode)
		assert.Equal(t, "VAT", doc.Transaction.Lines[0].TradeSettlement.ApplicableTradeTax[0].TaxType)
		assert.Equal(t, "S", doc.Transaction.Lines[0].TradeSettlement.ApplicableTradeTax[0].TaxCode)
		assert.Equal(t, "19", doc.Transaction.Lines[0].TradeSettlement.ApplicableTradeTax[0].TaxRatePercent)
		assert.Equal(t, "1800.00", doc.Transaction.Lines[0].TradeSettlement.Sum)
	})
}
