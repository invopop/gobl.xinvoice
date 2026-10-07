package xinvoice_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/invopop/gobl"
	cii "github.com/invopop/gobl.cii"
	xinvoice "github.com/invopop/gobl.de.xinvoice"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/convert"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// loadCIIEnvelope loads an envelope from test/data/cii.
func loadCIIEnvelope(t *testing.T, name string) *gobl.Envelope {
	t.Helper()
	return loadTestEnvelope(t, filepath.Join("test", "data", "cii", name))
}

// TestCIIFormatGolden converts each envelope in test/data/cii with its
// gobl.cii format and compares the XML against the golden files.
func TestCIIFormatGolden(t *testing.T) {
	for _, tc := range []struct {
		dir    string
		format cii.Format
	}{
		{"xrechnung", xinvoice.CIIFormatXRechnung},
		{"zugferd", xinvoice.CIIFormatZUGFeRD},
	} {
		t.Run(tc.dir, func(t *testing.T) {
			examples, err := filepath.Glob(filepath.Join("test", "data", "cii", tc.dir, "*.json"))
			require.NoError(t, err)
			require.NotEmpty(t, examples)
			for _, example := range examples {
				name := filepath.Base(example)
				t.Run(name, func(t *testing.T) {
					env := loadCIIEnvelope(t, filepath.Join(tc.dir, name))
					doc, err := cii.Export(env, cii.WithFormat(tc.format))
					require.NoError(t, err)
					data, err := cii.Encode(doc)
					require.NoError(t, err)

					golden := filepath.Join("test", "data", "cii", tc.dir, "out", strings.Replace(name, ".json", ".xml", 1))
					if *update {
						require.NoError(t, os.WriteFile(golden, data, 0644))
					}
					want, err := os.ReadFile(golden)
					require.NoError(t, err)
					assert.Equal(t, string(want), string(data))
				})
			}
		})
	}
}

func TestCIIFormats(t *testing.T) {
	t.Run("XRechnung identity headers", func(t *testing.T) {
		env := loadCIIEnvelope(t, "xrechnung/invoice-de-de.json")
		out, err := cii.ExportInvoice(env, cii.WithFormat(xinvoice.CIIFormatXRechnung))
		require.NoError(t, err)
		assert.Equal(t, xinvoice.CustomizationIDXRechnung, out.ExchangedContext.GuidelineContext.ID)
		assert.Equal(t, xinvoice.ProfileIDPeppolBilling, out.ExchangedContext.BusinessContext.ID)
	})

	t.Run("XRechnung found by its identifiers", func(t *testing.T) {
		f := cii.FindFormat(xinvoice.CustomizationIDXRechnung, xinvoice.ProfileIDPeppolBilling)
		require.NotNil(t, f)
		assert.Equal(t, xinvoice.KeyXRechnungCII, f.Key)
	})

	t.Run("ZUGFeRD profiles", func(t *testing.T) {
		for _, tc := range []struct {
			format    cii.Format
			guideline string
			vesID     string
		}{
			{xinvoice.CIIFormatZUGFeRDBasic, "urn:cen.eu:en16931:2017#compliant#urn:zugferd.de:2p0:basic", "de.zugferd:basic:2.5.2"},
			{xinvoice.CIIFormatZUGFeRD, "urn:cen.eu:en16931:2017", "de.zugferd:en16931:2.5.2"},
			{xinvoice.CIIFormatZUGFeRDExtended, "urn:cen.eu:en16931:2017#conformant#urn:zugferd.de:2p0:extended", "de.zugferd:extended:2.5.2"},
		} {
			t.Run(tc.format.Key.String(), func(t *testing.T) {
				assert.Equal(t, tc.guideline, tc.format.GuidelineID)
				assert.Equal(t, tc.vesID, tc.format.VESID)
				assert.Empty(t, tc.format.OutputGuidelineID)
			})
		}
	})

	t.Run("ZUGFeRD hybrid profiles found by their guideline", func(t *testing.T) {
		for _, want := range []cii.Format{xinvoice.CIIFormatZUGFeRDBasic, xinvoice.CIIFormatZUGFeRDExtended} {
			f := cii.FindFormat(want.GuidelineID, "")
			require.NotNil(t, f)
			assert.Equal(t, want.Key, f.Key)
		}
	})

	t.Run("ZUGFeRD ship-to from the customer", func(t *testing.T) {
		env := loadCIIEnvelope(t, "zugferd/standard-invoice.json")
		inv, ok := env.Extract().(*bill.Invoice)
		require.True(t, ok)
		require.Nil(t, inv.Delivery)
		out, err := cii.ExportInvoice(env, cii.WithFormat(xinvoice.CIIFormatZUGFeRD))
		require.NoError(t, err)
		require.NotNil(t, out.Transaction.Delivery.Receiver)
		assert.Equal(t, inv.Customer.Name, out.Transaction.Delivery.Receiver.Name)
	})

	t.Run("no ZUGFeRD ship-to with EN 16931", func(t *testing.T) {
		env := loadCIIEnvelope(t, "zugferd/standard-invoice.json")
		out, err := cii.ExportInvoice(env)
		require.NoError(t, err)
		assert.Nil(t, out.Transaction.Delivery.Receiver)
	})
}

