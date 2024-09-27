package xinvoice

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

// IssueDateFormat is the issue date format in the form YYYYMMDD
const IssueDateFormat = "102"

// Header a collection of data for a Cross Industry Invoice Header that is exchanged between two or more parties in written, printed or electronic form.
type Header struct {
	ID           string `xml:"ram:ID"`
	TypeCode     string `xml:"ram:TypeCode"`
	IssueDate    *Date  `xml:"ram:IssueDateTime>udt:DateTimeString"`
	IncludedNote *Note  `xml:"ram:IncludedNote,omitempty"`
}

// Note defines note in the RAM structure
type Note struct {
	Content     string `xml:"ram:Content"`
	SubjectCode string `xml:"ram:SubjectCode"`
}

// NewHeader creates the ExchangedDocument part of a EN 16931 compliant invoice
func NewHeader(inv *bill.Invoice) *Header {
	return &Header{
		ID:       invoiceNumber(inv.Series, inv.Code),
		TypeCode: invoiceTypeCode(inv),
		IssueDate: &Date{
			Date:   formatIssueDate(inv.IssueDate),
			Format: IssueDateFormat,
		},
	}
}

func formatIssueDate(date cal.Date) string {
	if date.IsZero() {
		return ""
	}
	t := date.Time()
	return t.Format("20060102")
}

func invoiceNumber(series cbc.Code, code cbc.Code) string {
	if series == "" {
		return series.String()
	}
	return fmt.Sprintf("%s-%s", series, code)
}

// For German suppliers, the element "Invoice type code" (BT-3) should only contain the
// following values from code list UNTDID 1001:
// - 326 (Partial invoice)
// - 380 (Commercial invoice)
// - 384 (Corrected invoice)
// - 389 (Self-billed invoice)
// - 381 (Credit note)
func invoiceTypeCode(inv *bill.Invoice) string {
	if isSelfBilledInvoice(inv) {
		return "389"
	}
	hash := map[cbc.Key]string{
		bill.InvoiceTypeStandard:   "380",
		bill.InvoiceTypeCorrective: "384",
		bill.InvoiceTypeCreditNote: "381",
	}
	return hash[inv.Type]
}

func isSelfBilledInvoice(inv *bill.Invoice) bool {
	return inv.Type == bill.InvoiceTypeStandard && inv.HasTags(tax.TagSelfBilled)
}
