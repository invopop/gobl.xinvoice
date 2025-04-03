// Package test provides tools for testing the library
package xinvoice_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/invopop/gobl"
	xinvoice "github.com/invopop/gobl.xinvoice"
	"github.com/invopop/gobl/addons/de/xrechnung"
	"github.com/invopop/gobl/addons/de/zugferd"
	"github.com/invopop/gobl/addons/fr/facturx"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	pathPatternXML  = "*.xml"
	pathPatternJSON = "*.json"
	pathConvert     = "gtox"
	pathParse       = "xtog"
	pathOut         = "out"

	staticUUID uuid.UUID = "0195ce71-dc9c-72c8-bf2c-9890a4a9f0a2"
)

// updateOut is a flag that can be set to update example files
var updateOut = flag.Bool("update", false, "Update the example files in test/data")

/*
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
*/

func TestConvertInvoice(t *testing.T) {
	examples := findSourceFiles(t, pathConvert, pathPatternJSON)
	formats := []string{
		"xrechnung-cii",
		"xrechnung-ubl",
		"facturx",
		"zugferd",
	}
	for _, example := range examples {
		inName := filepath.Base(example)
		env := loadEnvelope(t, inName)
		writeEnvelope(example, env) // Update source if required
		t.Run(inName, func(t *testing.T) {
			for _, format := range formats {
				env := loadEnvelope(t, inName)
				outName := strings.Replace(inName, ".json", "-"+format+".xml", 1)

				t.Run(format, func(t *testing.T) {
					var data []byte
					var err error
					switch format {
					case "xrechnung-cii":
						addInvoiceAddon(env, xrechnung.V3)
						data, err = xinvoice.ConvertToXRechnungCII(env)
					case "xrechnung-ubl":
						addInvoiceAddon(env, xrechnung.V3)
						data, err = xinvoice.ConvertToXRechnungUBL(env)
					case "facturx":
						addInvoiceAddon(env, facturx.V1)
						data, err = xinvoice.ConvertToFacturX(env)
					case "zugferd":
						addInvoiceAddon(env, zugferd.V2)
						data, err = xinvoice.ConvertToZUGFeRD(env)
					}
					require.NoError(t, err)

					if *updateOut {
						err = os.WriteFile(outputFilepath(pathConvert, outName), data, 0644)
						require.NoError(t, err)
					}

					output := loadOutputFile(t, pathConvert, outName)
					assert.Equal(t, string(output), string(data), "Output should match the expected XML. Update with --update flag.")
				})
			}
		})
	}
}

func addInvoiceAddon(env *gobl.Envelope, addon cbc.Key) {
	inv := env.Extract().(*bill.Invoice)
	inv.Addons.List = append(inv.Addons.List, addon)
	env.Calculate()
}

func TestParseInvoice(t *testing.T) {
	examples := findSourceFiles(t, pathParse, pathPatternXML)

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

			// Reset UUIDs
			env.Head.UUID = staticUUID
			if inv, ok := env.Extract().(*bill.Invoice); ok {
				inv.UUID = staticUUID
			}
			require.NoError(t, env.Calculate())

			writeEnvelope(dataPath(pathParse, pathOut, outName), env)

			// Extract the invoice from the envelope
			inv, ok := env.Extract().(*bill.Invoice)
			require.True(t, ok, "Document should be an invoice")

			// Marshal only the invoice
			data, err := json.MarshalIndent(inv, "", "\t")
			require.NoError(t, err)

			// Load the expected output
			output := loadOutputFile(t, pathParse, outName)

			// Parse the expected output to extract the invoice
			var expectedEnv gobl.Envelope
			err = json.Unmarshal(output, &expectedEnv)
			require.NoError(t, err)

			expectedInvoice, ok := expectedEnv.Extract().(*bill.Invoice)
			require.True(t, ok, "Expected document should be an invoice")

			// Marshal the expected invoice
			expectedData, err := json.MarshalIndent(expectedInvoice, "", "  ")
			require.NoError(t, err)

			assert.JSONEq(t, string(expectedData), string(data), "Invoice should match the expected JSON. Update with --update flag.")
		})
	}
}

func parseInvoiceFrom(t *testing.T, name string) (*gobl.Envelope, error) {
	t.Helper()
	path := dataPath(pathParse, name)
	src, err := os.Open(path)
	require.NoError(t, err)
	defer func() {
		if cerr := src.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	data, err := io.ReadAll(src)
	if err != nil {
		require.NoError(t, err)
	}
	return xinvoice.ConvertToGOBL(data)
}

// loadEnvelopeWithUpdate returns a GOBL Envelope from a file in the `test/data/convert` folder
func loadEnvelope(t *testing.T, name string) *gobl.Envelope {
	t.Helper()
	path := dataPath(pathConvert, name)

	src, _ := os.Open(path)
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(src)
	require.NoError(t, err)

	env := new(gobl.Envelope)
	require.NoError(t, json.Unmarshal(buf.Bytes(), env))

	// Clear the IDs
	env.Head.UUID = staticUUID
	if inv, ok := env.Extract().(*bill.Invoice); ok {
		inv.UUID = staticUUID
	}
	require.NoError(t, env.Calculate())
	require.NoError(t, env.Validate())

	return env
}

func writeEnvelope(path string, env *gobl.Envelope) {
	if !*updateOut {
		return
	}
	data, err := json.MarshalIndent(env, "", "\t")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		panic(err)
	}
}

func outputFilepath(path, name string) string {
	return filepath.Join(dataPath(path, pathOut, name))
}

func loadOutputFile(t *testing.T, path, name string) []byte {
	t.Helper()
	src, err := os.Open(outputFilepath(path, name))
	require.NoError(t, err)
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(src); err != nil {
		require.NoError(t, err)
	}
	return buf.Bytes()
}

func findSourceFiles(t *testing.T, path, pattern string) []string {
	path = filepath.Join(dataPath(), path, pattern)
	files, err := filepath.Glob(path)
	require.NoError(t, err)
	return files
}

func schemaPath() string {
	return filepath.Join(dataPath(), "schema")
}

func dataPath(files ...string) string {
	files = append([]string{rootFolder(), "test", "data"}, files...)
	return filepath.Join(files...)
}

func rootFolder() string {
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
