package xinvoice

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/de"
	"github.com/invopop/gobl/tax"
)

// Parses the XML information for a Party object
func ParseParty(party *TradeParty) *org.Party {
	p := &org.Party{
		Name: party.Name,
	}

	if party.DefinedTradeContact != nil && party.DefinedTradeContact.PersonName != nil {
		p.People = []*org.Person{
			{
				Name: &org.Name{
					Given: *party.DefinedTradeContact.PersonName,
				},
			},
		}
	}

	if party.PostalTradeAddress != nil {
		p.Addresses = []*org.Address{
			{
				Street:   party.PostalTradeAddress.LineOne,
				Locality: party.PostalTradeAddress.CityName,
				Code:     party.PostalTradeAddress.PostcodeCode,
				Country:  party.PostalTradeAddress.CountryID,
			},
		}
	}

	if party.DefinedTradeContact != nil && party.DefinedTradeContact.TelephoneUniversalCommunication != nil {
		p.Telephones = []*org.Telephone{
			{
				Number: party.DefinedTradeContact.TelephoneUniversalCommunication.CompleteNumber,
			},
		}
	}

	if party.DefinedTradeContact != nil && party.DefinedTradeContact.EmailURIUniversalCommunication != nil {
		p.Emails = []*org.Email{
			{
				Address: party.DefinedTradeContact.EmailURIUniversalCommunication.URIID,
			},
		}
	}
	if party.SpecifiedTaxRegistration != nil {
		for _, taxReg := range *party.SpecifiedTaxRegistration {
			if taxReg.ID != nil {
				switch taxReg.ID.SchemeID {
				case "VA":
					p.TaxID = &tax.Identity{
						Country: l10n.TaxCountryCode(party.PostalTradeAddress.CountryID),
						Code:    cbc.Code(regexp.MustCompile(`\D`).ReplaceAllString(taxReg.ID.Value, "")),
					}
				case "FC":
					identity := &org.Identity{
						Key:  de.IdentityKeyTaxNumber,
						Code: cbc.Code(taxReg.ID.Value),
					}
					p.Identities = append(p.Identities, identity)
				}
			}
		}
	}

	return p
}
