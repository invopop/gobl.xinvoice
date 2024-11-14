// Package xinvoice helps convert GOBL into XRechnung and Factur-X documents and vice versa.
package xinvoice

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl.cii"
	"github.com/invopop/gobl.cii/document"
	"github.com/invopop/gobl.ubl"
)

const (
	ciiHeader = "CrossIndustryInvoice"
	ublHeader = "Invoice"
)

// Currently Factur-X and Zugferd have the same context header, but
// keeping them separate to avoid confusion.
var mapFormatGuideline = map[string]string{
	"xrechnung": "urn:cen.eu:en16931:2017#compliant#urn:xeinkauf.de:kosit:xrechnung_3.0",
	"facturx":   "urn:cen.eu:en16931:2017#conformant#urn:factur-x.eu:1p0:extended",
	"zugferd":   "urn:cen.eu:en16931:2017#conformant#urn:factur-x.eu:1p0:extended",
}

// Convert converts a GOBL envelope into an XRechnung or Factur-X document.
func Convert(d []byte, f string) ([]byte, error) {
	j := json.Valid(d)
	var o []byte

	if j {
		env := new(gobl.Envelope)
		if err := json.Unmarshal(d, env); err != nil {
			return nil, fmt.Errorf("parsing input as GOBL Envelope: %w", err)
		}

		doc, err := cii.ToCII(env)
		if err != nil {
			return nil, fmt.Errorf("building XRechnung and Factur-X document: %w", err)
		}

		g, ok := mapFormatGuideline[f]
		if !ok {
			return nil, fmt.Errorf("invalid format %q - must be one of: xrechnung, facturx, zugferd", f)
		}

		// Add the guideline context
		doc.ExchangedContext.GuidelineContext.ID = g

		// Add particular fields required by the format
		switch f {
		case "xrechnung":
			doc.ExchangedContext.BusinessContext = &document.ExchangedContextParameter{
				ID: "urn:fdc:peppol.eu:2017:poacc:billing:01:1.0",
			}
		}

		o, err = doc.Bytes()
		if err != nil {
			return nil, fmt.Errorf("generating XRechnung and Factur-X xml: %w", err)
		}
	} else {
		// Assume XML if not JSON
		r, err := extractRootName(d)
		if err != nil {
			return nil, fmt.Errorf("extracting root name: %w", err)
		}

		var env *gobl.Envelope
		if r == ciiHeader {
			env, err = cii.ToGOBL(d)
			if err != nil {
				return nil, fmt.Errorf("converting CII to GOBL: %w", err)
			}
		} else if r == ublHeader {
			env, err = ubl.ToGOBL(d)
			if err != nil {
				return nil, fmt.Errorf("converting UBL to GOBL: %w", err)
			}
		} else {
			return nil, fmt.Errorf("unknown XML format: %s", r)
		}

		o, err = json.MarshalIndent(env, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("generating JSON output: %w", err)
		}
	}

	return o, nil
}

// Helper function to extract the root element name or specific header
func extractRootName(d []byte) (string, error) {
	dc := xml.NewDecoder(bytes.NewReader(d))
	for {
		tk, err := dc.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("error parsing XML: %w", err)
		}
		switch t := tk.(type) {
		case xml.StartElement:
			return t.Name.Local, nil
		}
	}
	return "", fmt.Errorf("no root element found")
}
