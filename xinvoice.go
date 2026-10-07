// Package xinvoice converts GOBL envelopes into the German electronic
// invoicing formats and back: XRechnung 3.0 in UBL and CII syntax, and
// ZUGFeRD's BASIC, EN 16931, and EXTENDED profiles.
package xinvoice

import (
	"errors"
	"fmt"
	"slices"

	"github.com/invopop/gobl"
	cii "github.com/invopop/gobl.cii"
	ubl "github.com/invopop/gobl.ubl"
	"github.com/invopop/gobl/addons/de/xrechnung"
	"github.com/invopop/gobl/addons/de/zugferd"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
)

// German document identifiers, written into the identity headers of the
// converted documents.
const (
	// CustomizationIDXRechnung identifies XRechnung 3.0 documents in
	// both syntaxes (BT-24).
	CustomizationIDXRechnung = "urn:cen.eu:en16931:2017#compliant#urn:xeinkauf.de:kosit:xrechnung_3.0"
	// GuidelineIDEN16931 identifies plain EN 16931 CII documents, the
	// guideline ZUGFeRD's EN 16931 profile uses.
	GuidelineIDEN16931 = "urn:cen.eu:en16931:2017"
	// GuidelineIDZUGFeRDBasic identifies ZUGFeRD's BASIC profile, a
	// restrictive CIUS of EN 16931 with its own reduced XSD.
	GuidelineIDZUGFeRDBasic = GuidelineIDEN16931 + "#compliant#urn:factur-x.eu:1p0:basic"
	// GuidelineIDZUGFeRDExtended identifies ZUGFeRD's EXTENDED profile,
	// a superset of EN 16931.
	GuidelineIDZUGFeRDExtended = GuidelineIDEN16931 + "#conformant#urn:factur-x.eu:1p0:extended"
	// ProfileIDPeppolBilling is the Peppol billing process identifier.
	ProfileIDPeppolBilling = "urn:fdc:peppol.eu:2017:poacc:billing:01:1.0"
)

// Format keys for the supported German document formats.
const (
	// FormatXRechnungUBL is XRechnung 3.0 in UBL syntax.
	FormatXRechnungUBL cbc.Key = "xrechnung-ubl-v3"
	// FormatXRechnungCII is XRechnung 3.0 in CII syntax.
	FormatXRechnungCII cbc.Key = "xrechnung-cii-v3"
	// FormatZUGFeRD is the XML side of ZUGFeRD's EN 16931 profile.
	// Readers locate the XML inside the PDF/A-3 by the exact file name
	// factur-x.xml.
	FormatZUGFeRD cbc.Key = "zugferd-v2"
	// FormatZUGFeRDBasic is the XML side of ZUGFeRD's BASIC profile. The
	// converter writes the same EN 16931 content under the BASIC
	// guideline; documents that carry data outside BASIC's reduced XSD
	// fail schematron validation instead of losing fields silently.
	FormatZUGFeRDBasic cbc.Key = "zugferd-basic-v2"
	// FormatZUGFeRDExtended is the XML side of ZUGFeRD's EXTENDED
	// profile, a superset of EN 16931.
	FormatZUGFeRDExtended cbc.Key = "zugferd-extended-v2"
)

// fileNameFacturX is the file name shared by all ZUGFeRD profiles.
// Readers locate the XML inside the PDF/A-3 by this exact name.
const fileNameFacturX = "factur-x.xml"

// ErrUnsupportedFormat is returned for an unknown format key.
var ErrUnsupportedFormat = errors.New("unsupported format")

// Format describes one supported German document format.
type Format struct {
	Key      cbc.Key
	Name     string
	FileName string

	syntax Syntax
	addons []cbc.Key
	// customizationID and profileID are the document identity headers.
	customizationID string
	profileID       string
	vesIDInvoice    string
	vesIDCreditNote string
}

// Addons returns a copy of the GOBL addon keys the format requires.
func (f *Format) Addons() []cbc.Key {
	return slices.Clone(f.addons)
}

