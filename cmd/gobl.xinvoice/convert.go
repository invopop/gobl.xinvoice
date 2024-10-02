package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/invopop/gobl"
	to_gobl "github.com/invopop/gobl.xinvoice/to_gobl"
	xinvoice "github.com/invopop/gobl.xinvoice/xinvoice"
	"github.com/spf13/cobra"
)

type convertOpts struct {
	*rootOpts
}

func convert(o *rootOpts) *convertOpts {
	return &convertOpts{rootOpts: o}
}

func (c *convertOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert <infile> <outfile>",
		Short: "Convert a GOBL JSON into a XRechnung & Factur-X XML document",
		RunE:  c.runE,
	}

	return cmd
}

func (c *convertOpts) runE(cmd *cobra.Command, args []string) error {
	if len(args) == 0 || len(args) > 2 {
		return fmt.Errorf("expected one or two arguments, the command usage is `gobl.xinvoice convert <infile> [outfile]`")
	}

	input, err := openInput(cmd, args)
	if err != nil {
		return err
	}
	defer input.Close() // nolint:errcheck

	out, err := c.openOutput(cmd, args)
	if err != nil {
		return err
	}
	defer out.Close() // nolint:errcheck

	inData, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	// Check if input is JSON or XML
	isJSON := json.Valid(inData)

	var outputData []byte

	if isJSON {
		env := new(gobl.Envelope)
		if err := json.Unmarshal(inData, env); err != nil {
			return fmt.Errorf("parsing input as GOBL Envelope: %w", err)
		}
		doc, err := xinvoice.NewDocument(env)
		if err != nil {
			return fmt.Errorf("building XRechnung and Factur-X document: %w", err)
		}

		outputData, err = doc.Bytes()
		if err != nil {
			return fmt.Errorf("generating XRechnung and Factur-X xml: %w", err)
		}
	} else {
		// Assume XML if not JSON
		doc := new(to_gobl.XMLDoc)
		if err := xml.Unmarshal(inData, doc); err != nil {
			return fmt.Errorf("parsing input document: %w", err)
		}

		env, err := to_gobl.NewDocumentGOBL(doc)
		if err != nil {
			return fmt.Errorf("building GOBL envelope: %w", err)
		}

		outputData, err = json.MarshalIndent(env, "", "  ")
		if err != nil {
			return fmt.Errorf("generating JSON output: %w", err)
		}
	}

	if _, err = out.Write(outputData); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	return nil
}
