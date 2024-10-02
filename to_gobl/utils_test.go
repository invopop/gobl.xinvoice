package to_gobl_test

import (
	"testing"

	"github.com/invopop/gobl.xinvoice/to_gobl"
	"github.com/invopop/gobl/cbc"
	"github.com/stretchr/testify/assert"
)

func TestParseDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Valid date", "20230515", "2023-05-15"},
		{"Invalid date", "20231345", "0001-01-01"},
		{"Empty string", "", "0001-01-01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := to_gobl.ParseDate(tt.input)
			assert.Equal(t, tt.expected, result.String())
		})
	}
}

func TestFindTaxKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Standard sales tax", "S", "standard"},
		{"Zero rated goods tax", "Z", "zero"},
		{"Tax exempt", "E", "exempt"},
		{"Unknown tax type", "X", "standard"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := to_gobl.FindTaxKey(tt.input)
			assert.Equal(t, tt.expected, string(result))
		})
	}
}

func TestTypeCodeParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Standard invoice", "380", "standard"},
		{"Credit note", "381", "credit-note"},
		{"Corrective invoice", "384", "corrective"},
		{"Proforma invoice", "325", "proforma"},
		{"Debit note", "383", "debit-note"},
		{"Unknown type code", "999", "other"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := to_gobl.TypeCodeParse(tt.input)
			assert.Equal(t, tt.expected, string(result))
		})
	}
}

func TestUnitFromUNECE(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Known UNECE code", "EA", "each"},
		{"Unknown UNECE code", "XYZ", "XYZ"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := to_gobl.UnitFromUNECE(cbc.Code(tt.input))
			assert.Equal(t, tt.expected, string(result))
		})
	}
}
