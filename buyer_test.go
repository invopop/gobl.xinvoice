package xinvoice_test

import (
	"testing"

	"github.com/invopop/gobl.xinvoice/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBuyer(t *testing.T) {
	t.Run("should contain the buyer info from standard invoice", func(t *testing.T) {
		doc, err := test.NewDocumentFrom("invoice-de-de.json")
		require.NoError(t, err)

		assert.Nil(t, err)
		assert.Equal(t, "DE282741168", doc.Transaction.Agreement.Buyer.ID)
		assert.Equal(t, "Sample Consumer", doc.Transaction.Agreement.Buyer.Name)
		assert.Equal(t, "80939", doc.Transaction.Agreement.Buyer.PostalTradeAddress.Postcode)
		assert.Equal(t, "Werner-Heisenberg-Allee", doc.Transaction.Agreement.Buyer.PostalTradeAddress.LineOne)
		assert.Equal(t, "München", doc.Transaction.Agreement.Buyer.PostalTradeAddress.City)
		assert.Equal(t, "DE", doc.Transaction.Agreement.Buyer.PostalTradeAddress.CountryID)
		assert.Equal(t, "email@sample.com", doc.Transaction.Agreement.Buyer.URIUniversalCommunication.URIID)
		assert.Equal(t, "EM", doc.Transaction.Agreement.Buyer.URIUniversalCommunication.SchemeID)
	})

	t.Run("should contain the buyer info from simplified invoice", func(t *testing.T) {
		doc, err := test.NewDocumentFrom("invoice-without-buyers-tax-id.json")
		require.NoError(t, err)

		assert.Nil(t, err)
		assert.Equal(t, "", doc.Transaction.Agreement.Buyer.ID)
		assert.Equal(t, "Sample Consumer", doc.Transaction.Agreement.Buyer.Name)
		assert.Equal(t, "80939", doc.Transaction.Agreement.Buyer.PostalTradeAddress.Postcode)
		assert.Equal(t, "Werner-Heisenberg-Allee", doc.Transaction.Agreement.Buyer.PostalTradeAddress.LineOne)
		assert.Equal(t, "München", doc.Transaction.Agreement.Buyer.PostalTradeAddress.City)
		assert.Equal(t, "DE", doc.Transaction.Agreement.Buyer.PostalTradeAddress.CountryID)
		assert.Equal(t, "email@sample.com", doc.Transaction.Agreement.Buyer.URIUniversalCommunication.URIID)
		assert.Equal(t, "EM", doc.Transaction.Agreement.Buyer.URIUniversalCommunication.SchemeID)
	})
}
