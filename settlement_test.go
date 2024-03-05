package xinvoice_test

import (
	"testing"

	xinvoice "github.com/invopop/gobl.xinvoice"
	"github.com/invopop/gobl.xinvoice/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSettlement(t *testing.T) {
	t.Run("should contain the transaction settlement", func(t *testing.T) {
		doc, err := test.NewDocumentFrom("invoice-de-de.json")
		require.NoError(t, err)

		assert.Nil(t, err)
		assert.Equal(t, "EUR", doc.Transaction.Settlement.Currency)
		assert.Equal(t, "1", doc.Transaction.Settlement.TypeCode)
		assert.Equal(t, "lorem ipsum", doc.Transaction.Settlement.Description)
		assert.Equal(t, "1800.00", doc.Transaction.Settlement.Summary.TotalAmount)
		assert.Equal(t, "1800.00", doc.Transaction.Settlement.Summary.TaxBasisTotalAmount)
		assert.Equal(t, "2142.00", doc.Transaction.Settlement.Summary.GrandTotalAmount)
		assert.Equal(t, "2142.00", doc.Transaction.Settlement.Summary.DuePayableAmount)
		assert.Equal(t, "342.00", doc.Transaction.Settlement.Summary.TaxTotalAmount.Amount)
		assert.Equal(t, "EUR", doc.Transaction.Settlement.Summary.TaxTotalAmount.Currency)
	})

	t.Run("should set ReferencedDocument for correction invoice", func(t *testing.T) {
		doc, err := test.NewDocumentFrom("correction-invoice.json")
		require.NoError(t, err)

		assert.Equal(t, "SAMPLE-001", doc.Transaction.Settlement.ReferencedDocument.IssuerAssignedID)
		assert.Equal(t, &xinvoice.Date{Date: "20240213", Format: "102"}, doc.Transaction.Settlement.ReferencedDocument.IssueDate)
	})
}
