// Package test provides tools for testing the library
package xinvoice_test

import (
	"flag"
	"path/filepath"
	"strings"
	"testing"

	"github.com/invopop/gobl.xinvoice/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var updateOut = flag.Bool("update", false, "Update the XML files in the test/data/out directory")

func TestNewDocument(t *testing.T) {
	schema, err := test.LoadSchema("schema.xsd")
	require.NoError(t, err)

	examples, err := test.GetDataGlob("*.json")
	require.NoError(t, err)

	for _, example := range examples {
		inName := filepath.Base(example)
		outName := strings.Replace(inName, ".json", ".xml", 1)

		t.Run(inName, func(t *testing.T) {
			doc, err := test.NewDocumentFrom(inName)
			require.NoError(t, err)

			data, err := doc.Bytes()
			require.NoError(t, err)

			err = test.ValidateXML(schema, data)
			require.NoError(t, err)

			output, err := test.LoadOutputFile(outName)
			assert.NoError(t, err)

			if *updateOut {
				err = test.SaveOutputFile(outName, data)
				require.NoError(t, err)
			} else {
				assert.Equal(t, output, data, "Output should match the expected XML. Update with --update flag.")
			}
		})
	}
}
