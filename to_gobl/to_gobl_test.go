package to_gobl_test

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"flag"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl.xinvoice/to_gobl"
	"github.com/invopop/gobl/bill"

	//"github.com/invopop/gobl.xinvoice/xinvoice/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var updateOut = flag.Bool("update", false, "Update the JSON files in the test/data/out directory")

func TestNewDocumentGOBL(t *testing.T) {
	examples, err := GetDataGlob("*.xml")
	require.NoError(t, err)

	for _, example := range examples {
		inName := filepath.Base(example)
		outName := strings.Replace(inName, ".xml", ".json", 1)

		t.Run(inName, func(t *testing.T) {
			// Load XML data into doc
			xmlData, err := os.ReadFile(example)
			require.NoError(t, err)

			doc := new(to_gobl.XMLDoc)
			err = xml.Unmarshal(xmlData, doc)
			require.NoError(t, err)

			// Now doc should be populated with data from the XML file
			goblEnv, err := to_gobl.NewDocumentGOBL(doc)
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

			if *updateOut {
				err = SaveOutputFile(outName, data)
				require.NoError(t, err)
			} else {
				assert.JSONEq(t, string(expectedData), string(data), "Invoice should match the expected JSON. Update with --update flag.")
			}
		})
	}
}

// LoadTestXMLDoc returns a to_gobl.XMLDoc from a file in the test data folder
func LoadTestXMLDoc(name string) (*to_gobl.XMLDoc, error) {
	src, err := os.Open(filepath.Join(GetDataPath(), name))
	if err != nil {
		return nil, err
	}
	defer src.Close()

	inData, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	doc := new(to_gobl.XMLDoc)
	if err := xml.Unmarshal(inData, doc); err != nil {
		return nil, err
	}

	return doc, nil
}

// GetDataGlob returns a list of files in the `test/data` folder that match the pattern
func GetDataGlob(pattern string) ([]string, error) {
	return filepath.Glob(filepath.Join(GetDataPath(), pattern))
}

// LoadOutputFile returns byte data from a file in the `test/data/out` folder
func LoadOutputFile(name string) ([]byte, error) {
	src, _ := os.Open(filepath.Join(GetOutPath(), name))

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(src); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// SaveOutputFile writes byte data to a file in the `test/data/out` folder
func SaveOutputFile(name string, data []byte) error {
	return os.WriteFile(filepath.Join(GetOutPath(), name), data, 0644)
}

// GetOutPath returns the path to the `test/data/out` folder
func GetOutPath() string {
	return filepath.Join(GetDataPath(), "out")
}

// GetDataPath returns the path to the `test/data` folder
func GetDataPath() string {
	return filepath.Join(GetTestPath(), "data")
}

// GetTestPath returns the path to the `test` folder
func GetTestPath() string {
	return filepath.Join(getRootFolder(), "test")
}

func getRootFolder() string {
	cwd, _ := os.Getwd()

	// for !isRootFolder(cwd) {
	// 	cwd = removeLastEntry(cwd)
	// }

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
	lastEntry := "\\" + filepath.Base(dir)
	i := strings.LastIndex(dir, lastEntry)
	return dir[:i]
}
