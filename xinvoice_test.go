// Package test provides tools for testing the library
package xinvoice_test

import (
	"testing"

	"github.com/invopop/gobl.xinvoice/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDocument(t *testing.T) {
	t.Run("should return bytes of the xml document", func(t *testing.T) {
		doc, err := test.NewDocumentFrom("invoice.json")
		require.NoError(t, err)

		data, err := doc.Bytes()
		require.NoError(t, err)

		output, err := test.LoadOutputFile("invoice-hardcoded.xml")
		require.NoError(t, err)

		assert.Equal(t, output, data)
	})
}
