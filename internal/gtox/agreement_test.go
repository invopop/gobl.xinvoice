package xinvoice_test

import (
	"testing"

	"github.com/invopop/gobl.xinvoice/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAgreement(t *testing.T) {
	t.Run("should contain the agreement info from standard invoice", func(t *testing.T) {
		doc, err := test.NewDocumentFrom("invoice-de-de.json")
		require.NoError(t, err)

		assert.Nil(t, err)
		assert.Equal(t, "XR-2024-2", doc.Transaction.Agreement.BuyerReference)
	})

	t.Run("should contain the agreement info from credit note", func(t *testing.T) {
		doc, err := test.NewDocumentFrom("credit-note.json")
		require.NoError(t, err)

		assert.Nil(t, err)
		assert.Equal(t, "XR-2024-4", doc.Transaction.Agreement.BuyerReference)
	})
}
