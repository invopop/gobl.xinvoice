# gobl.xinvoice

GOBL conversion into and from Factur-X (FR) and XRechnung/ZUGFeRD (DE) formats.

[![codecov](https://codecov.io/gh/invopop/gobl.xinvoice/graph/badge.svg?token=TMW8MWSZ9P)](https://codecov.io/gh/invopop/gobl.xinvoice)

Copyright [Invopop Ltd.](https://invopop.com) 2024. Released publicly under the [Apache License Version 2.0](LICENSE). For commercial licenses please contact the [dev team at invopop](mailto:dev@invopop.com). In order to accept contributions to this library we will require transferring copyrights to Invopop Ltd.

## Usage

### Go Package

Usage of the XInvoice conversion library is quite straight forward. You must first have a GOBL Envelope including an invoice ready to convert. There are some samples here in the test/data directory. You must then choose the conversion direction and output format. Supported formats:

- From GOBL:
    - To X-Rechnung (CII syntax)
    - To X-Rechnung (UBL syntax)
    - To ZUGFeRD (CII syntax)
    - To Factur-X (CII syntax)
- To GOBL:
    - From any of the above formats

The library uses the [gobl.cii](https://github.com/invopop/gobl.cii) and [gobl.ubl](https://github.com/invopop/gobl.ubl) libraries to perform the conversion.

Example of converting a GOBL invoice to XRechnung CII:
```go
package main

import (
    "os"

    "github.com/invopop/gobl"
    xinvoice "github.com/invopop/gobl.xinvoice"
)

func main {
    data, err := os.ReadFile("./test/data/invoice-de-de.json")
    if err != nil {
        panic(err)
    }

    env := new(gobl.Envelope)
    if err := json.Unmarshal(data, env); err != nil {
        panic(err)
    }

    doc, err := xinvoice.ConvertToXRechnungCII(env)
    if err != nil {
        panic(err)
    }

    err = os.WriteFile("invoice-de-de-xrechnung.xml", doc, 0644)
    if err != nil {
        panic(err)
    }
}
```
Example of converting a X-Invoice file (any format) to GOBL:
```go
package main

import (
    "os"

    "github.com/invopop/gobl"
    xinvoice "github.com/invopop/gobl.xinvoice"
)

func main {
    data, err := os.ReadFile("./test/data/invoice-xrechnung-cii.xml")
    if err != nil {
        panic(err)
    }

    env, err := xinvoice.ConvertToGOBL(data)
    if err != nil {
        panic(err)
    }

    output, err := json.MarshalIndent(env, "", "  ")
    if err != nil {
        panic(err)
    }
}
```


### Command Line

The GOBL to XInvoice tool also includes a command line helper. You can it install manually in your Go environment with:

```bash
go install ./cmd/gobl.xinvoice
```

Usage is very straightforward, with the only flag being the preferred output format: `xrechnung-cii`, `xrechnung-ubl`, `zugferd` or `facturx`, to be used in case of conversion from GOBL. It is also possible to save the output to a file by adding a path as its second argument:

```bash
# XInvoice to GOBL
gobl.xinvoice convert ./test/data/invoice-de-de.xml

# XInvoice to GOBL (with output file)
gobl.xinvoice convert ./test/data/invoice-de-de.xml ./test/data/out/invoice-de-de-gobl.json

# GOBL to XRechnung CII
gobl.xinvoice convert ./test/data/invoice-de-de.json --format xrechnung-cii

# GOBL to Factur-X CII
gobl.xinvoice convert ./test/data/invoice-de-de.json --format facturx
```

## Testing

### testify

The library uses testify for testing. To run the tests you can use the command:
```bash
go test
```

## Development

### Limitations

There are some limitations in the current conversion process to GOBL.

#### CII

1. GOBL does not currently support additional embedded documents, so the AdditionalReferencedDocument field (BG-24 in EN 16931) is not supported and lost in the conversion.
2. Payment advances do not include their own tax rate, they use the global tax rate of the invoice.
3. The field TotalPrepaidAmount (BT-113) in CII is not directly mapped to GOBL, so payment advances must be included in the SpecifiedAdvancePayment field in CII, or they will be lost in conversion.
4. The fields ReceivableSpecifiedTradeAccountingAccount (BT-133) and DesignatedProductClassification (BT-158) are added as a note to the line, with the type code as the key.

#### UBL

1. GOBL does not currently support additional embedded documents, so the AdditionalReferencedDocument field (BG-24 in EN 16931) is not supported and lost in the conversion.
2. GOBL only supports a single period in the ordering, so only the first InvoicePeriod (BG-14) in the UBL is taken.
3. Fields ProfileID (BT-23) and CustomizationID (BT-24) in UBL are not supported and lost in the conversion.
4. The AccountingCost (BT-19, BT-133) fields are added as notes.
5. Payment advances do not include their own tax rate, they use the global tax rate of the invoice.

### EN 16931 Compliance

This tool is compliant with the EN 16931 standard, which defines the semantic data model for electronic invoices in the European Union. X-Rechnung and Factur-X are in turn extensions of this standard.

### XRechnung

Useful links:

- [XStandard Einkauf](https://xeinkauf.de/): XRechnung organization.
- [eInvoicing in Germany](https://ec.europa.eu/digital-building-blocks/sites/display/DIGITAL/eInvoicing+in+Germany)

Specifications:

- [XRechnung 3.0.2](https://xeinkauf.de/app/uploads/2024/07/302-XRechnung-2024-06-20.pdf): Released 20/06/2024
- [English Summary of XRechnung 3.0.2](https://xeinkauf.de/app/uploads/2024/10/XRechnung-EnglishSummary-v302.pdf)

Authentication:

- [Federal Central Invoice Submission Portal (ZRE)](https://xrechnung.bund.de/prod/authenticate.do)
- [Online Access Act-compliant Invoice Submission Portal](https://xrechnung-bdr.de/edi/auth/login)

### ZUGFeRD

Useful links:

- [Forum elektronische Rechnung Deutschland](https://www.ferd-net.de/): ZUGFeRD organization.

Specifications:

- [ZUGFeRD 2.3.2](https://www.ferd-net.de/standards/zugferd-2.3.2/zugferd-2.3.2.html): Released 13/11/2024

### Factur-X

Useful links:

- [Le Forum National de la Facture Electronique](https://fnfe-mpe.org/): Factur-X organization.

Specifications:

- [Factur-X 1.07.2](https://fnfe-mpe.org/factur-x/): Released 13/11/2024