// formats defines the supported German document formats.
var formats = []*Format{
	{
		Key:             FormatXRechnungUBL,
		Name:            "XRechnung UBL Invoice/CreditNote V3",
		FileName:        "xrechnung-ubl.xml",
		syntax:          SyntaxUBL,
		addons:          []cbc.Key{xrechnung.V3},
		customizationID: CustomizationIDXRechnung,
		profileID:       ProfileIDPeppolBilling,
		vesIDInvoice:    "de.xrechnung:ubl-invoice:3.0.2",
		vesIDCreditNote: "de.xrechnung:ubl-creditnote:3.0.2",
	},
	{
		Key:             FormatXRechnungCII,
		Name:            "XRechnung CII Invoice/CreditNote V3",
		FileName:        "xrechnung-cii.xml",
		syntax:          SyntaxCII,
		addons:          []cbc.Key{xrechnung.V3},
		customizationID: CustomizationIDXRechnung,
		profileID:       ProfileIDPeppolBilling,
		vesIDInvoice:    vesIDXRechnungCII,
		vesIDCreditNote: vesIDXRechnungCII,
	},
	{
		Key:             FormatZUGFeRD,
		Name:            "ZUGFeRD V2 (CII)",
		FileName:        fileNameFacturX,
		syntax:          SyntaxCII,
		addons:          []cbc.Key{zugferd.V2},
		customizationID: GuidelineIDEN16931,
		vesIDInvoice:    vesIDZUGFeRD,
		vesIDCreditNote: vesIDZUGFeRD,
	},
	{
		Key:             FormatZUGFeRDBasic,
		Name:            "ZUGFeRD V2 BASIC (CII)",
		FileName:        fileNameFacturX,
		syntax:          SyntaxCII,
		addons:          []cbc.Key{zugferd.V2},
		customizationID: GuidelineIDZUGFeRDBasic,
		vesIDInvoice:    vesIDZUGFeRDBasic,
		vesIDCreditNote: vesIDZUGFeRDBasic,
	},
	{
		Key:             FormatZUGFeRDExtended,
		Name:            "ZUGFeRD V2 EXTENDED (CII)",
		FileName:        fileNameFacturX,
		syntax:          SyntaxCII,
		addons:          []cbc.Key{zugferd.V2},
		customizationID: GuidelineIDZUGFeRDExtended,
		vesIDInvoice:    vesIDZUGFeRDExtended,
		vesIDCreditNote: vesIDZUGFeRDExtended,
	},
}

// Formats returns copies of the supported German document formats.
func Formats() []*Format {
	out := make([]*Format, len(formats))
	for i, f := range formats {
		c := *f
		out[i] = &c
	}
	return out
}

// FormatFor returns a copy of the format with the given key, or nil.
func FormatFor(key cbc.Key) *Format {
	for _, f := range formats {
		if f.Key == key {
			c := *f
			return &c
		}
	}
	return nil
}

// Document carries the serialized XML and its transmission metadata.
type Document struct {
	// Data is the XML document.
	Data []byte
	// Format is the format the document was generated in.
	Format *Format
	// Namespace and Element identify the XML root.
	Namespace string
	Element   string
	// CustomizationID and ProfileID identify the document flavor and
	// business process.
	CustomizationID string
	ProfileID       string
	// Version is the syntax version (UBL "2.1", CII "D16B").
	Version string
	// VESID is the validation rule set the document must satisfy.
	VESID string
}

// Option configures a conversion.
type Option func(*options)

type options struct {
	attachments []BinaryAttachment
}

// WithAttachment embeds a file inside the generated XML.
func WithAttachment(a BinaryAttachment) Option {
	return func(o *options) {
		o.attachments = append(o.attachments, a)
	}
}

// Convert converts a GOBL envelope into the given German format. The
// invoice gains the format's addon when missing, so the German rules
// run before any mapping.
func Convert(env *gobl.Envelope, format cbc.Key, opts ...Option) (*Document, error) {
	f := FormatFor(format)
	if f == nil {
		return nil, fmt.Errorf("%w: '%s'", ErrUnsupportedFormat, format)
	}
	o := new(options)
	for _, opt := range opts {
		opt(o)
	}

	if env == nil {
		return nil, errors.New("nil envelope")
	}
	inv, ok := env.Extract().(*bill.Invoice)
	if !ok {
		return nil, errors.New("envelope does not contain an invoice")
	}
	if err := ensureAddons(env, inv, f.Addons()); err != nil {
		return nil, err
	}

	if f.syntax == SyntaxUBL {
		return convertUBL(env, inv, f, o)
	}
	return convertCII(env, inv, f, o)
}

