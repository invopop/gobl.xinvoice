// Package test provides tools for testing the library
package xinvoice_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/invopop/gobl"
	xinvoice "github.com/invopop/gobl.xinvoice"
	"github.com/invopop/gobl/bill"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	xmlPattern  = "*.xml"
	jsonPattern = "*.json"
)

func TestGtoX(t *testing.T) {

	examples, err := getDataGlob(jsonPattern)
	require.NoError(t, err)

	for _, example := range examples {
		inName := filepath.Base(example)
		for _, format := range []string{"xrechnung-cii", "xrechnung-ubl", "facturx", "zugferd"} {
			outName := strings.Replace(inName, ".json", "-"+format+".xml", 1)

			t.Run(inName+"-"+format, func(t *testing.T) {
				env, err := loadTestEnvelope(inName)
				require.NoError(t, err)

				var doc []byte
				switch format {
				case "xrechnung-cii":
					doc, err = xinvoice.ConvertToXRechnungCII(env)
				case "xrechnung-ubl":
					doc, err = xinvoice.ConvertToXRechnungUBL(env)
				case "facturx":
					doc, err = xinvoice.ConvertToFacturX(env)
				case "zugferd":
					doc, err = xinvoice.ConvertToZUGFeRD(env)
				}

				require.NoError(t, err)

				output, err := loadOutputFile(outName)
				assert.NoError(t, err)
				assert.Equal(t, output, doc, "Output should match the expected XML. Update with --update flag.")

			})
		}
	}
}

func TestXtoG(t *testing.T) {
	examples, err := getDataGlob("*.xml")
	require.NoError(t, err)

	for _, example := range examples {
		inName := filepath.Base(example)
		outName := strings.Replace(inName, ".xml", ".json", 1)

		t.Run(inName, func(t *testing.T) {
			// Load XML data
			xmlData, err := os.ReadFile(example)
			require.NoError(t, err)

			// Convert XML to GOBL
			env, err := xinvoice.ConvertToGOBL(xmlData)
			require.NoError(t, err)

			// Extract the invoice from the envelope
			invoice, ok := env.Extract().(*bill.Invoice)
			require.True(t, ok, "Document should be an invoice")
			// Remove UUID from the invoice
			invoice.UUID = ""

			// Marshal only the invoice
			data, err := json.MarshalIndent(invoice, "", "  ")
			require.NoError(t, err)

			// Load the expected output
			output, err := loadOutputFile(outName)
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

			assert.JSONEq(t, string(expectedData), string(data), "Invoice should match the expected JSON. Update with --update flag.")
		})
	}
}

// loadTestEnvelope returns a GOBL Envelope from a file in the `test/data` folder
func loadTestEnvelope(name string) (*gobl.Envelope, error) {
	src, _ := os.Open(filepath.Join(getConversionTypePath(jsonPattern), name))
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(src); err != nil {
		return nil, err
	}
	env := new(gobl.Envelope)
	if err := json.Unmarshal(buf.Bytes(), env); err != nil {
		return nil, err
	}

	return env, nil
}

func loadOutputFile(name string) ([]byte, error) {
	var pattern string
	if strings.HasSuffix(name, ".json") {
		pattern = xmlPattern
	} else {
		pattern = jsonPattern
	}
	src, _ := os.Open(filepath.Join(getOutPath(pattern), name))
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(src); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func getDataGlob(pattern string) ([]string, error) {
	return filepath.Glob(filepath.Join(getConversionTypePath(pattern), pattern))
}

func getOutPath(pattern string) string {
	return filepath.Join(getConversionTypePath(pattern), "out")
}

func getDataPath() string {
	return filepath.Join(getTestPath(), "data")
}

func getConversionTypePath(pattern string) string {
	if pattern == xmlPattern {
		return filepath.Join(getDataPath(), "xtog")
	}
	return filepath.Join(getDataPath(), "gtox")
}

func getTestPath() string {
	return filepath.Join(getRootFolder(), "test")
}

func getRootFolder() string {
	cwd, _ := os.Getwd()

	for !isRootFolder(cwd) {
		cwd = removeLastEntry(cwd)
	}

	return cwd
}

func isRootFolder(dir string) bool {
	files, _ := os.ReadDir(dir)

	for _, file := range files {
		if file.Name() == "go.mod" {
			return true
		}
	}

	return false
}

func removeLastEntry(dir string) string {
	lastEntry := "/" + filepath.Base(dir)
	i := strings.LastIndex(dir, lastEntry)
	return dir[:i]
}
