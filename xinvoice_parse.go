package xinvoice

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"

	"github.com/invopop/gobl"
	cii "github.com/invopop/gobl.cii"
	ubl "github.com/invopop/gobl.ubl"
	"github.com/invopop/gobl/addons/de/xrechnung"
	"github.com/invopop/gobl/addons/de/zugferd"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
)

// Syntax identifies which of the two German XML syntaxes a parsed
// document uses.
type Syntax string

// The two syntaxes German invoices arrive in.
const (
	SyntaxUBL Syntax = "ubl"
	SyntaxCII Syntax = "cii"
)

// ErrUnknownDocument is returned when the data is neither a UBL nor a
// CII invoice.
var ErrUnknownDocument = errors.New("not a recognized UBL or CII invoice")

// Parsed carries the outcome of parsing a German XML document.
type Parsed struct {
	// Envelope contains the GOBL invoice.
	Envelope *gobl.Envelope
	// Attachments are the binary attachments embedded in the XML.
	Attachments []BinaryAttachment
	// Syntax the document was written in.
	Syntax Syntax
}

// ParseOption configures parsing.
type ParseOption func(*parseOptions)

type parseOptions struct {
	from, to cbc.URI
}

// WithRouting records the transport direction on the parsed envelope's
// header. Both values must be set for the option to take effect.
func WithRouting(from, to cbc.URI) ParseOption {
	return func(o *parseOptions) {
		o.from = from
		o.to = to
	}
}

// Parse reads a German XML invoice in either syntax and converts it to
// a GOBL envelope. The syntax is decided by the document's root
// namespace.
func Parse(data []byte, opts ...ParseOption) (*Parsed, error) {
	o := new(parseOptions)
	for _, opt := range opts {
		opt(o)
	}

	ns, err := rootNamespace(data)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnknownDocument, err)
	}

	switch ns {
	case ubl.NamespaceUBLInvoice, ubl.NamespaceUBLCreditNote:
		return parseUBL(data, o)
	case cii.NamespaceRSM:
		return parseCII(data, o)
	default:
		return nil, fmt.Errorf("%w: unexpected root namespace %q", ErrUnknownDocument, ns)
	}
}

// parseUBL covers XRechnung UBL invoices and credit notes.
func parseUBL(data []byte, o *parseOptions) (*Parsed, error) {
	doc, err := ubl.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("parsing UBL document: %w", err)
	}
	in, ok := doc.(*ubl.Invoice)
	if !ok {
		return nil, fmt.Errorf("unsupported UBL document type %T", doc)
	}
	// Pin the plain EN 16931 format so the import does not depend
	// on gobl.ubl's own format detection.
	opts := []ubl.Option{ubl.WithFormat(ubl.FormatEN16931)}
	if o.from != "" && o.to != "" {
		opts = append(opts, ubl.WithRouting(o.from, o.to))
	}
	env, err := ubl.Import(in, opts...)
	if err != nil {
		return nil, fmt.Errorf("converting UBL document to GOBL: %w", err)
	}
	// Stamp the German addon on documents that identify as XRechnung.
	if in.CustomizationID == CustomizationIDXRechnung {
		if err := stampAddons(env, []cbc.Key{xrechnung.V3}); err != nil {
			return nil, fmt.Errorf("applying German addon to parsed document: %w", err)
		}
	}
	return &Parsed{
		Envelope:    env,
		Syntax:      SyntaxUBL,
		Attachments: ublAttachments(in),
	}, nil
}

// parseCII covers XRechnung CII and the XML inside ZUGFeRD PDFs.
func parseCII(data []byte, o *parseOptions) (*Parsed, error) {
	wire, err := cii.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("parsing CII document: %w", err)
	}
	in, ok := wire.(*cii.Invoice)
	if !ok {
		return nil, fmt.Errorf("unsupported CII document type %T", wire)
	}
	var opts []cii.Option
	if o.from != "" && o.to != "" {
		opts = append(opts, cii.WithRouting(o.from, o.to))
	}
	env, err := ciiImport(in, opts)
	if err != nil {
		return nil, fmt.Errorf("converting CII document to GOBL: %w", err)
	}
	// Stamp the German addon that matches the document's guideline. The
	// generic EN 16931 guideline also names non-German documents, so it
	// is left to gobl.cii's own detection.
	if addons := ciiAddonsFor(in); addons != nil {
		if err := stampAddons(env, addons); err != nil {
			return nil, fmt.Errorf("applying German addon to parsed document: %w", err)
		}
	}
	return &Parsed{
		Envelope:    env,
		Syntax:      SyntaxCII,
		Attachments: ciiAttachments(in),
	}, nil
}

