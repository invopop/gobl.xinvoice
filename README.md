# gobl.xinvoice

GOBL conversion into Factur-X (FR) and XRechnung/ZUGFeRD (DE) formats.

Copyright [Invopop Ltd.](https://invopop.com) 2023. Released publicly under the [Apache License Version 2.0](LICENSE). For commercial licenses please contact the [dev team at invopop](mailto:dev@invopop.com). In order to accept contributions to this library we will require transferring copyrights to Invopop Ltd.

## Usage

### Go Package

Usage of the XInvoice conversion library is quite straight forward. You must first have a GOBL Envelope including an invoice ready to convert. There are some samples here in the test/data directory.

```go
package main

import (
    "os"

    "github.com/invopop/gobl"
    xinvoice "github.com/invopop/gobl.xinvoice"
)

func main {
    data, _ := os.ReadFile("./test/data/invoice-de-de.json")

    env := new(gobl.Envelope)
    if err := json.Unmarshal(data, env); err != nil {
        panic(err)
    }

    // Prepare the CFDI document
    doc, err := xinvoice.NewDocument(env)
    if err != nil {
        panic(err)
    }

    // Create the XML output
    out, err := doc.Bytes()
    if err != nil {
        panic(err)
    }

    // TODO: do something with the output
}
```

### Command Line

The GOBL to XInvoice tool also includes a command line helper. You can it install manually in your Go environment with:

```bash
go install ./cmd/gobl.xinvoice
```

Usage is very straightforward:

```bash
gobl.xinvoice convert ./test/data/invoice-de-de.json
```

## Testing

### testify

The library uses testify for testing. To run the tests you can use the command:
```
go test
```

### Github action

There is a additional validator for XRechnung that is provided as a github action. You can use it to make sure that the created documents are valid. To run it go to:
```
https://github.com/invopop/gobl.xinvoice/actions/workflows/xrechnung_validator.yaml
```
Select a branch on which the file you want to test is on. Provide it with a file path to a file in the repository and click Run Workflow

## Development

### XRechnung

Useful links:

- [XStandard Einkauf - XRechnung](https://xeinkauf.de/xrechnung/)
- [eInvoicing in Germany](https://ec.europa.eu/digital-building-blocks/sites/display/DIGITAL/eInvoicing+in+Germany)

Specifications:

- [XRechnung 3.0.1 (DE)](https://xeinkauf.de/app/uploads/2023/09/301-XRechnung-2023-09-22.pdf)
- [English Summary of XRechnung 3.0.1](https://xeinkauf.de/app/uploads/2023/09/XRechnung-EnglishSummary-v301.pdf)
- [Document structure description](https://portal3.gefeg.com/invoice/tthome/index/617afdc4-623f-44e0-a05b-5b878840e508?page=1&useSelectedItemPosition=true)

Authentication:

- [Federal Central Invoice Submission Portal (ZRE)](https://xrechnung.bund.de/prod/authenticate.do)
- [Online Access Act-compliant Invoice Submission Portal](https://xrechnung-bdr.de/edi/auth/login)

