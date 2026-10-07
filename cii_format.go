package xinvoice

import (
	"github.com/invopop/gobl"
	cii "github.com/invopop/gobl.cii"
	"github.com/invopop/gobl/addons/de/xrechnung"
	"github.com/invopop/gobl/addons/de/zugferd"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/schema"
)

// Keys of the German CII formats in the GOBL convert register.
const (
	KeyXRechnungCII       cbc.Key = "cii+de-xrechnung-v3"
	KeyZUGFeRDCII         cbc.Key = "cii+de-zugferd-v2"
	KeyZUGFeRDBasicCII    cbc.Key = "cii+de-zugferd-v2+basic"
	KeyZUGFeRDExtendedCII cbc.Key = "cii+de-zugferd-v2+extended"
)

// VESIDs of the CII formats.
const (
	vesIDXRechnungCII    = "de.xrechnung:cii:3.0.2"
	vesIDZUGFeRD         = "de.zugferd:en16931:2.5.2"
	vesIDZUGFeRDBasic    = "de.zugferd:basic:2.5.2"
	vesIDZUGFeRDExtended = "de.zugferd:extended:2.5.2"
)

var ciiSchemas = []schema.ID{schema.Lookup(bill.Invoice{})}

// CIIFormatXRechnung is the gobl.cii format for XRechnung 3.0 in CII syntax.
var CIIFormatXRechnung = cii.Format{
	Key:         KeyXRechnungCII,
	Name:        i18n.NewString("CII XRechnung 3"),
	Countries:   []l10n.Code{l10n.DE},
	Schemas:     ciiSchemas,
	GuidelineID: CustomizationIDXRechnung,
	BusinessID:  ProfileIDPeppolBilling,
	Version:     cii.VersionD16B,
	Addons:      []cbc.Key{xrechnung.V3},
	VESID:       vesIDXRechnungCII,
}

// CIIFormatZUGFeRD is the gobl.cii format for the ZUGFeRD EN 16931
// (COMFORT) profile.
var CIIFormatZUGFeRD = cii.Format{
	Key:         KeyZUGFeRDCII,
	Name:        i18n.NewString("CII ZUGFeRD EN 16931"),
	Countries:   []l10n.Code{l10n.DE},
	Schemas:     ciiSchemas,
	GuidelineID: cii.GuidelineIDEN16931,
	Version:     cii.VersionD16B,
	Addons:      []cbc.Key{zugferd.V2},
	VESID:       vesIDZUGFeRD,
	ExportFuncs: []cii.ExportFunc{exportZUGFeRD},
}

// CIIFormatZUGFeRDBasic is the gobl.cii format for the ZUGFeRD BASIC profile,
// a CIUS of EN 16931.
var CIIFormatZUGFeRDBasic = cii.Format{
	Key:         KeyZUGFeRDBasicCII,
	Name:        i18n.NewString("CII ZUGFeRD BASIC"),
	Countries:   []l10n.Code{l10n.DE},
	Schemas:     ciiSchemas,
	GuidelineID: cii.GuidelineIDEN16931 + "#compliant#urn:zugferd.de:2p0:basic",
	Version:     cii.VersionD16B,
	Addons:      []cbc.Key{zugferd.V2},
	VESID:       vesIDZUGFeRDBasic,
	ExportFuncs: []cii.ExportFunc{exportZUGFeRD},
}

// CIIFormatZUGFeRDExtended is the gobl.cii format for the ZUGFeRD EXTENDED
// profile, which allows sub-lines.
var CIIFormatZUGFeRDExtended = cii.Format{
	Key:         KeyZUGFeRDExtendedCII,
	Name:        i18n.NewString("CII ZUGFeRD EXTENDED"),
	Countries:   []l10n.Code{l10n.DE},
	Schemas:     ciiSchemas,
	GuidelineID: cii.GuidelineIDEN16931 + "#conformant#urn:zugferd.de:2p0:extended",
	Version:     cii.VersionD16B,
	Addons:      []cbc.Key{zugferd.V2},
	VESID:       vesIDZUGFeRDExtended,
	ExportFuncs: []cii.ExportFunc{exportZUGFeRD, exportSubLines},
}

func init() {
	cii.RegisterFormats(CIIFormatXRechnung, CIIFormatZUGFeRD, CIIFormatZUGFeRDBasic, CIIFormatZUGFeRDExtended)
}

// exportZUGFeRD applies the ZUGFeRD rules to an exported invoice.
func exportZUGFeRD(_ *cii.Format, env *gobl.Envelope, doc cii.Document) error {
	inv, ok := env.Extract().(*bill.Invoice)
	out, ok2 := doc.(*cii.Invoice)
	if ok && ok2 {
		zugferdShipTo(inv, out)
	}
	return nil
}

// exportSubLines writes line breakdowns as sub-lines, which the EXTENDED
// profile allows.
func exportSubLines(_ *cii.Format, env *gobl.Envelope, doc cii.Document) error {
	inv, ok := env.Extract().(*bill.Invoice)
	out, ok2 := doc.(*cii.Invoice)
	if ok && ok2 {
		cii.ExpandGroupLines(inv, out)
	}
	return nil
}

// zugferdShipTo fills the ship-to party from the customer for rule
// BR-FX-EN-04 when the invoice gives no delivery details, as happens with
// imported invoices. Prepayment invoices (386) are left alone.
func zugferdShipTo(inv *bill.Invoice, out *cii.Invoice) {
	if inv.Delivery != nil || inv.Tax.GetExt(untdid.ExtKeyDocumentType) == "386" {
		return
	}
	if inv.Customer == nil || len(inv.Customer.Addresses) == 0 {
		return
	}
	if out.Transaction == nil || out.Transaction.Delivery == nil {
		return
	}
	out.Transaction.Delivery.Receiver = cii.NewDeliveryParty(inv.Customer)
}