// ensureAddons adds the missing addons, recalculating so their
// normalizations run, and validates the envelope. Validation always
// runs: a caller can declare an addon and still violate its rules.
func ensureAddons(env *gobl.Envelope, inv *bill.Invoice, required []cbc.Key) error {
	var missing []cbc.Key
	existing := inv.GetAddons()
	for _, a := range required {
		if !a.In(existing...) {
			missing = append(missing, a)
		}
	}
	if len(missing) > 0 {
		inv.SetAddons(append(existing, missing...)...)
		if err := env.Calculate(); err != nil {
			return fmt.Errorf("calculating envelope with addons: %w", err)
		}
	}
	if err := env.Validate(); err != nil {
		return fmt.Errorf("validating envelope with addons: %w", err)
	}
	return nil
}

// buildUBL builds the plain EN 16931 UBL document and reworks it into
// the German format.
func buildUBL(env *gobl.Envelope, f *Format) (*ubl.Invoice, error) {
	out, err := ubl.ExportInvoice(env, ubl.WithFormat(ubl.FormatEN16931))
	if err != nil {
		return nil, fmt.Errorf("converting to UBL: %w", err)
	}
	f.applyUBL(out)
	return out, nil
}

// applyUBL writes the format's identity headers onto the UBL document.
func (f *Format) applyUBL(out *ubl.Invoice) {
	out.CustomizationID = f.customizationID
	out.ProfileID = &ubl.IDType{Value: f.profileID}
}

// buildCII builds the plain EN 16931 CII document and reworks it into
// the German format.
func buildCII(env *gobl.Envelope, f *Format) (*cii.Invoice, error) {
	out, err := cii.ExportInvoice(env, cii.WithFormat(cii.FormatEN16931))
	if err != nil {
		return nil, fmt.Errorf("converting to CII: %w", err)
	}
	if inv, ok := env.Extract().(*bill.Invoice); ok && zugferd.V2.In(inv.GetAddons()...) {
		zugferdShipTo(inv, out)
	}
	f.applyCII(out)
	return out, nil
}

// applyCII writes the format's identity headers onto the CII document.
func (f *Format) applyCII(out *cii.Invoice) {
	out.ExchangedContext.GuidelineContext.ID = f.customizationID
	if f.profileID != "" {
		out.ExchangedContext.BusinessContext = &cii.ExchangedContextParameter{ID: f.profileID}
	}
}

// convertUBL produces the German UBL document with its transmission
// metadata.
func convertUBL(env *gobl.Envelope, inv *bill.Invoice, f *Format, o *options) (*Document, error) {
	out, err := buildUBL(env, f)
	if err != nil {
		return nil, err
	}

	for _, a := range o.attachments {
		out.AddBinaryAttachment(ubl.BinaryAttachment{
			ID:          a.ID,
			Description: a.Description,
			Data:        a.Data,
			MimeCode:    a.MimeCode,
			Filename:    a.Filename,
		})
	}

	data, err := ubl.Encode(out)
	if err != nil {
		return nil, fmt.Errorf("serializing UBL document: %w", err)
	}

	return &Document{
		Data:            data,
		Format:          f,
		Namespace:       out.UBLNamespace,
		Element:         out.XMLName.Local,
		CustomizationID: f.customizationID,
		ProfileID:       f.profileID,
		Version:         ubl.Version,
		VESID:           f.vesID(inv),
	}, nil
}

// convertCII produces the German CII document with its transmission
// metadata.
func convertCII(env *gobl.Envelope, inv *bill.Invoice, f *Format, o *options) (*Document, error) {
	out, err := buildCII(env, f)
	if err != nil {
		return nil, err
	}

	for _, a := range o.attachments {
		out.AddBinaryAttachment(cii.BinaryAttachment{
			ID:          a.ID,
			Description: a.Description,
			Data:        a.Data,
			MimeCode:    a.MimeCode,
			Filename:    a.Filename,
		})
	}

	data, err := cii.Encode(out)
	if err != nil {
		return nil, fmt.Errorf("serializing CII document: %w", err)
	}

	return &Document{
		Data:            data,
		Format:          f,
		Namespace:       cii.NamespaceRSM,
		Element:         "CrossIndustryInvoice",
		CustomizationID: f.customizationID,
		ProfileID:       f.profileID,
		Version:         cii.VersionD16B,
		VESID:           f.vesID(inv),
	}, nil
}

// vesID returns the validation rule set for the invoice's document type.
func (f *Format) vesID(inv *bill.Invoice) string {
	if inv.GetType().Has(bill.InvoiceTypeCreditNote) {
		return f.vesIDCreditNote
	}
	return f.vesIDInvoice
}
