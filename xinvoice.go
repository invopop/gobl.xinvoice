// Package xinvoice helps convert GOBL into XRechnung and Factur-X documents.
package xinvoice

import (
	"encoding/xml"

	"github.com/invopop/gobl"
)

// Document is a pseudo-model for containing the XML document being created
type Document struct {
	XMLName xml.Name `xml:"rsm:CrossIndustryInvoice"`
}

// NewDocument converts a GOBL envelope into a XRechnung and Factur-X document
func NewDocument(_ *gobl.Envelope) (*Document, error) {
	doc := Document{}
	return &doc, nil
}

// Bytes returns the XML representation of the document in bytes
func (d *Document) Bytes() ([]byte, error) {
	bytes, err := xml.MarshalIndent(d, "", "  ")
	if err != nil {
		return nil, err
	}

	return append([]byte(xml.Header), bytes...), nil
}
