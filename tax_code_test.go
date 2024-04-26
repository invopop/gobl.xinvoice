package xinvoice_test

import (
	"testing"

	xinvoice "github.com/invopop/gobl.xinvoice"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestFindTaxCode(t *testing.T) {
	t.Run("should return correct tax category", func(t *testing.T) {
		taxCode := xinvoice.FindTaxCode(tax.RateStandard)

		assert.Equal(t, xinvoice.StandardSalesTax, taxCode)
	})

	t.Run("should return zero tax category", func(t *testing.T) {
		taxCode := xinvoice.FindTaxCode(tax.RateZero)

		assert.Equal(t, xinvoice.ZeroRatedGoodsTax, taxCode)
	})

	t.Run("should return zero tax category", func(t *testing.T) {
		taxCode := xinvoice.FindTaxCode(tax.RateExempt)

		assert.Equal(t, xinvoice.TaxExempt, taxCode)
	})
}
