package xinvoice

import (
	"time"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
)

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
		ID:       inv.Series + "-" + inv.Code,
		TypeCode: invoiceTypeCode(inv.Type),
		IssueDate: &Date{
			Date:   formatIssueDate(inv.IssueDate),
			Format: "102",
		},
	}
}

func formatIssueDate(date cal.Date) string {
	if date.IsZero() {
		return ""
	}
	t := time.Date(date.Year, date.Month, date.Day, 0, 0, 0, 0, time.UTC)
	return t.Format("20060102")
}

// For German suppliers, the element "Invoice type code" (BT-3) should only contain the
// following values from code list UNTDID 1001:
// - 326 (Partial invoice)
// - 380 (Commercial invoice)
// - 384 (Corrected invoice)
// - 389 (Self-billed invoice)
// - 381 (Credit note)
func invoiceTypeCode(t cbc.Key) string {
	hash := map[cbc.Key]string{
		bill.InvoiceTypeStandard:   "380",
		bill.InvoiceTypeCorrective: "384",
		bill.InvoiceTypeCreditNote: "381",
	}
	return hash[t]
}
