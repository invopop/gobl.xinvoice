package xinvoice_test

import (
	"testing"

	gtox "github.com/invopop/gobl.xinvoice/internal/gtox"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestFindTaxCode(t *testing.T) {
	t.Run("should return correct tax category", func(t *testing.T) {
		taxCode := gtox.FindTaxCode(tax.RateStandard)

		assert.Equal(t, gtox.StandardSalesTax, taxCode)
	})

	t.Run("should return zero tax category", func(t *testing.T) {
		taxCode := gtox.FindTaxCode(tax.RateZero)

		assert.Equal(t, gtox.ZeroRatedGoodsTax, taxCode)
	})

	t.Run("should return zero tax category", func(t *testing.T) {
		taxCode := gtox.FindTaxCode(tax.RateExempt)

		assert.Equal(t, gtox.TaxExempt, taxCode)
	})
}