// breakdownEnvelope loads the fixture with a breakdown on the first line.
func breakdownEnvelope(t *testing.T, name string) *gobl.Envelope {
	t.Helper()
	env := loadCIIEnvelope(t, name)
	inv, ok := env.Extract().(*bill.Invoice)
	require.True(t, ok)

	design := num.MakeAmount(3000, 2)
	build := num.MakeAmount(6000, 2)
	inv.Lines[0].Discounts, inv.Lines[0].Charges = nil, nil
	inv.Lines[0].Breakdown = []*bill.SubLine{
		{Quantity: num.MakeAmount(1, 0), Item: &org.Item{Name: "Design", Price: &design}},
		{
			Quantity:  num.MakeAmount(2, 0),
			Item:      &org.Item{Name: "Build", Price: &build},
			Discounts: []*bill.LineDiscount{{Amount: num.MakeAmount(500, 2), Reason: "Promotion"}},
		},
		{Quantity: num.MakeAmount(1, 0), Item: &org.Item{Name: "Sample, no charge"}},
	}
	require.NoError(t, env.Calculate())
	return env
}

// TestCIIZUGFeRDExtendedSubLines writes a breakdown out as sub-invoice lines
// and reads it back: the breakdown and every total must survive.
func TestCIIZUGFeRDExtendedSubLines(t *testing.T) {
	env := breakdownEnvelope(t, "zugferd/standard-invoice.json")
	inv, ok := env.Extract().(*bill.Invoice)
	require.True(t, ok)

	doc, err := cii.Export(env, cii.WithFormat(xinvoice.CIIFormatZUGFeRDExtended))
	require.NoError(t, err)
	data, err := cii.Encode(doc)
	require.NoError(t, err)

	in, err := cii.Decode(data)
	require.NoError(t, err)
	parsed, err := cii.Import(in)
	require.NoError(t, err)
	out, ok := parsed.Extract().(*bill.Invoice)
	require.True(t, ok)

	assert.False(t, out.HasTags(tax.TagBypass))
	require.Len(t, out.Lines, len(inv.Lines))
	require.Len(t, out.Lines[0].Breakdown, len(inv.Lines[0].Breakdown))
	for j, sl := range inv.Lines[0].Breakdown {
		gs := out.Lines[0].Breakdown[j]
		assert.Equal(t, sl.Item.Name, gs.Item.Name)
		assert.True(t, sl.Quantity.Equals(gs.Quantity))
	}
	assert.Equal(t, inv.Totals.Sum.String(), out.Totals.Sum.String())
	assert.Equal(t, inv.Totals.Tax.String(), out.Totals.Tax.String())
	assert.Equal(t, inv.Totals.Payable.String(), out.Totals.Payable.String())
}

func TestCIIConvertRegister(t *testing.T) {
	t.Run("detect and import", func(t *testing.T) {
		data, err := os.ReadFile(filepath.Join("test", "data", "parse", "invoice-xrechnung-cii-v3.xml"))
		require.NoError(t, err)
		f, err := convert.Detect(data)
		require.NoError(t, err)
		assert.Equal(t, xinvoice.KeyXRechnungCII, f.Key)
	})

	t.Run("export", func(t *testing.T) {
		env := loadCIIEnvelope(t, "xrechnung/invoice-de-de.json")
		out, err := convert.Export(env, xinvoice.KeyXRechnungCII)
		require.NoError(t, err)
		f, err := convert.Detect(out.Data)
		require.NoError(t, err)
		assert.Equal(t, xinvoice.KeyXRechnungCII, f.Key, "detected again")
	})
}
