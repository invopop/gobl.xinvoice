// Package xinvoice helps convert GOBL into XRechnung and Factur-X documents and vice versa.
package xinvoice

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/invopop/gobl"
	cii "github.com/invopop/gobl.cii"
	ubl "github.com/invopop/gobl.ubl"
)

const (
	ciiHeader = "CrossIndustryInvoice"
	ublHeader = "Invoice"
)

// ConvertToGOBL converts an XRechnung, Factur-X or UBL document into a GOBL envelope
func ConvertToGOBL(d []byte) (*gobl.Envelope, error) {
	r, err := extractRootName(d)
	if err != nil {
		return nil, fmt.Errorf("extracting root name: %w", err)
	}

	var env *gobl.Envelope
	if r == ciiHeader {
		env, err = cii.ParseInvoice(d)
		if err != nil {
			return nil, fmt.Errorf("converting CII to GOBL: %w", err)
		}
	} else if r == ublHeader {
		env, err = ubl.ParseInvoice(d)
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
	doc, err := cii.ConvertInvoice(env, cii.WithContext(cii.ContextXRechnung))
	if err != nil {
		return nil, fmt.Errorf("convert invoice: %w", err)
	}
	return doc.Bytes()
}

// ConvertToXRechnungUBL converts a GOBL envelope into an XRechnung document
func ConvertToXRechnungUBL(env *gobl.Envelope) ([]byte, error) {
	doc, err := ubl.ConvertInvoice(env, ubl.WithContext(ubl.ContextXRechnung))
	if err != nil {
		return nil, fmt.Errorf("convert invoice: %w", err)
	}
	return doc.Bytes()
}

// ConvertToZUGFeRD converts a GOBL envelope into a ZUGFeRD document
func ConvertToZUGFeRD(env *gobl.Envelope) ([]byte, error) {
	doc, err := cii.ConvertInvoice(env, cii.WithContext(cii.ContextZUGFeRD))
	if err != nil {
		return nil, fmt.Errorf("convert invoice: %w", err)
	}
	return doc.Bytes()
}

// ConvertToFacturX converts a GOBL envelope into a Factur-X document
func ConvertToFacturX(env *gobl.Envelope) ([]byte, error) {
	doc, err := cii.ConvertInvoice(env, cii.WithContext(cii.ContextFacturX))
	if err != nil {
		return nil, fmt.Errorf("convert invoice: %w", err)
	}
	return doc.Bytes()
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
