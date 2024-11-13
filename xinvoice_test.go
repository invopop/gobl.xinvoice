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
	"github.com/invopop/gobl.xinvoice"
	"github.com/invopop/gobl/bill"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDocument(t *testing.T) {
	examples, err := getDataGlob("*.json")
	require.NoError(t, err)

	for _, example := range examples {
		inName := filepath.Base(example)

		for _, format := range []string{"xml", "json"} {
			outName := strings.Replace(inName, ".json", "."+format+".xml", 1)

			t.Run(inName, func(t *testing.T) {
				src, _ := os.Open(filepath.Join(getDataPath(), inName))

				buf := new(bytes.Buffer)
				_, err := buf.ReadFrom(src)
				require.NoError(t, err)

				doc, err := xinvoice.Convert(buf.Bytes(), format)
				require.NoError(t, err)

				output, err := LoadOutputFile(outName)
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
			goblEnv, err := xinvoice.Convert(xmlData, "")
			require.NoError(t, err)

			env := new(gobl.Envelope)
			err = json.Unmarshal(goblEnv, env)
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
			output, err := LoadOutputFile(outName)
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

// LoadOutputFile returns byte data from a file in the `test/data/out` folder
func LoadOutputFile(name string) ([]byte, error) {
	src, _ := os.Open(filepath.Join(getOutPath(), name))

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(src); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// getDataGlob returns a list of files in the `test/data` folder that match the pattern
func getDataGlob(pattern string) ([]string, error) {
	return filepath.Glob(filepath.Join(getDataPath(), pattern))
}

// getOutPath returns the path to the `test/data/out` folder
func getOutPath() string {
	return filepath.Join(getDataPath(), "out")
}

// getDataPath returns the path to the `test/data` folder
func getDataPath() string {
	return filepath.Join(getTestPath(), "data")
}

// getTestPath returns the path to the `test` folder
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
