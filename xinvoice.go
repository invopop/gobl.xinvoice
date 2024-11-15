// Package xinvoice helps convert GOBL into XRechnung and Factur-X documents and vice versa.
package xinvoice

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/invopop/gobl"
	cii "github.com/invopop/gobl.cii"
	"github.com/invopop/gobl.cii/document"
	ubl "github.com/invopop/gobl.ubl"
)

const (
	ciiHeader = "CrossIndustryInvoice"
	ublHeader = "Invoice"

	xRechnungGuideline = "urn:cen.eu:en16931:2017#compliant#urn:xeinkauf.de:kosit:xrechnung_3.0"
	xRechnungProfile   = "urn:fdc:peppol.eu:2017:poacc:billing:01:1.0"

	// Currently Factur-X and Zugferd have the same context header, but
	// keeping them separate to avoid confusion.
	facturXGuideline = "urn:cen.eu:en16931:2017#conformant#urn:factur-x.eu:1p0:extended"
	zugferdGuideline = "urn:cen.eu:en16931:2017#conformant#urn:factur-x.eu:1p0:extended"
)

// ConvertToGOBL converts an XRechnung, Factur-X or UBL document into a GOBL envelope
func ConvertToGOBL(d []byte) (*gobl.Envelope, error) {
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

	return env, nil
}

// ConvertToXRechnungCII converts a GOBL envelope into an XRechnung document
func ConvertToXRechnungCII(env *gobl.Envelope) ([]byte, error) {
	doc, err := cii.ToCII(env)
	if err != nil {
		return nil, fmt.Errorf("building XRechnung document: %w", err)
	}

	doc.ExchangedContext.GuidelineContext.ID = xRechnungGuideline
	doc.ExchangedContext.BusinessContext = &document.ExchangedContextParameter{
		ID: xRechnungProfile,
	}

	o, err := doc.Bytes()
	if err != nil {
		return nil, fmt.Errorf("generating XRechnung xml: %w", err)
	}

	return o, nil
}

// ConvertToXRechnungUBL converts a GOBL envelope into an XRechnung document
func ConvertToXRechnungUBL(env *gobl.Envelope) ([]byte, error) {
	doc, err := ubl.ToUBL(env)
	if err != nil {
		return nil, fmt.Errorf("building XRechnung document: %w", err)
	}

	doc.CustomizationID = xRechnungGuideline
	doc.ProfileID = xRechnungProfile

	o, err := doc.Bytes()
	if err != nil {
		return nil, fmt.Errorf("generating XRechnung xml: %w", err)
	}

	return o, nil
}

// ConvertToZUGFeRD converts a GOBL envelope into a ZUGFeRD document
func ConvertToZUGFeRD(env *gobl.Envelope) ([]byte, error) {
	doc, err := cii.ToCII(env)
	if err != nil {
		return nil, fmt.Errorf("building ZUGFeRD document: %w", err)
	}

	doc.ExchangedContext.GuidelineContext.ID = zugferdGuideline

	o, err := doc.Bytes()
	if err != nil {
		return nil, fmt.Errorf("generating ZUGFeRD xml: %w", err)
	}

	return o, nil
}

// ConvertToFacturX converts a GOBL envelope into a Factur-X document
func ConvertToFacturX(env *gobl.Envelope) ([]byte, error) {
	doc, err := cii.ToCII(env)
	if err != nil {
		return nil, fmt.Errorf("building Factur-X document: %w", err)
	}

	doc.ExchangedContext.GuidelineContext.ID = facturXGuideline

	o, err := doc.Bytes()
	if err != nil {
		return nil, fmt.Errorf("generating Factur-X xml: %w", err)
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
