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

func NewHeader(inv *bill.Invoice) *Header {
	id := inv.Series + "-" + inv.Code
	typeCode := invoiceTypeCode(inv.Type)
	date := formatIssueDate(inv.IssueDate)

	return &Header{
		ID:       id,
		TypeCode: typeCode,
		IssueDate: &Date{
			Date:   date,
			Format: "102",
		},
	}
}

func formatIssueDate(date cal.Date) string {
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
	hash := map[string]string{
		"standard":    "380",
		"corrective":  "384",
		"credit-note": "381",
	}
	return hash[t.String()]
}