// ublAttachments returns the files embedded in the UBL document.
func ublAttachments(in *ubl.Invoice) []BinaryAttachment {
	var atts []BinaryAttachment
	for _, ba := range in.ExtractBinaryAttachments() {
		atts = append(atts, BinaryAttachment{
			ID:          ba.ID,
			Description: ba.Description,
			Data:        ba.Data,
			MimeCode:    ba.MimeCode,
			Filename:    ba.Filename,
		})
	}
	return atts
}

// ciiAttachments returns the files embedded in the CII document.
func ciiAttachments(in *cii.Invoice) []BinaryAttachment {
	var atts []BinaryAttachment
	for _, ba := range in.ExtractBinaryAttachments() {
		atts = append(atts, BinaryAttachment{
			ID:          ba.ID,
			Description: ba.Description,
			Data:        ba.Data,
			MimeCode:    ba.MimeCode,
			Filename:    ba.Filename,
		})
	}
	return atts
}

// ciiAddonsFor returns the GOBL addons matching the CII document's
// guideline ID, or nil. Only guideline IDs that name a German
// specification match: the generic EN 16931 guideline (which ZUGFeRD's
// EN 16931 profile shares with plain EN 16931 documents) is left to
// gobl.cii's own detection.
func ciiAddonsFor(in *cii.Invoice) []cbc.Key {
	if in.ExchangedContext == nil || in.ExchangedContext.GuidelineContext == nil {
		return nil
	}
	switch in.ExchangedContext.GuidelineContext.ID {
	case CustomizationIDXRechnung:
		return []cbc.Key{xrechnung.V3}
	case GuidelineIDZUGFeRDBasic, GuidelineIDZUGFeRDExtended,
		guidelineIDZUGFeRD20Basic, guidelineIDZUGFeRD20Extended:
		return []cbc.Key{zugferd.V2}
	}
	return nil
}

// ZUGFeRD 2.0 wrote its own guideline IDs. ZUGFeRD 2.1 replaced them with
// the Factur-X 1.0 IDs the formats write today. Documents from 2.0
// software are still in circulation, so parsing recognizes both.
const (
	guidelineIDZUGFeRD20Basic    = GuidelineIDEN16931 + "#compliant#urn:zugferd.de:2p0:basic"
	guidelineIDZUGFeRD20Extended = GuidelineIDEN16931 + "#conformant#urn:zugferd.de:2p0:extended"
)

// stampAddons adds the missing addons to the envelope's invoice and
// recalculates. Parsing does not validate: a received document that
// violates the addon's rules must still produce an envelope.
func stampAddons(env *gobl.Envelope, required []cbc.Key) error {
	inv, ok := env.Extract().(*bill.Invoice)
	if !ok {
		return nil
	}
	var missing []cbc.Key
	existing := inv.GetAddons()
	for _, a := range required {
		if !a.In(existing...) {
			missing = append(missing, a)
		}
	}
	if len(missing) == 0 {
		return nil
	}
	inv.SetAddons(append(existing, missing...)...)
	return env.Calculate()
}

// ciiImport wraps gobl.cii's Import, turning its panics on structurally
// incomplete documents into errors.
func ciiImport(in *cii.Invoice, opts []cii.Option) (env *gobl.Envelope, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("incomplete document: %v", r)
		}
	}()
	return cii.Import(in, opts...)
}

// rootNamespace returns the namespace of the document's root element.
func rootNamespace(data []byte) (string, error) {
	dc := xml.NewDecoder(bytes.NewReader(data))
	for {
		tk, err := dc.Token()
		if err == io.EOF {
			return "", errors.New("no root element found")
		}
		if err != nil {
			return "", fmt.Errorf("reading XML: %w", err)
		}
		if t, ok := tk.(xml.StartElement); ok {
			return t.Name.Space, nil
		}
	}
}
