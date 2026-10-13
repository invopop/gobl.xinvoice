package xinvoice_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/invopop/gobl"
	xinvoice "github.com/invopop/gobl.de.xinvoice"
	ubl "github.com/invopop/gobl.ubl"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/convert"
	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const staticUUID uuid.UUID = "0195ce71-dc9c-72c8-bf2c-9890a4a9f0a2"

// loadUBLEnvelope loads, calculates, and validates a GOBL envelope from
// test/data/ubl, with fixed IDs so the output is stable.
func loadUBLEnvelope(t *testing.T, name string) *gobl.Envelope {
	t.Helper()
	return loadTestEnvelope(t, filepath.Join("test", "data", "ubl", name))
}

// loadTestEnvelope loads, calculates, and validates the GOBL envelope at
// path, with fixed IDs so the output is stable.
func loadTestEnvelope(t *testing.T, path string) *gobl.Envelope {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	env := new(gobl.Envelope)
	require.NoError(t, json.Unmarshal(data, env))
	env.Head.UUID = staticUUID
	if inv, ok := env.Extract().(*bill.Invoice); ok {
		inv.UUID = staticUUID
	}
	require.NoError(t, env.Calculate())
	require.NoError(t, env.Validate())
	if *update {
		data, err := json.MarshalIndent(env, "", "\t")
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(path, data, 0644))
	}
	return env
}

// TestUBLFormatXRechnungGolden converts each envelope in
// test/data/ubl/xrechnung with the gobl.ubl format and compares the XML
// against the golden files.
func TestUBLFormatXRechnungGolden(t *testing.T) {
	examples, err := filepath.Glob(filepath.Join("test", "data", "ubl", "xrechnung", "*.json"))
	require.NoError(t, err)
	require.NotEmpty(t, examples)
	for _, example := range examples {
		name := filepath.Base(example)
		t.Run(name, func(t *testing.T) {
			env := loadUBLEnvelope(t, filepath.Join("xrechnung", name))
			doc, err := ubl.ExportInvoice(env, ubl.WithFormat(xinvoice.UBLFormatXRechnung))
			require.NoError(t, err)
			data, err := ubl.Encode(doc)
			require.NoError(t, err)

			golden := filepath.Join("test", "data", "ubl", "xrechnung", "out", strings.Replace(name, ".json", ".xml", 1))
			if *update {
				require.NoError(t, os.WriteFile(golden, data, 0644))
			}
			want, err := os.ReadFile(golden)
			require.NoError(t, err)
			assert.Equal(t, string(want), string(data))
		})
	}
}

func TestUBLFormatXRechnung(t *testing.T) {
	t.Run("identity headers", func(t *testing.T) {
		env := loadUBLEnvelope(t, "xrechnung/invoice-xr-minimal.json")
		doc, err := ubl.Export(env, ubl.WithFormat(xinvoice.UBLFormatXRechnung))
		require.NoError(t, err)
		out, ok := doc.(*ubl.Invoice)
		require.True(t, ok)
		assert.Equal(t, xinvoice.CustomizationIDXRechnung, out.CustomizationID)
		assert.Equal(t, xinvoice.ProfileIDPeppolBilling, out.ProfileID.Value)
	})

	t.Run("found by its identifiers", func(t *testing.T) {
		f := ubl.FindFormat(xinvoice.CustomizationIDXRechnung, xinvoice.ProfileIDPeppolBilling)
		require.NotNil(t, f)
		assert.Equal(t, xinvoice.KeyXRechnungUBL, f.Key)
	})

	t.Run("VESIDs", func(t *testing.T) {
		env := loadUBLEnvelope(t, "xrechnung/invoice-xr-minimal.json")
		inv, ok := env.Extract().(*bill.Invoice)
		require.True(t, ok)
		assert.Equal(t, "de.xrechnung:ubl-invoice:3.0.2", xinvoice.UBLFormatXRechnung.VESID(inv))
	})
}

func TestConvertRegister(t *testing.T) {
	t.Run("detect and import", func(t *testing.T) {
		data, err := os.ReadFile(filepath.Join("test", "data", "parse", "invoice-xrechnung-ubl-v3.xml"))
		require.NoError(t, err)
		f, err := convert.Detect(data)
		require.NoError(t, err)
		assert.Equal(t, xinvoice.KeyXRechnungUBL, f.Key)

		env, err := convert.Import(data)
		require.NoError(t, err)
		inv, ok := env.Extract().(*bill.Invoice)
		require.True(t, ok)
		assert.Contains(t, inv.GetAddons(), xinvoice.UBLFormatXRechnung.Addons[0])
	})

	t.Run("export", func(t *testing.T) {
		env := loadUBLEnvelope(t, "xrechnung/invoice-xr-minimal.json")
		out, err := convert.Export(env, xinvoice.KeyXRechnungUBL, ubl.FormatEN16931.Key)
		require.NoError(t, err)
		assert.Equal(t, xinvoice.KeyXRechnungUBL, out.Format.Key)

		f, err := convert.Detect(out.Data)
		require.NoError(t, err)
		assert.Equal(t, xinvoice.KeyXRechnungUBL, f.Key, "detected again")
	})
}
