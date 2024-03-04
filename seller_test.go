package xinvoice_test

import (
	"testing"

	"github.com/invopop/gobl.xinvoice/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSeller(t *testing.T) {
	t.Run("should contain the seller info", func(t *testing.T) {
		doc, err := test.NewDocumentFrom("invoice-de-de.json")
		require.NoError(t, err)

		assert.Nil(t, err)
		assert.Equal(t, "Provide One GmbH", doc.Transaction.Agreement.Seller.Name)
		assert.Equal(t, "John Doe", doc.Transaction.Agreement.Seller.Contact.PersonName)
		assert.Equal(t, "+49100200300", doc.Transaction.Agreement.Seller.Contact.Phone)
		assert.Equal(t, "billing@example.com", doc.Transaction.Agreement.Seller.Contact.Email)
		assert.Equal(t, "69190", doc.Transaction.Agreement.Seller.PostalTradeAddress.Postcode)
		assert.Equal(t, "Dietmar-Hopp-Allee", doc.Transaction.Agreement.Seller.PostalTradeAddress.LineOne)
		assert.Equal(t, "Walldorf", doc.Transaction.Agreement.Seller.PostalTradeAddress.City)
		assert.Equal(t, "DE", doc.Transaction.Agreement.Seller.PostalTradeAddress.CountryID)
		assert.Equal(t, "billing@example.com", doc.Transaction.Agreement.Seller.URIUniversalCommunication.URIID)
		assert.Equal(t, "EM", doc.Transaction.Agreement.Seller.URIUniversalCommunication.SchemeID)
	})
}
