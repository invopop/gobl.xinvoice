package xinvoice_test

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/invopop/gobl"
	xinvoice "github.com/invopop/gobl.xinvoice"
	"github.com/invopop/gobl.xinvoice/test"
	"github.com/invopop/gobl/bill"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var updateOutXtoG = flag.Bool("update", false, "Update the JSON files in the test/data/out directory")
var updateOutGtoX = flag.Bool("update", false, "Update the XML files in the test/data/out directory")

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

			if *updateOutGtoX {
				err = test.SaveOutputFile(outName, data)
				require.NoError(t, err)
			} else {
				assert.Equal(t, output, data, "Output should match the expected XML. Update with --update flag.")
			}
		})
	}
}

func TestNewDocumentGOBL(t *testing.T) {
	examples, err := test.GetDataGlob("*.xml")
	require.NoError(t, err)

	for _, example := range examples {
		inName := filepath.Base(example)
		outName := strings.Replace(inName, ".xml", ".json", 1)

		t.Run(inName, func(t *testing.T) {
			// Load XML data into doc
			xmlData, err := os.ReadFile(example)
			require.NoError(t, err)

			doc := new(xinvoice.XMLDoc)
			err = xml.Unmarshal(xmlData, doc)
			require.NoError(t, err)

			// Now doc should be populated with data from the XML file
			goblEnv, err := xinvoice.NewDocumentGOBL(doc)
			require.NoError(t, err)

			// Extract the invoice from the envelope
			invoice, ok := goblEnv.Extract().(*bill.Invoice)
			require.True(t, ok, "Document should be an invoice")

			// Remove UUID from the invoice
			invoice.UUID = ""

			// Marshal only the invoice
			data, err := json.MarshalIndent(invoice, "", "  ")
			require.NoError(t, err)

			// Load the expected output
			output, err := test.LoadOutputFile(outName)
			assert.NoError(t, err)

			// Parse the expected output to extract the invoice
			var expectedEnv gobl.Envelope
			err = json.Unmarshal(output, &expectedEnv)
			require.NoError(t, err)

			expectedInvoice, ok := expectedEnv.Extract().(*bill.Invoice)
			require.True(t, ok, "Expected document should be an invoice")

			// Remove UUID from the expected invoice
			expectedInvoice.UUID = ""

			// Marshal the expected invoice
			expectedData, err := json.MarshalIndent(expectedInvoice, "", "  ")
			require.NoError(t, err)

			if *updateOutXtoG {
				err = test.SaveOutputFile(outName, data)
				require.NoError(t, err)
			} else {
				assert.JSONEq(t, string(expectedData), string(data), "Invoice should match the expected JSON. Update with --update flag.")
			}
		})
	}
}

func getRootFolder() string {
	cwd, _ := os.Getwd()

	// for !isRootFolder(cwd) {
	// 	cwd = removeLastEntry(cwd)
	// }

	return cwd
}
